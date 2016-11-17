package vars

type Err struct {
	Code int64
	Msg  string
}

var (
	ErrBindJSON     = &Err{100, "bind json error"}
	ErrLoginParams  = &Err{101, "login params err"}
	ErrUserNotFound = &Err{102, "user not found"}
	ErrLogin        = &Err{103, "user login err"}
)
