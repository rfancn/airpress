package api

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
)

type OptionHandler struct {
	OptionService service.OptionService
}

func NewOptionHandler(
	optionService service.OptionService,
) *OptionHandler {
	return &OptionHandler{
		OptionService: optionService,
	}
}

// CommentInput 评论选项查询无输入参数。
type CommentInput struct{}

// Comment 获取评论相关选项。
func (o *OptionHandler) Comment(ctx context.Context, _ *CommentInput) (*dto.HumaOut[map[string]interface{}], error) {
	result := make(map[string]interface{})

	result[property.CommentGravatarSource.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentGravatarSource)
	result[property.CommentGravatarDefault.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentGravatarDefault)
	result[property.CommentContentPlaceholder.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentContentPlaceholder)
	return dto.HumaOK(result)
}
