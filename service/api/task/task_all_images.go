package task

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type TaskImagesRes struct {
	Md5       string                 `json:"md5"`
	Result    *imageend.AllPointsRep `json:"results"`
	ThrResult []*imageend.ThrResRep  `json:"thr_rep"`
}

var (
	imagesDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func GetTaskAllImages(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	taskId := c.Query("task_id")
	pageIndex, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("rows"))
	if err != nil {
		log.Error(fmt.Sprintf("strconv Atoi err%v", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	if strings.EqualFold(taskId, "") {
		log.Error(fmt.Sprintf("parmars nil err"))
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

	task, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	taskImages, err := imagemodel.QueryTaskImages(task.TaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task images err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	var count int64 = 0
	smallTask, err := smalltaskmodel.QueryfineTuneTask(task.TaskId, "fineTune")
	if err != nil {
		log.Error(fmt.Sprintf("query fineTune small task not found err %s", err))
	}

	if smallTask != nil {
		count = getHasPointImages(taskImages, smallTask.SmallTaskId, smallTask.PointType)
	}

	var results []*imagemodel.ImageModel
	if len(taskImages) < pageIndex*pageSize && len(taskImages) > (pageIndex-1)*pageSize {
		results = taskImages[(pageIndex-1)*pageSize : len(taskImages)]
	} else if len(taskImages) > pageIndex*pageSize {
		results = taskImages[(pageIndex-1)*pageSize : pageIndex*pageSize]
	}
	//	fmt.Println(results)

	rep := make([]*TaskImagesRes, 0, 0)
	for _, image := range results {
		if image.Results[strconv.Itoa(int(task.PointType))] == nil {
			continue
		}
		thrRes := imageend.GetThrResults(image)
		res := imageend.GetTaskNotFineResults(image, task)

		imRes := &TaskImagesRes{
			Md5:       imagesDomain + image.Md5,
			Result:    res,
			ThrResult: thrRes,
		}
		rep = append(rep, imRes)
	}

	total := int(math.Ceil(float64(task.Count) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":       0,
		"images":     rep,
		"count":      count,
		"page":       pageIndex,
		"total":      total,
		"created_at": task.CreatedAt.Format("2006-01-02 03:04:05"),
		"records":    len(taskImages),
	})
}

func getHasPointImages(imageList []*imagemodel.ImageModel, stmId string, pointType int64) int64 {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.FineResults[strconv.Itoa(int(pointType))]

		if len(result) == 0 || result == nil {
			continue
		}

		for _, res := range result {
			if strings.EqualFold(res.SmallTaskId, stmId) {
				list = append(list, image)
				break
			}
		}

	}

	return int64(len(list))
}
