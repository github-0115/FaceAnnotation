package imagetaskmodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ImageTaskModel struct {
	ImageTaskId string    `bson:"image_task_id" json:"image_task_id" binding:"required"`
	TaskId      []string  `bson:"task_id" json:"task_id" binding:"required"`
	Total       int64     `bson:"total" json:"total" binding:"required"`
	Images      []string  `bson:"images" json:"images" binding:"required"`
	Introduce   string    `bson:"introduce" json:"introduce" binding:"required"`
	Status      int64     `bson:"status" json:"status" binding:"required"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at" binding:"required"`
}

var (
	ErrImageTaskModelNotFound = errors.New("image task Model not found")
	ErrImageTaskModelCursor   = errors.New("Cursor err")
)

func (itm *ImageTaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("image_task").Insert(&itm)
}

func QueryImageTask(id string) (*ImageTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	coll := new(ImageTaskModel)
	err := s.DB(db.Face.DB).C("image_task").Find(bson.M{
		"image_task_id": id,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find image task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrImageTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrImageTaskModelCursor
		}
		return nil, err
	}
	return coll, nil
}

func GetImageTasks(md5s []string) ([]*ImageTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*ImageTaskModel
	err := s.DB(db.Face.DB).C("image_task").Find(nil).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find image task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrImageTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrImageTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}
