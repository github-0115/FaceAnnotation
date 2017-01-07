package thrfacemodel

import (
	cfg "FaceAnnotation/config"
	"bytes"
	"encoding/json"
	//	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"

	//	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

var (
	apiUrl     = cfg.Cfg.FaceDetectUrl
	api_key    = cfg.Cfg.FaceApiKey
	api_secret = cfg.Cfg.FaceApiSecret

//	imageUrl   = "http://www.qq1234.org/uploads/allimg/120509/1_120509171458_7.jpg"
)

type ThrFaceModel struct {
	ImgId     string       `json:"img_id"`
	SessionId string       `json:"session_id"`
	ImgHeight int64        `json:"img_height"`
	ImgWidth  int64        `json:"img_width"`
	Face      []*FaceModel `json:"face"`
}

type FaceModel struct {
	FaceId    string     `json:"face_id"`
	Attribute *Attribute `json:"attribute"`
	Position  *Position  `json:"position"`
}

type Attribute struct {
	Age     *Age     `json:"age"`
	Gender  *Gender  `json:"gender"`
	Race    *Race    `json:"race"`
	Smiling *Smiling `json:"smiling"`
}

type Age struct {
	Range float32 `json:"range"`
	Value float32 `json:"value"`
}

type Gender struct {
	Confidence float32 `json:"confidence"`
	Value      string  `json:"value"`
}

type Race struct {
	Confidence float32 `json:"confidence"`
	Value      string  `json:"value"`
}

type Smiling struct {
	Value float32 `json:"value"`
}

type Position struct {
	Center     *Point  `json:"center"`
	EyeLeft    *Point  `json:"eye_left"`
	EyeRight   *Point  `json:"eye_right"`
	MouthLeft  *Point  `json:"mouth_left"`
	MouthRight *Point  `json:"mouth_right"`
	Nose       *Point  `json:"nose"`
	Height     float32 `json:"height"`
	Width      float32 `json:"width"`
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func ThrFaceFileRes(fileName string, fileBytes []byte) (*ThrFaceModel, error) {

	extraParams := map[string]string{
		"api_key":    api_key,
		"api_secret": api_secret,
	}

	request, err := newfileUploadRequest("http://apicn.faceplusplus.com/v2/detection/detect", extraParams, "img", fileName, fileBytes)
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

	rep := new(ThrFaceModel)
	if err := json.Unmarshal(body.Bytes(), &rep); err != nil {
		log.Error(fmt.Sprintf("json unmarshal err=%s", err))
		return nil, err
	}

	defer resp.Body.Close()
	fmt.Println(body)

	return rep, nil
}

func newfileUploadRequest(uri string, params map[string]string, paramName, fileName string, fileBytes []byte) (*http.Request, error) {

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

func ThrFaceRes(imageUrl string) (*ThrFaceModel, error) {

	extraParams := map[string]string{
		"api_key":    api_key,
		"api_secret": api_secret,
	}

	request, err := newurlfileUploadRequest("http://apicn.faceplusplus.com/v2/detection/detect", extraParams, "img", imageUrl)
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

	rep := new(ThrFaceModel)
	if err := json.Unmarshal(body.Bytes(), &rep); err != nil {
		log.Error(fmt.Sprintf("json unmarshal err=%s", err))
		return nil, err
	}

	defer resp.Body.Close()
	fmt.Println(rep)
	fmt.Println("====+++")
	fmt.Println(body)

	return rep, nil
}

// Creates a new file upload http request with optional extra params
func newurlfileUploadRequest(uri string, params map[string]string, paramName, imageUrl string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(imageUrl))
	if err != nil {
		return nil, err
	}

	fileBytes, err := getLocalFile(imageUrl)
	if err != nil {
		log.Error(fmt.Sprintf("get image err = %s", err))
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

func getLocalFile(path string) ([]byte, error) {
	fileByte, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(fmt.Sprintf("ioutil read image file err" + err.Error()))
		return nil, err
	}
	return fileByte, nil
}

func getImg(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Error(fmt.Sprintf("get url=%s pic err=%s", url, err))
		return nil, err
	}
	defer resp.Body.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	return pix, err
}
