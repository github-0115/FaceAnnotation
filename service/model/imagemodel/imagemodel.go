package imagemodel

import (
	db "FaceAnnotation/service/db"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	//	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ImageModel struct {
	TaskId    []string                        `bson:"task_id" json:"task_id" binding:"required"` //与 taskmodel 关联
	Md5       string                          `bson:"md5" json:"md5" binding:"required"`
	Results   map[string]map[string][]*Points `bson:"results" json:"results" binding:"required"`
	ThrFaces  map[string][]interface{}        `binding:"required" bson:"thr_faces" json:"thr_faces"`
	Url       string                          `binding:"required" bson:"url" json:"url"`
	CreatedAt time.Time                       `bson:"created_at" json:"created_at" binding:"required"`
}

type ThrFaces struct {
	FaceFivePoint       thrfacemodel.ThrFaceModel        `bson:"face_five_point" json:"face_five_point" binding:"required"`
	FaceEightThreePoint thrfacemodel.EightThreeFaceModel `bson:"face_eight_thr_point" json:"face_eight_thr_point" binding:"required"`
	ImportPoint         *ImportPoint                     `bson:"import_point" json:"import_point" binding:"required"`
}

type Result struct {
	Faces     map[string][]*Points `binding:"required" bson:"faces" json:"faces"`
	Status    int64                `bson:"status" json:"status" binding:"required"` //0:一般标图结果 1:微调结果
	CreatedAt time.Time            `bson:"created_at" json:"created_at" binding:"required"`
}

type ImportPoint struct {
	Name   string    `binding:"required" bson:"name" json:"name"`
	Points []float64 `binding:"required" bson:"points" json:"points"`
}

type Points struct {
	SmallTaskId string   `bson:"small_task_id" json:"small_task_id" binding:"required"`
	User        string   `bson:"user" json:"user" binding:"required"`
	Points      []*Point `binding:"required" bson:"points" json:"points"`
}

type Point struct {
	Y float64 `binding:"required" bson:"y" json:"y"`
	X float64 `binding:"required" bson:"x" json:"x"`
}

var (
	ErrImageModelNotFound = errors.New("image Model not found")
	ErrImageModelCursor   = errors.New("Cursor err")
)

func (im *ImageModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("image").Insert(&im)
}

func QueryImage(md5 string) (*ImageModel, error) {
	coll := new(ImageModel)
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image").Find(bson.M{
		"md5": md5,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrImageModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrImageModelCursor
		}
		return nil, err
	}
	return coll, nil
}

func GetSmallTaskImages(md5s []string) ([]*ImageModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*ImageModel
	err := s.DB(db.Face.DB).C("image").Find(bson.M{
		"md5": bson.M{"$in": md5s},
	}).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small task image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrImageModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrImageModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QueryTaskImages(taskId string) ([]*ImageModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*ImageModel
	err := s.DB(db.Face.DB).C("image").Find(bson.M{
		"task_id": taskId,
	}).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrImageModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrImageModelCursor
		}
		return nil, err
	}
	return results, nil
}

func UpdateImageModel(md5 string, taskId string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image").Update(bson.M{
		"md5": md5,
	}, bson.M{"$push": bson.M{"task_id": taskId}})
	if err != nil {
		log.Error(fmt.Sprintf("upsert image point err ", err))
		if err == mgo.ErrNotFound {
			return ErrImageModelNotFound
		}
		return ErrImageModelCursor
	}
	return nil
}

func UpsertImageModel(res *ImageModel) (bool, error) {
	s := db.Face.GetSession()
	defer s.Close()

	_, err := s.DB(db.Face.DB).C("image").Upsert(bson.M{
		"md5": res.Md5,
	}, res)
	if err != nil {
		log.Error(fmt.Sprintf("upsert image point err ", err))
		if err == mgo.ErrNotFound {
			return false, ErrImageModelNotFound
		}
		return false, ErrImageModelCursor
	}
	return true, nil
}

func GetImageList(dirPth string) (files []string, err error) {
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}

		urls := strings.Split(filename, "\\")
		files = append(files, urls[len(urls)-1])

		return nil
	})
	return files, err
}
