package apm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apm "github.com/stackify/stackify-go-apm"
)

func TestStackifyAPM(t *testing.T) {
	stackifyapm, err := apm.NewStackifyAPM()

	assert.Nil(t, err)
	assert.NotNil(t, stackifyapm.Context)
	assert.NotNil(t, stackifyapm.Tracer)
	assert.NotNil(t, stackifyapm.TraceProvider)
}

func TestStackifyAPMShutDown(t *testing.T) {
	stackifyapm, err := apm.NewStackifyAPM()

	stackifyapm.Shutdown()

	assert.Nil(t, err)
}
