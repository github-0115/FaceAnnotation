package image_task

import (
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

type ImageTaskParmars struct {
	Introduce string `json:"introduce"` //本次说明
}

func CreateImageTask(c *gin.Context) {
	//	var imageTaskParmars ImageTaskParmars
	//	if err := c.BindJSON(&imageTaskParmars); err != nil {
	//		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
	//		c.JSON(400, gin.H{
	//			"code":    vars.ErrBindJSON.Code,
	//			"message": vars.ErrBindJSON.Msg,
	//		})
	//		return
	//	}

	//	introduce := imageTaskParmars.Introduce

	imageTaskId := uuid.NewV4().String()
	imagetaskColl, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		if err != imagetaskmodel.ErrImageTaskModelNotFound {
			log.Error(fmt.Sprintf("image task exist err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageTaskExist.Code,
				"message": vars.ErrImageTaskExist.Msg,
			})
			return
		}
	}

	if imagetaskColl != nil {
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskExist.Code,
			"message": vars.ErrImageTaskExist.Msg,
		})
		return
	}

	imageTaskColl := &imagetaskmodel.ImageTaskModel{
		ImageTaskId: imageTaskId,
		Introduce:   time.Now().Format("2006-01-02 03:04:05"),
		Status:      0,
		CreatedAt:   time.Now(),
	}

	err = imageTaskColl.Save()
	if err != nil {
		log.Error(fmt.Sprintf("image task save err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskSave.Code,
			"message": vars.ErrTaskSave.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":          0,
		"image_task_id": imageTaskId,
	})
}
