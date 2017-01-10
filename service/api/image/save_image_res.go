package image

import (
	cfg "FaceAnnotation/config"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	timeoutmodel "FaceAnnotation/service/model/timeoutmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

type ImageResParams struct {
	Md5         string                 `json:"md5"`
	SmallTaskId string                 `json:"small_task_id"`
	Points      *imagemodel.Points     `json:"points"`
	FineRes     *imagemodel.FineResult `json:"fine_res"`
}

func SaveImageRes(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)

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
	fineRes := imageResParams.FineRes

	userColl, err := usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

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

	//fineTune
	if strings.EqualFold(userColl.Identity, usermodel.UserIdentity.FineTune) && strings.EqualFold(smallTaskModel.Areas, usermodel.UserIdentity.FineTune) {
		//
		if fineRes == nil || fineRes.Result == nil {
			log.Error(fmt.Sprintf(" user parmar nil fineRes=%s ,fineRes.Result=%s error:", fineRes, fineRes.Result))
			c.JSON(400, gin.H{
				"code":    vars.ErrLoginParams.Code,
				"message": vars.ErrLoginParams.Msg,
			})
			return
		}
		if imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))] == nil {
			imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))] = make([]*imagemodel.FineResult, 0, 0)
		}

		timeOutImages, err := timeoutmodel.QuerySmallTaskImage(imageModel.Md5, smallTaskId)
		if err != nil {
			if err != timeoutmodel.ErrTimeOutModelNotFound {
				log.Error(fmt.Sprintf("time out image query err", err.Error()))
			}
		}
		if timeOutImages != nil && len(timeOutImages) != 0 {
			fineRes.CreatedAt = timeOutImages[0].CreatedAt.Format("2006-01-02 03:04:05")
		} else {
			fineRes.CreatedAt = time.Now().Add(-10 * time.Minute).Format("2006-01-02 03:04:05")
		}
		fineRes.FinishedAt = time.Now().Format("2006-01-02 03:04:05")
		imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))] = append(imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))], fineRes)

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
		for _, res := range imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))] {
			if strings.EqualFold(res.SmallTaskId, smallTaskId) {
				count += 1
			}
		}

		if count >= smallTaskModel.LimitCount {
			err := updateTaskStatus(smallTaskModel)
			if err != nil {
				log.Error(fmt.Sprintf("task status update err", err.Error()))
			}
		}

		c.JSON(200, gin.H{
			"code":    0,
			"message": "save success",
		})
		return
	}

	if points == nil || points.Points == nil {
		log.Error(fmt.Sprintf(" user parmar points=%s,points.Points=%s nil error:", points, points.Points))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	if imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))] == nil {
		imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))] = make(map[string][]*imagemodel.Points)
	}
	timeOutImages, err := timeoutmodel.QuerySmallTaskImage(imageModel.Md5, smallTaskId)
	if err != nil {
		if err != timeoutmodel.ErrTimeOutModelNotFound {
			log.Error(fmt.Sprintf("time out image query err", err.Error()))
		}
	}
	if points != nil && len(points.Points) != 0 {
		points.CreatedAt = timeOutImages[0].CreatedAt.Format("2006-01-02 03:04:05")
	} else {
		points.CreatedAt = time.Now().Add(-10 * time.Minute).Format("2006-01-02 03:04:05")
	}
	points.FinishedAt = time.Now().Format("2006-01-02 03:04:05")
	imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas] = append(imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas], points)

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
	for _, res := range imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas] {
		if strings.EqualFold(res.SmallTaskId, smallTaskId) {
			count += 1
		}
	}

	if count >= smallTaskModel.LimitCount {
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
		for _, res := range image.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas] {
			if strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
				count += 1
			}
		}

		if count < smallTaskModel.LimitCount {
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
		if err != smalltaskmodel.ErrSmallTaskModelNotFound {
			log.Error(fmt.Sprintf("small task query err", err.Error()))
			return err
		}
	}
	if len(notList) == 0 {
		allimageList, err := imagemodel.QueryTaskImages(smallTaskModel.TaskId)
		md5s := make([]string, 0, 0)
		for _, image := range allimageList {
			md5s = append(md5s, image.Md5)
		}
		stm := &smalltaskmodel.SmallTaskModel{
			TaskId:          smallTaskModel.TaskId,
			SmallTaskId:     uuid.NewV4().String(),
			SmallTaskImages: md5s,
			PointType:       smallTaskModel.PointType,
			Areas:           "fineTune",
			LimitCount:      int64(cfg.Cfg.FinetuneUserCount),
			Status:          0,
			CreatedAt:       time.Now(),
		}

		err = stm.Save()
		if err != nil {
			log.Error(fmt.Sprintf(" finetune task save err", err.Error()))
			return err
		}

	}
	return nil
}

func updateTaskStatus(smallTaskModel *smalltaskmodel.SmallTaskModel) error {
	imageList, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		return err
	}
	for _, image := range imageList {
		var (
			count int64 = 0
		)
		for _, res := range image.FineResults[strconv.Itoa(int(smallTaskModel.PointType))] {
			if strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
				count += 1
			}
		}

		if count < smallTaskModel.LimitCount {
			return err
		}
	}

	err = smalltaskmodel.UpdateSmallTasks(smallTaskModel.SmallTaskId, smalltaskmodel.TaskStatus.Success)
	if err != nil {
		log.Error(fmt.Sprintf("fine small task status update err", err.Error()))
		return err
	}

	err = taskmodel.UpdateTaskStatus(smallTaskModel.TaskId, taskmodel.TaskStatus.Success)
	if err != nil {
		log.Error(fmt.Sprintf(" task status update err", err.Error()))
		return err
	}

	return nil
}
