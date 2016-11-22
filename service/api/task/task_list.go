package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetTaskList(c *gin.Context) {
	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Error(fmt.Sprintf("get task list parmars error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskParmars.Code,
			"message": vars.ErrTaskParmars.Msg,
		})
		return
	}

	task_list, err := taskmodel.QueryTaskList(int64(status))
	if err != nil {
		log.Error(fmt.Sprintf("get task list error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskListNotFound.Code,
			"message": vars.ErrTaskListNotFound.Msg,
		})
		return
	}

	if task_list == nil {
		task_list = make([]*taskmodel.TaskModel, 0, 0)
	}

	c.JSON(200, gin.H{
		"code":      0,
		"task_list": task_list,
	})
}
