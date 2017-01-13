package uploadmodel

import (
	cfg "FaceAnnotation/config"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	//	taskmodel "FaceAnnotation/service/model/taskmodel"
	//	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	//	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/inconshreveable/log15"
)

var (
	bucket               *oss.Bucket
	ErrClient            = errors.New("oss New client err")
	ErrIsBucketExist     = errors.New("oss IsBucketExist err")
	ErrCreateBucket      = errors.New("oss CreateBucket err")
	ErrBucket            = errors.New("oss get Bucket err")
	ErrPutObjectFromFile = errors.New("oss PutObjectFromFile err")
	ErrReadAllFile       = errors.New("oss ErrReadAllFile err")
)

//func UploadFile(photoName string, []byte) (string, error) {
////	fileByte, err := ioutil.ReadAll(file)
////	if err != nil {
////		log.Error(fmt.Sprintf("ioutil  ReadAll file err" + err.Error()))
////		return "", ErrReadAllFile
////	}

//	//	h := md5.New()
//	//	h.Write(fileByte)
//	//	photoName := hex.EncodeToString(h.Sum(nil))

//	_, err = Ossload(photoName, fileByte)
//	if err != nil {
//		log.Error(fmt.Sprintf("user oss photo failed. err=%#v", err))
//		if err == ErrClient {
//			return photoName, ErrClient
//		} else if err == ErrIsBucketExist {
//			return photoName, ErrIsBucketExist
//		} else if err == ErrCreateBucket {
//			return "", ErrCreateBucket
//		} else if err == ErrBucket {
//			return "", ErrBucket
//		}
//		return "", ErrPutObjectFromFile
//	}

//	return photoName, nil
//}

func UploadUrlsFile(urls []string, taskid string, points map[string][]interface{}) (string, error) {
	bucket, err := ossBucket()
	if err != nil {
		return "", err
	}

	send := make(chan string, len(urls))
	getchan := make(chan string, len(urls))

	for i := 0; i < 100; i++ {
		go ossWorker(bucket, taskid, points, send, getchan)
	}

	for _, res := range urls {
		send <- res
	}
	close(send)

	return "", nil
}

func UploadLocalFiles(urls []string, taskid string, points map[string][]interface{}) (string, error) {
	log.Info(fmt.Sprintf("upload start%s", time.Now()))
	bucket, err := ossBucket()
	if err != nil {
		return "", err
	}

	send := make(chan string, len(urls))
	getchan := make(chan string, len(urls))

	for i := 0; i < 100; i++ {
		go ossWorker(bucket, taskid, points, send, getchan)
	}

	for _, res := range urls {
		send <- res
	}
	close(send)

	results := make([]string, 0, 0)
	for {
		resStr := <-getchan
		results = append(results, resStr)

		if len(results) == len(urls) {
			break
		}
	}
	log.Info(fmt.Sprintf("upload end%s", time.Now()))
	log.Info(fmt.Sprintf("upload res%s", results))
	return "import success", nil
}

func ossWorker(bucket *oss.Bucket, taskid string, points map[string][]interface{}, send, getchan chan string) {

	for res := range send {

		fileByte, err := ioutil.ReadFile(res)
		if err != nil {
			getchan <- "n"
			continue
		}
		h := md5.New()
		h.Write(fileByte)
		photoName := hex.EncodeToString(h.Sum(nil))

		names := strings.Split(res, "\\")

		imageColl, err := imagemodel.QueryImage(photoName)
		if imageColl == nil {
			imageColl = &imagemodel.ImageModel{
				TaskId:    []string{taskid}, //taskid
				Md5:       photoName,
				Url:       names[len(names)-1],
				CreatedAt: time.Now(),
			}

			err = bucket.PutObject(photoName, bytes.NewReader(fileByte))
			if err != nil {
				log.Error(fmt.Sprintf("client Bucket PutObject err" + err.Error()))
			}
		} else {
			imageColl.TaskId = append(imageColl.TaskId, taskid)
		}

		if imageColl.ThrFaces == nil {
			imageColl.ThrFaces = make(map[string]map[string]interface{})
		}
		imageColl.ThrFaces["deepir_import"]["95"] = points[names[len(names)-1]]

		_, err = imagemodel.UpsertImageModel(imageColl)
		if err != nil {
			log.Error(fmt.Sprintf("image=%s update err=%s", err.Error()))
		}

		getchan <- "y"
	}
}

func ossBucket() (*oss.Bucket, error) {
	client, err := oss.New(cfg.Cfg.ALAkDomian, cfg.Cfg.ALAkId, cfg.Cfg.ALAkSecret)
	if err != nil {
		log.Error(fmt.Sprintf("oss.New client err" + err.Error()))
		return nil, ErrClient
	}

	isExist, err := client.IsBucketExist("faceannotation")
	if err != nil {
		log.Error(fmt.Sprintf("oss client.IsBucketExist  err" + err.Error()))
		return nil, ErrIsBucketExist
	}

	if isExist {
		bucket, err = client.Bucket("faceannotation")
		if err != nil {
			log.Error(fmt.Sprintf("oss client Bucket err" + err.Error()))
			return nil, ErrBucket
		}
	} else {
		err = client.CreateBucket("faceannotation")
		if err != nil {
			log.Error(fmt.Sprintf("oss client CreateBucket err" + err.Error()))
			return nil, ErrCreateBucket
		}
	}

	// 设置Bucket ACL
	err = client.SetBucketACL("faceannotation", oss.ACLPublicRead)
	if err != nil {
		log.Error(fmt.Sprintf("oss set Bucket ACL err" + err.Error()))
	}

	return bucket, nil
}

func UploadFile(photoName string, fileByte []byte) (string, error) {

	_, err := Ossload(photoName, fileByte)
	if err != nil {
		log.Error(fmt.Sprintf("user oss photo failed. err=%#v", err))
		if err == ErrClient {
			return photoName, ErrClient
		} else if err == ErrIsBucketExist {
			return photoName, ErrIsBucketExist
		} else if err == ErrCreateBucket {
			return "", ErrCreateBucket
		} else if err == ErrBucket {
			return "", ErrBucket
		}
		return "", ErrPutObjectFromFile
	}

	return photoName, nil
}

func Ossload(fileName string, fileByte []byte) (string, error) {
	client, err := oss.New(cfg.Cfg.ALAkDomian, cfg.Cfg.ALAkId, cfg.Cfg.ALAkSecret)
	if err != nil {
		log.Error(fmt.Sprintf("oss.New client err" + err.Error()))
		return "", ErrClient
	}

	isExist, err := client.IsBucketExist("faceannotation")
	if err != nil {
		log.Error(fmt.Sprintf("oss client.IsBucketExist  err" + err.Error()))
		return "", ErrIsBucketExist
	}

	if isExist {
		bucket, err = client.Bucket("faceannotation")
		if err != nil {
			log.Error(fmt.Sprintf("oss client Bucket err" + err.Error()))
			return "", ErrBucket
		}
	} else {
		err = client.CreateBucket("faceannotation")
		if err != nil {
			log.Error(fmt.Sprintf("oss client CreateBucket err" + err.Error()))
			return "", ErrCreateBucket
		}
	}

	// 设置Bucket ACL
	err = client.SetBucketACL("faceannotation", oss.ACLPublicRead)
	if err != nil {
		log.Error(fmt.Sprintf("oss set Bucket ACL err" + err.Error()))
	}

	err = bucket.PutObject(fileName, bytes.NewReader(fileByte))
	if err != nil {
		log.Error(fmt.Sprintf("client Bucket PutObject err" + err.Error()))
		return "", ErrPutObjectFromFile
	}

	fileUrl := "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/" + fileName

	return fileUrl, nil
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
