package u

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
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
	return err == nil && st.Mode().IsRegular()
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
	Must(err)
	return dir
}

// CreateDirMust creates a directory. Panics on error
func CreateDirMust(path string) string {
	err := CreateDirIfNotExists(path)
	Must(err)
	return path
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
	Must(err)
	return dir
}

// WriteFileCreateDirMust is like ioutil.WriteFile() but also creates
// intermediary directories
func WriteFileCreateDirMust(d []byte, path string) error {
	if err := CreateDirIfNotExists(filepath.Dir(path)); err != nil {
		return err
	}
	return ioutil.WriteFile(path, d, 0644)
}

func WriteFileMust(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	Must(err)
}

func ReadFileMust(path string) []byte {
	d, err := ioutil.ReadFile(path)
	Must(err)
	return d
}

// like io.Closer Close() but ignores an error so better to use as
// defer CloseNoError(f)
func CloseNoError(f io.Closer) {
	_ = f.Close()
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

func RemoveFilesInDirMust(dir string) {
	if !DirExists(dir) {
		return
	}
	files, err := ioutil.ReadDir(dir)
	Must(err)
	for _, fi := range files {
		if !fi.Mode().IsRegular() {
			continue
		}
		path := filepath.Join(dir, fi.Name())
		err = os.Remove(path)
		Must(err)
	}
}

func RemoveFileLogged(path string) {
	err := os.Remove(path)
	if err == nil {
		Logf("RemoveFileLogged('%s')\n", path)
		return
	}
	if os.IsNotExist(err) {
		// TODO: maybe should print note
		return
	}
	Logf("os.Remove('%s') failed with '%s'\n", path, err)
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

func CopyFileMust(dst, src string) {
	Must(CopyFile(dst, src))
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

// absolute path of the current directory
func CurrDirAbsMust() string {
	dir, err := filepath.Abs(".")
	Must(err)
	return dir
}

// we are executed for do/ directory so top dir is parent dir
func CdUpDir(dirName string) {
	startDir := CurrDirAbsMust()
	dir := startDir
	for {
		// we're already in top directory
		if filepath.Base(dir) == dirName && DirExists(dir) {
			err := os.Chdir(dir)
			Must(err)
			return
		}
		parentDir := filepath.Dir(dir)
		PanicIf(dir == parentDir, "invalid startDir: '%s', dir: '%s'", startDir, dir)
		dir = parentDir
	}
}

func FmtSizeHuman(size int64) string {
	return humanize.Bytes(uint64(size))
}

func PrintFileSize(path string) {
	st, err := os.Stat(path)
	if err != nil {
		fmt.Printf("File '%s' doesn't exist\n", path)
		return
	}
	fmt.Printf("'%s': %s\n", path, FmtSizeHuman(st.Size()))
}

func AreFilesEuqalMust(path1, path2 string) bool {
	d1 := ReadFileMust(path1)
	d2 := ReadFileMust(path2)
	return bytes.Equal(d1, d2)
}

func FilesSameSize(path1, path2 string) bool {
	s1, err := GetFileSize(path1)
	if err != nil {
		return false
	}
	s2, err := GetFileSize(path2)
	if err != nil {
		return false
	}
	return s1 == s2
}

func DirCopyRecur(dstDir, srcDir string, shouldCopyFn func(path string) bool) ([]string, error) {
	err := CreateDirIfNotExists(dstDir)
	if err != nil {
		return nil, err
	}
	fileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}
	var allCopied []string
	for _, fi := range fileInfos {
		name := fi.Name()
		if fi.IsDir() {
			dst := filepath.Join(dstDir, name)
			src := filepath.Join(srcDir, name)
			copied, err := DirCopyRecur(dst, src, shouldCopyFn)
			if err != nil {
				return nil, err
			}
			allCopied = append(allCopied, copied...)
			continue
		}

		src := filepath.Join(srcDir, name)
		dst := filepath.Join(dstDir, name)
		shouldCopy := true
		if shouldCopyFn != nil {
			shouldCopy = shouldCopyFn(src)
		}
		if !shouldCopy {
			continue
		}
		CopyFileMust(dst, src)
		allCopied = append(allCopied, src)
	}
	return allCopied, nil
}

func DirCopyRecurMust(dstDir, srcDir string, shouldCopyFn func(path string) bool) []string {
	copied, err := DirCopyRecur(dstDir, srcDir, shouldCopyFn)
	Must(err)
	return copied
}
