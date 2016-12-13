package taskmodel

import (
	db "FaceAnnotation/service/db"
	facemodel "FaceAnnotation/service/model/facemodel"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TaskModel struct {
	Title     string        `bson:"title" json:"title" binding:"required"`
	Count     int64         `bson:"count" json:"count" binding:"required"`
	Images    []*ImageModel `bson:"images" json:"images" binding:"required"`
	Status    int64         `bson:"status" json:"status" binding:"required"` // 0 未开始  1 正在进行  2 已完成
	CreatedAt time.Time     `bson:"created_at" json:"created_at" binding:"required"`
}

type ImageModel struct {
	Url    string `bson:"url" json:"url" binding:"required"` // 0 未标  1 已标注完成
	Status int64  `bson:"status" json:"status" binding:"required"`
}

type ImportModel struct {
	Url         string                 `bson:"url" json:"url"`
	OriginFaces []*facemodel.Landmarks `binding:"required" bson:"origin_faces" json:"origin_faces"`
}

var (
	ErrTaskModelNotFound = errors.New("Face Model not found")
	ErrTaskModelCursor   = errors.New("Cursor err")
	ErrDirNotFound       = errors.New("dir ads not found")
	ErrFileNotFound      = errors.New("file ads not found")
	ErrReadFile          = errors.New("read file err")
)

func (t *TaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("task").Insert(&t)
}

func QueryTaskList(status int64) ([]*TaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()
	var result []*TaskModel
	err := s.DB(db.Face.DB).C("task").Find(bson.M{
		"status": status,
	}).All(&result)

	if err != nil {
		log.Error(fmt.Sprintf("find task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrTaskModelCursor
		}
		return nil, err
	}
	return result, nil
}

func QueryTask(title string) (*TaskModel, error) {
	coll := new(TaskModel)
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("task").Find(bson.M{
		"title": title,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrTaskModelCursor
		}
		return nil, err
	}
	return coll, nil
}

func UpdateTaskStatus(title string, status int64) error {

	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("task").Update(bson.M{
		"title": title,
	}, bson.M{"$set": bson.M{"status": status}})

	if err != nil {
		log.Error(fmt.Sprintf("update task status err ", err))
		if err == mgo.ErrNotFound {
			return ErrTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return ErrTaskModelCursor
		}
		return err
	}
	return nil
}

func UpdateTaskImageStatus(title string, url string, status int64) error {

	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("task").Update(bson.M{
		"title":      title,
		"images.url": url,
	}, bson.M{"$set": bson.M{"images.$.status": status}})

	if err != nil {
		log.Error(fmt.Sprintf("update task image status err ", err))
		if err == mgo.ErrNotFound {
			return ErrTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return ErrTaskModelCursor
		}
		return err
	}
	return nil
}

func ByteToTaskModel(filebyte []byte) (*TaskModel, []*ImportModel, error) {

	urlStr := string(filebyte)
	url_lines := strings.Split(urlStr, "\n")

	url_list := make([]string, 0, 0)
	for _, url_line := range url_lines {
		if url_line != "" {

			url_list = append(url_list, url_line)
		}
	}

	importTask, importModelList, err := importFileToTaskModel(url_list)
	if err != nil {
		log.Error(fmt.Sprintf("json marshal err==%s", err))
		return nil, nil, err
	}
	return importTask, importModelList, nil
}

func importFileToTaskModel(s []string) (*TaskModel, []*ImportModel, error) {
	list := make([]*ImportModel, 0, len(s))
	taskModel := &TaskModel{
		Count:     int64(len(s)),
		Status:    0,
		CreatedAt: time.Now(),
	}
	images := make([]*ImageModel, 0, len(s))
	for _, res := range s {
		rep := new(ImportModel)
		if err := json.Unmarshal([]byte(res), &rep); err != nil {
			log.Error(fmt.Sprintf("json unmarshal err=%s", err))

			return nil, nil, err
		}
		list = append(list, rep)
		image := &ImageModel{
			Url:    rep.Url,
			Status: 0,
		}
		images = append(images, image)
	}
	taskModel.Images = images

	return taskModel, list, nil
}
