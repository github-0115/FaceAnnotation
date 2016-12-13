package face

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	redismodel "FaceAnnotation/utils/redisclient"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type FaceParams struct {
	Title       string                   `json:"title"`
	User        string                   `json:"user"`
	Url         string                   `json:"url"`
	Faces       []*facemodel.FacesPoints `json:"faces"`
	OriginFaces []*facemodel.Landmarks   `json:"origin_faces"`
}

func UpsertFacePoint(c *gin.Context) {

	var faceParams FaceParams
	if err := c.BindJSON(&faceParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	title := faceParams.Title
	user := faceParams.User
	url := faceParams.Url
	faces := faceParams.Faces
	originFaces := faceParams.OriginFaces

	faceModel := &facemodel.FaceModel{
		Title:       title,
		User:        user,
		Url:         url,
		Faces:       faces,
		OriginFaces: originFaces,
	}

	if len(faceModel.Faces) == 0 {
		log.Error(fmt.Sprintf("face points len = 0 "))
	}

	_, err := facemodel.UpsertFaceResult(faceModel)
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
		url_strs := strings.Split(faceModel.Url, "/")
		err = taskmodel.UpdateTaskImageStatus(faceModel.Title, url_strs[len(url_strs)-1], 1)
		if err != nil {
			log.Error(fmt.Sprintf("update task image status 1 error:%s", err.Error()))
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
