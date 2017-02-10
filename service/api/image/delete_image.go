package image

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type DeleteParams struct {
	Md5         []string `json:"md5"`
	ImageTaskId string   `json:"image_task_id"`
}

func DeleteImage(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)

	var deleteParams DeleteParams
	if err := c.BindJSON(&deleteParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	urls := deleteParams.Md5
	imageTaskId := deleteParams.ImageTaskId
	if len(urls) == 0 || imageTaskId == "" {
		log.Error(fmt.Sprintf("parmar nil err%v"))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	_, err := usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	imageTaskModel, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	for _, url := range urls {
		urlStrs := strings.Split(url, "/")
		md5 := urlStrs[len(urlStrs)-1]

		for _, task := range imageTaskModel.TaskId {
			taskModel, err := taskmodel.QueryTask(task)
			if err != nil {
				if err != taskmodel.ErrTaskModelNotFound {
					log.Error(fmt.Sprintf("query task  err", err.Error()))
					c.JSON(400, gin.H{
						"code":    vars.ErrTaskNotFound.Code,
						"message": vars.ErrTaskNotFound.Msg,
					})
					return
				}
			}

			if taskModel == nil {
				continue
			}
			smallTaskModels, err := smalltaskmodel.QueryTaskImageSmallTask(taskModel.TaskId, md5)
			if err != nil {
				log.Error(fmt.Sprintf("query small task err %s", err))
				c.JSON(400, gin.H{
					"code":    vars.ErrSmallTaskNotFound.Code,
					"message": vars.ErrSmallTaskNotFound.Msg,
				})
				return
			}

			for _, smallTask := range smallTaskModels {
				err := smalltaskmodel.RemoveSmallTaskImage(smallTask.SmallTaskId, md5)
				if err != nil {
					if err != smalltaskmodel.ErrSmallTaskModelNotFound {
						log.Error(fmt.Sprintf("query small task err %s", err))
						c.JSON(400, gin.H{
							"code":    vars.ErrSmallTaskNotFound.Code,
							"message": vars.ErrSmallTaskNotFound.Msg,
						})
						return
					}
				}
			}

			_, err = imagemodel.DeleteTaskImage(taskModel.TaskId, md5)
			if err != nil {
				if err != imagemodel.ErrImageModelNotFound {
					log.Error(fmt.Sprintf("image query err", err.Error()))
					c.JSON(400, gin.H{
						"code":    vars.ErrNotImage.Code,
						"message": vars.ErrNotImage.Msg,
					})
					return
				}
			} else {
				err = taskmodel.UpdateTaskCount(taskModel.TaskId, taskModel.Count-1)
				if err != nil {
					if err != taskmodel.ErrTaskModelNotFound {
						log.Error(fmt.Sprintf("task query err", err.Error()))
						c.JSON(400, gin.H{
							"code":    vars.ErrTaskNotFound.Code,
							"message": vars.ErrTaskNotFound.Msg,
						})
						return
					}
				}
			}
		}

		err = imagetaskmodel.PullImageTaskImages(imageTaskModel.ImageTaskId, md5)
		if err != nil {
			log.Error(fmt.Sprintf("image task not found err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageTaskNotFound.Code,
				"message": vars.ErrImageTaskNotFound.Msg,
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "delete image success",
	})
}
