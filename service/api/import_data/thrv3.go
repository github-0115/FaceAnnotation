package import_data

import (
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	"fmt"
	//	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ThrFaceV3(c *gin.Context) {

	file, fileHeader, err := c.Request.FormFile("image")
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(fmt.Sprintf("ioutil ReadAll file err" + err.Error()))
		c.JSON(400, gin.H{
			"code":    0,
			"message": "read file err",
		})
		return
	}

	thr, err := thrfacemodel.ThrFaceFileResV3(fileHeader.Filename, fileByte)
	if err != nil {
		log.Error(fmt.Sprintf("face ReadAll file err" + err.Error()))
		c.JSON(400, gin.H{
			"code":    0,
			"message": "face=+=",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":  0,
		"v3res": thr,
	})
}
