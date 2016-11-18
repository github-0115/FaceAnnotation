package image

import (
	cfg "FaceAnnotation/config"
	facemodel "FaceAnnotation/service/model/facemodel"
	//	imagemodel "FaceAnnotation/service/model/imagemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	redismodel "FaceAnnotation/utils/redisclient"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

var (
	completeStr  = "completed"
	image_domain = cfg.Cfg.Domian + "origin_images/"
)

func GetOneImage(c *gin.Context) {
	title := c.Query("title")
	log.Error(fmt.Sprintf("get local task title err %S", title))
	taskModel, err := taskmodel.QueryTask(title)
	if err != nil {
		log.Error(fmt.Sprintf("get task error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	//	image_list, err := imagemodel.GetImageList("./origin_images")
	//	if err != nil {
	//		log.Error(fmt.Sprintf("get local image list err %S", err.Error()))
	//		c.JSON(400, gin.H{
	//			"code":    vars.ErrJsonUnmarshal.Code,
	//			"message": vars.ErrJsonUnmarshal.Msg,
	//		})
	//		return
	//	}

	already_list, err := facemodel.QueryAll()
	if err != nil {
		log.Error(fmt.Sprintf("get already image list err %S", err.Error()))
		if err != facemodel.ErrFaceModelNotFound {
			c.JSON(400, gin.H{
				"code":    vars.ErrFaceCursor.Code,
				"message": vars.ErrFaceCursor.Msg,
			})
			return
		}
	}

	var not_list []string
	if already_list == nil {

		not_list = taskModel.Url

	} else {
		not_list = RemoveDuplicatesAndEmpty(taskModel, already_list)
	}

	if not_list == nil {
		log.Error(fmt.Sprintf("The task has been completed !"))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskCompleted.Code,
			"message": vars.ErrTaskCompleted.Msg,
		})
		return
	}

	image_url := getNotCompletedImage(not_list)
	if strings.EqualFold(image_url, completeStr) {
		log.Error(fmt.Sprintf("The task has been completed !"))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskCompleted.Code,
			"message": vars.ErrTaskCompleted.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":      0,
		"image_url": image_domain + image_url,
	})
}

func getNotCompletedImage(not_list []string) string {

	for _, res := range not_list {
		_, err := redismodel.GetCheckEmailStr(res)
		if err != nil {
			if err == redismodel.RedisNotFound {
				err = redismodel.SetCheckEmailStr(res, "yes")
				if err != nil {
					log.Error(fmt.Sprintf("image set redis err %S", err.Error()))
				}

				return res
			}

			log.Error(fmt.Sprintf("image get %s redis err %S", res, err.Error()))
		}

	}
	return completeStr
}
