package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func MakeTmpDir() (dir string, err error) {
	username := "unknown_user"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	// Remove path seperators in usernames.
	username = strings.Replace(username, `\`, "", -1)
	username = strings.Replace(username, "/", "", -1)

	return ioutil.TempDir("", fmt.Sprintf("courier_%s_", username))
}

func CopyDirContents(srcDir, dstDir, ignoreDir string) error {

	LogDebug(`Copying directory contents from %q to %q`, srcDir, dstDir)

	// Ignore the error, since the directory may not exist. If there are any
	// 'stat' type errors, then they will be caught when we start copying.
	_ = os.RemoveAll(dstDir)

	// Copy the files from src to dst. This is a bit painful in Go...
	return filepath.Walk(srcDir, func(p string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ignoreDir {
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(srcDir, p)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		src, err := filepath.EvalSymlinks(p)
		if err != nil {
			return err
		}
		info, err = os.Stat(src)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Create dir.
			LogDebug(`Creating dir %q`, dst)
			if err := os.MkdirAll(dst, info.Mode()); err != nil {
				return err
			}
		} else {
			// Copy the file. Yes this is ineffictient on Linux. But it's the
			// only way to do things so that it works for both Linux and
			// Windows without modification.
			LogDebug(`Copying file from %q to %q`, src, dst)
			if buf, err := ioutil.ReadFile(src); err != nil {
				return err
			} else if err := ioutil.WriteFile(dst, buf, info.Mode()); err != nil {
				return err
			}
		}
		return nil
	})
}

func CreateDirHash(dir string, ignoreDir string) ([]byte, error) {
	LogDebug(`Creating directory hash for %q`, dir)

	hash := sha1.New()
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		src, err := filepath.EvalSymlinks(p)
		if err != nil {
			return err
		}
		info, err = os.Stat(src)
		if err != nil {
			return err
		}

		LogDebug(`Base path: %q`, info.Name())
		if info.IsDir() && info.Name() == ignoreDir {
			return filepath.SkipDir
		}

		// Get the path relative to the folder we're hashing
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		// Hash the filename
		LogDebug(`Hashing filename: %q`, rel)
		hash.Write([]byte(rel))
		// Hash the file mode.
		LogDebug(`Hashing the mode: %q`, info.Mode().String())
		hash.Write([]byte(info.Mode().String()))
		// Hash the file contents.
		switch {
		case info.Mode().IsDir():
			return nil
		case info.Mode().IsRegular():
			fp, err := os.Open(p)
			if err != nil {
				return err
			}
			defer fp.Close()

			LogDebug(`Hashing contents of %q`, rel)
			if _, err := io.Copy(hash, fp); err != nil {
				return err
			}
		// Not a directory or reg file -> probably a symlink
		default:
			return fmt.Errorf("%q has unexpected FileMode %q.  If it is a symlink please delete it manually and try again.",
				dir, info.Mode().String())
		}

		return nil
	})

	LogDebug(`Hash for %q: %x`, dir, hash.Sum(nil))
	return hash.Sum(nil), err
}
