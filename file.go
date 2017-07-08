package u

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// PathExists returns true if a filesystem path exists
// Treats any error (e.g. lack of access due to permissions) as non-existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileExists returns true if a given path exists and is a file
func FileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir() && st.Mode().IsRegular()
}

// DirExists returns true if a given path exists and is a directory
func DirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// PathIsDir returns true if a path exists and is a directory
// Returns false, nil if a path exists and is not a directory (e.g. a file)
// Returns undefined, error if there was an error e.g. because a path doesn't exists
func PathIsDir(path string) (isDir bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// GetFileSize returns size of the file
func GetFileSize(path string) (int64, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// CreateDirIfNotExists creates a directory if it doesn't exist
func CreateDirIfNotExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// CreateDirIfNotExistsMust creates a directory. Panics on error
func CreateDirIfNotExistsMust(dir string) string {
	err := os.MkdirAll(dir, 0755)
	PanicIfErr(err)
	return dir
}

// CreateDirMust creates a directory. Panics on error
func CreateDirMust(path string) {
	err := CreateDirIfNotExists(path)
	PanicIfErr(err)
}

// CreateDirForFile creates intermediary directories for a file
func CreateDirForFile(path string) error {
	dir := filepath.Dir(path)
	return CreateDirIfNotExists(dir)
}

// CreateDirForFileMust is like CreateDirForFile. Panics on error.
func CreateDirForFileMust(path string) string {
	dir := filepath.Dir(path)
	err := CreateDirIfNotExists(dir)
	PanicIfErr(err)
	return dir
}

// WriteBytesToFile is like ioutil.WriteFile() but also creates intermediary
// directories
func WriteBytesToFile(d []byte, path string) error {
	if err := CreateDirIfNotExists(filepath.Dir(path)); err != nil {
		return err
	}
	return ioutil.WriteFile(path, d, 0644)
}

// ListFilesInDir returns a list of files in a directory
func ListFilesInDir(dir string, recursive bool) []string {
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		isDir, err := PathIsDir(path)
		if err != nil {
			return err
		}
		if isDir {
			if recursive || path == dir {
				return nil
			}
			return filepath.SkipDir
		}
		files = append(files, path)
		return nil
	})
	return files
}

// CopyFile copies a file
func CopyFile(dst, src string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()
	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()
	if _, err = io.Copy(fdst, fsrc); err != nil {
		return err
	}
	return nil
}

// ReadLinesFromReader reads all lines from io.Reader. Newlines are not included.
func ReadLinesFromReader(r io.Reader) ([]string, error) {
	res := make([]string, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}

// ReadLinesFromFile reads all lines from a file. Newlines are not included.
func ReadLinesFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadLinesFromReader(f)
}

// Sha1OfFile returns 20-byte sha1 of file content
func Sha1OfFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		//fmt.Printf("os.Open(%s) failed with %s\n", path, err.Error())
		return nil, err
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		//fmt.Printf("io.Copy() failed with %s\n", err.Error())
		return nil, err
	}
	return h.Sum(nil), nil
}

// Sha1HexOfFile returns 40-byte hex sha1 of file content
func Sha1HexOfFile(path string) (string, error) {
	sha1, err := Sha1OfFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1), nil
}

// PathMatchesExtensions returns true if path matches any of the extensions
func PathMatchesExtensions(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}
	ext := strings.ToLower(filepath.Ext(path))
	for _, allowed := range extensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// DeleteFilesIf deletes a files in a given directory if shouldDelete callback
// returns true
func DeleteFilesIf(dir string, shouldDelete func(os.FileInfo) bool) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		if fi.IsDir() || !fi.Mode().IsRegular() {
			continue
		}
		if shouldDelete(fi) {
			path := filepath.Join(dir, fi.Name())
			err = os.Remove(path)
			// Maybe: keep deleting?
			if err != nil {
				return err
			}
		}
	}
	return nil
}
