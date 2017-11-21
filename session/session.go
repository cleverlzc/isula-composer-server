package session

import (
	"errors"
	"fmt"
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

func Run(task models.Task) error {
	if task.OutputFile != "" {
		// TODO, if we want to support it, we need to check env in 'init'
		return errors.New("builder is not implemented")
	}

	cacheDir, err := GetCacheDir(task)
	logs.Debug("Run task on %s\n", cacheDir)
	if err != nil {
		return err
	}

	logs.Debug("the script is '%s'\n", task.Scripts)
	cmd := exec.Command("/bin/sh", "-c", task.Scripts)
	cmd.Dir = cacheDir
	err = cmd.Run()
	return err
}
