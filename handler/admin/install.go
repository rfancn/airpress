package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/rfancn/airpress/handler/trans"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util/xerr"
)

type InstallHandler struct {
	InstallService service.InstallService
}

func NewInstallHandler(installService service.InstallService) *InstallHandler {
	return &InstallHandler{
		InstallService: installService,
	}
}

func (i *InstallHandler) InstallBlog(ctx *gin.Context) (interface{}, error) {
	var installParam param.Install
	err := ctx.ShouldBindJSON(&installParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	err = i.InstallService.InstallBlog(ctx, installParam)
	if err != nil {
		return nil, err
	}
	return "安装完成", nil
}
