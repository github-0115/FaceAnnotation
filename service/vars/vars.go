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
	ErrJsonUnmarshal    = &Err{200, "json unmarshal err"}
	ErrFaceModelUpsert  = &Err{201, "face points upsert error"}
	ErrFaceCursor       = &Err{202, "find face url cursor err"}
	ErrLocalImageGet    = &Err{300, "get local image error"}
	ErrTaskListNotFound = &Err{400, "task list not found error"}
	ErrTaskNotFound     = &Err{401, "task not found error"}
	ErrTaskCursor       = &Err{402, "find task cursor err"}
	ErrTaskCompleted    = &Err{403, "The task has been completed !"}
)
