package image

import (
	cfg "FaceAnnotation/config"
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
	if title == "" {
		log.Error(fmt.Sprintf("get task parmars error:%s"))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageParmars.Code,
			"message": vars.ErrImageParmars.Msg,
		})
		return
	}

	taskModel, err := taskmodel.QueryTask(title)
	if err != nil {
		log.Error(fmt.Sprintf("get task error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	var not_list []string
	if taskModel.Status != 2 {

		not_list = getNotAnnotationList(taskModel)

	} else {

		not_list = make([]string, 0, 0)
	}

	if len(not_list) == 0 {
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

func getNotAnnotationList(a *taskmodel.TaskModel) (ret []string) {
	for _, res := range a.Images {
		if res.Status == 0 {
			ret = append(ret, res.Url)
		}
	}
	return
}

func getNotCompletedImage(not_list []string) string {

	for _, res := range not_list {
		_, err := redismodel.GetCheckImageStr(res)
		if err != nil {
			if err == redismodel.RedisNotFound {
				err = redismodel.SetCheckImageStr(res, "yes")
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
