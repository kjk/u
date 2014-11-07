package u

import (
	"archive/zip"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// treats any error (e.g. lack of access due to permissions) as non-existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Returns true, nil if a path exists and is a directory
// Returns false, nil if a path exists and is not a directory (e.g. a file)
// Returns undefined, error if there was an error e.g. because a path doesn't exists
func PathIsDir(path string) (isDir bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func GetFileSize(path string) (int64, error) {
	if fi, err := os.Lstat(path); err != nil {
		return 0, err
	} else {
		return fi.Size(), nil
	}
}

func DirifyFileName(fn string) string {
	sep := string(os.PathSeparator)
	res := fn[:2] + sep + fn[2:4] + sep
	res += fn[4:6] + sep + fn[6:8] + sep + fn[8:10]
	return res + sep + fn[10:]
}

func CreateDirIfNotExists(path string) error {
	if !PathExists(path) {
		return os.MkdirAll(path, 0777)
	}
	return nil
}

func CreateDirIfNotExistsMust(dir string) string {
	if err := os.MkdirAll(dir, 0755); err != nil {
		PanicIfErr(err)
	}
	return dir
}

func CreateDirMust(path string) {
	err := CreateDirIfNotExists(path)
	PanicIfErr(err)
}

func CreateDirForFile(path string) error {
	dir := filepath.Dir(path)
	return CreateDirIfNotExists(dir)
}

func CreateDirForFileMust(path string) string {
	dir := filepath.Dir(path)
	err := CreateDirIfNotExists(dir)
	PanicIfErr(err)
	return dir
}

func WriteBytesToFile(d []byte, path string) error {
	if err := CreateDirIfNotExists(filepath.Dir(path)); err != nil {
		return err
	}
	return ioutil.WriteFile(path, d, 0644)
}

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

// the names of files inside the zip file are relatitve to dirToZip e.g.
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

func ReadLinesFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadLinesFromReader(f)
}
