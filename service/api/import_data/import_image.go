package import_data

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	uploadmodel "FaceAnnotation/service/model/uploadmodel"
	vars "FaceAnnotation/service/vars"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ImportImage(c *gin.Context) {
	isRes := c.PostForm("isRes")
	taskId := c.PostForm("task_id")
	imageFile, imageFileHeader, err := c.Request.FormFile("image")
	if err != nil || taskId == "" {
		log.Error(fmt.Sprintf("import image parmars err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImportImageParmars.Code,
			"message": vars.ErrImportImageParmars.Msg,
		})
		return
	}

	_, err = taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
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
			TaskId:    []string{taskId},
			Md5:       photoName,
			Url:       imageFileHeader.Filename,
			CreatedAt: time.Now(),
		}

		_, err = uploadmodel.UploadFile(photoName, fileByte)
		if err != nil {
			log.Error(fmt.Sprintf("image %s upload oss err %s", imageFileHeader.Filename, err.Error()))
		}

	} else {
		imageColl.TaskId = append(imageColl.TaskId, taskId)
	}
	//face++ res
	var thrRes *thrfacemodel.EightThreeFaceModel
	if strings.EqualFold(isRes, "not") {
		five, _ := thrfacemodel.ThrFaceFileRes(photoName, fileByte)
		if err != nil {
			log.Error(fmt.Sprintf("get face++ five res fail err:%s", err))
		}
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

	_, err = imagemodel.UpsertImageModel(imageColl)
	if err != nil {
		log.Error(fmt.Sprintf("image update err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelUpdate.Code,
			"message": vars.ErrImageModelUpdate.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"task_id": taskId,
	})
}
