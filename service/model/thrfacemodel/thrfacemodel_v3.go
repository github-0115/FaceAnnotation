package thrfacemodel

import (
	cfg "FaceAnnotation/config"
	"bytes"
	"encoding/json"

	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	log "github.com/inconshreveable/log15"
)

type FaceModelV3 struct {
	ImageId   string  `json:"image_id"`
	RequestId string  `json:"requst_id"`
	TimeUsed  int64   `json:"time_used"`
	Faces     []*Face `json:"faces"`
}

type Face struct {
	FaceToken     string            `json:"face_token"`
	FaceRectangle *FaceRectangle    `json:"face_rectangle"`
	Landmark      map[string]*Point `json:"Landmark"`
}

type FaceRectangle struct {
	ImageWidth  float64 `json:"width"`
	ImageHeight float64 `json:"height"`
	Top         float64 `json:"top"`
	Left        float64 `json:"left"`
}

var (
	api_url   = cfg.Cfg.FaceDetectUrl
	apiKey    = cfg.Cfg.FaceApiKey
	apiSecret = cfg.Cfg.FaceApiSecret
)

func ThrFaceFileResV3(fileName string, fileBytes []byte) (*FaceModelV3, error) {

	extraParams := map[string]string{
		"api_key":         apiKey,
		"api_secret":      apiSecret,
		"return_landmark": "1",
	}

	request, err := newfileUploadRequest(api_url, extraParams, "image_file", fileName, fileBytes)
	if err != nil {
		log.Error(fmt.Sprintf("new req.. err, err=%#v", err))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(fmt.Sprintf("client.Do err, err=%#v", err))
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("read form err, err=%#v", err))
		return nil, err
	}

	rep := new(FaceModelV3)
	if err := json.Unmarshal(body.Bytes(), &rep); err != nil {
		log.Error(fmt.Sprintf("json unmarshal err=%s", err))
		return nil, err
	}

	defer resp.Body.Close()
	fmt.Println(body)

	return rep, nil
}

func newfileUploadRequestV3(uri string, params map[string]string, paramName, fileName string, fileBytes []byte) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(fileBytes))

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}
