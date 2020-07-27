package service

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/panshul007/apideps/logger"
	"github.com/stretchr/testify/assert"
)

const apiRepoPath = "git@bitbucket.org:egym-com/apis.git"

func TestListTags(t *testing.T) {
	fs := memfs.New()
	_ = logger.SetupLogger(true)
	eLog := logger.Logger()
	d := NewDepLoader(eLog)
	repo, err := d.cloneRepo(apiRepoPath, fs)
	assert.Nil(t, err)
	d.listTags(repo)
}

func TestListRefs(t *testing.T) {
	fs := memfs.New()
	_ = logger.SetupLogger(true)
	eLog := logger.Logger()
	d := NewDepLoader(eLog)
	repo, err := d.cloneRepo(apiRepoPath, fs)
	assert.Nil(t, err)
	err = d.listReferences(repo)
	assert.Nil(t, err)
}

func TestResolveTag(t *testing.T) {
	fs := memfs.New()
	_ = logger.SetupLogger(true)
	eLog := logger.Logger()
	d := NewDepLoader(eLog)
	repo, err := d.cloneRepo(apiRepoPath, fs)
	assert.Nil(t, err)
	hash, err := d.resolveTag("measurement-service-v1.0.1", repo)
	assert.Nil(t, err)
	assert.NotNil(t, hash)
	assert.NotEqual(t, plumbing.ZeroHash, hash)
	assert.Equal(t, "4b00f9948a6aec2caa4689c58a826243a4bee1aa", hash.String())
}
