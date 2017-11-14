package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetVersionInfo(t *testing.T) {
	commit := "test"
	version := "0.0.1"
	SetVersionInfo(commit, version)
	assert.Equal(t, commit, versionInfo.GitCommit)
	assert.Equal(t, version, versionInfo.Version)
}
