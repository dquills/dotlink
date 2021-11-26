package linker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Overwrite bool              `yaml:"overwrite-existing"`
	Mkdirs    bool              `yaml:"make-dirs"`
	Backup    bool              `yaml:"backup-existing"`
	Paths     map[string]string `yaml:"paths"`
}

func (c *Config) LinkAll() {
	errors := []error{}
	for path, link := range c.Paths {
		err := c.LinkOne(path, link)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		fmt.Printf("\n")
		fmt.Printf("%d error(s) found:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())

		}
	}
}

func (c *Config) LinkOne(file, linkPath string) error {
	// TODO: linkPath must not be the same directory to prevent recursive symlinks
	if file[0] == '/' || file[0] == '~' {
		// We should only link things from our cwd or the passed in path
		return fmt.Errorf("absolute path found: %s (source files should exist in cwd or the directory passed in with '-d')",
			file,
		)
	}

	// clean the filepath just in case...
	fileCleaned, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("unable to get absolute path for  %s: %s", file,
			err.Error(),
		)
	}

	// Check to see if the file to link actually exists...
	if _, err = os.Stat(fileCleaned); os.IsNotExist(err) {
		return fmt.Errorf("unable to link %s: file does not exist", file)
	}

	// if a dir is provided, the symlink path should end in the original file/folder name
	if linkPath[len(linkPath)-1] == '/' {
		linkPath += filepath.Base(file)
	}
	lpCleaned := GetFullPath(linkPath)

	// Stop yourself from borking your system
	err = guardSymlink(fileCleaned, lpCleaned)
	if err != nil {
		return err
	}

	// check if target dir exists, otherwise mkdir -p
	linkDir := filepath.Dir(lpCleaned)
	if _, err := os.Stat(linkDir); errors.Is(err, os.ErrNotExist) {
		if c.Mkdirs {
			err := mkAllDirs(linkDir)
			if err != nil {
				return fmt.Errorf("unable to make dirs %s: %s", linkDir,
					err.Error(),
				)
			}
		} else {
			return fmt.Errorf("path %s does not exist: set 'make-dirs' to true to create it", linkDir)
		}
	}

	// if target dest exists, delete it or back it up
	// TODO: eval the symlink and do nothing if it's the same link
	if _, err := os.Stat(lpCleaned); !os.IsNotExist(err) {
		if c.Backup {
			err := backup(lpCleaned)
			if err != nil {
				return fmt.Errorf("unable to overwrite %s: %s", linkPath,
					err.Error(),
				)
			}
		} else if c.Overwrite {
			err := os.Remove(lpCleaned)
			if err != nil {
				return fmt.Errorf("unable to overwrite %s: %s", linkPath,
					err.Error(),
				)
			}
		}
	}
	err = newSymlink(fileCleaned, lpCleaned)
	if err != nil {
		return fmt.Errorf("unable to link %s to %s: %s",
			file,
			linkPath,
			err.Error(),
		)
	}

	return nil
}
