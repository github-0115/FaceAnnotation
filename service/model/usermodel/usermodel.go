package usermodel

import (
	db "FaceAnnotation/service/db"
	"errors"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserColl struct {
	UserId    string    `bson:"user_id" binding:"required"`
	Username  string    `binding:"required" bson:"username" json:"username"`
	Password  string    `binding:"required" bson:"password" json:"password"`
	Identity  string    `binding:"required" bson:"identity" json:"identity"` //admin、fineTune、normal
	Status    int64     `binding:"required" bson:"status" json:"status"`
	CreatedAt time.Time `bson:"created_at" binding:"required" json:"created_at"`
}

type userIdentity struct {
	Admin    string
	Normal   string
	FineTune string
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserCursor   = errors.New("Cursor err")
	UserIdentity    = &userIdentity{"admin", "normal", "fineTune"}
)

func (u *UserColl) Save() error {
	s := db.User.GetSession()
	defer s.Close()
	return s.DB(db.User.DB).C("user").Insert(&u)
}

func QueryUser(username string) (*UserColl, error) {
	coll := new(UserColl)
	s := db.User.GetSession()
	defer s.Close()
	err := s.DB(db.User.DB).C("user").Find(bson.M{
		"username": username,
	}).One(coll)

	if err != nil {
		log.Error(fmt.Sprintf("find user err ", err))
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
	count, err := s.DB(db.User.DB).C("user").Find(bson.M{
		"user_id": userId,
	}).Count()

	if err != nil {
		log.Error(fmt.Sprintf("find user err ", err))
		if err == mgo.ErrNotFound {
			return 0, ErrUserNotFound
		} else if err == mgo.ErrCursor {
			return 0, ErrUserCursor
		}
		return 0, err
	}
	return count, nil
}
