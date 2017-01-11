package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

type TaskParmars struct {
	Count int64 `json:"count"`
	//	Area      []string `json:"area"`
	PointType int64  `json:"point_type"`
	MinUnit   int64  `json:"min_unit"`
	LimitUser int64  `json:"limit_user"`
	Introduce string `json:"introduce"` //本次说明
}

func CreateTask(c *gin.Context) {

	var taskParmars TaskParmars
	if err := c.BindJSON(&taskParmars); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	count := taskParmars.Count
	//	area := taskParmars.Area
	minUnit := taskParmars.MinUnit
	pointType := taskParmars.PointType
	limitUser := taskParmars.LimitUser
	introduce := taskParmars.Introduce

	taskId := uuid.NewV4().String()
	taskColl, err := taskmodel.QueryTask(taskId)
	if err != nil {
		if err != taskmodel.ErrTaskModelNotFound {
			log.Error(fmt.Sprintf("task create err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrTaskExist.Code,
				"message": vars.ErrTaskExist.Msg,
			})
			return
		}
	}

	var areas []string
	switch pointType {
	case 5:
		areas = []string{"leftEye", "rightEye", "mouth", "nouse"}
	case 27:
		areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "mouth", "nouse", "face"}
	case 68:
		areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "mouth", "nouse", "face"}
	case 83:
		areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "mouth", "nouse", "face"}
	case 95:
		areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
	default:
		areas = []string{"leftEyebrow", "rightEyebrow", "leftEye", "rightEye", "leftEar", "rightEar", "mouth", "nouse", "face"}
	}

	taskColl = &taskmodel.TaskModel{
		TaskId:    taskId,
		Area:      areas,
		MinUnit:   minUnit,
		PointType: pointType,
		LimitUser: limitUser,
		Count:     count,
		Introduce: introduce,
		Status:    0,
		CreatedAt: time.Now(),
	}

	err = taskColl.Save()
	if err != nil {
		log.Error(fmt.Sprintf("task save err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskSave.Code,
			"message": vars.ErrTaskSave.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"task_id": taskId,
	})
}
