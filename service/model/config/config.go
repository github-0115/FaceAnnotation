package config

import (
	db "FaceAnnotation/service/db"
)

var BanTokens map[string]bool

type Config struct {
	BanTokens []string `bson:"ban_tokens"`
}

func GetBanTokens() ([]string, error) {
	config := new(Config)
	s := db.User.GetSession()
	defer s.Close()
	err := s.DB(db.User.DB).C("config").Find(nil).One(config)
	if err != nil {
		return nil, err
	}
	return config.BanTokens, nil
}

func LoadBanTokens() {
	BanTokens = make(map[string]bool)

	banTokens, err := GetBanTokens()
	if err != nil {
		return
	}
	for _, token := range banTokens {
		BanTokens[token] = true
	}
}
