package zerror

type ErrorLevel = string

const (
	Debug ErrorLevel = "Debug"
	Info  ErrorLevel = "Info"
	Warn  ErrorLevel = "Warn"
	Error ErrorLevel = "Error"
	Fatal ErrorLevel = "Fatal"
)

type Err struct {
	Code  int
	Msg   string
	Level string
}

func (e Err) Error() string {
	return e.Msg
}

func (e *Err) SetErrorMsg(msg string) {
	e.Msg = msg
}

func (e *Err) SetErrorLevel(level ErrorLevel)  {
	e.Level = level
}

var ErrServerDown = Err{10, "服务进程异常", Error}
var ErrSignalDone = Err{11, "系统关闭服务进程", Warn}

var ErrUnknown = Err{100, "服务未知错误", Error}
var ErrInnerServer = Err{101, "服务内部错误", Error}
var ErrUserInfo = Err{102, "用户信息错误", Error}
var ErrParamParse = Err{103, "请求参数异常", Warn}
var ErrRequestTooFrequent = Err{104, "操作过于频繁，请稍后再试", Warn}
