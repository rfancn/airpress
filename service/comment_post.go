package service

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
)

type PostCommentService interface {
	BaseCommentService
	CreateBy(ctx context.Context, commentParam *param.Comment) (*entity.Comment, error)
	CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error)
	UpdateBy(ctx context.Context, commentID int32, commentParam *param.Comment) (*entity.Comment, error)
}
