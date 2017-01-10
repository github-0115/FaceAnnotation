package smalltaskmodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SmallTaskModel struct {
	TaskId          string    `bson:"task_id" json:"task_id" binding:"required"`             //与 taskmodel 关联
	SmallTaskId     string    `bson:"small_task_id" json:"small_task_id" binding:"required"` //与 taskmodel 关联
	SmallTaskImages []string  `bson:"small_task_images" json:"small_task_images" binding:"required"`
	PointType       int64     `bson:"point_type" json:"point_type" binding:"required"`
	Areas           string    `bson:"areas" json:"areas" binding:"required"`             //标识标注的哪个部位
	LimitCount      int64     `bson:"limit_count" json:"limit_count" binding:"required"` //人数限制
	Status          int64     `bson:"status" json:"status" binding:"required"`           // 0 创建成功  1  正在进行中 2 已标注完成
	CreatedAt       time.Time `bson:"created_at" json:"created_at" binding:"required"`
}

type taskStatus struct {
	NotStart int64
	Start    int64
	Success  int64
}

type ImageModel struct {
	User   []string `binding:"required" bson:"user" json:"user"`
	Md5    string   `bson:"md5" json:"md5" binding:"required"`
	Status int64    `bson:"status" json:"status" binding:"required"` // 0 未标  1 已标注完成
}

var (
	ErrSmallTaskModelNotFound = errors.New("small task Model not found")
	ErrSmallTaskModelCursor   = errors.New("Cursor err")
	TaskStatus                = &taskStatus{0, 1, 2}
)

func (stm *SmallTaskModel) Save() error {
	s := db.Face.GetSession()
	defer s.Close()
	return s.DB(db.Face.DB).C("small_task").Insert(&stm)
}

func QuerySmallTask(id string) (*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	coll := new(SmallTaskModel)
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"small_task_id": id,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find small task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return coll, nil
}

func QueryNotSmallTask() ([]*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": bson.M{"$ne": TaskStatus.Success},
	}).Sort("areas").All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small_task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QueryNorNotSmallTask() ([]*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"areas":  bson.M{"$ne": "fineTune"},
		"status": bson.M{"$ne": TaskStatus.Success},
	}).Sort("areas").All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small_task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QueryAreaNotSmallTask(area string) ([]*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"areas":  area,
		"status": bson.M{"$ne": TaskStatus.Success},
	}).Sort("areas").All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small_task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QuerySmallTasks(status int64) ([]*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": status,
	}).Sort("areas").All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}

func QueryPageNotSmallTask(pageIndex int, pageSize int) ([]*SmallTaskModel, int, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": bson.M{"$ne": TaskStatus.Success},
	}).Sort("areas").Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&results)

	count, err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": bson.M{"$ne": TaskStatus.Success},
	}).Count()

	if err != nil {
		log.Error(fmt.Sprintf("find small_task err ", err))
		if err == mgo.ErrNotFound {
			return nil, 0, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, 0, ErrSmallTaskModelCursor
		}
		return nil, 0, err
	}
	return results, count, nil
}

func QueryPageSmallTasks(status int64, pageIndex int, pageSize int) ([]*SmallTaskModel, int, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": status,
	}).Sort("areas").Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&results)

	count, err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"status": status,
	}).Count()

	if err != nil {
		log.Error(fmt.Sprintf("find small task err ", err))
		if err == mgo.ErrNotFound {
			return nil, 0, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, 0, ErrSmallTaskModelCursor
		}
		return nil, 0, err
	}
	return results, count, nil
}

func QueryTaskSmallTasks(taskId string) ([]*SmallTaskModel, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*SmallTaskModel
	err := s.DB(db.Face.DB).C("small_task").Find(bson.M{
		"task_id": taskId,
		"status":  bson.M{"$ne": TaskStatus.Success},
	}).Sort("areas").All(&results)

	if err != nil {
		log.Error(fmt.Sprintf("find small task err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrSmallTaskModelCursor
		}
		return nil, err
	}
	return results, nil
}

func UpdateSmallTasks(id string, value int64) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("small_task").Update(bson.M{
		"small_task_id": id,
	}, bson.M{"$set": bson.M{"status": value}})

	if err != nil {
		log.Error(fmt.Sprintf("find small task err ", err))
		if err == mgo.ErrNotFound {
			return ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return ErrSmallTaskModelCursor
		}
		return err
	}
	return nil
}

func RemoveSmallTask(taskId string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("small_task").Remove(bson.M{
		"task_id": taskId,
	})
	if err != nil {
		log.Error(fmt.Sprintf("remove small task err ", err))
		if err == mgo.ErrNotFound {
			return ErrSmallTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return ErrSmallTaskModelCursor
		}
		return err
	}
	return nil
}
