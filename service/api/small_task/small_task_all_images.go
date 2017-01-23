package small_task

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type SmallTaskImagesRes struct {
	Md5       string                 `json:"md5"`
	Result    *imageend.AllPointsRep `json:"results"`
	ThrResult []*imageend.ThrResRep  `json:"thr_rep"`
}

var (
	imagesDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func GetSmallTaskAllImages(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	smallTaskId := c.Query("task_id")
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

	if strings.EqualFold(smallTaskId, "") {
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

	smallTaskModel, err := smalltaskmodel.QuerySmallTask(smallTaskId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	smallTaskImages, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("query small task images err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	images := getHasPointImages(smallTaskImages, smallTaskModel) // area string, limitCount int64, stmId string, pointType int64

	var results []*imagemodel.ImageModel
	if len(images) < pageIndex*pageSize && len(images) > (pageIndex-1)*pageSize {
		results = images[(pageIndex-1)*pageSize : len(images)]
	} else if len(images) > pageIndex*pageSize {
		results = images[(pageIndex-1)*pageSize : pageIndex*pageSize]
	}

	rep := make([]*SmallTaskImagesRes, 0, 0)
	for _, image := range results {
		if image.Results[strconv.Itoa(int(smallTaskModel.PointType))] == nil {
			continue
		}
		thrRes := imageend.GetThrResults(image)
		res := imageend.GetSmallTaskNotFineResults(image, smallTaskModel)

		imRes := &SmallTaskImagesRes{
			Md5:       imagesDomain + image.Md5,
			Result:    res,
			ThrResult: thrRes,
		}
		rep = append(rep, imRes)
	}
	notCount := getNotCompleteImageCount(smallTaskImages, smallTaskModel.Areas, smallTaskModel.LimitCount, smallTaskModel.SmallTaskId, smallTaskModel.PointType)
	total := int(math.Ceil(float64(len(smallTaskModel.SmallTaskImages)) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":       0,
		"page":       pageIndex,
		"total":      total,
		"count":      len(smallTaskModel.SmallTaskImages) - int(notCount),
		"records":    len(smallTaskModel.SmallTaskImages),
		"images":     rep,
		"point_type": smallTaskModel.PointType,
		"area":       smallTaskModel.Areas,
	})
}

func getHasPointImages(imageList []*imagemodel.ImageModel, smallTaskModel *smalltaskmodel.SmallTaskModel) []*imagemodel.ImageModel {
	list := make([]*imagemodel.ImageModel, 0, 0)

	for _, image := range imageList {
		result := image.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas]

		if len(result) == 0 || result == nil {
			continue
		}

		for _, res := range result {
			if strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
				list = append(list, image)
				break
			}
		}

	}

	return list
}

//func getNotCompleteImageCount(imageList []*imagemodel.ImageModel, area string, limitCount int64, stmId string, pointType int64) int64 {
//	list := make([]*imagemodel.ImageModel, 0, 0)
//	for _, image := range imageList {
//		result := image.Results[strconv.Itoa(int(pointType))][area]

//		if len(result) == 0 || result == nil {
//			list = append(list, image)
//			continue
//		}

//		var (
//			count int64 = 0
//		)
//		for _, res := range result {

//			if strings.EqualFold(res.SmallTaskId, stmId) {
//				count += 1
//			}
//		}

//		if count < limitCount {
//			list = append(list, image)
//		}
//	}
//	return int64(len(list))
//}
