package linker

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// expandHome obtains the user's home directory from '~'
func expandHome(p string) string {
	if p[0] != '~' {
		return p
	}

	// TODO: this doesn't need to be done multiple times...
	usr, _ := user.Current()
	home := usr.HomeDir

	if len(p) == 1 {
		return home
	} else {
		return home + p[1:]
	}
}

func GetFullPath(p string) string {
	var finalPath string

	pathExpanded := filepath.Clean(p)
	if pathExpanded[0] == '~' {
		finalPath = expandHome(pathExpanded)
	} else {
		finalPath = pathExpanded
	}

	return finalPath
}

func newSymlink(source, dest string) error {
	err := os.Symlink(source, dest)
	if err != nil {
		return err
	}
	return nil
}

func mkAllDirs(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}

func backup(filePath string) error {
	backupPath := filePath + ".bak"
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		// remove the old backup for now...
		// TODO: add a number at the end or something
		err := os.Remove(backupPath)
		if err != nil {
			return err
		}
	}
	err := os.Rename(filePath, backupPath)
	if err != nil {
		return err
	}

	return nil
}

func guardSymlink(file, dest string) error {
	// guard that the file/dest aren't the same
	if file == dest {
		return fmt.Errorf("source and destination file cannot be the same")
	}

	// guard that the dest is not a parent folder of the 'file'
	// this is probably quite naive...but it works for my purposes
	fs := strings.Split(file, "/")
	ds := strings.Split(dest, "/")
	if len(ds) >= len(fs) {
		return nil
	}

	lastDS := len(ds) - 1
	for i := 0; i <= len(fs); i++ {
		if i > lastDS {
			// error out because lastFS is a parent
			return fmt.Errorf("unable to symlink %[1]s to %[2]s: %[2]s is a parent of %[1]s",
				file,
				dest,
			)
		}

		if fs[i] != ds[i] {
			break
		}
	}

	// guard that we aren't killing '~/'
	usr, _ := user.Current()
	if usr.HomeDir == dest {
		return fmt.Errorf("unable to symlink %s to %s: that would nuke your home directory",
			file,
			dest,
		)
	}
	return nil
}
