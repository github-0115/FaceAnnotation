package task

import (
	//	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type TaskRep struct {
	Key         int64          `json:"key"`
	Mark        string         `json:"mark"`
	TaskId      string         `json:"task_id"`
	Children    []*TwoChildren `json:"children"`
	Situation   int64          `json:"situation"`
	Description string         `json:"description"`
	Status      int64          `json:"status"`
	CreatedAt   string         `json:"created_at"`
}

type TwoChildren struct {
	Key         int64            `json:"key"`
	Mark        string           `json:"mark"`
	TaskId      string           `json:"task_id"`
	Children    []*ThreeChildren `json:"children"`
	Situation   int64            `json:"situation"`
	Description string           `json:"description"`
	LimitCount  int64            `json:"limit_count"`
	Status      int64            `json:"status"`
	CreatedAt   string           `json:"created_at"`
}

type ThreeChildren struct {
	Key         int64  `json:"key"`
	Mark        string `json:"mark"`
	TaskId      string `json:"task_id"`
	Situation   int64  `json:"situation"`
	Description string `json:"description"`
	Area        string `json:"area"`
	LimitCount  int64  `json:"limit_count"`
	Status      int64  `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func TaskList(c *gin.Context) {
	pageIndex, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("rows"))
	//	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Error(fmt.Sprintf("strconv Atoi err%v", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	taskList := make([]*TaskRep, 0, 0)

	tasks, records, err := taskmodel.QueryPageTasks(pageIndex, pageSize)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskListNotFound.Code,
			"message": vars.ErrTaskListNotFound.Msg,
		})
		return
	}
	if tasks == nil {

		c.JSON(200, gin.H{
			"code":    0,
			"tasks":   taskList,
			"page":    pageIndex,
			"total":   0,
			"records": 0,
		})
		return
	}
	//	taskRep := &TaskRep{
	//		Key:  0,
	//		Mark: "one",
	//	}
	//	taskRep.Description = tasks[0].CreatedAt.Format("2006-01-02 03:04:05") + "导入数据"
	//	taskRep.TaskId = tasks[0].TaskId
	//	taskRep.Status = tasks[0].Status
	//	taskRep.Situation = tasks[0].Count
	//	taskRep.CreatedAt = tasks[0].CreatedAt.Format("2006-01-02 03:04:05")
	var key int64 = 1
	for i := 0; i < len(tasks); i++ {
		key += 1
		taskRep := &TaskRep{
			Key:         key,
			Mark:        "one",
			TaskId:      tasks[i].TaskId,
			Situation:   tasks[i].Count,
			Description: strconv.Itoa(int(tasks[i].PointType)),
			Status:      tasks[i].Status,
			CreatedAt:   tasks[i].CreatedAt.Format("2006-01-02 03:04:05"),
		}
		//		taskRep.Description = tasks[0].CreatedAt.Format("2006-01-02 03:04:05") + "导入数据"
		//		taskRep.TaskId = tasks[0].TaskId
		//		taskRep.Status = tasks[0].Status
		//		taskRep.Situation = tasks[0].Count
		//		taskRep.CreatedAt = tasks[0].CreatedAt.Format("2006-01-02 03:04:05")

		key += 1
		twoChildren := &TwoChildren{
			Key:         key,
			Mark:        "two",
			LimitCount:  tasks[i].LimitUser,
			TaskId:      tasks[i].TaskId,
			Situation:   tasks[i].Count,
			Description: strconv.Itoa(int(tasks[i].PointType)),
			Status:      tasks[i].Status,
			CreatedAt:   tasks[i].CreatedAt.Format("2006-01-02 03:04:05"),
		}

		smallTasks, err := smalltaskmodel.QueryTaskAllSmallTasks(tasks[i].TaskId)
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			//			c.JSON(400, gin.H{
			//				"code":    vars.ErrSmallTaskNotFound.Code,
			//				"message": vars.ErrSmallTaskNotFound.Msg,
			//			})
			//			return
		}

		if smallTasks == nil {
			threeChildren := &ThreeChildren{}
			twoChildren.Children = append(twoChildren.Children, threeChildren)
			taskRep.Children = append(taskRep.Children, twoChildren)
			continue
		}

		for j := 0; j < len(smallTasks); j++ {
			key += 1
			threeChildren := &ThreeChildren{
				Key:         key,
				Mark:        "three",
				LimitCount:  smallTasks[j].LimitCount,
				TaskId:      smallTasks[j].SmallTaskId,
				Situation:   int64(len(smallTasks[j].SmallTaskImages)),
				Area:        smallTasks[j].Areas,
				Description: strconv.Itoa(int(smallTasks[j].PointType)) + smallTasks[j].Areas,
				Status:      smallTasks[j].Status,
				CreatedAt:   smallTasks[j].CreatedAt.Format("2006-01-02 03:04:05"),
			}

			twoChildren.Children = append(twoChildren.Children, threeChildren)
		}

		taskRep.Children = append(taskRep.Children, twoChildren)
		taskList = append(taskList, taskRep)
	}

	total := int(math.Ceil(float64(records) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":    0,
		"tasks":   taskList,
		"page":    pageIndex,
		"total":   total,
		"records": records,
	})
}
