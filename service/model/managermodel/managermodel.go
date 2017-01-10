package managermodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ManagerColl struct {
	ManagerId   string    `bson:"manager_id" json:"manager_id" binding:"required"`
	Managername string    `binding:"required" bson:"managername" json:"managername"`
	Password    string    `binding:"required" bson:"password" json:"password"`
	CreatedAt   time.Time `bson:"created_at" binding:"required" json:"created_at"`
}

var (
	ErrUserNotFound = errors.New("manager not found")
	ErrUserCursor   = errors.New("Cursor err")
)

func (u *ManagerColl) Save() error {
	s := db.User.GetSession()
	defer s.Close()
	return s.DB(db.User.DB).C("manager").Insert(&u)
}

func QueryUser(username string) (*ManagerColl, error) {
	coll := new(ManagerColl)
	s := db.User.GetSession()
	defer s.Close()
	err := s.DB(db.User.DB).C("manager").Find(bson.M{
		"managername": username,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find manager err ", err))
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		} else if err == mgo.ErrCursor {
			return nil, ErrUserCursor
		}
		return nil, err
	}
	return coll, nil
}

func QueryUserExist(userId string) (int, error) {

	s := db.User.GetSession()
	defer s.Close()
	count, err := s.DB(db.User.DB).C("manager").Find(bson.M{
		"manager_id": userId,
	}).Count()

	if err != nil {
		log.Error(fmt.Sprintf("find manager err ", err))
		if err == mgo.ErrNotFound {
			return 0, ErrUserNotFound
		} else if err == mgo.ErrCursor {
			return 0, ErrUserCursor
		}
		return 0, err
	}
	return count, nil
}
