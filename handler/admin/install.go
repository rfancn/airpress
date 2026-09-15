package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type InstallHandler struct {
	InstallService service.InstallService
}

func NewInstallHandler(installService service.InstallService) *InstallHandler {
	return &InstallHandler{
		InstallService: installService,
	}
}

// InstallBlogInput 安装博客输入。
type InstallBlogInput struct {
	Body param.Install `doc:"安装参数"`
}

// InstallBlog 安装博客（huma 风格，挂在无鉴权的 admin public group 下）。
func (i *InstallHandler) InstallBlog(ctx context.Context, in *InstallBlogInput) (*dto.HumaOut[string], error) {
	if err := i.InstallService.InstallBlog(ctx, in.Body); err != nil {
		return dto.HumaErr[string](err)
	}
	return dto.HumaOK("安装完成")
}
