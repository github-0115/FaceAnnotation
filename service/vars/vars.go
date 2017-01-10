package vars

type Err struct {
	Code int64
	Msg  string
}

var (
	ErrBindJSON           = &Err{100, "bind json error"}
	ErrUserNotFound       = &Err{101, "user not found"}
	ErrUserCursor         = &Err{102, "find user cursor err"}
	ErrUserNameExist      = &Err{103, "user name exist err"}
	ErrUserIdExist        = &Err{104, "user id exist err"}
	ErrUserSave           = &Err{105, "user save err"}
	ErrLoginParams        = &Err{106, "user  params err"}
	ErrLogin              = &Err{107, "user login err"}
	ErrNeedToken          = &Err{301, "need auth token"}
	ErrInvalidToken       = &Err{301, "invalid token"}
	ErrIncompleteToken    = &Err{301, "token incomplete"}
	ErrTaskParmars        = &Err{400, "task list parmars err"}
	ErrTaskCompleted      = &Err{401, "The task has been completed !"}
	ErrTaskSave           = &Err{402, "task save err"}
	ErrTaskCursor         = &Err{403, "find task cursor err"}
	ErrTaskNotFound       = &Err{404, "task not found error"}
	ErrTaskListNotFound   = &Err{405, "task list not found error"}
	ErrImportTaskParmars  = &Err{406, "import task parmars err"}
	ErrTaskExist          = &Err{407, "task exist err"}
	ErrReadImportFile     = &Err{408, "read import task err"}
	ErrJsonUnmarshal      = &Err{600, "json unmarshal err"}
	ErrFaceParmars        = &Err{601, "face parmars is nil err"}
	ErrFaceModelUpsert    = &Err{602, "face points upsert error"}
	ErrFaceCursor         = &Err{603, "find face url cursor err"}
	ErrLocalImageGet      = &Err{700, "get local image error"}
	ErrImageParmars       = &Err{701, "image list parmars err"}
	ErrImageModelNotFound = &Err{702, "image model not found err"}
	ErrImageModelSave     = &Err{703, "image model save err"}
	ErrImageModelUpdate   = &Err{704, "image model update err"}
	ErrImportImageParmars = &Err{705, "import image parmars err"}
	ErrReadImage          = &Err{706, "read image error"}
	ErrNotImage           = &Err{707, "Under the task gets no pictures error"}
	ErrSmallTaskSave      = &Err{800, "create small task save err"}
	ErrSmallTaskNotFound  = &Err{801, "create small task not found err"}
	ErrNotSmallTask       = &Err{802, "not small task err"}
)
