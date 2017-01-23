package task

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
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

func GetTaskImages(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	taskId := c.Query("task_id")
	flag := c.DefaultQuery("flag", "not")
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

	if strings.EqualFold(flag, "all") {
		//
		imageTaskColl, err := imagetaskmodel.QueryImageTask(taskId)
		if err != nil {
			log.Error(fmt.Sprintf("image task not found err", err.Error()))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageTaskNotFound.Code,
				"message": vars.ErrImageTaskNotFound.Msg,
			})
			return
		}

		tasks, err := taskmodel.GetImageTasks(imageTaskColl.TaskId)
		if err != nil {
			log.Error(fmt.Sprintf("query task err %s", err))
		}

		if tasks == nil {
			rep := make([]*TaskImagesRep, 0, 0)
			c.JSON(200, gin.H{
				"code":       0,
				"tasks":      rep,
				"count":      0,
				"page":       pageIndex,
				"total":      len(imageTaskColl.Images),
				"created_at": imageTaskColl.CreatedAt.Format("2006-01-02 03:04:05"),
				"records":    0,
			})
			return
		}
		//		pointTypes := make([]int64, 0, 0)
		var pointType int64 = 5
		for _, task := range tasks {
			//			pointTypes = append(pointTypes, task.PointType)
			if task.PointType > pointType {
				pointType = task.PointType
			}
		}

		images, err := imagemodel.GetSmallTaskImages(imageTaskColl.Images)
		if err != nil {
			log.Error(fmt.Sprintf("query task images err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrImageModelNotFound.Code,
				"message": vars.ErrImageModelNotFound.Msg,
			})
			return
		}

		var results []*imagemodel.ImageModel
		if len(images) < pageIndex*pageSize && len(images) > (pageIndex-1)*pageSize {
			results = images[(pageIndex-1)*pageSize : len(images)]
		} else if len(images) > pageIndex*pageSize {
			results = images[(pageIndex-1)*pageSize : pageIndex*pageSize]
		}

		rep := make([]*TaskImagesRep, 0, 0)
		for _, image := range results {

			imRep := &TaskImagesRep{
				TaskId:    image.TaskId,
				Md5:       imagesDomain + image.Md5,
				PointType: strconv.Itoa(int(pointType)),
				Status:    imageTaskColl.Status,
				CreatedAt: imageTaskColl.CreatedAt,
			}
			rep = append(rep, imRep)

		}
		total := int(math.Ceil(float64(len(imageTaskColl.Images)) / float64(pageSize)))
		c.JSON(200, gin.H{
			"code":       0,
			"tasks":      rep,
			"count":      len(imageTaskColl.Images),
			"page":       pageIndex,
			"total":      total,
			"created_at": imageTaskColl.CreatedAt.Format("2006-01-02 03:04:05"),
			"records":    len(images),
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

	taskImages, err := imagemodel.QueryTaskImages(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query task images err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	smallTask, err := smalltaskmodel.QueryfineTuneTask(taskId, "fineTune")
	if err != nil {
		log.Error(fmt.Sprintf("query fineTune small task err %s", err))
		rep := make([]*TaskImagesRep, 0, 0)
		c.JSON(200, gin.H{
			"code":       0,
			"tasks":      rep,
			"count":      task.Count,
			"page":       pageIndex,
			"total":      1,
			"created_at": task.CreatedAt.Format("2006-01-02 03:04:05"),
			"records":    0,
		})
		return
	}

	images := taskGetFineCompleteImages(taskImages, smallTask.SmallTaskId, task.PointType)

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
			Md5:       imagesDomain + image.Md5,
			PointType: strconv.Itoa(int(smallTask.PointType)),
			Result:    taskSwitchPoints(image, task.PointType),
			ThrResult: thr_Res,
			Status:    smallTask.Status,
			CreatedAt: task.CreatedAt,
		}
		rep = append(rep, imRep)
	}

	total := int(math.Ceil(float64(len(images)) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":       0,
		"tasks":      rep,
		"count":      len(taskImages),
		"page":       pageIndex,
		"total":      total,
		"created_at": task.CreatedAt.Format("2006-01-02 03:04:05"),
		"records":    len(images),
	})
}

func taskSwitchPoints(image *imagemodel.ImageModel, pointType int64) []*imagemodel.Points {
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

func taskGetFineCompleteImages(imageList []*imagemodel.ImageModel, stmId string, pointType int64) []*imagemodel.ImageModel {
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
