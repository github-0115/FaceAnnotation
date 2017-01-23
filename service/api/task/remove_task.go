package task

import (
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func RemoveTask(c *gin.Context) {
	taskId := c.PostForm("task_id")

	taskColl, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	if taskColl != nil {

		smallTasks, err := smalltaskmodel.QueryTaskAllSmallTasks(taskColl.TaskId)
		if err != nil {
			if err != smalltaskmodel.ErrSmallTaskModelNotFound {
				log.Error(fmt.Sprintf("query small tasks err %s", err))
				c.JSON(400, gin.H{
					"code":    vars.ErrSmallTaskNotFound.Code,
					"message": vars.ErrSmallTaskNotFound.Msg,
				})
				return
			}
		}

		for _, smallTask := range smallTasks {

			err := smalltaskmodel.RemoveSmallTask(smallTask.SmallTaskId)
			if err != nil {
				log.Error(fmt.Sprintf("remove small tasks err %s", err))
			}

		}
	}

	err = taskmodel.RemoveTask(taskColl.TaskId)
	if err != nil {
		log.Error(fmt.Sprintf("remove tasks err %s", err))
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": " remove task success",
	})
}
