package vars

type Err struct {
	Code int64
	Msg  string
}

var (
	ErrBindJSON         = &Err{100, "bind json error"}
	ErrUserNotFound     = &Err{101, "user not found"}
	ErrUserCursor       = &Err{102, "find user cursor err"}
	ErrUserNameExist    = &Err{103, "user name exist err"}
	ErrUserIdExist      = &Err{104, "user id exist err"}
	ErrUserSave         = &Err{105, "user save err"}
	ErrLoginParams      = &Err{106, "user login params err"}
	ErrLogin            = &Err{107, "user login err"}
	ErrNeedToken        = &Err{301, "need auth token"}
	ErrInvalidToken     = &Err{301, "invalid token"}
	ErrIncompleteToken  = &Err{301, "token incomplete"}
	ErrTaskParmars      = &Err{400, "task list parmars err"}
	ErrTaskCompleted    = &Err{401, "The task has been completed !"}
	ErrTaskNotFound     = &Err{402, "task not found error"}
	ErrTaskCursor       = &Err{403, "find task cursor err"}
	ErrTaskListNotFound = &Err{404, "task list not found error"}
	ErrJsonUnmarshal    = &Err{600, "json unmarshal err"}
	ErrFaceParmars      = &Err{601, "face parmars is nil err"}
	ErrFaceModelUpsert  = &Err{602, "face points upsert error"}
	ErrFaceCursor       = &Err{603, "find face url cursor err"}
	ErrLocalImageGet    = &Err{700, "get local image error"}
	ErrImageParmars     = &Err{701, "image list parmars err"}
)
