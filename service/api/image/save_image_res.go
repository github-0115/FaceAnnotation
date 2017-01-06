package image

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type ImageResParams struct {
	Md5         string             `json:"md5"`
	SmallTaskId string             `json:"small_task_id"`
	Points      *imagemodel.Points `json:"points"`
}

func SaveImageRes(c *gin.Context) {
	//	name, _ := c.Get("username")
	//	username := name.(string)

	var imageResParams ImageResParams
	if err := c.BindJSON(&imageResParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}

	md5Str := imageResParams.Md5
	smallTaskId := imageResParams.SmallTaskId
	points := imageResParams.Points

	smallTaskModel, err := smalltaskmodel.QuerySmallTask(smallTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	imageModel, err := imagemodel.QueryImage(md5Str)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	if imageModel.Results["deepir"] == nil {
		imageModel.Results["deepir"] = make(map[string][]*imagemodel.Points)
	}

	imageModel.Results["deepir"][smallTaskModel.Areas] = append(imageModel.Results["deepir"][smallTaskModel.Areas], points)

	_, err = imagemodel.UpsertImageModel(imageModel)
	if err != nil {
		log.Error(fmt.Sprintf("image update err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelUpdate.Code,
			"message": vars.ErrImageModelUpdate.Msg,
		})
		return
	}

	var count int64 = 0
	for _, res := range imageModel.Results["deepir"][smallTaskModel.Areas] {
		if strings.EqualFold(res.SmallTaskId, smallTaskId) {
			count += 1
		}
	}

	if count == smallTaskModel.LimitCount {
		err := updateStatus(smallTaskModel)
		if err != nil {
			log.Error(fmt.Sprintf("task status update err", err.Error()))
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "save success",
	})
}

func updateStatus(smallTaskModel *smalltaskmodel.SmallTaskModel) error {
	imageList, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		return err
	}
	for _, image := range imageList {
		var (
			count int64 = 0
		)
		for _, res := range image.Results["deepir"][smallTaskModel.Areas] {
			if strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
				count += 1
			}
		}

		if count != smallTaskModel.LimitCount {
			return err
		}
	}

	err = smalltaskmodel.UpdateSmallTasks(smallTaskModel.SmallTaskId, smalltaskmodel.TaskStatus.Success)
	if err != nil {
		log.Error(fmt.Sprintf("small task status update err", err.Error()))
		return err
	}

	notList, err := smalltaskmodel.QueryTaskSmallTasks(smallTaskModel.TaskId)
	if err != nil {
		log.Error(fmt.Sprintf("small task query err", err.Error()))
		return err
	}
	if len(notList) == 0 {
		err := taskmodel.UpdateTaskStatus(smallTaskModel.TaskId, taskmodel.TaskStatus.Success)
		if err != nil {
			log.Error(fmt.Sprintf(" task status update err", err.Error()))
			return err
		}
	}
	return nil
}
