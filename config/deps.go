package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ApiDependencies struct {
	Dependencies map[string]Dependency `yaml:"dependencies"`
}

type Dependency struct {
	Repo       string `yaml:"repo"`
	RepoFolder string `yaml:"repofolder"`
	Commit     string `yaml:"commit,omitempty"`
	Tag        string `yaml:"tag,omitempty"`
	TargetPath string `yaml:"targetpath"`
}

func LoadDeps(depsFile string) (*ApiDependencies, error) {
	depsBytes, err := ioutil.ReadFile(depsFile)
	if err != nil {
		return nil, fmt.Errorf("error while reading apideps file: %w", err)
	}
	apiDeps := &ApiDependencies{}
	err = yaml.Unmarshal(depsBytes, apiDeps)
	if err != nil {
		return nil, fmt.Errorf("error while parsing api deps file: %w", err)
	}
	return apiDeps, nil
}

func ListDeps(depsFile string) error {
	apiDeps, err := LoadDeps(depsFile)
	if err != nil {
		return fmt.Errorf("error while loading dependencies: %w", err)
	}
	apiDeps.prettyPrint()
	return nil
}

func (d *ApiDependencies) prettyPrint() {
	for name, dep := range d.Dependencies {
		fmt.Printf("%s :\n", name)
		fmt.Printf("\t repo: %s\n", dep.Repo)
		fmt.Printf("\t repofolder: %s\n", dep.RepoFolder)
		if dep.Commit != "" {
			fmt.Printf("\t commit: %s\n", dep.Commit)
		}
		if dep.Tag != "" {
			fmt.Printf("\t tag: %s\n", dep.Tag)
		}
		fmt.Printf("\t targetpath: %s\n", dep.TargetPath)
	}
}
