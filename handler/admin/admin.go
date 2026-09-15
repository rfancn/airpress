package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
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

// IsInstalledInput 检查安装状态无输入参数。
type IsInstalledInput struct{}

// IsInstalled 检查博客是否已安装。
func (a *AdminHandler) IsInstalled(ctx context.Context, _ *IsInstalledInput) (*dto.HumaOut[interface{}], error) {
	data, err := a.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](data)
}

// AuthPreCheckInput 登录预检输入。
type AuthPreCheckInput struct {
	Body param.LoginParam `doc:"登录参数"`
}

// AuthPreCheck 登录预检（判断是否需要 MFA 验证码）。
func (a *AdminHandler) AuthPreCheck(ctx context.Context, in *AuthPreCheckInput) (*dto.HumaOut[*dto.LoginPreCheckDTO], error) {
	user, err := a.AdminService.Authenticate(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[*dto.LoginPreCheckDTO](err)
	}
	return dto.HumaOK(&dto.LoginPreCheckDTO{NeedMFACode: a.TwoFactorMFAService.UseMFA(user.MfaType)})
}

// AuthInput 登录输入。
type AuthInput struct {
	Body param.LoginParam `doc:"登录参数"`
}

// Auth 管理员登录。
func (a *AdminHandler) Auth(ctx context.Context, in *AuthInput) (*dto.HumaOut[*dto.AuthTokenDTO], error) {
	data, err := a.AdminService.Auth(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[*dto.AuthTokenDTO](err)
	}
	return dto.HumaOK(data)
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

// RefreshTokenInput 刷新令牌输入。
type RefreshTokenInput struct {
	RefreshToken string `path:"refreshToken" doc:"刷新令牌"`
}

// RefreshToken 刷新访问令牌。
func (a *AdminHandler) RefreshToken(ctx context.Context, in *RefreshTokenInput) (*dto.HumaOut[*dto.AuthTokenDTO], error) {
	data, err := a.AdminService.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return dto.HumaErr[*dto.AuthTokenDTO](err)
	}
	return dto.HumaOK(data)
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
