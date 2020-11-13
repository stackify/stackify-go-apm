package transport_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
	"github.com/stackify/stackify-go-apm/transport"
)

func TestNewTransportWithDefaultConfig(t *testing.T) {
	c := config.NewConfig(config.WithTransportType("default"))
	c.LogPath = "/usr/local/stackify/stackify-python-apm/log/"
	c.ProcessID = "ProcessID"
	c.HostName = "HostName"

	transport := transport.NewTransport(c)
	ttype := reflect.TypeOf(transport)

	if ttype.String() != "*transport.defaultTransport" {
		t.Errorf("Error transport is not type of defaultTransport\n")
	}
}

func TestNewTransportWithConfigValue(t *testing.T) {
	c := config.NewConfig(config.WithTransportType("default"))

	transport := transport.NewTransport(c)
	ttype := reflect.TypeOf(transport)

	if ttype.String() != "*transport.defaultTransport" {
		t.Errorf("Error transport is not type of defaultTransport\n")
	}
}

func TestDefaultTransportMethods(t *testing.T) {
	c := config.NewConfig(config.WithTransportType("default"))
	fileNameFormat := fmt.Sprintf("%s%s#%s-", c.LogPath, c.HostName, c.ProcessID) + "%d.log"
	fileName := fmt.Sprintf(fileNameFormat, 1)

	transport := transport.NewTransport(c)
	transport.SendAll()
	transport.HandleTrace(new(span.StackifySpan))

	file, _ := os.Open(fileName)
	fileStat, _ := file.Stat()
	if fileStat.Size() == 0 {
		t.Errorf("Error transport is not logging traces to logfile\n")
	}
}

func TestDefaultTransportDeletedLogFile(t *testing.T) {
	c := config.NewConfig(config.WithTransportType("default"))
	fileNameFormat := fmt.Sprintf("%s%s#%s-", c.LogPath, c.HostName, c.ProcessID) + "%d.log"
	fileName := fmt.Sprintf(fileNameFormat, 1)

	transport := transport.NewTransport(c)
	os.Remove(fileName)
	transport.HandleTrace(new(span.StackifySpan))

	file, _ := os.Open(fileName)
	fileStat, _ := file.Stat()
	if fileStat.Size() == 0 {
		t.Errorf("Error transport not creating deleted log file\n")
	}
}

func TestDefaultTransportRollOver(t *testing.T) {
	c := config.NewConfig(
		config.WithTransportType("default"),
		config.WithLogFileThresholdSize(150),
	)
	c.ProcessID = "test"
	fileNameFormat := fmt.Sprintf("%s%s#%s-", c.LogPath, c.HostName, c.ProcessID) + "%d.log"
	expectedFileName := fmt.Sprintf(fileNameFormat, 2)

	transport := transport.NewTransport(c)
	transport.HandleTrace(new(span.StackifySpan))
	transport.HandleTrace(new(span.StackifySpan))

	_, err := os.Stat(expectedFileName)
	if os.IsNotExist(err) {
		t.Errorf("Error transport not incrementing log files\n")
	}
}

func TestDefaultTransportMaxLogFiles(t *testing.T) {
	c := config.NewConfig(config.WithTransportType("default"))
	fileNameFormat := fmt.Sprintf("%s%s#%s-", c.LogPath, c.HostName, c.ProcessID) + "%d.log"
	fileName := fmt.Sprintf(fileNameFormat, 1)
	var logFiles []string

	for i := 1; i < 15; i++ {
		c.ProcessID = strconv.Itoa(i)
		transport.NewTransport(c)
	}

	dir, err := filepath.Abs(filepath.Dir(fileName))
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, _ error) error {
		if err == nil && !fi.IsDir() && filepath.Ext(path) == ".log" {
			logFiles = append(logFiles, fi.Name())
		}
		return nil
	})

	if len(logFiles) != config.MaxLogFilesCount {
		t.Errorf("Error transport not deleting old log files\n")
	}
}
