package task

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type TaskRep struct {
	Title     string                  `json:"title"`
	AllCount  int64                   `json:"all_count"`
	NotCount  int64                   `json:"not_count"`
	Images    []*taskmodel.ImageModel `json:"images"`
	Status    int64                   `json:"status"` // 0 未开始  1 正在进行  2 已完成
	CreatedAt string                  `json:"created_at"`
}

func GetTaskList(c *gin.Context) {
	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Error(fmt.Sprintf("get task list parmars error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskParmars.Code,
			"message": vars.ErrTaskParmars.Msg,
		})
		return
	}

	task_list, err := taskmodel.QueryTaskList(int64(status))
	if err != nil {
		log.Error(fmt.Sprintf("get task list error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskListNotFound.Code,
			"message": vars.ErrTaskListNotFound.Msg,
		})
		return
	}

	var taskRep_list []*TaskRep
	if task_list == nil {
		taskRep_list = make([]*TaskRep, 0, 0)
	} else {
		taskRep_list = GetTaskRepList(task_list)
	}

	c.JSON(200, gin.H{
		"code":      0,
		"task_list": taskRep_list,
	})
}

func GetTaskRepList(t []*taskmodel.TaskModel) []*TaskRep {
	taskRep_list := make([]*TaskRep, 0, 0)
	for _, res := range t {
		rep := &TaskRep{
			Title:     res.Title,
			AllCount:  res.Count,
			Images:    res.Images,
			Status:    res.Status,
			CreatedAt: res.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		var count int64 = 0
		for _, res := range res.Images {
			if res.Status == 0 {
				count = count + 1
			}
		}
		rep.NotCount = count
		taskRep_list = append(taskRep_list, rep)
	}

	return taskRep_list
}
