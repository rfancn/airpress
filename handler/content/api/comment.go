package api

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/service"
)

type CommentHandler struct {
	BaseCommentService service.BaseCommentService
}

func NewCommentHandler(baseCommentService service.BaseCommentService) *CommentHandler {
	return &CommentHandler{
		BaseCommentService: baseCommentService,
	}
}

// LikeInput 评论点赞输入。
type LikeInput struct {
	CommentID int32 `path:"commentID" doc:"评论ID"`
}

// Like 给评论点赞。
func (c *CommentHandler) Like(ctx context.Context, in *LikeInput) (*dto.HumaOut[interface{}], error) {
	err := c.BaseCommentService.IncreaseLike(ctx, in.CommentID)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}
