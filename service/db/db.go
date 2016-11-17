package db

import (
	c "FaceAnnotation/config"
	"fmt"

	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
)

var (
	Session   *mgo.Session
	DBUserUri string
	User      *MongoSession
)

func init() {
	DBUserUri = c.DBCfg.UserCenterMongoTask.DB
	mongoUrl := c.DBCfg.UserCenterMongoTask.String()

	var err error
	Session, err = mgo.Dial(mongoUrl)
	if err != nil {
		log.Error(fmt.Sprintf("connect Mongo error", err))
		panic(err)
	}

	User = new(MongoSession).Init(
		c.DBCfg.UserCenterMongoTask,
		"user",
	)
}
