package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/panshul007/apideps/logger"
)

type FileOps struct {
	log logger.GenericLogger
}

func NewFileOpsInstance(log logger.GenericLogger) *FileOps {
	return &FileOps{
		log: log,
	}
}

func (d *FileOps) ensureTargetFolder(targetPath string) error {
	d.log.Debugf("ensuring folder exists: %s", targetPath)
	return os.MkdirAll(targetPath, os.ModePerm)
}

func (d *FileOps) copyFolderPath(commit *object.Tree, srcFolderPath string, targetFolderPath string) error {
	srcFolder := filepath.Clean(srcFolderPath)
	dstFolder := filepath.Clean(targetFolderPath)
	d.log.Infof("copying folder: %s to: %s", srcFolder, dstFolder)
	// ensure source folder is present in commit object as a folder
	err := d.ensureSourceFolder(commit, srcFolder)
	if err != nil {
		return err
	}
	// ensure target folder exists or create
	err = d.ensureTargetFolder(targetFolderPath)
	if err != nil {
		return fmt.Errorf("error while creating target folder: %s, error: %w", dstFolder, err)
	}
	// list the files to be copied
	subTree, err := commit.Tree(srcFolder)
	if err != nil {
		return fmt.Errorf("error while getting the subtree for source folder: %s, error: %w", srcFolder, err)
	}
	err = d.listTreeFiles(subTree)
	if err != nil {
		return fmt.Errorf("error while listing tree file for src: %s, error: %w", srcFolder, err)
	}
	// iterate and copy the files
	files := subTree.Files()
	for {
		file, err := files.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		// check for directory
		// if file is folder make recursive call
		dir := filepath.Dir(file.Name)
		if dir != "" && dir != "." {
			err := d.copyFolderPath(subTree, dir, filepath.Join(dstFolder, dir))
			if err != nil {
				return err
			}
		}
		// copy file
		dstFilePath := filepath.Join(dstFolder, filepath.Base(file.Name))
		err = d.copyFileObject(file, dstFilePath)
		if err != nil {
			return fmt.Errorf("error while copying file path: %s, base name: %s to: %s, error: %w", file.Name, filepath.Base(file.Name), dstFilePath, err)
		}
	}
	return nil
}

func (d *FileOps) ensureSourceFolder(commit *object.Tree, srcFolder string) error {
	files := commit.Files()
	for {
		file, err := files.Next()
		if err != nil {
			return fmt.Errorf("source folder: %s not found in source commit tree", srcFolder)
		}
		if srcFolder == filepath.Dir(file.Name) {
			return nil
		}
	}
}

func (d *FileOps) copyFileObject(srcObject *object.File, dstFilePath string) error {
	d.log.Debugf("copying file: %s to: %s", srcObject.Name, dstFilePath)

	// opening source file object from memory filesystem
	srcFile, err := srcObject.Reader()
	if err != nil {
		return fmt.Errorf("error opening file: %s, error: %w", srcObject.Name, err)
	}
	// closing source file with defer
	defer func(log logger.GenericLogger, srcfile io.ReadCloser, srcFilename string) {
		err := srcfile.Close()
		if err != nil {
			log.Errorf("error while closing file: %s, error: %v", srcFilename, err)
		}
	}(d.log, srcFile, srcObject.Name)

	// creating destination file on real file system
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("error while creating destination file: %s, error: %w", dstFilePath, err)
	}
	// closing destination file
	defer func(log logger.GenericLogger, dstfile *os.File) {
		err := dstfile.Close()
		if err != nil {
			log.Errorf("error while closing file: %s, error: %v", dstfile.Name(), err)
		}
	}(d.log, dstFile)

	// copying to destination file on real filesystem from source file on memory filesystem
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error while copying srcFile: %s to dstFile: %s, error: %w", srcObject.Name, dstFilePath, err)
	}
	// flushing the destination file writer buffer
	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("error while finishing to copy srcFile: %s to dstFile: %s, error: %w", srcObject.Name, dstFilePath, err)
	}

	// copying file permissions
	srcFilePerm, err := srcObject.Mode.ToOSFileMode()
	if err != nil {
		return fmt.Errorf("error while extracting file mode from src file: %s, error: %w", srcObject.Name, err)
	}
	err = os.Chmod(dstFilePath, srcFilePerm)
	if err != nil {
		return fmt.Errorf("error while setting file permissions from source file to destination file: %w", err)
	}
	return nil
}

func (d *FileOps) listTreeFiles(commitTree *object.Tree) error {
	return commitTree.Files().ForEach(func(f *object.File) error {
		base := filepath.Base(f.Name)
		dir, fname := filepath.Split(f.Name)
		fmt.Printf("file name: %s, type: %v, dir: %v, base:%s, splitDir: %s, splitFName:%s\n", f.Name, f.Type().String(), filepath.Dir(f.Name), base, dir, fname)
		return nil
	})
}
