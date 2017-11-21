package session

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/isula/isula-composer-server/models"
)

var rootCacheDir string

func init() {
	// TODO init from config
	rootCacheDir = "/tmp/build-cache"
}

// GetCacheDir returns the cache directory
func GetCacheDir(task models.Task) (string, error) {
	cacheDir := filepath.Join(rootCacheDir, fmt.Sprintf("%d", task.ID))
	fi, err := os.Stat(cacheDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	if fi != nil && fi.IsDir() {
		return cacheDir, nil
	}
	err = os.MkdirAll(cacheDir, 0700)
	if err != nil {
		return "", err
	}

	return cacheDir, nil
}

// Run starts to run the building task
func Run(task models.Task) error {
	task.Status = models.TaskStatusRunning
	models.UpdateTask(&task)

	if task.OutputFile != "" {
		// TODO, if we want to support it, we need to check env in 'init'
		task.Status = models.TaskStatusInnerError
		models.UpdateTask(&task)
		return errors.New("builder is not implemented")
	}

	cacheDir, err := GetCacheDir(task)
	logs.Debug("Run task on %s\n", cacheDir)
	if err != nil {
		task.Status = models.TaskStatusInnerError
		models.UpdateTask(&task)
		return err
	}

	logs.Debug("the script is '%s'\n", task.Scripts)
	cmd := exec.Command("/bin/sh", "-c", task.Scripts)
	cmd.Dir = cacheDir
	// TODO: run it at background
	err = cmd.Run()
	if err != nil {
		task.Status = models.TaskStatusFailed
	} else {
		task.Status = models.TaskStatusFinish
	}
	models.UpdateTask(&task)
	return err
}

// GetFileStat returns the status of the output file
func GetFileStat(task models.Task, url string) (os.FileInfo, error) {
	cacheDir, err := GetCacheDir(task)
	logs.Debug("GetFileStat on %s/%s\n", cacheDir, url)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(cacheDir, url)
	return os.Stat(file)
}

// ReadDir gets the files of an output directory
func ReadDir(task models.Task, url string) ([]string, error) {
	var files []string
	cacheDir, err := GetCacheDir(task)
	logs.Debug("ReadDir on %s/%s\n", cacheDir, url)
	if err != nil {
		return files, err
	}

	file := filepath.Join(cacheDir, url)
	fis, err := ioutil.ReadDir(file)
	if err != nil {
		return files, err
	}

	for _, fi := range fis {
		files = append(files, fi.Name())
	}

	return files, nil
}

// ReadFile reads the content of an output file
func ReadFile(task models.Task, url string) ([]byte, error) {
	cacheDir, err := GetCacheDir(task)
	logs.Debug("ReadFile on %s/%s\n", cacheDir, url)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(cacheDir, url)
	return ioutil.ReadFile(file)
}

// Remove removes all the task outputs
func Remove(task models.Task) error {
	cacheDir, err := GetCacheDir(task)
	logs.Debug("Remove %s\n", cacheDir)
	if err != nil {
		return err
	}

	return os.RemoveAll(cacheDir)
}
