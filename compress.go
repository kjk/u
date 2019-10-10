package u

import (
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// implement io.ReadCloser over os.File wrapped with io.Reader.
// io.Closer goes to os.File, io.Reader goes to wrapping reader
type readerWrappedFile struct {
	f *os.File
	r io.Reader
}

func (rc *readerWrappedFile) Close() error {
	return rc.f.Close()
}

func (rc *readerWrappedFile) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}

// OpenFileMaybeCompressed opens a file that might be compressed with gzip
// or bzip2.
// TODO: could sniff file content instead of checking file extension
func OpenFileMaybeCompressed(path string) (io.ReadCloser, error) {
	ext := strings.ToLower(filepath.Ext(path))
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if ext == ".gz" {
		r, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		rc := &readerWrappedFile{
			f: f,
			r: r,
		}
		return rc, nil
	}
	if ext == ".bz2" {
		r := bzip2.NewReader(f)
		rc := &readerWrappedFile{
			f: f,
			r: r,
		}
		return rc, nil
	}
	return f, nil
}

// ReadFileMaybeCompressed reads file. Ungzips if it's gzipped.
func ReadFileMaybeCompressed(path string) ([]byte, error) {
	r, err := OpenFileMaybeCompressed(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// WriteFileGzipped writes data to a path, using best gzip compression
func WriteFileGzipped(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		f.Close()
		os.Remove(path)
		return err
	}
	err = w.Close()
	if err != nil {
		f.Close()
		os.Remove(path)
		return err
	}
	err = f.Close()
	if err != nil {
		os.Remove(path)
		return err
	}
	return nil
}

// GzipFile compresses srcPath with gzip and saves as dstPath
func GzipFile(dstPath, srcPath string) error {
	fSrc, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fSrc.Close()
	fDst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fDst.Close()
	w, err := gzip.NewWriterLevel(fDst, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, fSrc)
	if err != nil {
		return err
	}
	return w.Close()
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
	return err
}

func ReadZipFileMust(path string) map[string][]byte {
	r, err := zip.OpenReader(path)
	Must(err)
	defer CloseNoError(r)
	res := map[string][]byte{}
	for _, f := range r.File {
		rc, err := f.Open()
		Must(err)
		d, err := ioutil.ReadAll(rc)
		Must(err)
		_ = rc.Close()
		res[f.Name] = d
	}
	return res
}

func zipAddFile(zw *zip.Writer, zipName string, path string) {
	zipName = filepath.ToSlash(zipName)
	d, err := ioutil.ReadFile(path)
	Must(err)
	w, err := zw.Create(zipName)
	Must(err)
	_, err = w.Write(d)
	Must(err)
	fmt.Printf("  added %s from %s\n", zipName, path)
}

func zipDirRecur(zw *zip.Writer, baseDir string, dirToZip string) {
	dir := filepath.Join(baseDir, dirToZip)
	files, err := ioutil.ReadDir(dir)
	Must(err)
	for _, fi := range files {
		if fi.IsDir() {
			zipDirRecur(zw, baseDir, filepath.Join(dirToZip, fi.Name()))
		} else if fi.Mode().IsRegular() {
			zipName := filepath.Join(dirToZip, fi.Name())
			path := filepath.Join(baseDir, zipName)
			zipAddFile(zw, zipName, path)
		} else {
			path := filepath.Join(baseDir, fi.Name())
			s := fmt.Sprintf("%s is not a dir or regular file", path)
			panic(s)
		}
	}
}

// toZip is a list of files and directories in baseDir
// Directories are added recursively
func CreateZipFile(dst string, baseDir string, toZip ...string) {
	os.Remove(dst)
	if len(toZip) == 0 {
		panic("must provide toZip args")
	}
	fmt.Printf("Creating zip file %s\n", dst)
	w, err := os.Create(dst)
	Must(err)
	defer CloseNoError(w)
	zw := zip.NewWriter(w)
	Must(err)
	for _, name := range toZip {
		path := filepath.Join(baseDir, name)
		fi, err := os.Stat(path)
		Must(err)
		if fi.IsDir() {
			zipDirRecur(zw, baseDir, name)
		} else if fi.Mode().IsRegular() {
			zipAddFile(zw, name, path)
		} else {
			s := fmt.Sprintf("%s is not a dir or regular file", path)
			panic(s)
		}
	}
	err = zw.Close()
	Must(err)
}
