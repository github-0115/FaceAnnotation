package small_task

import (
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	vars "FaceAnnotation/service/vars"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetSmallTasks(c *gin.Context) {
	area := c.Query("area")

	var (
		smalltaskList []*smalltaskmodel.SmallTaskModel
		err           = errors.New("")
	)
	if area == "" {
		smalltaskList, err = smalltaskmodel.QueryNotSmallTask()
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}

	} else {
		smalltaskList, err = smalltaskmodel.QueryAreaNotSmallTask(area)
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}
	}

	if smalltaskList == nil {
		log.Error(fmt.Sprintf("not small task to allot err %s", err))
		c.JSON(200, gin.H{
			"code":    vars.ErrNotSmallTask.Code,
			"message": vars.ErrNotSmallTask.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":        0,
		"small_tasks": smalltaskList,
		"small_task":  smalltaskList[0].SmallTaskId,
	})
}
