package Response

import (
	"github.com/kataras/iris"
	"week13/pkg/zerror"
)

// ErrResponse 失败时回包
func ErrResponse(ctx iris.Context, err error) {
	// 处理error
	custom := zerror.HandleError(err)

	// 返回结果
	resp := iris.Map{
		"data":    nil,
		"code":    custom.Code,
		"message": custom.Msg,
	}
	ctx.JSON(resp)
}

// SuccResponse 成功时回包
func SuccResponse(ctx iris.Context, data interface{}, message string) {
	ctx.JSON(iris.Map{
		"code":    0,
		"message": message,
		"data":    data,
	})
}

