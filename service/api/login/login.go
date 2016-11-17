package login

import (
	user "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	security "FaceAnnotation/utils/security"
	"fmt"
	//	"strings"

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
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	if !security.CheckPasswordHash(password, userColl.Password) {
		log.Error(fmt.Sprintf("password error.username=%s, password=%s, saved_password=%s", username, password, userColl.Password))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "login success !",
	})

}
