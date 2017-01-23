package task

import (
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func StartTask(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	taskId := c.PostForm("task_id")

	_, err := usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	taskColl, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return

	}

	smallTasks, err := smalltaskmodel.QueryTaskAllSmallTasks(taskColl.TaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small tasks err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	for _, smallTask := range smallTasks {
		if smallTask.Status == smalltaskmodel.TaskStatus.Stop {
			err := smalltaskmodel.UpdateSmallTasks(smallTask.SmallTaskId, smalltaskmodel.TaskStatus.Start)
			if err != nil {
				log.Error(fmt.Sprintf("start small tasks  update status err %s", err))
			}
		}
	}

	err = taskmodel.UpdateTaskStatus(taskColl.TaskId, taskmodel.TaskStatus.Start)
	if err != nil {
		log.Error(fmt.Sprintf("start task update status err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return

	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "start success",
	})
}
