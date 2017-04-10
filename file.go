package u

import (
	"archive/zip"
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// PathExists returns true if a filesystem path exists
// Treats any error (e.g. lack of access due to permissions) as non-existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

// DirifyFileName converts a file name into sub-directories
func DirifyFileName(fn string) string {
	sep := string(os.PathSeparator)
	res := fn[:2] + sep + fn[2:4] + sep
	res += fn[4:6] + sep + fn[6:8] + sep + fn[8:10]
	return res + sep + fn[10:]
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

// CreateZipWithDirContent creates a zip file with the content of a directory.
// The names of files inside the zip file are relatitve to dirToZip e.g.
// if dirToZip is foo and there is a file foo/bar.txt, the name in the zip
// will be bar.txt
func CreateZipWithDirContent(zipFilePath, dirToZip string) error {
	if isDir, err := PathIsDir(dirToZip); err != nil || !isDir {
		// TODO: should return an error if err == nil && !isDir
		return err
	}
	zf, err := os.Create(zipFilePath)
	if err != nil {
		//fmt.Printf("Failed to os.Create() %s, %s\n", zipFilePath, err.Error())
		return err
	}
	defer zf.Close()
	zipWriter := zip.NewWriter(zf)
	// TODO: is the order of defer here can create problems?
	// TODO: need to check error code returned by Close()
	defer zipWriter.Close()

	//fmt.Printf("Walk root: %s\n", config.LocalDir)
	err = filepath.Walk(dirToZip, func(pathToZip string, info os.FileInfo, err error) error {
		if err != nil {
			//fmt.Printf("WalkFunc() received err %s from filepath.Wath()\n", err.Error())
			return err
		}
		//fmt.Printf("%s\n", path)
		isDir, err := PathIsDir(pathToZip)
		if err != nil {
			//fmt.Printf("PathIsDir() for %s failed with %s\n", pathToZip, err.Error())
			return err
		}
		if isDir {
			return nil
		}
		toZipReader, err := os.Open(pathToZip)
		if err != nil {
			//fmt.Printf("os.Open() %s failed with %s\n", pathToZip, err.Error())
			return err
		}
		defer toZipReader.Close()

		zipName := pathToZip[len(dirToZip)+1:] // +1 for '/' in the path
		inZipWriter, err := zipWriter.Create(zipName)
		if err != nil {
			//fmt.Printf("Error in zipWriter(): %s\n", err.Error())
			return err
		}
		_, err = io.Copy(inZipWriter, toZipReader)
		if err != nil {
			return err
		}
		//fmt.Printf("Added %s to zip file\n", pathToZip)
		return nil
	})
	return nil
}

// ReadLinesFromReader reads all lines from io.Reader
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

// ReadLinesFromFile reads all lines from a file
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
