package import_data

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	vars "FaceAnnotation/service/vars"
	"encoding/json"
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
			image.ThrFaces["deepir_import"] = map[string]interface{}{}
		}

		if importPoints[image.Url] == nil {
			//face++ res
			var thrRes *thrfacemodel.EightThreeFaceModel
			five, _ := thrfacemodel.ThrFaceFileRes(image.Url, fileByte)
			if err != nil {
				log.Error(fmt.Sprintf("get face++ five res fail err:%s", err))
			}
			if five.Face[0] != nil {
				thrRes, err = thrfacemodel.EightThreeFace(five.Face[0].FaceId)
				if err != nil {
					log.Error(fmt.Sprintf("get face++ 83 res fail err:%s", err))
				}
				thrRes.Result[0].FaceHeight = five.Face[0].Position.Height
				thrRes.Result[0].FaceWidth = five.Face[0].Position.Width
				thrRes.Result[0].ImageWidth = five.ImgWidth
				thrRes.Result[0].ImageHeight = five.ImgHeight
				res1B, _ := json.Marshal(thrRes)

				var result interface{}
				if err := json.Unmarshal(res1B, &result); err != nil {
					fmt.Println("json unmarshal err=%s", err)
				}
				fmt.Println("-----result%s-----", result)
				if image.ThrFaces == nil {
					image.ThrFaces["face++"] = make(map[string]interface{})
				}
				image.ThrFaces["face++"] = make(map[string]interface{})
				image.ThrFaces["face++"]["83"] = result
			}

		} else {
			image.ThrFaces["deepir_import"]["95"] = importPoints[image.Url]
		}

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
