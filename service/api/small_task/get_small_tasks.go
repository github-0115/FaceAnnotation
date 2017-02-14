package small_task

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	smalltaskmodel "FaceAnnotation/service/model/smalltaskmodel"
	timeoutmodel "FaceAnnotation/service/model/timeoutmodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

var (
	s = rand.NewSource(time.Now().Unix())
	r = rand.New(s)
)

func GetSmallTasks(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	area := c.Query("area")

	userColl, err := usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	//fineTune
	if strings.EqualFold(userColl.Identity, usermodel.UserIdentity.FineTune) {

		timeOutTask, err := getTimeOutFineImage(username)
		if timeOutTask != "" {
			c.JSON(200, gin.H{
				"code":          0,
				"small_task_id": timeOutTask,
			})
			return
		}

		smalltaskList, err := smalltaskmodel.QueryAreaNotSmallTask("fineTune")
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}

		smallTaskId, err := getFineTune(username, smalltaskList)
		if strings.EqualFold(smallTaskId, "") {
			log.Error(fmt.Sprintf("not fintune small task to allot err %s", err))
			c.JSON(200, gin.H{
				"code":    vars.ErrNotSmallTask.Code,
				"message": vars.ErrNotSmallTask.Msg,
			})
			return
		}

		c.JSON(200, gin.H{
			"code":          0,
			"small_task_id": smallTaskId,
		})
		return
	}

	//	timeOutTask, err := getTimeOutImage(username)
	//	if timeOutTask != "" {
	//		c.JSON(200, gin.H{
	//			"code":          0,
	//			"small_task_id": timeOutTask,
	//		})
	//		return
	//	}

	var (
		smalltaskList []*smalltaskmodel.SmallTaskModel
		//		err           = errors.New("")
	)
	if area == "" {
		smalltaskList, err = smalltaskmodel.QueryNorNotSmallTask()
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}

	} else {
		smalltaskList, err = smalltaskmodel.QueryAreaNotSmallTask(area)
		if err != nil {
			log.Error(fmt.Sprintf("query small task err %s", err))
			c.JSON(400, gin.H{
				"code":    vars.ErrSmallTaskNotFound.Code,
				"message": vars.ErrSmallTaskNotFound.Msg,
			})
			return
		}
	}

	if smalltaskList == nil {
		log.Error(fmt.Sprintf("not small task to allot err %s", err))
		c.JSON(200, gin.H{
			"code":    vars.ErrNotSmallTask.Code,
			"message": vars.ErrNotSmallTask.Msg,
		})
		return
	}

	var (
		smallTaskId string
		count       int = 0
	)
	for {
		smallTaskId, err = getTasksId(username, smalltaskList[r.Intn(len(smalltaskList))])
		if err != nil {
			log.Error(fmt.Sprintf("not small task to allot err %s", err))
		}
		count = count + 1
		if !strings.EqualFold(smallTaskId, "") {
			break
		} else {
			if count >= len(smalltaskList) {
				break
			}
			//			log.Error(fmt.Sprintf("not small task to allot err %s", err))
			//			c.JSON(200, gin.H{
			//				"code":    vars.ErrNotSmallTask.Code,
			//				"message": vars.ErrNotSmallTask.Msg,
			//			})
			//			return
		}
	}

	if strings.EqualFold(smallTaskId, "") {
		log.Error(fmt.Sprintf("not small task to allot err %s", err))
		c.JSON(200, gin.H{
			"code":    vars.ErrNotSmallTask.Code,
			"message": vars.ErrNotSmallTask.Msg,
		})
		return
	}

	//	smallTaskId, err := getSmallTasksId(username, smalltaskList)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("not small task to allot err %s", err))
	//		c.JSON(200, gin.H{
	//			"code":    vars.ErrNotSmallTask.Code,
	//			"message": vars.ErrNotSmallTask.Msg,
	//		})
	//		return
	//	}

	//	if strings.EqualFold(smallTaskId, "") {
	//		log.Error(fmt.Sprintf("not small task to allot err %s", err))
	//		c.JSON(200, gin.H{
	//			"code":    vars.ErrNotSmallTask.Code,
	//			"message": vars.ErrNotSmallTask.Msg,
	//		})
	//		return
	//	}

	c.JSON(200, gin.H{
		"code":          0,
		"small_task_id": smallTaskId,
	})
}

func getFineTune(username string, smalltaskList []*smalltaskmodel.SmallTaskModel) (string, error) {
	//
	for _, smallTask := range smalltaskList {
		imageList, err := imagemodel.GetSmallTaskImages(smallTask.SmallTaskImages)
		if err != nil {
			log.Error(fmt.Sprintf("image query err", err.Error()))
			continue
		}

		notImages := getFineNotImageList(imageList, username, smallTask.LimitCount, smallTask.SmallTaskId, smallTask.PointType)
		if len(notImages) == 0 {
			log.Error(fmt.Sprintf("Under the task gets no pictures err"))
			continue
		}

		for _, image := range notImages {
			timeOutModels, err := timeoutmodel.QuerySmallTaskImage(image.Md5, smallTask.SmallTaskId)
			if err != nil {
				if err != timeoutmodel.ErrTimeOutModelNotFound {
					log.Error(fmt.Sprintf("image query err", err.Error()))
					continue
				}
			}
			if timeOutModels == nil {
				return smallTask.SmallTaskId, nil
			}
			var flag bool = false
			for _, res := range timeOutModels {
				if strings.EqualFold(res.User, username) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			//继续判断是否超出数量限制
			result := image.FineResults[strconv.Itoa(int(smallTask.PointType))]
			var count int64 = 0
			for _, res := range result {
				if strings.EqualFold(res.SmallTaskId, smallTask.SmallTaskId) {
					count += 1
				}
			}

			if int64(len(timeOutModels))+count == smallTask.LimitCount {
				continue
			}

			return smallTask.SmallTaskId, nil
		}

	}

	return "", nil
}

func getTasksId(username string, smallTask *smalltaskmodel.SmallTaskModel) (string, error) {

	imageList, err := imagemodel.GetSmallTaskImages(smallTask.SmallTaskImages)
	if err != nil {
		log.Error(fmt.Sprintf("image query err", err.Error()))
		return "", nil
	}

	notImages := getNotImageList(imageList, username, smallTask.Areas, smallTask.LimitCount, smallTask.SmallTaskId, smallTask.PointType)
	if len(notImages) == 0 {
		log.Error(fmt.Sprintf("Under the task gets no pictures err"))
		return "", nil
	}

	for _, image := range notImages {
		timeOutModels, err := timeoutmodel.QuerySmallTaskImage(image.Md5, smallTask.SmallTaskId)
		if err != nil {
			if err != timeoutmodel.ErrTimeOutModelNotFound {
				log.Error(fmt.Sprintf("image query err", err.Error()))
				continue
			}
		}
		if timeOutModels == nil {
			return smallTask.SmallTaskId, nil
		}
		var flag bool = false
		for _, res := range timeOutModels {
			if strings.EqualFold(res.User, username) {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		//继续判断是否超出数量限制
		result := image.Results[strconv.Itoa(int(smallTask.PointType))][smallTask.Areas]
		var count int64 = 0
		for _, res := range result {
			if strings.EqualFold(res.SmallTaskId, smallTask.SmallTaskId) {
				count += 1
			}
		}

		if int64(len(timeOutModels))+count >= smallTask.LimitCount {
			continue
		}

		return smallTask.SmallTaskId, nil
	}

	return "", nil
}

//筛选可以获取到image的taskid
func getSmallTasksId(username string, smalltaskList []*smalltaskmodel.SmallTaskModel) (string, error) {
	timeOutImages, err := timeoutmodel.QueryUserImages(username)
	if err != nil {
		if err != timeoutmodel.ErrTimeOutModelNotFound {
			log.Error(fmt.Sprintf("time out image query err", err.Error()))
			return "", err
		}
	}

	for _, smallTask := range smalltaskList {
		for _, timeOutimage := range timeOutImages {
			if strings.EqualFold(timeOutimage.SmallTaskId, smallTask.SmallTaskId) {
				continue
			}
		}

		imageList, err := imagemodel.GetSmallTaskImages(smallTask.SmallTaskImages)
		if err != nil {
			log.Error(fmt.Sprintf("image query err", err.Error()))
			continue
		}

		notImages := getNotImageList(imageList, username, smallTask.Areas, smallTask.LimitCount, smallTask.SmallTaskId, smallTask.PointType)
		if len(notImages) == 0 {
			log.Error(fmt.Sprintf("Under the task gets no pictures err"))
			continue
		}

		for _, image := range notImages {
			timeOutModels, err := timeoutmodel.QuerySmallTaskImage(image.Md5, smallTask.SmallTaskId)
			if err != nil {
				if err != timeoutmodel.ErrTimeOutModelNotFound {
					log.Error(fmt.Sprintf("image query err", err.Error()))
					continue
				}
			}
			if timeOutModels == nil {
				return smallTask.SmallTaskId, nil
			}
			var flag bool = false
			for _, res := range timeOutModels {
				if strings.EqualFold(res.User, username) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			//继续判断是否超出数量限制
			result := image.Results[strconv.Itoa(int(smallTask.PointType))][smallTask.Areas]
			var count int64 = 0
			for _, res := range result {
				if strings.EqualFold(res.SmallTaskId, smallTask.SmallTaskId) {
					count += 1
				}
			}

			if int64(len(timeOutModels))+count == smallTask.LimitCount {
				continue
			}

			return smallTask.SmallTaskId, nil
		}
	}
	return "", nil
}

//查出结果中没有自己的且未标完的
func getFineNotImageList(imageList []*imagemodel.ImageModel, username string, limitCount int64, stmId string, pointType int64) []*imagemodel.ImageModel {

	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {

		result := image.FineResults[strconv.Itoa(int(pointType))]
		//		log.Info(fmt.Sprintf("image query result %s err", result))
		if len(result) == 0 || result == nil {
			list = append(list, image)
			continue
		}

		var (
			count int64 = 0
			flag  bool  = true
		)
		for _, res := range result {
			if res.User == username {
				if strings.EqualFold(res.SmallTaskId, stmId) {
					log.Info(fmt.Sprintf("image query result user %s %s err", res.User, username))
					flag = false
					break
				}
				continue
			}

			if strings.EqualFold(res.SmallTaskId, stmId) {
				count += 1
			}

		}

		if count < limitCount && flag {
			list = append(list, image)
		}
	}
	return list
}

//查出结果中没有自己的且未标完的
func getNotImageList(imageList []*imagemodel.ImageModel, username string, area string, limitCount int64, stmId string, pointType int64) []*imagemodel.ImageModel {
	list := make([]*imagemodel.ImageModel, 0, 0)
	for _, image := range imageList {
		result := image.Results[strconv.Itoa(int(pointType))][area]
		//		log.Info(fmt.Sprintf("image query result %s err", result))
		if len(result) == 0 || result == nil {
			list = append(list, image)
			continue
		}

		var (
			count int64 = 0
			flag  bool  = true
		)
		for _, res := range result {
			if res.User == username {
				if strings.EqualFold(res.SmallTaskId, stmId) {
					log.Info(fmt.Sprintf("image query result user %s %s err", res.User, username))
					flag = false
					break
				}
				continue
			}

			if strings.EqualFold(res.SmallTaskId, stmId) {
				count += 1
			}

		}

		if count < limitCount && flag {
			list = append(list, image)
		}
	}
	return list
}

//首先给出自己未完成的
func getTimeOutImage(username string) (string, error) {
	timeOutImages, err := timeoutmodel.QueryUserImages(username)
	if err != nil {
		if err != timeoutmodel.ErrTimeOutModelNotFound {
			log.Error(fmt.Sprintf("time out image query err", err.Error()))
			return "", err
		}
	}

	//	var imageModel *imagemodel.ImageModel
	if timeOutImages != nil {
		for _, timeOut := range timeOutImages {
			imageModel, err := imagemodel.QueryImage(timeOut.Md5)
			if err != nil {
				log.Error(fmt.Sprintf("image query err", err.Error()))
				continue
			}
			smallTaskModel, err := smalltaskmodel.QuerySmallTask(timeOut.SmallTaskId)
			if err != nil {
				log.Error(fmt.Sprintf("query small task err %s", err))
			}

			result := imageModel.Results[strconv.Itoa(int(smallTaskModel.PointType))][smallTaskModel.Areas]
			var flag bool = false
			for _, res := range result {
				if strings.EqualFold(res.User, username) && strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			return timeOut.SmallTaskId, nil
		}
	}

	return "", nil
}

func getTimeOutFineImage(username string) (string, error) {
	timeOutImages, err := timeoutmodel.QueryUserImages(username)
	if err != nil {
		if err != timeoutmodel.ErrTimeOutModelNotFound {
			log.Error(fmt.Sprintf("time out image query err", err.Error()))
			return "", err
		}
	}

	//	var imageModel *imagemodel.ImageModel
	if timeOutImages != nil {
		for _, timeOut := range timeOutImages {
			imageModel, err := imagemodel.QueryImage(timeOut.Md5)
			if err != nil {
				log.Error(fmt.Sprintf("image query err", err.Error()))
				continue
			}
			smallTaskModel, err := smalltaskmodel.QuerySmallTask(timeOut.SmallTaskId)
			if err != nil {
				log.Error(fmt.Sprintf("query small task err %s", err))
			}

			result := imageModel.FineResults[strconv.Itoa(int(smallTaskModel.PointType))]
			var flag bool = false
			for _, res := range result {
				if strings.EqualFold(res.User, username) && strings.EqualFold(res.SmallTaskId, smallTaskModel.SmallTaskId) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			return timeOut.SmallTaskId, nil
		}
	}

	return "", nil
}
