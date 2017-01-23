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

func GetImageTasks() ([]*ImageTaskModel, error) {
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

func QueryPageImageTasks(pageIndex int, pageSize int) ([]*ImageTaskModel, int, error) {
	s := db.Face.GetSession()
	defer s.Close()

	var results []*ImageTaskModel
	err := s.DB(db.Face.DB).C("image_task").Find(nil).Sort("created_at").Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(&results)

	count, err := s.DB(db.Face.DB).C("image_task").Find(nil).Count()

	if err != nil {
		log.Error(fmt.Sprintf("find image task err ", err))
		if err == mgo.ErrNotFound {
			return nil, 0, ErrImageTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return nil, 0, ErrImageTaskModelCursor
		}
		return nil, 0, err
	}
	return results, count, nil
}

func UpdateImageTaskImages(taskId string, md5 string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image_task").Update(bson.M{
		"image_task_id": taskId,
	}, bson.M{"$push": bson.M{"images": md5}})
	if err != nil {
		log.Error(fmt.Sprintf("upsert image point err ", err))
		if err == mgo.ErrNotFound {
			return ErrImageTaskModelNotFound
		}
		return ErrImageTaskModelCursor
	}
	return nil
}

func UpdateImageTaskTaskId(imagetaskId string, taskId string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image_task").Update(bson.M{
		"image_task_id": imagetaskId,
	}, bson.M{"$push": bson.M{"task_id": taskId}})
	if err != nil {
		log.Error(fmt.Sprintf("upsert image point err ", err))
		if err == mgo.ErrNotFound {
			return ErrImageTaskModelNotFound
		}
		return ErrImageTaskModelCursor
	}
	return nil
}

func UpdateImageTaskIntroduce(imagetaskId string, introduce string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image_task").Update(bson.M{
		"image_task_id": imagetaskId,
	}, bson.M{"$set": bson.M{"introduce": introduce}})
	if err != nil {
		log.Error(fmt.Sprintf("upsert image point err ", err))
		if err == mgo.ErrNotFound {
			return ErrImageTaskModelNotFound
		}
		return ErrImageTaskModelCursor
	}
	return nil
}

func RemoveImageTask(taskId string) error {
	s := db.Face.GetSession()
	defer s.Close()

	err := s.DB(db.Face.DB).C("image_task").Remove(bson.M{
		"image_task_id": taskId,
	})

	if err != nil {
		log.Error(fmt.Sprintf("remove image task err ", err))
		if err == mgo.ErrNotFound {
			return ErrImageTaskModelNotFound
		} else if err == mgo.ErrCursor {
			return ErrImageTaskModelCursor
		}
		return err
	}
	return nil
}
