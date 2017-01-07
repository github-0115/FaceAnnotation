package timeoutmodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TimeOutModel struct {
	SmallTaskId string    `bson:"small_task_id" json:"small_task_id" binding:"required"`
	Md5         string    `bson:"md5" json:"md5" binding:"required"`
	User        string    `binding:"required" bson:"user" json:"user"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at" binding:"required"` //超时自动删除该条记录
}

var (
	ErrTimeOutModelNotFound = errors.New("task image Model not found")
	ErrTimeOutModelCursor   = errors.New("Cursor err")
)

func (tom *TimeOutModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("timeout").Insert(&tom)
}

func QueryUserImages(username string) ([]*TimeOutModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*TimeOutModel
	err := s.DB(db.Face.DB).C("timeout").Find(bson.M{
		"user": username,
	}).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find task image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrTimeOutModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrTimeOutModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QueryUserTsakImage(username string, smId string) ([]*TimeOutModel, error) {

	s := db.Face.GetSession()
	defer s.Close()

	var results []*TimeOutModel
	err := s.DB(db.Face.DB).C("timeout").Find(bson.M{
		"user":          username,
		"small_task_id": smId,
	}).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find task image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrTimeOutModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrTimeOutModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QuerySmallTaskImage(md5 string, smallTaskId string) ([]*TimeOutModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*TimeOutModel
	err := s.DB(db.Face.DB).C("timeout").Find(bson.M{
		"md5":           md5,
		"small_task_id": smallTaskId,
	}).All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find task image err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrTimeOutModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrTimeOutModelCursor
		}
		return nil, err
	}
	return results, nil
}
