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
	Title       string         `bson:"title" json:"title" binding:"required"`
	User        string         `bson:"user" json:"user" binding:"required"`
	Url         string         `bson:"url" json:"url" binding:"required"`
	Faces       []*FacesPoints `binding:"required" bson:"faces" json:"faces"`
	OriginFaces []*Landmarks   `binding:"required" bson:"origin_faces" json:"origin_faces"`
}

type FaceUrl struct {
	Url string `bson:"url" json:"url" binding:"required"`
}

type FacesPoints struct {
	LeftEyeBrow  []*Point `binding:"required" bson:"left_eye_brow" json:"left_eye_brow"`
	RightEyeBrow []*Point `binding:"required" bson:"right_eye_brow" json:"right_eye_brow"`
	LeftEye      []*Point `binding:"required" bson:"left_eye" json:"left_eye"`
	RightEye     []*Point `binding:"required" bson:"right_eye" json:"right_eye"`
	LeftEar      []*Point `binding:"required" bson:"left_ear" json:"left_ear"`
	RightEar     []*Point `binding:"required" bson:"right_ear" json:"right_ear"`
	Mouth        []*Point `binding:"required" bson:"mouth" json:"mouth"`
	Nouse        []*Point `binding:"required" bson:"nouse" json:"nouse"`
	Face         []*Point `binding:"required" bson:"face" json:"face"`
}

type Landmarks struct {
	Landmarks []*Landmark `binding:"required" bson:"landmarks" json:"landmarks"`
}

type LeftEyeBrow struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type RightEyeBrow struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type LeftEye struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type RightEye struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type LeftEar struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type RightEar struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Mouth struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Nouse struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Face struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Points struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Point struct {
	Y float32 `binding:"required" bson:"y" json:"y"`
	X float32 `binding:"required" bson:"x" json:"x"`
}

type Landmark struct {
	Landmark []*OnePoint `binding:"required" bson:"landmark" json:"landmark"`
}

type OnePoint struct {
	OnePoint float32 `binding:"required" bson:"one_point" json:"one_point"`
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
