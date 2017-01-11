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

type ImageRep struct {
	TaskId      []string             `json:"task_id"`
	SmallTaskId string               `json:"small_task_id"`
	Md5         string               `json:"md5"`
	Area        string               `json:"area"`
	PointType   string               `json:"point_type"`
	Result      []*imagemodel.Points `json:"rep"`
	ThrResult   *imageend.PointsRep  `json:"thr_rep"`
}

var (
	imageDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func GetSmallTaskImages(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	stId := c.Query("small_task_id")
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

	smallTaskModel, err := smalltaskmodel.QuerySmallTask(stId)
	if err != nil {
		log.Error(fmt.Sprintf("query small task err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	imageList, err := imagemodel.GetSmallTaskImages(smallTaskModel.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrSmallTaskNotFound.Code,
			"message": vars.ErrSmallTaskNotFound.Msg,
		})
		return
	}

	images := getCompleteImages(imageList, smallTaskModel.Areas, stId, smallTaskModel.PointType)
	var results []*imagemodel.ImageModel
	if len(images) < pageIndex*pageSize && len(images) > (pageIndex-1)*pageSize {
		results = images[(pageIndex-1)*pageSize : len(images)]
	} else if len(images) > pageIndex*pageSize {
		results = images[(pageIndex-1)*pageSize : pageIndex*pageSize]
	}

	rep := make([]*ImageRep, 0, 0)
	for _, image := range results {
		thr_Res := imageend.SwitchPoint(smallTaskModel.PointType, image)
		imRep := &ImageRep{
			TaskId:      image.TaskId,
			SmallTaskId: stId,
			Area:        smallTaskModel.Areas,
			Md5:         imageDomain + image.Md5,
			PointType:   strconv.Itoa(int(smallTaskModel.PointType)),
			Result:      image.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas],
			ThrResult:   thr_Res,
		}
		rep = append(rep, imRep)
	}

	total := int(math.Ceil(float64(len(smallTaskModel.SmallTaskImages)) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":     0,
		"res":      rep,
		"complete": len(images),
		"page":     pageIndex,
		"total":    total,
		"records":  len(smallTaskModel.SmallTaskImages),
	})
}

func getCompleteImages(imageList []*imagemodel.ImageModel, area string, stmId string, pointType int64) []*imagemodel.ImageModel {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.Results[strconv.Itoa(int(pointType))][area]

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
