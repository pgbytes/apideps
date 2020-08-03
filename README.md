# apideps

A tool to fetch files (proto files) from an API repository, with multiple apis stored as a mono repo. 
Every folder from the api repo can be fetched at a specific version of the repo, buy commit hash or tag.  

The output will be a folder structure with copied files as defined by the `targetpath` option in config.
 
# WIP

This tool is under active development and at a very early stage. 

## Install

On a system with Go installed and GO Modules turned on:

```
go get -u github.com/panshul007/apideps
```

## Config

The tool reads the api dependencies from the config yaml file. (Default: `apideps.yaml` in the execution folder)

Each dependency can be defined as:

```yaml
  dependency1:
    repo: ""
    repofolder: ""
    commit: ""
    tag: ""
    targetpath: ""
```

Where:
    
- dependency1  -> is the name for dependency unique in the config
- repo -> git clone repository url eg: `git@bitbucket.org:foobar-company/apis.git`
- repofolder -> the folder within the repo to be extracted eg: `service1/v1`. This folder path should be relative to the root of repo. 
- commit -> the complete commit hash of repo from which the API folder is to be extracted
- tag -> the tag name of repo from which the API folder is to be extracted
- targetpath -> the folder path to which the `repofolder` will be copied recursively. eg: `api/service1/v1`
    
## Usage

```
apideps --help
```

## In memory FS options:
* https://github.com/spf13/afero
* https://github.com/go-git/go-billy 
    * https://pkg.go.dev/gopkg.in/src-d/go-billy.v4/memfs?tab=doc

## Credits

* This tool uses [go-git](https://github.com/go-git/go-git) internally for all git operations.
