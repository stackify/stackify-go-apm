package config

const (
	StackifyInstrumentationName string = "stackifyapm_tracer"
	DefaultLogFileThresholdSize int64  = 50000000
	DefaultTransportType        string = "default"
	MaxLogFilesCount            int    = 10
)

var (
	TransportTypes = map[string]bool{
		DefaultTransportType: true,
	}
)
