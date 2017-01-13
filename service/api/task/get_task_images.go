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
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type TaskImagesRep struct {
	TaskId    []string              `json:"task_id"`
	Md5       string                `json:"md5"`
	Area      string                `json:"area"`
	PointType string                `json:"point_type"`
	Result    []*imagemodel.Points  `json:"rep"`
	ThrResult []*imageend.ThrResRep `json:"thr_rep"`
	Status    int64                 `json:"status"`
	CreatedAt time.Time             `json:"created_at" `
}

var (
	imageDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func GetTaskImages(c *gin.Context) {
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
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	smallTask, err := smalltaskmodel.QueryfineTuneTask(taskId, "fineTune")
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	taskImages, err := imagemodel.QueryTaskImages(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	images := getFineCompleteImages(taskImages, smallTask.TaskId, task.PointType)

	var results []*imagemodel.ImageModel
	if len(images) < pageIndex*pageSize && len(images) > (pageIndex-1)*pageSize {
		results = images[(pageIndex-1)*pageSize : len(images)]
	} else if len(images) > pageIndex*pageSize {
		results = images[(pageIndex-1)*pageSize : pageIndex*pageSize]
	}

	rep := make([]*TaskImagesRep, 0, 0)
	for _, image := range results {
		thr_Res := imageend.ThrResults(task.PointType, image)
		imRep := &TaskImagesRep{
			TaskId:    image.TaskId,
			Area:      smallTask.Areas,
			Md5:       imageDomain + image.Md5,
			PointType: strconv.Itoa(int(smallTask.PointType)),
			Result:    switchPoints(image, task.PointType),
			ThrResult: thr_Res,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
		}
		rep = append(rep, imRep)
	}

	total := int(math.Ceil(float64(len(taskImages)) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":     0,
		"tasks":    rep,
		"complete": len(images),
		"page":     pageIndex,
		"total":    total,
		"records":  len(taskImages),
	})
}

func switchPoints(image *imagemodel.ImageModel, pointType int64) []*imagemodel.Points {
	fineRes := make([]*imagemodel.Points, 0, 0)
	if image.FineResults[strconv.Itoa(int(pointType))] == nil {
		log.Info(fmt.Sprintf("image  = %s switch point", image.Md5))
		return fineRes
	}

	fines := image.FineResults[strconv.Itoa(int(pointType))]
	for _, fine := range fines {
		points := make([]*imagemodel.Point, 0, 0)
		//		var p *imagemodel.Point
		for _, point := range fine.Result {
			if point != nil {

				for _, res := range point {
					points = append(points, res)
				}

			}
		}

		pointsRes := &imagemodel.Points{
			SmallTaskId: fine.SmallTaskId,
			User:        fine.User,
			Points:      points,
			Sys:         fine.Sys,
			CreatedAt:   fine.CreatedAt,
			FinishedAt:  fine.FinishedAt,
		}

		fineRes = append(fineRes, pointsRes)
	}

	return fineRes
}

func getFineCompleteImages(imageList []*imagemodel.ImageModel, stmId string, pointType int64) []*imagemodel.ImageModel {
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
	return list
}
