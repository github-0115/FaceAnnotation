package usertaskmodel

import (
	db "FaceAnnotation/service/db"
	"fmt"
	"time"
)

type UserTaskModel struct {
	TaskId         string        `bson:"task_id" json:"task_id" binding:"required"`           //与 taskmodel 关联
	UserTaskId     string        `bson:"user_task_id" json:"user_task_id" binding:"required"` //与 taskmodel 关联
	User           string        `binding:"required" bson:"user" json:"user"`
	UserTaskImages []*ImageModel `bson:"user_task_images" json:"user_task_images" binding:"required"`
	PointType      int64         `bson:"point_type" json:"point_type" binding:"required"`
	Areas          []string      `bson:"areas" json:"areas" binding:"required"`             //标识标注的哪个部位
	LimitCount     int64         `bson:"limit_count" json:"limit_count" binding:"required"` //标识标注的哪个部位
	Status         int64         `bson:"status" json:"status" binding:"required"`           // 0 未标  1 已标注完成
	CreatedAt      time.Time     `bson:"created_at" json:"created_at" binding:"required"`
}

type SmallTaskModel struct {
	SmallTaskImages []string    `bson:"small_task_images" json:"small_task_images" binding:"required"`
	LeftEyeBrow     []*UserTask `binding:"required" bson:"left_eye_brow" json:"left_eye_brow"`
	RightEyeBrow    []*UserTask `binding:"required" bson:"right_eye_brow" json:"right_eye_brow"`
	LeftEye         []*UserTask `binding:"required" bson:"left_eye" json:"left_eye"`
	RightEye        []*UserTask `binding:"required" bson:"right_eye" json:"right_eye"`
	LeftEar         []*UserTask `binding:"required" bson:"left_ear" json:"left_ear"`
	RightEar        []*UserTask `binding:"required" bson:"right_ear" json:"right_ear"`
	Mouth           []*UserTask `binding:"required" bson:"mouth" json:"mouth"`
	Nouse           []*UserTask `binding:"required" bson:"nouse" json:"nouse"`
	Face            []*UserTask `binding:"required" bson:"face" json:"face"`
	All             []*UserTask `binding:"required" bson:"face" json:"face"`
}

type ImageModel struct {
	Md5    string `bson:"md5" json:"md5" binding:"required"`
	Status int64  `bson:"status" json:"status" binding:"required"` // 0 未标  1 已标注完成
}

func (utk *UserTaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("user_task").Insert(&utk)
}
