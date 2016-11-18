package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetTaskList(c *gin.Context) {
	status := c.Query("status")

	task_list, err := taskmodel.QueryTaskList(status)
	if err != nil {
		log.Error(fmt.Sprintf("get task list error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskListNotFound.Code,
			"message": vars.ErrTaskListNotFound.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":      0,
		"task_list": task_list,
	})
}
