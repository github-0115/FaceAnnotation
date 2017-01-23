package export_data

import (
	exportmodel "FaceAnnotation/service/model/exportmodel"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func ExportData(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
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

	smallTask, err := smalltaskmodel.QueryfineTuneTask(task.TaskId, "fineTune")
	if err != nil {
		log.Error(fmt.Sprintf("query fineTune small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrNoValidDataExport.Code,
			"message": vars.ErrNoValidDataExport.Msg,
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

	images := getFineCompleteImages(taskImages, smallTask.SmallTaskId, task.PointType)

	if len(images) == 0 {
		log.Error(fmt.Sprintf("no fineTune res err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrNoValidDataExport.Code,
			"message": vars.ErrNoValidDataExport.Msg,
		})
		return
	}

	//	exportmodel.ImageZip(images)
	res := make([]*exportmodel.Res, 0, 0)
	for _, image := range images {
		//daochu
		fmt.Println("export start......")
		points := switchPoints(image, smallTask.SmallTaskId, smallTask.PointType)
		if points == nil {
			continue
		}
		eRes := &exportmodel.Res{
			Name:   image.Url,
			Points: points,
		}
		res = append(res, eRes)

	}
	resName := "images(" + username + time.Now().Format("20060102030405") + ").zip"
	dataUrl, err := exportmodel.ImageDataZip(images, res, resName)
	if err != nil {
		log.Error(fmt.Sprintf("save export data err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrNoValidDataExport.Code,
			"message": vars.ErrNoValidDataExport.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":     0,
		"data_url": dataUrl,
		"message":  "expoet data success",
	})
}

func switchPoints(image *imagemodel.ImageModel, stmId string, pointType int64) []*imagemodel.Point {

	if image.FineResults[strconv.Itoa(int(pointType))] == nil {
		log.Info(fmt.Sprintf("image  = %s switch point", image.Md5))
		return nil
	}

	fines := image.FineResults[strconv.Itoa(int(pointType))]
	for _, fine := range fines {

		if strings.EqualFold(fine.SmallTaskId, stmId) {

			points := make([]*imagemodel.Point, 0, 0)
			for _, point := range fine.Result {
				if point != nil && len(point) != 0 {

					for _, res := range point {
						points = append(points, res)
					}
				}
			}
			return points
		}

	}

	return nil
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

/*
	for _, image := range images {
		//daochu
		fmt.Println("export start......")
		points := switchPoints(image, smallTask.SmallTaskId, smallTask.PointType)
		if points == nil {
			continue
		}
		eRes := &exportmodel.Res{
			Points: points,
		}

				err := exportmodel.ImageZip(image.Md5, image.Url)
				if err != nil {
					log.Info(fmt.Sprintf("image  = %s export save err = %s", image.Md5, err))
				}

				err = exportmodel.SaveImageRes(image.Url, points)
				if err != nil {
					log.Info(fmt.Sprintf("image  = %s export res save err = %s", image.Md5, err))
				}

	}*/
