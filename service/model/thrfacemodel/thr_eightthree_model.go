package thrfacemodel

import (
	cfg "FaceAnnotation/config"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

var (
	landUrl = cfg.Cfg.FaceLandmarkUrl
)

type EightThreeFaceModel struct {
	SessionId string    `json:"session_id"`
	Result    []*Result `json:"result"`
}

type Result struct {
	FaceId   string            `json:"face_id"`
	Landmark map[string]*Point `json:"Landmark"`
}

func EightThreeFace(c *gin.Context) {
	extraParams := map[string]string{
		"api_key":    api_key,
		"api_secret": api_secret,
		"type":       "83p",
		"face_id":    "78197afb212ee433d08a9a1961f64fda", //face_id
	}
	fmt.Println(extraParams)
	request, err := newMultipartRequest(landUrl, extraParams)
	if err != nil {
		log.Error(fmt.Sprintf("new 83 req.. err, err=%#v", err))
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(fmt.Sprintf("83 client.Do err, err=%#v", err))
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("83 read form err, err=%#v", err))
	}

	rep := new(EightThreeFaceModel)
	if err := json.Unmarshal(body.Bytes(), &rep); err != nil {
		log.Error(fmt.Sprintf("json unmarshal err=%s", err))
	}

	defer resp.Body.Close()

	c.JSON(200, gin.H{
		"code": 0,
		"rep":  rep,
	})
}

// Creates a new file upload http request with optional extra params
func newMultipartRequest(url string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}
