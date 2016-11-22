package redisclient

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	cfg "FaceAnnotation/config"

	log "github.com/inconshreveable/log15"
	"gopkg.in/redis.v3"
)

const imageurlKey = "face_image_url:"

var redisAddr = cfg.DBCfg.RedisBroker.Host + ":" + strconv.Itoa(cfg.DBCfg.RedisBroker.Port)
var RedisNotFound = errors.New("redis not found err")
var client = redis.NewClient(&redis.Options{
	Addr:     redisAddr,
	Password: cfg.DBCfg.RedisBroker.Password,
	DB:       cfg.DBCfg.RedisBroker.DB,
})

func SetCheckImageStr(imageName string, isYes string) error {
	err := client.Set(imageurlKey+imageName, isYes, time.Duration(20)*time.Minute).Err()
	if err != nil {
		log.Error(fmt.Sprintf("imageurlKey: imageName=%s, isYes=%s, err=%#v.", imageName, isYes, err))
		return err
	}
	return nil
}

func GetCheckImageStr(imageName string) (string, error) {
	isYes, err := client.Get(imageurlKey + imageName).Result()
	if err == redis.Nil {
		log.Error(fmt.Sprintf("imageurlKey: imageName=%s, isYes=%s, err=%#v.", imageName, isYes, err))
		return "", RedisNotFound
	} else if err != nil {

		return "", err
	}
	return isYes, nil
}

func DeleteCheckImageStr(imageName string) error {
	err := client.Del(imageurlKey + imageName).Err()
	if err == redis.Nil {
		log.Error(fmt.Sprintf("del imageurlKey: imageName=%s, err=%#v.", imageName, err))
		return RedisNotFound
	} else if err != nil {

		return err
	}
	return nil
}
