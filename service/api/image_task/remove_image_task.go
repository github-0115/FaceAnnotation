package image_task

import (
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func RemoveImageTask(c *gin.Context) {
	imageTaskId := c.PostForm("task_id")

	imagetaskColl, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	if imagetaskColl == nil {
		log.Error(fmt.Sprintf("image task not found  err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	if len(imagetaskColl.TaskId) != 0 {
		taskColls, err := taskmodel.GetImageTasks(imagetaskColl.TaskId)
		if err != nil {
			if err != taskmodel.ErrTaskModelNotFound {
				log.Error(fmt.Sprintf("task not found err", err.Error()))
				c.JSON(400, gin.H{
					"code":    vars.ErrTaskNotFound.Code,
					"message": vars.ErrTaskNotFound.Msg,
				})
				return
			}
		}

		for _, taskColl := range taskColls {
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

				err = taskmodel.RemoveTask(taskColl.TaskId)
				if err != nil {
					log.Error(fmt.Sprintf("remove tasks err %s", err))
				}

			}
		}
	}

	err = imagetaskmodel.RemoveImageTask(imageTaskId)
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
		"message": " remove image task success",
	})
}
