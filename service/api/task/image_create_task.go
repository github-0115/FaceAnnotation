package task

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

type ImageTaskParmars struct {
	ImageTaskId string `json:"image_task_id"`
	PointType   int64  `json:"point_type"`
	MinUnit     int64  `json:"min_unit"`
	LimitUser   int64  `json:"limit_user"`
	Introduce   string `json:"introduce"` //本次说明
}

func ImageCreateTask(c *gin.Context) {

	var imageTaskParmars ImageTaskParmars
	if err := c.BindJSON(&imageTaskParmars); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	imageTaskId := imageTaskParmars.ImageTaskId
	minUnit := imageTaskParmars.MinUnit
	pointType := imageTaskParmars.PointType
	limitUser := imageTaskParmars.LimitUser
	introduce := imageTaskParmars.Introduce

	imagetaskColl, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

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
		Count:     int64(len(imagetaskColl.Images)),
		Introduce: introduce,
		Status:    0,
		CreatedAt: time.Now(),
	}

	//create smalltask
	stms := make([]*smalltaskmodel.SmallTaskModel, 0, 0)
	fmt.Println(math.Ceil(float64(len(imagetaskColl.Images)) / float64(taskColl.MinUnit)))
	for i := 0; i < int(math.Ceil(float64(len(imagetaskColl.Images))/float64(taskColl.MinUnit))); i++ {
		for _, res := range taskColl.Area {
			stm := &smalltaskmodel.SmallTaskModel{
				TaskId:      taskId,
				SmallTaskId: uuid.NewV4().String(),
				PointType:   taskColl.PointType,
				Areas:       res,
				LimitCount:  taskColl.LimitUser,
				Status:      0,
				CreatedAt:   time.Now(),
			}

			if (i+1)*int(taskColl.MinUnit) > len(imagetaskColl.Images) {
				stm.SmallTaskImages = imagetaskColl.Images[i*int(taskColl.MinUnit) : len(imagetaskColl.Images)]
			} else {
				stm.SmallTaskImages = imagetaskColl.Images[i*int(taskColl.MinUnit) : (i+1)*int(taskColl.MinUnit)]
			}

			stms = append(stms, stm)
		}
	}

	err = saveSmallTask(stms)
	if err != nil {
		log.Error(fmt.Sprintf("create small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskSave.Code,
			"message": vars.ErrSmallTaskSave.Msg,
		})
		return
	}

	for _, image := range imagetaskColl.Images {
		err := imagemodel.UpdateImageModel(image, taskId)
		if err != nil {
			log.Error(fmt.Sprintf("update image=%s taskid err=%s", image, err.Error()))
		}
	}

	err = imagetaskmodel.UpdateImageTaskTaskId(imagetaskColl.ImageTaskId, taskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task update taskId err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
	}

	err = taskColl.Save()
	if err != nil {
		log.Error(fmt.Sprintf("task save err", err.Error()))
		smalltaskmodel.RemoveSmallTask(taskColl.TaskId)
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskSave.Code,
			"message": vars.ErrTaskSave.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"task_id": taskId,
		"message": "create " + strconv.Itoa(len(stms)) + "small task success",
	})
}

func saveSmallTask(stms []*smalltaskmodel.SmallTaskModel) error {
	for _, res := range stms {
		err := res.Save()
		if err != nil {
			smalltaskmodel.RemoveSmallTask(res.TaskId)
			return err
		}
	}
	return nil
}
