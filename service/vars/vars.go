package vars

type Err struct {
	Code int64
	Msg  string
}

var (
	ErrBindJSON      = &Err{100, "bind json error"}
	ErrUserNotFound  = &Err{101, "user not found"}
	ErrUserCursor    = &Err{102, "find user cursor err"}
	ErrUserNameExist = &Err{103, "user name exist err"}
	ErrUserIdExist   = &Err{104, "user id exist err"}
	ErrUserSave      = &Err{105, "user save err"}
	ErrLoginParams   = &Err{106, "user login params err"}
	ErrLogin         = &Err{107, "user login err"}
)
