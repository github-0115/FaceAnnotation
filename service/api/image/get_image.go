package image

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	timeoutmodel "FaceAnnotation/service/model/timeoutmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type ImageRep struct {
	SmallTaskId string                 `json:"small_task_id"`
	Image       *imagemodel.ImageModel `json:"image"`
	PointType   int64                  `json:"point_type"`
	Areas       string                 `json:"areas"`
}

func GetImage(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	smallTaskId := c.Query("small_task_id")

	smallTaskModel, err := smalltaskmodel.QuerySmallTask(smallTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}
	//All images to get the small task
	imageList, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}
	//All images not complete
	notImages := getNotImageList(imageList, username, smallTaskModel.Areas, smallTaskModel.LimitCount, smallTaskId)
	var md5Str *imagemodel.ImageModel
	for _, image := range notImages {
		timeOutModels, err := timeoutmodel.QuerySmallTaskImage(image.Md5, smallTaskId)
		if err != nil {
			if err != timeoutmodel.ErrTimeOutModelNotFound {
				log.Error(fmt.Sprintf("image query err", err.Error()))
				c.JSON(400, gin.H{
					"code":    vars.ErrImageModelNotFound.Code,
					"message": vars.ErrImageModelNotFound.Msg,
				})
				return
			}
		}
		if timeOutModels == nil {
			md5Str = image
			break
		}

		result := image.Results["deepir"][smallTaskModel.Areas]
		var count int64 = 0
		for _, res := range result {
			if strings.EqualFold(res.SmallTaskId, smallTaskId) {
				count += 1
			}
		}

		if int64(len(timeOutModels))+count == smallTaskModel.LimitCount {
			continue
		}

		md5Str = image
		break

	}

	timeOutModel := timeoutmodel.TimeOutModel{
		SmallTaskId: smallTaskId,
		Md5:         md5Str.Md5,
		User:        username,
		CreatedAt:   time.Now(),
	}

	err = timeOutModel.Save()
	if err != nil {
		log.Error(fmt.Sprintf("user=%s get image =%s timeOutModel save err%s", username, md5Str, err.Error()))
	}

	rep := &ImageRep{
		SmallTaskId: smallTaskId,
		Image:       md5Str,
		PointType:   smallTaskModel.PointType,
		Areas:       smallTaskModel.Areas,
	}

	if smallTaskModel.Status != smalltaskmodel.TaskStatus.Start {
		err = smalltaskmodel.UpdateSmallTasks(smallTaskModel.SmallTaskId, smalltaskmodel.TaskStatus.Start)
		if err != nil {
			log.Error(fmt.Sprintf("small task status update err", err.Error()))
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"res":  rep,
	})
}

func getNotImageList(imageList []*imagemodel.ImageModel, username string, area string, limitCount int64, stmId string) []*imagemodel.ImageModel {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.Results["deepir"][area]
		if len(result) == 0 {
			list = append(list, image)
			continue
		}
		var count int64 = 0
		for _, res := range result {
			if res.User == username {
				if strings.EqualFold(res.SmallTaskId, stmId) {
					break
				}
				continue
			}
			if strings.EqualFold(res.SmallTaskId, stmId) {
				count += 1
			}
			if count == limitCount {
				break
			}
		}
		list = append(list, image)
	}
	return list
}

//	timeOutModel, err := timeoutmodel.QueryUserImage(username)
//	if err != nil {
//		if err != timeoutmodel.ErrTimeOutModelNotFound {
//			log.Error(fmt.Sprintf("image query err", err.Error()))
//			c.JSON(400, gin.H{
//				"code":    vars.ErrImageModelNotFound.Code,
//				"message": vars.ErrImageModelNotFound.Msg,
//			})
//			return
//		}
//	}

//	if timeOutModel != nil {
//		//
//	}
