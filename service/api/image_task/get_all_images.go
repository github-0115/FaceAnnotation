package image_task

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	imagetaskmodel "FaceAnnotation/service/model/imagetaskmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type ImagesRep struct {
	Md5       string                 `json:"md5"`
	Results   *imageend.AllPointsRep `json:"results"`
	ThrResult []*imageend.ThrResRep  `json:"thr_rep"`
}

var (
	imagesDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func GetAllImages(c *gin.Context) {
	imageTaskId := c.Query("task_id")
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
	if strings.EqualFold(imageTaskId, "") {
		log.Error(fmt.Sprintf("parmars nil err"))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	imageTaskColl, err := imagetaskmodel.QueryImageTask(imageTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("image task not found err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageTaskNotFound.Code,
			"message": vars.ErrImageTaskNotFound.Msg,
		})
		return
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

	rep := make([]*ImagesRep, 0, 0)
	for _, image := range results {
		thrRes := imageend.GetThrResults(image)
		res := imageend.GetAllResults(image)
		imagesRep := &ImagesRep{
			Md5:       imagesDomain + image.Md5,
			Results:   res,
			ThrResult: thrRes,
		}

		rep = append(rep, imagesRep)
	}

	count := getCompleteImageCount(images)
	total := int(math.Ceil(float64(len(imageTaskColl.Images)) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":       0,
		"page":       pageIndex,
		"total":      total,
		"count":      count,
		"records":    len(imageTaskColl.Images),
		"images":     rep,
		"created_at": imageTaskColl.CreatedAt.Format("2006-01-02 03:04:05"),
	})
}

func getCompleteImageCount(imageList []*imagemodel.ImageModel) int64 {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {

		if image.FineResults["95"] != nil || len(image.FineResults["95"]) == 0 {
			list = append(list, image)
			continue
		}

		if image.FineResults["83"] != nil || len(image.FineResults["83"]) == 0 {
			list = append(list, image)
			continue
		}

		if image.FineResults["68"] != nil || len(image.FineResults["68"]) == 0 {
			list = append(list, image)
			continue
		}

		if image.FineResults["27"] != nil || len(image.FineResults["27"]) == 0 {
			list = append(list, image)
			continue
		}

		if image.FineResults["5"] != nil || len(image.FineResults["5"]) == 0 {
			list = append(list, image)
			continue
		}

	}
	return int64(len(list))
}
