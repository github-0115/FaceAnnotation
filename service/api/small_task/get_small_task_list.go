package small_task

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type SmallTaskRep struct {
	TaskId         string    `json:"task_id"`       //与 taskmodel 关联
	SmallTaskId    string    `json:"small_task_id"` //与 taskmodel 关联
	CompleteImages int64     `json:"complete_images"`
	TotalImages    int64     `json:"total_images"`
	PointType      int64     `json:"point_type"`
	Area           string    `json:"area"`        //标识标注的哪个部位
	LimitCount     int64     `json:"limit_count"` //人数限制
	Status         int64     `json:"status"`      // 0 创建成功  1  正在进行中 2 已标注完成
	CreatedAt      time.Time `json:"created_at"`
}

func SmallTaskList(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	pageIndex, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("rows"))
	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Error(fmt.Sprintf("strconv Atoi err%v", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}
	_, err = usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	var (
		smallTasks []*smalltaskmodel.SmallTaskModel
		records    int = 0
	)
	if status == 2 {
		smallTasks, records, err = smalltaskmodel.QueryPageSmallTasks(int64(status), pageIndex, pageSize)
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}
	} else {
		smallTasks, records, err = smalltaskmodel.QueryPageNotSmallTask(pageIndex, pageSize)
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}
	}
	if smallTasks == nil {
		smallTasks = make([]*smalltaskmodel.SmallTaskModel, 0, 0)
	}

	smallTaskList := make([]*SmallTaskRep, 0, 0)
	for _, smallTask := range smallTasks {

		rep := &SmallTaskRep{
			TaskId:      smallTask.TaskId,
			SmallTaskId: smallTask.SmallTaskId,
			TotalImages: int64(len(smallTask.SmallTaskImages)),
			PointType:   smallTask.PointType,
			Area:        smallTask.Areas,
			LimitCount:  smallTask.LimitCount,
			Status:      smallTask.Status,
			CreatedAt:   smallTask.CreatedAt,
		}

		if smallTask.Status == smalltaskmodel.TaskStatus.Success {
			rep.CompleteImages = int64(len(smallTask.SmallTaskImages))
			smallTaskList = append(smallTaskList, rep)
			continue
		}

		imageList, err := imagemodel.GetSmallTaskImages(smallTask.SmallTaskImages)
		if err != nil {
			log.Error(fmt.Sprintf("image query err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}

		count := getNotCompleteImageCount(imageList, smallTask.Areas, smallTask.LimitCount, smallTask.SmallTaskId, smallTask.PointType)
		rep.CompleteImages = int64(len(smallTask.SmallTaskImages)) - count
		smallTaskList = append(smallTaskList, rep)
	}

	total := int(math.Ceil(float64(records) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":    0,
		"tasks":   smallTaskList,
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
