package import_data

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	uploadmodel "FaceAnnotation/service/model/uploadmodel"
	vars "FaceAnnotation/service/vars"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ImportImage(c *gin.Context) {
	imageTaskId := c.PostForm("image_task_id")
	introduce := c.PostForm("introduce")
	isRes := c.PostForm("isRes")
	resFile, _, err := c.Request.FormFile("res")
	imageFile, imageFileHeader, err := c.Request.FormFile("image")
	if err != nil {
		log.Error(fmt.Sprintf("import image parmars err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImportImageParmars.Code,
			"message": vars.ErrImportImageParmars.Msg,
		})
		return
	}

	imagetaskColl, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	importPoints := make(map[string][]interface{})
	if resFile != nil && strings.EqualFold(isRes, "yes") {

		resByte, err := ioutil.ReadAll(resFile)
		if err != nil {
			log.Error(fmt.Sprintf("ioutil read res file err" + err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrReadImage.Code,
				"message": vars.ErrReadImage.Msg,
			})
			return
		}

		importPoints = readPoint(resByte)
		fmt.Println(importPoints)
	}

	fileByte, err := ioutil.ReadAll(imageFile)
	if err != nil {
		log.Error(fmt.Sprintf("ioutil read image file err" + err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrReadImage.Code,
			"message": vars.ErrReadImage.Msg,
		})
		return
	}

	h := md5.New()
	h.Write(fileByte)
	photoName := hex.EncodeToString(h.Sum(nil))

	imageColl, err := imagemodel.QueryImage(photoName)
	if err != nil {
		if err != imagemodel.ErrImageModelNotFound {
			log.Error(fmt.Sprintf("query image err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageModelNotFound.Code,
				"message": vars.ErrImageModelNotFound.Msg,
			})
			return
		}
	}

	if imageColl == nil {
		imageColl = &imagemodel.ImageModel{
			Md5:       photoName,
			Url:       imageFileHeader.Filename,
			ThrFaces:  make(map[string]map[string]interface{}),
			CreatedAt: time.Now(),
		}

		_, err = uploadmodel.UploadFile(photoName, fileByte)
		if err != nil {
			log.Error(fmt.Sprintf("image %s upload oss err %s", imageFileHeader.Filename, err.Error()))
		}

	}

	if strings.EqualFold(isRes, "not") {
		//face++ res
		imageColl, err = faceRes(photoName, fileByte, imageColl)
	} else {
		if imageColl.ThrFaces["deepir_import"] == nil {
			imageColl.ThrFaces["deepir_import"] = map[string]interface{}{}
		}
		// res?=nil
		if importPoints[imageColl.Url] == nil {
			//face++ res
			imageColl, err = faceRes(photoName, fileByte, imageColl)
		} else {
			imageColl.ThrFaces["deepir_import"]["95"] = importPoints[imageColl.Url]
		}
	}

	_, err = imagemodel.UpsertImageModel(imageColl)
	if err != nil {
		log.Error(fmt.Sprintf("image update err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelUpdate.Code,
			"message": vars.ErrImageModelUpdate.Msg,
		})
		return
	}

	err = imagetaskmodel.UpdateImageTaskImages(imagetaskColl.ImageTaskId, imageColl.Md5)
	if err != nil {
		log.Error(fmt.Sprintf("update image task err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}
	if !strings.EqualFold(imagetaskColl.Introduce, introduce) {
		err = imagetaskmodel.UpdateImageTaskIntroduce(imagetaskColl.ImageTaskId, introduce)
		if err != nil {
			log.Error(fmt.Sprintf("update image task introduce err", err.Error()))
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "import image success",
	})
}

/*
func v2() {
	five, _ := thrfacemodel.ThrFaceFileRes(photoName, fileByte)
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
		if imageColl.ThrFaces == nil {
			imageColl.ThrFaces["face++"] = make(map[string]interface{})
		}
		imageColl.ThrFaces["face++"] = make(map[string]interface{})
		imageColl.ThrFaces["face++"]["83"] = result
	}

	//face++ res  V3
	var thrRes *thrfacemodel.FaceModelV3
	if strings.EqualFold(isRes, "not") {

		thrRes, err = thrfacemodel.ThrFaceFileResV3(photoName, fileByte)
		if err != nil {
			log.Error(fmt.Sprintf("get face++ res fail err:%s", err))
		}
		res1B, _ := json.Marshal(thrRes)

		var result interface{}
		if err := json.Unmarshal(res1B, &result); err != nil {
			fmt.Println("json unmarshal err=%s", err)
		}
		fmt.Println("-----result%s-----", result)
		if imageColl.ThrFaces == nil {
			imageColl.ThrFaces["face++"] = make(map[string]interface{})
		}
		imageColl.ThrFaces["face++"] = make(map[string]interface{})
		imageColl.ThrFaces["face++"]["83"] = result
	}
}
*/
