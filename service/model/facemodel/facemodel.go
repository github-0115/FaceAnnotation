package facemodel

import (
	db "FaceAnnotation/service/db"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FaceModel struct {
	Title string    `bson:"title" json:"title" binding:"required"`
	User  string    `bson:"user" json:"user" binding:"required"`
	Url   string    `bson:"url" json:"url" binding:"required"`
	Faces []*Points `binding:"required" bson:"faces" json:"faces"`
}

type FaceUrl struct {
	Url string `bson:"url" json:"url" binding:"required"`
}

type Points struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Point struct {
	Y float32 `binding:"required" bson:"y" json:"y"`
	X float32 `binding:"required" bson:"x" json:"x"`
}

var (
	ErrFaceModelNotFound = errors.New("Face Model not found")
	ErrFaceModelCursor   = errors.New("Cursor err")
)

func QueryAll() ([]*FaceUrl, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var result []*FaceUrl
	err := s.DB(db.Face.DB).C("face").Find(nil).Select(bson.M{"url": 1}).All(&result)
	if err != nil {
		log.Error(fmt.Sprintf("upsert face point err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrFaceModelNotFound
		}
		return nil, ErrFaceModelCursor
	}
	return result, nil
}

func UpsertFaceResult(res *FaceModel) (bool, error) {
	s := db.Face.GetSession()
	defer s.Close()

	_, err := s.DB(db.Face.DB).C("face").Upsert(bson.M{
		"url": res.Url,
	}, res)
	if err != nil {
		log.Error(fmt.Sprintf("upsert face point err ", err))
		if err == mgo.ErrNotFound {
			return false, ErrFaceModelNotFound
		}
		return false, ErrFaceModelCursor
	}
	return true, nil
}

func StringToJson(res string) (*FaceModel, error) {

	rep := new(FaceModel)
	if err := json.Unmarshal([]byte(res), &rep); err != nil {
		log.Error(fmt.Sprintf("face json unmarshal err=%s", err))

		return nil, err
	}

	return rep, nil
}
