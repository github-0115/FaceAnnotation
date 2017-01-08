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
	SmallTaskId string `json:"small_task_id"`
	//	Image       *imagemodel.ImageModel `json:"image"`
	Md5       string     `json:"md5"`
	PointsRep *PointsRep `json:"points_rep"`
	PointType int64      `json:"point_type"`
	Areas     string     `json:"areas"`
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

	md5Str, err := getTimeOutImage(username, smallTaskId, smallTaskModel.Areas)

	if md5Str != nil {
		pRep := SwitchPoint(md5Str)
		rep := &ImageRep{
			SmallTaskId: smallTaskId,
			Md5:         md5Str.Md5,
			PointsRep:   pRep,
			PointType:   smallTaskModel.PointType,
			Areas:       smallTaskModel.Areas,
		}

		c.JSON(200, gin.H{
			"code": 0,
			"res":  rep,
		})

		return
	}

	//All images to get the small task
	imageList, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrNotImage.Code,
			"message": vars.ErrNotImage.Msg,
		})
		return
	}
	//All images not complete
	notImages := getNotImageList(imageList, username, smallTaskModel.Areas, smallTaskModel.LimitCount, smallTaskId)
	//	log.Info(fmt.Sprintf("image query notImages %s err", notImages))
	if len(notImages) == 0 {
		log.Error(fmt.Sprintf("Under the task gets no pictures err"))
		c.JSON(400, gin.H{
			"code":    vars.ErrNotImage.Code,
			"message": vars.ErrNotImage.Msg,
		})
		return
	}
	//	var md5Str *imagemodel.ImageModel
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
		var flag bool = false
		for _, res := range timeOutModels {
			if strings.EqualFold(res.User, username) {
				flag = true
				break
			}
		}
		if flag {
			continue
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

	if md5Str == nil {
		log.Error(fmt.Sprintf("Under the task gets no pictures err"))
		c.JSON(400, gin.H{
			"code":    vars.ErrNotImage.Code,
			"message": vars.ErrNotImage.Msg,
		})
		return
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

	pRep := SwitchPoint(md5Str)
	rep := &ImageRep{
		SmallTaskId: smallTaskId,
		//		Image:       md5Str,
		Md5:       md5Str.Md5,
		PointsRep: pRep,
		PointType: smallTaskModel.PointType,
		Areas:     smallTaskModel.Areas,
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

func getTimeOutImage(username string, smallTaskId string, area string) (*imagemodel.ImageModel, error) {
	timeOutImages, err := timeoutmodel.QueryUserTsakImage(username, smallTaskId)
	if err != nil {
		if err != timeoutmodel.ErrTimeOutModelNotFound {
			log.Error(fmt.Sprintf("time out image query err", err.Error()))
			return nil, err
		}
	}

	var imageModel *imagemodel.ImageModel
	if timeOutImages != nil {
		for _, timeOut := range timeOutImages {
			imageModel, err = imagemodel.QueryImage(timeOut.Md5)
			if err != nil {
				log.Error(fmt.Sprintf("image query err", err.Error()))
				continue
			}

			result := imageModel.Results["deepir"][area]
			//		log.Info(fmt.Sprintf("image query result %s err", result))
			var flag bool = false
			for _, res := range result {
				if strings.EqualFold(res.User, username) && strings.EqualFold(res.SmallTaskId, smallTaskId) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			return imageModel, nil
		}
	}

	return nil, nil
}

func getNotImageList(imageList []*imagemodel.ImageModel, username string, area string, limitCount int64, stmId string) []*imagemodel.ImageModel {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.Results["deepir"][area]
		//		log.Info(fmt.Sprintf("image query result %s err", result))
		if len(result) == 0 || result == nil {
			list = append(list, image)
			continue
		}

		var (
			count int64 = 0
			flag  bool  = true
		)
		for _, res := range result {
			if res.User == username {
				if strings.EqualFold(res.SmallTaskId, stmId) {
					log.Info(fmt.Sprintf("image query result user %s %s err", res.User, username))
					flag = false
					break
				}
				continue
			}

			if strings.EqualFold(res.SmallTaskId, stmId) {
				count += 1
			}

		}

		if count < limitCount && flag {
			list = append(list, image)
		}
	}
	return list
}
