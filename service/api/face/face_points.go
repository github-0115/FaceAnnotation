package face

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	redismodel "FaceAnnotation/utils/redisclient"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func UpsertFacePoint(c *gin.Context) {
	face_str := c.PostForm("face")

	if face_str == "" {
		log.Error(fmt.Sprintf("face parmars error: face point is nil"))
		c.JSON(400, gin.H{
			"code":    vars.ErrFaceParmars.Code,
			"message": vars.ErrFaceParmars.Msg,
		})
		return
	}

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

	if faceModel.Title != "" {
		taskModel, err := taskmodel.QueryTask(faceModel.Title)
		if err != nil {
			log.Error(fmt.Sprintf("get task error:%s", err.Error()))
		}

		err = taskmodel.UpdateTaskImageStatus(faceModel.Title, faceModel.Url, 1)
		if err != nil {
			log.Error(fmt.Sprintf("update task status 1 error:%s", err.Error()))
		}

		if taskModel.Status == 0 {
			err := taskmodel.UpdateTaskStatus(faceModel.Title, 1)
			if err != nil {
				log.Error(fmt.Sprintf("update task status 1 error:%s", err.Error()))
			}
		}

		if isAllCompleted(taskModel) == "yes" {
			err := taskmodel.UpdateTaskStatus(faceModel.Title, 2)
			if err != nil {
				log.Error(fmt.Sprintf("update task status 2 error:%s", err.Error()))
			}
		}
	}

	err = redismodel.DeleteCheckImageStr(faceModel.Url)
	if err != nil {
		log.Error(fmt.Sprintf("image del redis err %S", err.Error()))
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "face points upsert success !",
	})
}

func isAllCompleted(a *taskmodel.TaskModel) string {
	for _, res := range a.Images {
		if res.Status == 0 {
			return "not"
		}
	}
	return "yes"
}
