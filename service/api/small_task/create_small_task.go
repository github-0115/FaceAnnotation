package small_task

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
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

func CreateSmallTask(c *gin.Context) {
	taskId := c.PostForm("task_id")

	taskColl, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	imageList, err := imagemodel.QueryTaskImages(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	imageMd5s := make([]string, 0, 0)
	for i := 0; i < len(imageList); i++ {
		imageMd5s = append(imageMd5s, imageList[i].Md5)
	}

	//create smalltask
	stms := make([]*smalltaskmodel.SmallTaskModel, 0, 0)
	for i := 0; i < int(math.Ceil(float64(len(imageMd5s))/float64(taskColl.MinUnit))); i++ {
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

			if (i+1)*int(taskColl.MinUnit) > len(imageMd5s) {
				stm.SmallTaskImages = imageMd5s[i*int(taskColl.MinUnit) : len(imageMd5s)]
			} else {
				stm.SmallTaskImages = imageMd5s[i*int(taskColl.MinUnit) : (i+1)*int(taskColl.MinUnit)]
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

	c.JSON(200, gin.H{
		"code":    0,
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
