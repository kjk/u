package u

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
)

// FileWalkEntry describes a single file from FileWalk
type FileWalkEntry struct {
	Dir      string
	FileInfo os.FileInfo
}

// Path returns full path of the file
func (e *FileWalkEntry) Path() string {
	return filepath.Join(e.Dir, e.FileInfo.Name())
}

// FileWalk describes a file traversal
type FileWalk struct {
	startDir    string
	FilesChan   chan *FileWalkEntry
	askedToStop int32
}

// Stop stops file traversal
func (ft *FileWalk) Stop() {
	atomic.StoreInt32(&ft.askedToStop, 1)
	// drain the channel
	for range ft.FilesChan {
	}
}

func fileWalkWorker(ft *FileWalk) {
	toVisit := []string{ft.startDir}
	defer close(ft.FilesChan)

	for len(toVisit) > 0 {
		shouldStop := atomic.LoadInt32(&ft.askedToStop)
		if shouldStop > 0 {
			return
		}
		// would be more efficient to shift by one and
		// chop off at the end
		dir := toVisit[0]
		toVisit = StringsRemoveFirst(toVisit)

		files, err := ioutil.ReadDir(dir)
		// TODO: should I send errors as well?
		if err != nil {
			continue
		}
		for _, fi := range files {
			path := filepath.Join(dir, fi.Name())
			mode := fi.Mode()
			if mode.IsDir() {
				toVisit = append(toVisit, path)
			} else if mode.IsRegular() {
				fte := &FileWalkEntry{
					Dir:      dir,
					FileInfo: fi,
				}
				ft.FilesChan <- fte
			}
		}
	}
}

// StartFileWalk starts a file traversal from startDir
func StartFileWalk(startDir string) *FileWalk {
	// buffered channel so that
	ch := make(chan *FileWalkEntry, 1024*64)
	ft := &FileWalk{
		startDir:  startDir,
		FilesChan: ch,
	}
	go fileWalkWorker(ft)
	return ft
}
