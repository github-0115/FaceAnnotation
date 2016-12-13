package task

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ImportTask(c *gin.Context) {

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Error(fmt.Sprintf("import parmar is err%v", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImportTaskParmars.Code,
			"message": vars.ErrImportTaskParmars.Msg,
		})
		return
	}
	fileName := fileHeader.Filename

	filebyte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(fmt.Sprintf("read import task file err%v", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrReadImportFile.Code,
			"message": vars.ErrReadImportFile.Msg,
		})
		return
	}

	taskModel, importModelList, err := taskmodel.ByteToTaskModel(filebyte)
	if err != nil {
		log.Error(fmt.Sprintf("read import task file error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrReadImportFile.Code,
			"message": vars.ErrReadImportFile.Msg,
		})
		return
	}
	taskModel.Title = "import_task_" + fileName

	err = taskModel.Save()
	if err != nil {
		log.Error(fmt.Sprintf("import task file save error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskSave.Code,
			"message": vars.ErrTaskSave.Msg,
		})
		return
	}

	for _, res := range importModelList {
		faceModel := &facemodel.FaceModel{
			Title:       taskModel.Title,
			Url:         res.Url,
			OriginFaces: res.OriginFaces,
		}
		_, err := facemodel.UpsertFaceResult(faceModel)
		if err != nil {
			log.Error(fmt.Sprintf("import task face third points upsert error:%s", err.Error()))
		}
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "import task success",
	})
}
