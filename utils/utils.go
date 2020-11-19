package utils

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// TimeToTimestamp function converts time.Time to string
func TimeToTimestamp(t time.Time) string {
	return fmt.Sprintf("%.4f", float64(t.UnixNano())/1e9*1000)
}

// TraceIdToUUID function converts TraceID to UUID string format
func TraceIdToUUID(s []byte) string {
	value, _ := uuid.FromBytes(s)
	return value.String()
}

// SpanIdToString function converts SpanId to string
func SpanIdToString(s []byte) string {
	value := binary.BigEndian.Uint32(s)
	return strconv.FormatUint(uint64(value), 10)
}
