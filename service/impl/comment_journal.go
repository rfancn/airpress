package impl

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/dal"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type journalCommentServiceImpl struct {
	service.BaseCommentService
}

func NewJournalCommentService(baseCommentService service.BaseCommentService) service.JournalCommentService {
	return &journalCommentServiceImpl{
		BaseCommentService: baseCommentService,
	}
}

func (j *journalCommentServiceImpl) CountByStatusAndJournalID(ctx context.Context, status consts.CommentStatus, journalIDs []int32) (map[int32]int64, error) {
	return j.CountByStatusAndContentIDs(ctx, status, journalIDs)
}

func (j *journalCommentServiceImpl) UpdateBy(ctx context.Context, commentID int32, commentParam *param.Comment) (*entity.Comment, error) {
	if commentID == 0 {
		return nil, nil
	}
	comment := j.BaseCommentService.ConvertParam(commentParam)
	comment.ID = commentID
	return j.BaseCommentService.Update(ctx, comment)
}

func (j *journalCommentServiceImpl) CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	count, err := commentDAL.WithContext(ctx).Where(commentDAL.Type.Eq(consts.CommentTypeJournal), commentDAL.Status.Eq(status)).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}
