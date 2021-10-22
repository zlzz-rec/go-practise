package zerror

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.feedtoken.tech/ztouch/zglib/zstring"
	"strings"
	"week13/internal/myapp/config"
	"week13/pkg/log"
)

func fmtError(err error) string {
	var length int
	arr := zstring.Split(fmt.Sprintf("%+v", err), "\n")
	if len(arr) > 15 {
		length = 15
	} else {
		length = len(arr)
	}
	return strings.Join(arr[:length], "\n")
}

func HandleError(err error) Err {
	// 获取error中的string
	final := fmtError(err)
	if config.Opts.Debug {
		fmt.Println(final)
	}

	// 获取根因
	root := errors.Cause(err)
	custom, ok := root.(Err)
	if !ok {
		custom = ErrUnknown
	}

	// 处理日志
	switch custom.Level {
	case Debug:
		log.Logger.Debug(final)
	case Info:
		log.Logger.Info(final)
	case Warn:
		log.Logger.Warn(final)
	case Error:
		log.Logger.Error(final)
	case Fatal:
		log.Logger.Fatal(final)
	}
	return custom
}
