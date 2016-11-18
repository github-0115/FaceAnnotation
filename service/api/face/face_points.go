package face

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func UpsertFacePoint(c *gin.Context) {
	face_str := c.PostForm("face")

	faceModel, err := facemodel.StringToJson(face_str)
	if err != nil {
		log.Error(fmt.Sprintf("face json unmarshal error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrJsonUnmarshal.Code,
			"message": vars.ErrJsonUnmarshal.Msg,
		})
		return
	}

	_, err = facemodel.UpsertFaceResult(faceModel)
	if err != nil {
		log.Error(fmt.Sprintf("face points upsert error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrFaceModelUpsert.Code,
			"message": vars.ErrFaceModelUpsert.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "face points upsert success !",
	})
}
