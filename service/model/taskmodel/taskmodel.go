package taskmodel

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
	TaskId    string    `bson:"task_id" json:"task_id" binding:"required"`        //与 imagemodel 关联
	Area      []string  `bson:"area" json:"area"  binding:"required"`             //标注区域
	PointType int64     `bson:"point_type" json:"point_type"  binding:"required"` //5、27、68、95点
	MinUnit   int64     `bson:"min_unit" json:"min_unit"  binding:"required"`     //拆分最小单元
	LimitUser int64     `bson:"limit_user" json:"limit_user"  binding:"required"` //标注人员数量限制
	Count     int64     `bson:"count" json:"count" binding:"required"`            //本次任务数量
	Introduce string    `bson:"introduce" json:"introduce" binding:"required"`    //本次说明
	Status    int64     `bson:"status" json:"status" binding:"required"`          // taskStatus 0 创建成功  1 导入图片成功  2 导入图片失败 3 正在进行中 4 标注完成
	CreatedAt time.Time `bson:"created_at" json:"created_at" binding:"required"`
}

type taskStatus struct {
	NotStart      int64 // 0 创建成功
	ImportSuccess int64 // 1 导入图片成功
	ImportFail    int64 // 2 导入图片失败
	Start         int64 // 3 正在进行中
	Success       int64 // 4 标注完成
}

type ImportModel struct {
	Url         string   `bson:"url" json:"url"`
	OriginFaces []string `binding:"required" bson:"origin_faces" json:"origin_faces"`
}

var (
	ErrTaskModelNotFound = errors.New("task Model not found")
	ErrTaskModelCursor   = errors.New("Cursor err")
	ErrDirNotFound       = errors.New("dir ads not found")
	ErrFileNotFound      = errors.New("file ads not found")
	ErrReadFile          = errors.New("read file err")
	TaskStatus           = &taskStatus{0, 1, 2, 3, 4}
)

func (t *TaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("task").Insert(&t)
}

func QueryTask(taskId string) (*TaskModel, error) {
	coll := new(TaskModel)
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("task").Find(bson.M{
		"task_id": taskId,
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

func UpdateTaskStatus(taskId string, status int64) error {

	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("task").Update(bson.M{
		"task_id": taskId,
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
