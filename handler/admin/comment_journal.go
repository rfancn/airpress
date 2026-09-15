package admin

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/service/impl"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type JournalCommentHandler struct {
	JournalCommentService   service.JournalCommentService
	OptionService           service.OptionService
	JournalService          service.JournalService
	JournalCommentAssembler assembler.JournalCommentAssembler
}

func NewJournalCommentHandler(journalCommentService service.JournalCommentService, optionService service.OptionService, journalService service.JournalService, journalCommentAssembler assembler.JournalCommentAssembler) *JournalCommentHandler {
	return &JournalCommentHandler{
		JournalCommentService:   journalCommentService,
		OptionService:           optionService,
		JournalService:          journalService,
		JournalCommentAssembler: journalCommentAssembler,
	}
}

// ListJournalCommentInput 日志评论列表查询输入。
type ListJournalCommentInput struct {
	Page     int      `query:"page" doc:"页码"`
	Size     int      `query:"size" doc:"每页数量"`
	Sort     []string `query:"sort" doc:"排序字段"`
	Keyword  string   `query:"keyword" doc:"关键词"`
	Status   string   `query:"status" doc:"评论状态"`
	ParentID int32    `query:"parentID" doc:"父评论ID"`
}

// ListJournalComment 获取日志评论分页列表。
func (j *JournalCommentHandler) ListJournalComment(ctx context.Context, in *ListJournalCommentInput) (*dto.HumaOut[*dto.Page], error) {
	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.Size},
	}
	// huma 不支持指针 query 参数，空串/0 视为未提供
	if in.Keyword != "" {
		commentQuery.Keyword = &in.Keyword
	}
	if in.Status != "" {
		status, err := consts.CommentStatusFromString(in.Status)
		if err != nil {
			return dto.HumaErr[*dto.Page](err)
		}
		commentQuery.CommentStatus = &status
	}
	if in.ParentID != 0 {
		commentQuery.ParentID = &in.ParentID
	}
	commentQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	comments, totalCount, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentDTOs, err := j.JournalCommentAssembler.ConvertToWithJournal(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentDTOs, totalCount, commentQuery.Page))
}

// ListJournalCommentLatestInput 最新日志评论查询输入。
type ListJournalCommentLatestInput struct {
	Top int32 `query:"top" doc:"获取数量"`
}

// ListJournalCommentLatest 获取最新日志评论列表。
func (j *JournalCommentHandler) ListJournalCommentLatest(ctx context.Context, in *ListJournalCommentLatestInput) (*dto.HumaOut[[]*vo.JournalCommentWithJournal], error) {
	if in.Top == 0 {
		return dto.HumaErr[[]*vo.JournalCommentWithJournal](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("top is required"))
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(in.Top)},
	}
	comments, _, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[[]*vo.JournalCommentWithJournal](err)
	}
	commentDTOs, err := j.JournalCommentAssembler.ConvertToWithJournal(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*vo.JournalCommentWithJournal](err)
	}
	return dto.HumaOK(commentDTOs)
}

// ListJournalCommentAsTreeInput 日志评论树形查询输入。
type ListJournalCommentAsTreeInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	Page      int   `query:"page" doc:"页码"`
}

// ListJournalCommentAsTree 获取日志评论树形分页列表。
func (j *JournalCommentHandler) ListJournalCommentAsTree(ctx context.Context, in *ListJournalCommentAsTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := j.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: in.Page}

	allComments, err := j.JournalCommentService.GetByContentID(ctx, in.JournalID, consts.CommentTypeJournal, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}

	commentVOs, totalCount, err := j.JournalCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, totalCount, page))
}

// ListJournalCommentWithParentInput 日志评论父级视图查询输入。
type ListJournalCommentWithParentInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	Page      int   `query:"page" doc:"页码"`
}

// ListJournalCommentWithParent 获取日志评论带父级的分页列表。
func (j *JournalCommentHandler) ListJournalCommentWithParent(ctx context.Context, in *ListJournalCommentWithParentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := j.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: in.Page}

	comments, totalCount, err := j.JournalCommentService.Page(ctx, param.CommentQuery{
		ContentID: &in.JournalID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}

	commentsWithParent, err := j.JournalCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentsWithParent, totalCount, page))
}

// CreateJournalCommentInput 创建日志评论输入。
type CreateJournalCommentInput struct {
	Body param.AdminComment `doc:"评论内容"`
}

// CreateJournalComment 创建日志评论。
func (j *JournalCommentHandler) CreateJournalComment(ctx context.Context, in *CreateJournalCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil || user == nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	blogURL, err := j.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	commentParam := in.Body
	commonParam := param.Comment{
		Author:            user.Username,
		Email:             user.Email,
		AuthorURL:         blogURL,
		Content:           commentParam.Content,
		PostID:            commentParam.PostID,
		ParentID:          commentParam.ParentID,
		AllowNotification: true,
		CommentType:       consts.CommentTypeJournal,
	}
	comment, err := j.JournalCommentService.CreateBy(ctx, &commonParam)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	result, err := j.JournalCommentAssembler.ConvertToDTO(ctx, comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// UpdateJournalCommentStatusInput 更新评论状态输入。
type UpdateJournalCommentStatusInput struct {
	CommentID int32  `path:"commentID" doc:"评论ID"`
	Status    string `path:"status" doc:"评论状态"`
}

// UpdateJournalCommentStatus 更新日志评论状态。
func (j *JournalCommentHandler) UpdateJournalCommentStatus(ctx context.Context, in *UpdateJournalCommentStatusInput) (*dto.HumaOut[*entity.Comment], error) {
	status, err := consts.CommentStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	comment, err := j.JournalCommentService.UpdateStatus(ctx, in.CommentID, status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	return dto.HumaOK(comment)
}

// UpdateJournalCommentInput 更新日志评论输入。
type UpdateJournalCommentInput struct {
	CommentID int32         `path:"commentID" doc:"评论ID"`
	Body      param.Comment `doc:"评论内容"`
}

// UpdateJournalComment 更新日志评论。
func (j *JournalCommentHandler) UpdateJournalComment(ctx context.Context, in *UpdateJournalCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	commentParam := in.Body
	if commentParam.AuthorURL != "" {
		err := util.Validate.Var(commentParam.AuthorURL, "url")
		if err != nil {
			return dto.HumaErr[*dto.Comment](xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("url is not available"))
		}
	}
	comment, err := j.JournalCommentService.UpdateBy(ctx, in.CommentID, &commentParam)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	result, err := j.JournalCommentAssembler.ConvertToDTO(ctx, comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// UpdateJournalStatusBatchInput 批量更新评论状态输入。
type UpdateJournalStatusBatchInput struct {
	Status int32   `path:"status" doc:"评论状态"`
	Body   []int32 `doc:"评论ID列表"`
}

// UpdateJournalStatusBatch 批量更新日志评论状态。
func (j *JournalCommentHandler) UpdateJournalStatusBatch(ctx context.Context, in *UpdateJournalStatusBatchInput) (*dto.HumaOut[[]*dto.Comment], error) {
	comments, err := j.JournalCommentService.UpdateStatusBatch(ctx, in.Body, consts.CommentStatus(in.Status))
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	commentDTOs, err := j.JournalCommentAssembler.ConvertToDTOList(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(commentDTOs)
}

// DeleteJournalCommentInput 删除日志评论输入。
type DeleteJournalCommentInput struct {
	CommentID int32 `path:"commentID" doc:"评论ID"`
}

// DeleteJournalComment 删除日志评论。
func (j *JournalCommentHandler) DeleteJournalComment(ctx context.Context, in *DeleteJournalCommentInput) (*dto.HumaOut[any], error) {
	err := j.JournalCommentService.Delete(ctx, in.CommentID)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// DeleteJournalCommentBatchInput 批量删除日志评论输入。
type DeleteJournalCommentBatchInput struct {
	Body []int32 `doc:"评论ID列表"`
}

// DeleteJournalCommentBatch 批量删除日志评论。
func (j *JournalCommentHandler) DeleteJournalCommentBatch(ctx context.Context, in *DeleteJournalCommentBatchInput) (*dto.HumaOut[any], error) {
	err := j.JournalCommentService.DeleteBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
