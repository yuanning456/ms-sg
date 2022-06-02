package controller

import (
	"log"
	"ms-sg/constant"
	"ms-sg/server/common"
	"ms-sg/server/web/logic"
	"ms-sg/server/web/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {

}

func (a *AccountController) Register(ctx *gin.Context) {
	req := &model.RegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Println("参数传递错误", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数错误"))
		return
	}

	err = logic.DefaultAccountLogic.Register(req)
	if err != nil {
		log.Println("业务注册出错", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.Success(constant.OK, map[string]string{
		"msg": "注册成功",
	}))
}