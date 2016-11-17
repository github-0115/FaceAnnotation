package user

import (
	user "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	security "FaceAnnotation/utils/security"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/satori/go.uuid"
)

type AddUserParams struct {
	Managername string `json:"username"`
	Password    string `json:"password"`
}

func AddUser(c *gin.Context) {

	var addUserParams AddUserParams
	if err := c.BindJSON(&addUserParams); err != nil {
		log.Error(fmt.Sprintf("bind json error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrBindJSON.Code,
			"message": vars.ErrBindJSON.Msg,
		})
		return
	}
	username := addUserParams.Managername
	password := addUserParams.Password

	userColl, err := user.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		if err != user.ErrUserNotFound {
			c.JSON(400, gin.H{
				"code":    vars.ErrUserCursor.Code,
				"message": vars.ErrUserCursor.Msg,
			})
			return
		}
	}

	if userColl != nil {
		log.Error(fmt.Sprintf("user name exist err"))
		c.JSON(400, gin.H{
			"code":    vars.ErrUserNameExist.Code,
			"message": vars.ErrUserNameExist.Msg,
		})
		return
	}

	user_id := uuid.NewV4().String()
	existUserNum, err := user.QueryUserExist(user_id)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		if err != user.ErrUserNotFound {
			c.JSON(400, gin.H{
				"code":    vars.ErrUserCursor.Code,
				"message": vars.ErrUserCursor.Msg,
			})
			return
		}
	}

	if existUserNum >= 1 {
		log.Error(fmt.Sprintf("user num >= 1, %i", existUserNum, ", user_id = %v", user_id))
		c.JSON(400, gin.H{
			"code":    vars.ErrUserIdExist.Code,
			"message": vars.ErrUserIdExist.Msg,
		})
		return
	}

	savedPassword := security.GeneratePasswordHash(password)
	user_new := &user.UserColl{
		UserId:    user_id,
		Username:  username,
		Password:  savedPassword,
		CreatedAt: time.Now(),
	}

	err = user_new.Save()
	if err != nil {
		log.Error(fmt.Sprintf("user save error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrUserSave.Code,
			"message": vars.ErrUserSave.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "add user success !",
	})
}
