package import_data

import (
	taskmodel "FaceAnnotation/service/model/taskmodel"
	uploadmodel "FaceAnnotation/service/model/uploadmodel"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

func ImportData(c *gin.Context) {
	//
	pointFile, _, err := c.Request.FormFile("point")
	imagePath := c.PostForm("image_path")
	//	areas := c.PostForm("areas")
	unitstr := c.PostForm("unit")
	limit_user := c.PostForm("limit_user")
	point_type := c.PostForm("point_type")
	unit, err := strconv.Atoi(unitstr)
	limitUser, err := strconv.Atoi(limit_user)
	pointType, err := strconv.Atoi(point_type)
	fileByte, err := ioutil.ReadAll(pointFile)
	if err != nil {
		log.Error(fmt.Sprintf("ioutil ReadAll file err" + err.Error()))
		c.JSON(400, gin.H{
			"code":    0,
			"message": "read file err",
		})
		return
	}

	importPoints := readPoint(fileByte)

	var imageList []string
	if imagePath != "" {
		imageList = getFilelist(imagePath)
	}
	log.Error(fmt.Sprintf("imageList=%s,importPoints=%s", imageList, importPoints))
	if len(imageList) != len(importPoints) {
		log.Error(fmt.Sprintf("imageList=%d,importPoints=%d", len(imageList), len(importPoints)))
		c.JSON(400, gin.H{
			"code":    0,
			"message": "image count != point count",
		})
		return
	}
	taskId := uuid.NewV4().String()
	taskModel := &taskmodel.TaskModel{
		TaskId:    taskId,
		Count:     int64(len(imageList)),
		PointType: int64(pointType),
		MinUnit:   int64(unit),
		LimitUser: int64(limitUser),
		Area:      []string{"left_eye_brow", "right_eye_brow", "left_eye", "right_eye", "left_ear", "right_ear", "mouth", "nouse", "face"},
		Introduce: "test",
		Status:    0,
		CreatedAt: time.Now(),
	}

	ss, err := uploadmodel.UploadLocalFiles(imageList, taskId, importPoints)
	if err != nil {
		log.Error(fmt.Sprintf("upload file err = %s", err))
	}

	taskModel.Save()

	c.JSON(200, gin.H{
		"code":      0,
		"taskModel": taskModel,
		"res":       ss,
	})
}

func getFilelist(path string) []string {
	list := make([]string, 0, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		list = append(list, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return list
}

func readPoint(dat []byte) map[string][]interface{} {
	urlStr := strings.Replace(string(dat), " ", "", -1)
	urls := strings.Split(urlStr, "\n")

	m := make(map[string][]interface{})
	for _, url := range urls {
		if url != "" {

			rr := strings.Split(url, "\t")
			pointStrs := strings.Split(rr[1], ",")

			points := []interface{}{}
			for _, pointStr := range pointStrs {
				f, _ := strconv.ParseFloat(pointStr, 32)
				points = append(points, f)
			}

			m[rr[0]] = points
		}
	}
	return m
}
