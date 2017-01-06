package import_data

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	uploadmodel "FaceAnnotation/service/model/uploadmodel"
	vars "FaceAnnotation/service/vars"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ImportImage(c *gin.Context) {

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

		err := imageColl.Save()
		if err != nil {
			log.Error(fmt.Sprintf("image save err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageModelSave.Code,
				"message": vars.ErrImageModelSave.Msg,
			})
			return
		}

		_, err = uploadmodel.UploadFile(photoName, fileByte)
		if err != nil {
			log.Error(fmt.Sprintf("image %s upload oss err %s", imageFileHeader.Filename, err.Error()))
		}

	} else {
		err = imagemodel.UpdateImageModel(photoName, taskId)
		if err != nil {
			log.Error(fmt.Sprintf("image update taskid err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageModelUpdate.Code,
				"message": vars.ErrImageModelUpdate.Msg,
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": taskId,
	})
}
