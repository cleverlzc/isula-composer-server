package session

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/isula/isula-composer-server/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	conn := os.Getenv("TESTCONN")
	models.InitDB(conn, "mysql", "default")
}

func TestGetCacheDir(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "isula-composer-server")
	os.MkdirAll(filepath.Join(tmpDir, "2"), 0700)
	fi, _ := os.Create(filepath.Join(tmpDir, "3"))
	fi.Close()
	rootCacheDir = tmpDir
	defer os.RemoveAll(tmpDir)

	cases := []struct {
		id       int64
		dir      string
		expected bool
	}{
		{1, filepath.Join(tmpDir, "1"), true},
		{2, filepath.Join(tmpDir, "2"), true},
		{3, "", false},
	}

	for _, c := range cases {
		var task models.Task
		task.ID = c.id
		cd, err := GetCacheDir(task)
		assert.Equal(t, c.dir, cd)
		assert.Equal(t, c.expected, err == nil)
	}
}

func TestRun(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "isula-composer-server")
	rootCacheDir = tmpDir
	defer os.RemoveAll(tmpDir)

	var task models.Task
	task.ID = 1
	task.OutputFile = "output"
	err := Run(task)
	assert.NotNil(t, err)

	task.OutputFile = ""
	task.Scripts = "touch a; touch b"
	err = Run(task)
	assert.Nil(t, err)
	//FIXME
	_, err = os.Stat(filepath.Join(tmpDir, "1", "a"))
	assert.Nil(t, err)
	_, err = os.Stat(filepath.Join(tmpDir, "1", "b"))
	assert.Nil(t, err)

	os.MkdirAll(filepath.Join(tmpDir, "2"), 0700)
	task.ID = 2
	err = Run(task)
	assert.Nil(t, err)

	fi, _ := os.Create(filepath.Join(tmpDir, "3"))
	fi.Close()
	task.ID = 3
	err = Run(task)
	assert.NotNil(t, err)
}

func TestGetFileStat(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "isula-composer-server")
	rootCacheDir = tmpDir
	defer os.RemoveAll(tmpDir)

	var task models.Task
	var url string
	task.ID = 1
	fi, _ := os.Create(filepath.Join(tmpDir, "1"))
	fi.Close()
	_, err := GetFileStat(task, url)
	assert.NotNil(t, err)

	task.ID = 2
	url = "/"
	_, err = GetFileStat(task, url)
	assert.Nil(t, err)

	url = "/notfound"
	_, err = GetFileStat(task, url)
	assert.NotNil(t, err)
}

func TestReadDir(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "isula-composer-server")
	rootCacheDir = tmpDir
	defer os.RemoveAll(tmpDir)

	var task models.Task
	var url string
	task.ID = 1
	fi, _ := os.Create(filepath.Join(tmpDir, "1"))
	fi.Close()
	_, err := ReadDir(task, url)
	assert.NotNil(t, err)

	task.ID = 2
	url = "/"
	os.MkdirAll(filepath.Join(tmpDir, "2", "sub"), 0700)
	files, err := ReadDir(task, url)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(files))

	url = "/notfound"
	_, err = ReadDir(task, url)
	assert.NotNil(t, err)
}

func TestReadFile(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "isula-composer-server")
	rootCacheDir = tmpDir
	defer os.RemoveAll(tmpDir)

	var task models.Task
	var url string
	task.ID = 1
	fi, _ := os.Create(filepath.Join(tmpDir, "1"))
	fi.Close()
	_, err := ReadFile(task, url)
	assert.NotNil(t, err)

	task.ID = 2
	url = "/notfound"
	_, err = ReadFile(task, url)
	assert.NotNil(t, err)

	task.ID = 3
	os.MkdirAll(filepath.Join(tmpDir, "3"), 0700)
	fi, _ = os.Create(filepath.Join(tmpDir, "3", "data"))
	fi.Close()
	url = "data"
	_, err = ReadFile(task, url)
	assert.Nil(t, err)
}
