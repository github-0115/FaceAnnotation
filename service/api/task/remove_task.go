package task

import (
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type DeleteParams struct {
	TaskId      string `json:"task_id"`
	ImageTaskId string `json:"image_task_id"`
}

func RemoveTask(c *gin.Context) {
	//	taskId := c.PostForm("task_id")

	var deleteParams DeleteParams
	if err := c.BindJSON(&deleteParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	taskId := deleteParams.TaskId
	imageTaskId := deleteParams.ImageTaskId

	imageTaskModel, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
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

	err = imagetaskmodel.PullImageTaskId(imageTaskModel.ImageTaskId, taskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": " remove task success",
	})
}
