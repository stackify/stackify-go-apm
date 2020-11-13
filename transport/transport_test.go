package transport_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
	"github.com/stackify/stackify-go-apm/transport"
)

func TestNewTransportWithDefaultConfig(t *testing.T) {
	c := new(config.Config)

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
		t.Errorf("Error transport is not logging traces to logfile\n")
	}
}
