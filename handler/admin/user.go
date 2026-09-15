package admin

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/impl"
	"github.com/rfancn/airpress/util/xerr"
)

type UserHandler struct {
	UserService         service.UserService
	TwoFactorMFAService service.TwoFactorTOTPMFAService
}

func NewUserHandler(userService service.UserService, twoFactorMFAService service.TwoFactorTOTPMFAService) *UserHandler {
	return &UserHandler{
		UserService:         userService,
		TwoFactorMFAService: twoFactorMFAService,
	}
}

// GetCurrentUserProfileInput 当前用户资料查询无输入参数。
type GetCurrentUserProfileInput struct{}

// GetCurrentUserProfile 获取当前用户资料。
func (u *UserHandler) GetCurrentUserProfile(ctx context.Context, _ *GetCurrentUserProfileInput) (*dto.HumaOut[*dto.User], error) {
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil {
		return dto.HumaErr[*dto.User](err)
	}
	return dto.HumaOK(u.UserService.ConvertToDTO(ctx, user))
}

// UpdateUserProfileInput 更新用户资料输入。
type UpdateUserProfileInput struct {
	Body param.User `doc:"用户参数"`
}

// UpdateUserProfile 更新用户资料。
func (u *UserHandler) UpdateUserProfile(ctx context.Context, in *UpdateUserProfileInput) (*dto.HumaOut[*dto.User], error) {
	user, err := u.UserService.Update(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.User](err)
	}
	return dto.HumaOK(u.UserService.ConvertToDTO(ctx, user))
}

// UpdatePasswordInput 更新密码输入。
type UpdatePasswordInput struct {
	Body UpdatePasswordBody `doc:"密码参数"`
}

// UpdatePasswordBody 密码请求体。
type UpdatePasswordBody struct {
	OldPassword string `json:"oldPassword" doc:"旧密码"`
	NewPassword string `json:"newPassword" doc:"新密码"`
}

// UpdatePassword 更新密码。
func (u *UserHandler) UpdatePassword(ctx context.Context, in *UpdatePasswordInput) (*dto.HumaOut[any], error) {
	err := u.UserService.UpdatePassword(ctx, in.Body.OldPassword, in.Body.NewPassword)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// GenerateMFAQRCodeInput 生成MFA二维码输入。
type GenerateMFAQRCodeInput struct {
	Body GenerateMFAQRCodeBody `doc:"MFA参数"`
}

// GenerateMFAQRCodeBody MFA请求体。
type GenerateMFAQRCodeBody struct {
	MFAType *consts.MFAType `json:"mfaType" doc:"MFA类型"`
}

// GenerateMFAQRCode 生成MFA二维码。
func (u *UserHandler) GenerateMFAQRCode(ctx context.Context, in *GenerateMFAQRCodeInput) (*dto.HumaOut[*vo.MFAFactorAuth], error) {
	if in.Body.MFAType == nil {
		return dto.HumaErr[*vo.MFAFactorAuth](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("parameter error"))
	}
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil {
		return dto.HumaErr[*vo.MFAFactorAuth](err)
	}

	mfaFactorAuthDTO := &vo.MFAFactorAuth{}
	if *in.Body.MFAType == consts.MFATFATotp {
		key, url, err := u.TwoFactorMFAService.GenerateOTPKey(ctx, user.Nickname)
		if err != nil {
			return dto.HumaErr[*vo.MFAFactorAuth](err)
		}
		mfaFactorAuthDTO.MFAType = consts.MFATFATotp
		mfaFactorAuthDTO.OptAuthURL = url
		mfaFactorAuthDTO.MFAKey = key
		qrCode, err := u.TwoFactorMFAService.GenerateMFAQRCode(ctx, url)
		if err != nil {
			return dto.HumaErr[*vo.MFAFactorAuth](err)
		}
		mfaFactorAuthDTO.QRImage = qrCode
		return dto.HumaOK(mfaFactorAuthDTO)
	} else {
		return dto.HumaErr[*vo.MFAFactorAuth](xerr.WithMsg(nil, "Not supported authentication").WithStatus(xerr.StatusBadRequest))
	}
}

// UpdateMFAInput 更新MFA输入。
type UpdateMFAInput struct {
	Body UpdateMFABody `doc:"MFA参数"`
}

// UpdateMFABody MFA更新请求体。
type UpdateMFABody struct {
	MFAType  *consts.MFAType `json:"mfaType" doc:"MFA类型"`
	MFAKey   string          `json:"mfaKey" doc:"MFA密钥"`
	AuthCode string          `json:"authcode" doc:"验证码"`
}

// UpdateMFA 更新MFA。
func (u *UserHandler) UpdateMFA(ctx context.Context, in *UpdateMFAInput) (*dto.HumaOut[any], error) {
	if in.Body.MFAType == nil {
		return dto.HumaErr[any](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("parameter error"))
	}
	err := u.UserService.UpdateMFA(ctx, in.Body.MFAKey, *in.Body.MFAType, in.Body.AuthCode)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
