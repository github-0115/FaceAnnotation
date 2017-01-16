package export_data

import (
	exportmodel "FaceAnnotation/service/model/exportmodel"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	//	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	//	"io/ioutil"
	//	"os"
	//	"path/filepath"
	"strconv"
	"strings"
	//	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ExportData(c *gin.Context) {
	//
	taskId := c.PostForm("task_id")

	task, err := taskmodel.QueryTask(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrTaskNotFound.Code,
			"message": vars.ErrTaskNotFound.Msg,
		})
		return
	}

	taskImages, err := imagemodel.QueryTaskImages(task.TaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	//	smallTask, err := smalltaskmodel.QueryfineTuneTask(task.TaskId, "fineTune")
	//	if err != nil {
	//		log.Error(fmt.Sprintf("query small task err %s", err))
	//		c.JSON(400, gin.H{
	//			"code":    vars.ErrSmallTaskNotFound.Code,
	//			"message": vars.ErrSmallTaskNotFound.Msg,
	//		})
	//		return
	//	}

	//	images := getFineCompleteImages(taskImages, smallTask.TaskId, task.PointType)

	for _, image := range taskImages {
		//daochu
		fmt.Println("export start......")
		//				exportmodel.SaveResFile(strings.Split(image.Url, ".")[0], image.FineResults[strconv.Itoa(int(smallTask.PointType))][0])
		//		exportmodel.SaveResFile(strings.Split(image.Url, ".")[0], image.ThrFaces["face++"])
		exportmodel.SaveResFile(strings.Split(image.Url, ".")[0], image)
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "expoet data success",
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
