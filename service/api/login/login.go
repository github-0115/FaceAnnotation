package login

import (
	user "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	security "FaceAnnotation/utils/security"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

type LoginParams struct {
	Managername string `json:"username"`
	Password    string `json:"password"`
}

func Login(c *gin.Context) {

	var loginParams LoginParams
	if err := c.BindJSON(&loginParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	username := loginParams.Managername
	password := loginParams.Password

	userColl, err := user.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrUserNotFound.Code,
			"message": vars.ErrUserNotFound.Msg,
		})
		return
	}

	savedPassword := security.GeneratePasswordHash(password)

	if !strings.EqualFold(userColl.Password, savedPassword) {
		log.Error(fmt.Sprintf("user login password error"))
		c.JSON(400, gin.H{
			"code":    vars.ErrUserNotFound.Code,
			"message": vars.ErrUserNotFound.Msg,
		})
		return
	}

	c.JSON(400, gin.H{
		"code":    0,
		"message": "login success !",
	})

}
