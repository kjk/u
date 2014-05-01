package u

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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

func CreateDirIfNotExists(path string) error {
	if !PathExists(path) {
		return os.MkdirAll(path, 0777)
	}
	return nil
}

func FileSha1(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		//fmt.Printf("os.Open(%s) failed with %s\n", path, err.Error())
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		//fmt.Printf("io.Copy() failed with %s\n", err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Sha1StringOfBytes(data []byte) string {
	return fmt.Sprintf("%x", Sha1OfBytes(data))
}

func Sha1OfBytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// TODO: remove in favor of Sha1OfBytes
func DataSha1(data []byte) (string, error) {
	return Sha1StringOfBytes(data), nil
}

func TimeSinceNowAsString(t time.Time) string {
	d := time.Now().Sub(t)
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	days := hours / 24
	hours = hours % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dhr", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dhr %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
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
