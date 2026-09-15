package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type EmailHandler struct {
	EmailService service.EmailService
}

func NewEmailHandler(emailService service.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailService: emailService,
	}
}

// TestEmailInput 发送测试邮件输入。
type TestEmailInput struct {
	Body param.TestEmail `doc:"测试邮件参数"`
}

// Test 发送测试邮件。
func (e *EmailHandler) Test(ctx context.Context, in *TestEmailInput) (*dto.HumaOut[any], error) {
	err := e.EmailService.SendTextEmail(ctx, in.Body.To, in.Body.Subject, in.Body.Content)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
