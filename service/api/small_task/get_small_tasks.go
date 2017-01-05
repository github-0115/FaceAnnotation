package small_task

import (
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetSmallTasks(c *gin.Context) {

	smalltaskList, err := smalltaskmodel.QueryNotSmallTask()
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	if smalltaskList == nil {
		smalltaskList = make([]*smalltaskmodel.SmallTaskModel, 0, 0)
	}

	c.JSON(200, gin.H{
		"code":       0,
		"small_task": smalltaskList,
	})
}
