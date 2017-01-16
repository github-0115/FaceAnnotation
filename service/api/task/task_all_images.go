package task

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	taskmodel "FaceAnnotation/service/model/taskmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
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

	taskImages, err := imagemodel.QueryPageTaskImages(taskId, pageIndex, pageSize)
	if err != nil {
		log.Error(fmt.Sprintf("query task images err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	rep := make([]*TaskImagesRep, 0, 0)
	for _, image := range taskImages {
		thr_Res := imageend.ThrResults(task.PointType, image)
		results := taskSwitchPoints(image, task.PointType)
		var status int64 = 1
		if results == nil {
			results = make([]*imagemodel.Points, 0, 0)
			status = 0
		}
		imRep := &TaskImagesRep{
			TaskId:    image.TaskId,
			Area:      "fineTune",
			Md5:       imagesDomain + image.Md5,
			PointType: strconv.Itoa(int(task.PointType)),
			Result:    results,
			ThrResult: thr_Res,
			Status:    status,
			CreatedAt: task.CreatedAt,
		}
		rep = append(rep, imRep)
	}

	total := int(math.Ceil(float64(task.Count) / float64(pageSize)))

	c.JSON(200, gin.H{
		"code":       0,
		"tasks":      rep,
		"count":      task.Count,
		"page":       pageIndex,
		"total":      total,
		"created_at": task.CreatedAt.Format("2006-01-02 03:04:05"),
		"records":    len(taskImages),
	})
}
