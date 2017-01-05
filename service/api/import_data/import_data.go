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

	if len(imageList) != len(importPoints) {
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
		Introduce: "title",
		Status:    0,
		CreatedAt: time.Now(),
	}

	//	url, err := uploadmodel.UploadFile(file)
	//	if err != nil {
	//		log.Error(fmt.Sprintf("upload file err = %s", err))
	//	}

	ss, err := uploadmodel.UploadLocalFiles(imageList, importPoints)
	if err != nil {
		log.Error(fmt.Sprintf("upload file err = %s", err))
	}

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

func readPoint(dat []byte) map[string]interface{} {
	urlStr := strings.Replace(string(dat), " ", "", -1)
	urls := strings.Split(urlStr, "\n")

	m := make(map[string]interface{})
	for _, url := range urls {
		if url != "" {

			rr := strings.Split(url, "\t")
			pointStrs := strings.Split(rr[1], ",")

			points := make([]float64, 0, 0)
			for _, pointStr := range pointStrs {
				f, _ := strconv.ParseFloat(pointStr, 32)
				points = append(points, f)
			}

			m[rr[0]] = points
		}
	}
	return m
}
