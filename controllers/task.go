package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/isula/isula-composer-server/models"
	"github.com/isula/isula-composer-server/session"
)

// Task defines the task operations
type Task struct {
	beego.Controller
}

// Create adds a new task
func (t *Task) Create() {
	user := t.Ctx.Input.Param(":user")
	output := t.Ctx.Input.Query("output")

	logs.Debug("Create task for '%s', output is '%s'", user, output)

	var config, scripts string

	file, _, err := t.Ctx.Request.FormFile("file")
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot find the upload file '%s'.", user))
		return
	}
	data, _ := ioutil.ReadAll(file)
	file.Close()

	if output != "" {
		config = string(data)
	} else {
		scripts = string(data)
	}

	id, err := models.AddTaskFull(user, output, config, scripts)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create a task for '%s'.", user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, id, nil)
}

// List lists the tasks
func (t *Task) List() {
	user := t.Ctx.Input.Param(":user")

	logs.Debug("List tasks of '%s'", user)

	tasks, err := models.QueryTaskListByUser(user)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the tasks from '%s'.", user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, tasks, nil)
}

// Get returns the task detail
func (t *Task) Get() {
	user := t.Ctx.Input.Param(":user")
	idStr := t.Ctx.Input.Param(":id")

	logs.Debug("Get task %s from '%s'", idStr, user)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Invalid id detected '%s': %v.", idStr, err))
		return
	}

	task, err := models.QueryTaskByID(int64(id))
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the task '%d' from '%s'.", id, user))
		return
	} else if task == nil {
		CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to find the task '%d' from '%s'.", id, user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, task, nil)
}

// Delete deletes the task and the task's data
func (t *Task) Delete() {
	user := t.Ctx.Input.Param(":user")
	idStr := t.Ctx.Input.Param(":id")

	logs.Debug("Delete task %s from '%s'", idStr, user)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Invalid id detected '%s': %v.", idStr, err))
		return
	}

	if num, err := models.RemoveTask(int64(id)); err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete the task '%d' from '%s'.", id, user))
		return
	} else if num == 0 {
		CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to delete the task '%d' from '%s', cannot find it.", id, user))
		return
	}

	var task models.Task
	task.ID = int64(id)
	err = session.Remove(task)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete the task '%d' from '%s'.", id, user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, "success to remove task", nil)
}

// Trigger calls the task
// TODO: we don't allow to edit a task now
func (t *Task) Trigger() {
	user := t.Ctx.Input.Param(":user")
	idStr := t.Ctx.Input.Param(":id")

	logs.Debug("Trigger task %s from '%s'", idStr, user)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Invalid id detected '%s': %v.", idStr, err))
		return
	}

	task, err := models.QueryTaskByID(int64(id))
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the task '%d' from '%s'.", id, user))
		return
	} else if task == nil {
		CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to find the task '%d' from '%s'.", id, user))
		return
	}

	err = session.Run(*task)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusExpectationFailed, err, fmt.Sprintf("Failed to run the task '%s/%d' from '%v'.", user, id, err))
		return
	}
	CtxSuccessWrap(t.Ctx, http.StatusOK, "ok", nil)
}

// Files returns the output files
func (t *Task) Files() {
	user := t.Ctx.Input.Param(":user")
	idStr := t.Ctx.Input.Param(":id")
	url := t.Ctx.Input.Param(":splat")

	logs.Debug("Get task files %s from '%s:%s'", url, user, idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Invalid id detected '%s': %v.", idStr, err))
		return
	}

	task, err := models.QueryTaskByID(int64(id))
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the task '%d' from '%s'.", id, user))
		return
	} else if task == nil {
		CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to find the task '%d' from '%s'.", id, user))
		return
	}

	fi, err := session.GetFileStat(*task, url)
	if err != nil {
		if os.IsNotExist(err) {
			CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to find the file '%s' from '%s:%d'.", url, user, id))
			return
		}
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the file '%s' from '%s:%d'.", url, user, id))
		return
	}

	if fi.IsDir() {
		logs.Debug("%s is a directory", url)
		files, _ := session.ReadDir(*task, url)
		t.TplName = "task/files.tpl"
		t.Data["files"] = files
		t.Render()
		return
	}

	data, err := session.ReadFile(*task, url)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the file content '%s' from '%s:%d'.", url, user, id))
		return
	}

	header := make(map[string]string)
	header["Content-Length"] = fmt.Sprint(len(data))
	CtxDataWrap(t.Ctx, http.StatusOK, data, header)
}
