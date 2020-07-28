package service

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/panshul007/apideps/config"
	"github.com/panshul007/apideps/logger"
)

type DepLoader struct {
	log     logger.GenericLogger
	fileOps *FileOps
	repoMap map[string]*git.Repository
}

func NewDepLoader(log logger.GenericLogger) *DepLoader {
	return &DepLoader{
		log:     log,
		fileOps: NewFileOpsInstance(log),
		repoMap: make(map[string]*git.Repository),
	}
}

func (d *DepLoader) FetchDeps(deps *config.ApiDependencies) error {
	d.log.Infof("Fetching dependencies...")
	fs := memfs.New()

	for api, dep := range deps.Dependencies {
		d.log.Infof("processing %s...", api)
		if dep.Repo == "" {
			continue
		}
		repo, err := d.cloneRepo(dep.Repo, fs)
		if err != nil {
			return fmt.Errorf("error while cloning repo: %s, with error: %w", dep.Repo, err)
		}

		ref, err := repo.Head()
		if err != nil {
			return fmt.Errorf("error while fetching head ref of repo: %w", err)
		}
		fmt.Printf("head ref: %v\n", ref.Hash())

		commitRef, err := d.resolveCommitRef(dep, repo)
		if err != nil {
			return fmt.Errorf("error while resolving commit ref for api: %s, with error: %w", api, err)
		}

		commit, err := repo.CommitObject(commitRef)
		if err != nil {
			return fmt.Errorf("could not load commit object of api: %s, for commit: %s, error: %w", api, commitRef.String(), err)
		}

		tree, err := commit.Tree()
		if err != nil {
			return fmt.Errorf("could not load tree from commit: %s, error: %w", commitRef.String(), err)
		}

		//err = d.fileOps.listTreeFiles(tree)
		//if err != nil {
		//	return err
		//}

		err = d.fileOps.copyFolderPath(tree, dep.RepoFolder, dep.TargetPath)
		if err != nil {
			return fmt.Errorf("error while copying folder: %w", err)
		}
	}
	return nil
}

func (d *DepLoader) cloneRepo(repoPath string, fs billy.Filesystem) (*git.Repository, error) {
	// check if repo is present in cache
	repoCache, found := d.repoMap[repoPath]
	if found {
		return repoCache, nil
	}
	d.log.Infof("cloning repo: %s", repoPath)
	store := memory.NewStorage()
	agent, err := ssh.NewSSHAgentAuth("bitbucket-pipelines")
	if err != nil {
		return nil, fmt.Errorf("error while getting ssh agent for clone, error: %w", err)
	}
	repo, err := git.Clone(store, fs, &git.CloneOptions{
		URL:      repoPath,
		Progress: os.Stdout,
		Auth:     agent,
	})
	if err != nil {
		return nil, err
	}
	// put in cache for reuse
	d.repoMap[repoPath] = repo
	return repo, nil
}

func (d *DepLoader) resolveCommitRef(dep config.Dependency, repo *git.Repository) (plumbing.Hash, error) {
	// check if at-least one of commit hash or tag name is provided
	if dep.Tag == "" && dep.Commit == "" {
		return plumbing.ZeroHash, fmt.Errorf("no commit ref provided, tag or commit hash required")
	}
	if dep.Tag != "" {
		return d.resolveTag(dep.Tag, repo)
	}
	return plumbing.NewHash(dep.Commit), nil
}

func (d *DepLoader) resolveTag(tag string, repo *git.Repository) (plumbing.Hash, error) {
	ref, err := repo.Tag(tag)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("could not resolve tag name: %s, with error: %w", tag, err)
	}
	tagObject, err := repo.TagObject(ref.Hash())
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("could not resolve ref name: %s, of ref hash: %v to tag object with error: %w", tag, ref.Hash(), err)
	}
	d.log.Infof("resolved tag: %s to commit: %v", tag, tagObject.Target)
	return tagObject.Target, nil
}

func (d *DepLoader) listTags(repo *git.Repository) {
	tags, err := repo.TagObjects()
	if err != nil {
		d.log.Errorf("error while getting tag objects: %v", err)
	}
	err = tags.ForEach(processTag)
	if err != nil {
		d.log.Errorf("error while iterating tag objects: %v", err)
	}
}

func (d *DepLoader) listReferences(repo *git.Repository) error {
	refs, err := repo.References()
	if err != nil {
		return err
	}
	err = refs.ForEach(func(r *plumbing.Reference) error {
		target := ""
		if r.Name().IsTag() {
			tag, err := d.resolveTag(r.Name().Short(), repo)
			if err == nil {
				target = tag.String()
			}
		}
		fmt.Printf("ref: %v, name: %s, type: %v, isRefTag: %v, short: %v, target: %v\n", r.String(), r.Name().String(), r.Type(), r.Name().IsTag(), r.Name().Short(), target)
		return nil
	})
	return err
}

func processTag(tag *object.Tag) error {
	fmt.Printf("tag: %v, commit: %v, target: %v\n", tag.Name, tag.Hash.String(), tag.Target)
	return nil
}
