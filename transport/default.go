package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"bitbucket.org/stackify/stackify-go-apm/config"
	"bitbucket.org/stackify/stackify-go-apm/trace/span"
)

const (
	defaultThreshold int64 = 50000000
	loggerPrefixSize int64 = 30
)

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
		fmt.Println(">>> Closing file: ", sil.getFileName())
		sil.file.Close()
	}

	if increment {
		sil.fileNameCounter += 1
	}
	fileName := sil.getFileName()

	fmt.Println(">>> Creating log file: ", fileName)
	var err error = nil
	sil.file, err = os.Create(fileName)
	if err != nil {
		fmt.Println(">>> Error creating file: "+fileName+". Err: ", err)
		return
	}

	os.Chmod(fileName, os.ModePerm)
	sil.log = log.New(sil.file, "", log.LUTC)
}

func (sil *stackifyIncrementalLogger) println(ssjson string) {
	sil.mutex.RLock()
	defer sil.mutex.RUnlock()

	if sil.file == nil {
		fmt.Println(">>> Log file doesn't exist.")
		return
	}

	if _, err := os.Stat(sil.getFileName()); os.IsNotExist(err) {
		fmt.Println(">>> Log file " + sil.getFileName() + "doesn't exist or deleted")
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
	fmt.Println(">>>>> Logging trace: ", stackifySpan.Id)
	stackifySpanJSON, _ := json.Marshal(stackifySpan)
	dt.println(string(stackifySpanJSON))
	// fmt.Println(string(stackifySpanJSON))
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
			defaultThreshold,
		),
	}
}
