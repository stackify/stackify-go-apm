package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
)

const (
	defaultThreshold int64 = 50000000
	loggerPrefixSize int64 = 30
)

type File struct {
	Path string
	os.FileInfo
}

type Files []*File

func (a Files) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Files) Len() int {
	return len(a)
}

type ByMtime struct {
	Files
}

func (a ByMtime) Less(i, j int) bool {
	return a.Files[i].ModTime().After(a.Files[j].ModTime())
}

type stackifyIncrementalLogger struct {
	fileNameFormat  string
	threshold       int64
	fileNameCounter int
	log             *log.Logger
	file            *os.File
	mutex           sync.RWMutex
}

func (sil *stackifyIncrementalLogger) getFileName() string {
	return fmt.Sprintf(sil.fileNameFormat, sil.fileNameCounter)
}

func (sil *stackifyIncrementalLogger) createFileOrRollOver(increment bool) {
	sil.mutex.RLock()
	defer sil.mutex.RLock()

	if sil.file != nil {
		sil.file.Close()
	}

	if increment {
		sil.fileNameCounter += 1
	}
	fileName := sil.getFileName()

	var err error = nil
	sil.file, err = os.Create(fileName)
	if err != nil {
		return
	}

	os.Chmod(fileName, os.ModePerm)
	sil.deleteOldLogFiles()
	sil.log = log.New(sil.file, "", log.LUTC)
}

func (sil *stackifyIncrementalLogger) deleteOldLogFiles() {
	files := sil.getLogFiles()
	if len(files) > config.MaxLogFilesCount {
		filesToDelete := files[config.MaxLogFilesCount:]
		for _, f := range filesToDelete {
			os.Remove(f.Path)
		}
	}
}

func (sil *stackifyIncrementalLogger) getLogFiles() Files {
	var logFiles Files

	dir, err := filepath.Abs(filepath.Dir(sil.getFileName()))
	if err != nil {
		return logFiles
	}

	err = filepath.Walk(dir, func(path string, fi os.FileInfo, _ error) error {
		if err == nil && !fi.IsDir() && filepath.Ext(path) == ".log" {
			logFiles = append(logFiles, &File{
				Path:     path,
				FileInfo: fi,
			})
		}
		return nil
	})
	sort.Sort(ByMtime{logFiles})
	return logFiles
}

func (sil *stackifyIncrementalLogger) println(ssjson string) {
	sil.mutex.RLock()
	defer sil.mutex.RUnlock()

	if sil.file == nil {
		return
	}

	if _, err := os.Stat(sil.getFileName()); os.IsNotExist(err) {
		sil.createFileOrRollOver(false)
	}

	if sil.shouldRollOver(ssjson) {
		sil.createFileOrRollOver(true)
	}
	sil.log.SetPrefix(time.Now().UTC().Format("2006-01-02, 15:04:05.000000") + "> ")
	sil.log.Println(ssjson)
}

func (sil *stackifyIncrementalLogger) shouldRollOver(ssjson string) bool {
	fileStat, _ := sil.file.Stat()
	fileSize := fileStat.Size()
	return fileSize+int64(len(ssjson))+loggerPrefixSize > sil.threshold
}

func newStackifyIncrementalLogger(fileNameFormat string, threshold int64) *stackifyIncrementalLogger {
	stackifyIncrementalLogger := &stackifyIncrementalLogger{
		fileNameFormat: fileNameFormat,
		threshold:      threshold,
	}
	stackifyIncrementalLogger.createFileOrRollOver(true)
	return stackifyIncrementalLogger
}

type defaultTransport struct {
	*stackifyIncrementalLogger

	config          *config.Config
	fileNameCounter int
	log             *log.Logger
	file            *os.File
	threshold       int64
	mutex           sync.RWMutex
}

func (dt *defaultTransport) HandleTrace(stackifySpan *span.StackifySpan) {
	stackifySpanJSON, _ := json.Marshal(stackifySpan)
	dt.println(string(stackifySpanJSON))
}

func (dt *defaultTransport) SendAll() {
	// nothing to flush for default transport
}

func newDefaultTransport(c *config.Config) Transport {
	fileNameFormat := fmt.Sprintf("%s%s#%s-", c.LogPath, c.HostName, c.ProcessID) + "%d.log"

	return &defaultTransport{
		config: c,
		stackifyIncrementalLogger: newStackifyIncrementalLogger(
			fileNameFormat,
			c.LogFileThresholdSize,
		),
	}
}
