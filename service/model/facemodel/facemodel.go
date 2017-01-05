package facemodel

import (
	db "FaceAnnotation/service/db"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FaceUrl struct {
	Url string `bson:"url" json:"url" binding:"required"`
}

type FaceModel struct {
	TaskId      string        `bson:"task_id" json:"task_id" binding:"required"` //与 taskmodel 关联
	Md5         string        `bson:"md5" json:"md5" binding:"required"`
	FiveResult  []*FiveResult `bson:"five_result" json:"five_result" binding:"required"`
	TwoResult   []*Result     `bson:"two_result" json:"two_result" binding:"required"`
	SixResult   []*Result     `bson:"six_result" json:"six_result" binding:"required"`
	EightResult []*Result     `bson:"eight_result" json:"eight_result" binding:"required"`
	NineResult  []*NineResult `bson:"nine_result" json:"nine_result" binding:"required"`
	ThrFaces    []*ThrFaces   `binding:"required" bson:"thr_faces" json:"thr_faces"`
}

type ImageModel struct {
	TaskId   string      `bson:"task_id" json:"task_id" binding:"required"` //与 taskmodel 关联
	Md5      string      `bson:"md5" json:"md5" binding:"required"`
	Results  *Result     `bson:"results" json:"results" binding:"required"`
	ThrFaces []*ThrFaces `binding:"required" bson:"thr_faces" json:"thr_faces"`
}

type ThrFaces struct {
	faceFivePoint       thrfacemodel.ThrFaceModel        `bson:"face_five_point" json:"face_five_point" binding:"required"`
	faceEightThreePoint thrfacemodel.EightThreeFaceModel `bson:"face_eight_thr_point" json:"face_eight_thr_point" binding:"required"`
}

type Result struct {
	Faces     map[string][]*Points `binding:"required" bson:"faces" json:"faces"`
	Status    int64                `bson:"status" json:"status" binding:"required"` //0:一般标图结果 1:微调结果
	CreatedAt time.Time            `bson:"created_at" json:"created_at" binding:"required"`
}

type Points struct {
	Points []*Point `binding:"required" bson:"points" json:"points"`
}

type Point struct {
	SmallTaskId []string `bson:"small_task_id" json:"small_task_id" binding:"required"`
	User        string   `bson:"user" json:"user" binding:"required"`
	Area        string   `bson:"area" json:"area" binding:"required"`
	Y           float32  `binding:"required" bson:"y" json:"y"`
	X           float32  `binding:"required" bson:"x" json:"x"`
}

type FiveResult struct {
	User      string        `bson:"user" json:"user" binding:"required"`
	Faces     []*FivePoints `binding:"required" bson:"faces" json:"faces"`
	Remark    int64         `bson:"remark" json:"remark" binding:"required"` //0:一般标图结果 1:微调结果
	CreatedAt time.Time     `bson:"created_at" json:"created_at" binding:"required"`
}

type NineResult struct {
	User      string        `bson:"user" json:"user" binding:"required"`
	Faces     []*NinePoints `binding:"required" bson:"faces" json:"faces"`
	Remark    int64         `bson:"remark" json:"remark" binding:"required"` //0:一般标图结果 1:微调结果
	CreatedAt time.Time     `bson:"created_at" json:"created_at" binding:"required"`
}

type NinePoints struct {
	LeftEyeBrow  *LeftEyeBrow  `binding:"required" bson:"left_eye_brow" json:"left_eye_brow"`
	RightEyeBrow *RightEyeBrow `binding:"required" bson:"right_eye_brow" json:"right_eye_brow"`
	LeftEye      *LeftEye      `binding:"required" bson:"left_eye" json:"left_eye"`
	RightEye     *RightEye     `binding:"required" bson:"right_eye" json:"right_eye"`
	LeftEar      *LeftEar      `binding:"required" bson:"left_ear" json:"left_ear"`
	RightEar     *RightEar     `binding:"required" bson:"right_ear" json:"right_ear"`
	Mouth        *Mouth        `binding:"required" bson:"mouth" json:"mouth"`
	Nouse        *Nouse        `binding:"required" bson:"nouse" json:"nouse"`
	Face         *Face         `binding:"required" bson:"face" json:"face"`
}

type FivePoints struct {
	EyeLeft    *Point `bson:"eye_left" json:"eye_left"`
	EyeRight   *Point `bson:"eye_right" json:"eye_right"`
	MouthLeft  *Point `bson:"mouth_left" json:"mouth_left"`
	MouthRight *Point `bson:"mouth_right" json:"mouth_right"`
	Nose       *Point `bson:"nose" json:"nose"`
}

type TwoPoints struct {
	LeftEyeBrow  *LeftEyeBrow  `binding:"required" bson:"left_eye_brow" json:"left_eye_brow"`
	RightEyeBrow *RightEyeBrow `binding:"required" bson:"right_eye_brow" json:"right_eye_brow"`
	LeftEye      *LeftEye      `binding:"required" bson:"left_eye" json:"left_eye"`
	RightEye     *RightEye     `binding:"required" bson:"right_eye" json:"right_eye"`
	Mouth        *Mouth        `binding:"required" bson:"mouth" json:"mouth"`
	Nouse        *Nouse        `binding:"required" bson:"nouse" json:"nouse"`
	Face         *Face         `binding:"required" bson:"face" json:"face"`
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

//func UpsertFaceResult(res *FaceModel) (bool, error) {
//	s := db.Face.GetSession()
//	defer s.Close()

//	_, err := s.DB(db.Face.DB).C("face").Upsert(bson.M{
//		"url": res.Url,
//	}, res)
//	if err != nil {
//		log.Error(fmt.Sprintf("upsert face point err ", err))
//		if err == mgo.ErrNotFound {
//			return false, ErrFaceModelNotFound
//		}
//		return false, ErrFaceModelCursor
//	}
//	return true, nil
//}

func StringToJson(res string) (*FaceModel, error) {

	rep := new(FaceModel)
	if err := json.Unmarshal([]byte(res), &rep); err != nil {
		log.Error(fmt.Sprintf("face json unmarshal err=%s", err))

		return nil, err
	}

	return rep, nil
}
