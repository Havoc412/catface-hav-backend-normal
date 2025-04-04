package users

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Destroy struct {
	// 表单参数验证结构体支持匿名结构体嵌套、以及匿名结构体与普通字段组合
	Id
}

// 验证器语法，参见 Register.go文件，有详细说明

func (d Destroy) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&d); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(d, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserShow表单参数验证器json化失败", "")
		return
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Users{}).Destroy(extraAddBindDataContext)

		// 以下代码为模拟 前置、后置函数的回调代码
		/*
			func(before_callback_fn func(context *gin.Context) bool, after_callback_fn func(context *gin.Context)) {
				if before_callback_fn(extraAddBindDataContext) {
					defer after_callback_fn(extraAddBindDataContext)
					(&Web.Users{}).Destroy(extraAddBindDataContext)
				} else {
					// 这里编写前置函数验证不通过的相关返回提示逻辑...

				}
			}((&Users.DestroyBefore{}).Before, (&Users.DestroyAfter{}).After)
		*/
	}
}
