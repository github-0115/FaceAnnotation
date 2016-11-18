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
		return
	} else {
		not_list = RemoveDuplicatesAndEmpty(taskModel, already_list)
	}

	c.JSON(200, gin.H{
		"code": 0,
		"url":  not_list,
	})
}

func RemoveDuplicatesAndEmpty(a *taskmodel.TaskModel, b []*facemodel.FaceUrl) (ret []string) {

	a_len := len(a.Url)
	b_len := len(b)
	for i := 0; i < a_len; i++ {
		for j := 0; j < b_len; j++ {
			if strings.EqualFold(a.Url[i], b[j].Url) {
				continue
			}
			ret = append(ret, a.Url[i])
		}
	}

	return
}
