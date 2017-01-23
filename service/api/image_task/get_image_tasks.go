package image_task

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"strings"

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
	Situation   string           `json:"situation"`
	Description string           `json:"description"`
	LimitCount  int64            `json:"limit_count"`
	Status      int64            `json:"status"`
	CreatedAt   string           `json:"created_at"`
}

type ThreeChildren struct {
	Key         int64  `json:"key"`
	Mark        string `json:"mark"`
	TaskId      string `json:"task_id"`
	Situation   string `json:"situation"`
	Description string `json:"description"`
	Area        string `json:"area"`
	LimitCount  int64  `json:"limit_count"`
	Status      int64  `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func ImageTaskList(c *gin.Context) {
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

	imageTasks, records, err := imagetaskmodel.QueryPageImageTasks(pageIndex, pageSize)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskListNotFound.Code,
			"message": vars.ErrTaskListNotFound.Msg,
		})
		return
	}
	if imageTasks == nil {

		c.JSON(200, gin.H{
			"code":    0,
			"tasks":   taskList,
			"page":    pageIndex,
			"total":   0,
			"records": 0,
		})
		return
	}

	var key int64 = 1
	for i := 0; i < len(imageTasks); i++ {
		key += 1
		taskRep := &TaskRep{
			Key:         key,
			Mark:        "one",
			TaskId:      imageTasks[i].ImageTaskId,
			Situation:   int64(len(imageTasks[i].Images)),
			Description: imageTasks[i].Introduce,
			Status:      imageTasks[i].Status,
			CreatedAt:   imageTasks[i].CreatedAt.Format("2006-01-02 03:04:05"),
		}

		tasks, err := taskmodel.GetImageTasks(imageTasks[i].TaskId)
		if err != nil {
			log.Error(fmt.Sprintf("query task err %s", err))
		}

		if tasks == nil {
			//			twoChildren := &TwoChildren{}
			//			threeChildren := &ThreeChildren{}
			//			twoChildren.Children = append(twoChildren.Children, threeChildren)
			//			taskRep.Children = append(taskRep.Children, twoChildren)
			taskList = append(taskList, taskRep)
			continue
		}
		for _, task := range tasks {

			key += 1
			twoChildren := &TwoChildren{
				Key:         key,
				Mark:        "two",
				LimitCount:  task.LimitUser,
				TaskId:      task.TaskId,
				Situation:   strconv.Itoa(0) + "/" + strconv.Itoa(len(imageTasks[i].Images)),
				Description: strconv.Itoa(int(task.PointType)),
				Status:      task.Status,
				CreatedAt:   task.CreatedAt.Format("2006-01-02 03:04:05"),
			}

			smallTasks, err := smalltaskmodel.QueryTaskAllSmallTasks(task.TaskId)
			if err != nil {
				log.Error(fmt.Sprintf("query small task err %s", err))
			}

			if smallTasks == nil {
				//				threeChildren := &ThreeChildren{}
				//				twoChildren.Children = append(twoChildren.Children, threeChildren)
				taskRep.Children = append(taskRep.Children, twoChildren)
				continue
			}
			for j := 0; j < len(smallTasks); j++ {
				var count int64 = 0
				if smallTasks[j].Status == smalltaskmodel.TaskStatus.Success {
					count = int64(len(smallTasks[j].SmallTaskImages))
				} else {
					imageList, err := imagemodel.GetSmallTaskImages(smallTasks[j].SmallTaskImages)
					if err != nil {
						log.Error(fmt.Sprintf("task images query err", err.Error()))
						continue
					}
					notCount := getNotCompleteImageCount(imageList, smallTasks[j].Areas, smallTasks[j].LimitCount, smallTasks[j].SmallTaskId, smallTasks[j].PointType)
					count = int64(len(smallTasks[j].SmallTaskImages)) - notCount
				}

				key += 1
				threeChildren := &ThreeChildren{
					Key:         key,
					Mark:        "three",
					LimitCount:  smallTasks[j].LimitCount,
					TaskId:      smallTasks[j].SmallTaskId,
					Description: strconv.Itoa(int(smallTasks[j].PointType)) + smallTasks[j].Areas,
					Situation:   strconv.Itoa(int(count)) + "/" + strconv.Itoa(len(smallTasks[j].SmallTaskImages)),
					Area:        smallTasks[j].Areas,
					Status:      smallTasks[j].Status,
					CreatedAt:   smallTasks[j].CreatedAt.Format("2006-01-02 03:04:05"),
				}

				if strings.EqualFold(smallTasks[j].Areas, "fineTune") {
					twoChildren.Situation = strconv.Itoa(int(count)) + "/" + strconv.Itoa(len(smallTasks[j].SmallTaskImages))
				}
				twoChildren.Children = append(twoChildren.Children, threeChildren)
			}

			taskRep.Children = append(taskRep.Children, twoChildren)
		}

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

func getNotCompleteImageCount(imageList []*imagemodel.ImageModel, area string, limitCount int64, stmId string, pointType int64) int64 {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.Results[strconv.Itoa(int(pointType))][area]

		if len(result) == 0 || result == nil {
			list = append(list, image)
			continue
		}

		var (
			count int64 = 0
		)
		for _, res := range result {

			if strings.EqualFold(res.SmallTaskId, stmId) {
				count += 1
			}
		}

		if count < limitCount {
			list = append(list, image)
		}
	}
	return int64(len(list))
}
