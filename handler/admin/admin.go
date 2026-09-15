package admin

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/rfancn/airpress/handler/trans"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util/xerr"
)

type AdminHandler struct {
	OptionService       service.OptionService
	AdminService        service.AdminService
	TwoFactorMFAService service.TwoFactorTOTPMFAService
}

func NewAdminHandler(optionService service.OptionService, adminService service.AdminService, twoFactorMFA service.TwoFactorTOTPMFAService) *AdminHandler {
	return &AdminHandler{
		OptionService:       optionService,
		AdminService:        adminService,
		TwoFactorMFAService: twoFactorMFA,
	}
}

func (a *AdminHandler) IsInstalled(ctx *gin.Context) (interface{}, error) {
	return a.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
}

func (a *AdminHandler) AuthPreCheck(ctx *gin.Context) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.ShouldBindJSON(&loginParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "")
	}

	user, err := a.AdminService.Authenticate(ctx, loginParam)
	if err != nil {
		return nil, err
	}
	return &dto.LoginPreCheckDTO{NeedMFACode: a.TwoFactorMFAService.UseMFA(user.MfaType)}, nil
}

func (a *AdminHandler) Auth(ctx *gin.Context) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.ShouldBindJSON(&loginParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "").WithStatus(xerr.StatusBadRequest)
	}

	return a.AdminService.Auth(ctx, loginParam)
}

// LogOutInput 登出无输入参数。
type LogOutInput struct{}

// LogOut 管理员登出。
func (a *AdminHandler) LogOut(ctx context.Context, _ *LogOutInput) (*dto.HumaOut[interface{}], error) {
	err := a.AdminService.ClearToken(ctx)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// SendResetCodeInput 发送重置密码验证码输入。
type SendResetCodeInput struct {
	Body param.ResetPasswordParam `doc:"重置密码参数"`
}

// SendResetCode 发送重置密码验证码。
func (a *AdminHandler) SendResetCode(ctx context.Context, in *SendResetCodeInput) (*dto.HumaOut[interface{}], error) {
	err := a.AdminService.SendResetPasswordCode(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

func (a *AdminHandler) RefreshToken(ctx *gin.Context) (interface{}, error) {
	refreshToken := ctx.Param("refreshToken")
	if refreshToken == "" {
		return nil, xerr.BadParam.New("refreshToken参数为空").WithStatus(xerr.StatusBadRequest).
			WithMsg("refreshToken 参数不能为空")
	}
	return a.AdminService.RefreshToken(ctx, refreshToken)
}

// GetEnvironmentsInput 获取环境信息无输入参数。
type GetEnvironmentsInput struct{}

// GetEnvironments 获取环境信息。
func (a *AdminHandler) GetEnvironments(ctx context.Context, _ *GetEnvironmentsInput) (*dto.HumaOut[*dto.EnvironmentDTO], error) {
	return dto.HumaOK(a.AdminService.GetEnvironments(ctx))
}

// GetLogFilesInput 获取日志文件输入。
type GetLogFilesInput struct {
	Lines int64 `query:"lines" required:"true" doc:"日志行数"`
}

// GetLogFiles 获取日志文件内容。
func (a *AdminHandler) GetLogFiles(ctx context.Context, in *GetLogFilesInput) (*dto.HumaOut[string], error) {
	result, err := a.AdminService.GetLogFiles(ctx, in.Lines)
	if err != nil {
		return dto.HumaErr[string](err)
	}
	return dto.HumaOK(result)
}
