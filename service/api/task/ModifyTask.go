package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type ModifyTaskParmars struct {
	TaskId    string `json:"task_id"`
	Introduce string `json:"introduce"` //本次说明
}

func ModifyTask(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)

	var taskParmars ModifyTaskParmars
	if err := c.BindJSON(&taskParmars); err != nil {
		log.Error(fmt.Sprintf("modify task bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	taskId := taskParmars.TaskId
	introduce := taskParmars.Introduce

	if introduce == "" || taskId == "" {
		log.Error(fmt.Sprintf("modify task parmars nil"))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	_, err := usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	taskColl, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return

	}

	err = taskmodel.UpdateTaskIntroduce(taskColl.TaskId, introduce)
	if err != nil {
		log.Error(fmt.Sprintf("stop task update status err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "modify task success",
	})
}
