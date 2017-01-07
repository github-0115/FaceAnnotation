package import_data

import (
	//	imagemodel "FaceAnnotation/service/model/imagemodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	//	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetThrResult(c *gin.Context) {
	//	taskId := c.Query("task_id")
	imagePath := c.Query("image")

	//	imageList, err := imagemodel.QueryTaskImages(taskId)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("query task images err", err.Error()))
	//		c.JSON(400, gin.H{
	//			"code":    vars.ErrImageModelNotFound.Code,
	//			"message": vars.ErrImageModelNotFound.Msg,
	//		})
	//		return
	//	}

	//for _,image:=range imageList{

	//}
	res, err := thrfacemodel.ThrFaceRes(imagePath)
	if err != nil {
		log.Error(fmt.Sprintf("get face++ res fail err:%s", err))
	}

	rep, err := thrfacemodel.EightThreeFace(res.Face[0].FaceId)
	if err != nil {
		log.Error(fmt.Sprintf("get face++ res fail err:%s", err))
	}

	c.JSON(200, gin.H{
		"code":   0,
		"thrres": res,
		"eigres": rep,
	})
}
