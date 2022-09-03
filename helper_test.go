package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetObjects(t *testing.T) {
	files := []string{
		"testdata/content.yaml",
		"testdata/content-multi-layers.yaml",
	}
	objs, err := GetObjects(files...)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(objs))
}
