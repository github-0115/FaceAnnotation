package import_data

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ImportResult(c *gin.Context) {
	taskId := c.PostForm("task_id")
	pointFile, _, err := c.Request.FormFile("point")
	fileByte, err := ioutil.ReadAll(pointFile)
	if err != nil {
		log.Error(fmt.Sprintf("ioutil ReadAll file err" + err.Error()))
		c.JSON(400, gin.H{
			"code":    0,
			"message": "read file err",
		})
		return
	}

	importPoints := readPoint(fileByte)

	imageList, err := imagemodel.QueryTaskImages(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task image err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	if len(imageList) != len(importPoints) {
		log.Error(fmt.Sprintf("image count != point count"))
	}

	for _, image := range imageList {
		if image.ThrFaces["deepir_import"] == nil {
			image.ThrFaces["deepir_import"] = []interface{}{}
		}
		//		image.ThrFaces["deepir_import"] = append(image.ThrFaces["deepir_import"], []interface{}{importPoints[image.Url]})
		image.ThrFaces["deepir_import"] = importPoints[image.Url]
		//		fmt.Println(importPoints[image.Url])
		_, err = imagemodel.UpsertImageModel(image)
		if err != nil {
			log.Error(fmt.Sprintf("image%s update thrres err", image.Url, err.Error()))
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "import res success",
	})
}
