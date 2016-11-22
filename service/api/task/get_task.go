package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetTask(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		log.Error(fmt.Sprintf("get task parmars error:%s"))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskParmars.Code,
			"message": vars.ErrTaskParmars.Msg,
		})
		return
	}

	taskModel, err := taskmodel.QueryTask(title)
	if err != nil {
		log.Error(fmt.Sprintf("get task error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	if taskModel == nil {
		taskModel = &taskmodel.TaskModel{}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"task": taskModel,
	})
}
