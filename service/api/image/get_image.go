package image

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetImageList(c *gin.Context) {
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

	c.JSON(200, gin.H{
		"code": 0,
		"url":  not_list,
	})
}

func RemoveDuplicatesAndEmpty(a *taskmodel.TaskModel, b []*facemodel.FaceUrl) (ret []string) {

	a_len := len(a.Images)
	b_len := len(b)
	for i := 0; i < a_len; i++ {
		for j := 0; j < b_len; j++ {
			if strings.EqualFold(a.Images[i].Url, b[j].Url) {
				continue
			}
			ret = append(ret, a.Images[i].Url)
		}
	}

	return
}
