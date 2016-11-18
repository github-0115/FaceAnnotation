package facemodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TaskModel struct {
	Title     string    `bson:"title" json:"title" binding:"required"`
	Count     int64     `bson:"count" json:"count" binding:"required"`
	Url       []string  `bson:"url" json:"url" binding:"required"`
	Status    int64     `bson:"status" json:"status" binding:"required"`
	CreatedAt time.Time `bson:"created_at" binding:"required" json:"created_at"`
}

var (
	ErrTaskModelNotFound = errors.New("Face Model not found")
	ErrTaskModelCursor   = errors.New("Cursor err")
)

func (t *TaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("task").Insert(&t)
}

func QueryTaskList(status string) ([]*TaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()
	var result []*TaskModel
	err := s.DB(db.Face.DB).C("task").Find(bson.M{
		"status": status,
	}).All(result)

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
