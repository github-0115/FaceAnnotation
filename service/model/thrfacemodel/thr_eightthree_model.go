package thrfacemodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	//	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

var (
	landUrl = "http://apicn.faceplusplus.com/v2/detection/landmark"
)

type EightThreeFaceModel struct {
	SessionId string    `json:"session_id"`
	Result    []*Result `json:"result"`
}

type Result struct {
	FaceId      string            `json:"face_id"`
	ImageWidth  float64           `json:"image_width"`
	ImageHeight float64           `json:"image_height"`
	FaceWidth   float64           `json:"face_width"`
	FaceHeight  float64           `json:"face_height"`
	Landmark    map[string]*Point `json:"Landmark"`
}

func EightThreeFace(face_id string) (*EightThreeFaceModel, error) {
	extraParams := map[string]string{
		"api_key":    api_key,
		"api_secret": api_secret,
		"type":       "83p",
		"face_id":    face_id, //face_id
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
		//		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("83 read form err, err=%#v", err))
		//		return nil, err
	}

	rep := new(EightThreeFaceModel)
	if err := json.Unmarshal(body.Bytes(), &rep); err != nil {
		log.Error(fmt.Sprintf("json unmarshal err=%s", err))
		//		return nil, err
	}

	defer resp.Body.Close()
	fmt.Println(body.String())
	return rep, nil
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
