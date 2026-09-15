package api

import (
	"context"
	"html/template"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type JournalHandler struct {
	JournalService          service.JournalService
	JournalCommentService   service.JournalCommentService
	OptionService           service.ClientOptionService
	JournalCommentAssembler assembler.JournalCommentAssembler
}

func NewJournalHandler(
	journalService service.JournalService,
	journalCommentService service.JournalCommentService,
	optionService service.ClientOptionService,
	journalCommentAssembler assembler.JournalCommentAssembler,
) *JournalHandler {
	return &JournalHandler{
		JournalService:          journalService,
		JournalCommentService:   journalCommentService,
		OptionService:           optionService,
		JournalCommentAssembler: journalCommentAssembler,
	}
}

// ListJournalInput 日志列表查询输入。
type ListJournalInput struct {
	Page     int    `query:"page" doc:"页码"`
	PageSize int    `query:"size" doc:"每页数量"`
	Keyword  string `query:"keyword" doc:"关键词"`
}

// ListJournal 获取日志列表。
func (j *JournalHandler) ListJournal(ctx context.Context, in *ListJournalInput) (*dto.HumaOut[*dto.Page], error) {
	journalQuery := param.JournalQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.PageSize},
	}
	if in.Keyword != "" {
		journalQuery.Keyword = &in.Keyword
	}
	journalQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	journalQuery.JournalType = consts.JournalTypePublic.Ptr()
	journals, totalCount, err := j.JournalService.ListJournal(ctx, journalQuery)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, journals)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(journalDTOs, totalCount, journalQuery.Page))
}

// GetJournalInput 日志详情查询输入。
type GetJournalInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
}

// GetJournal 获取日志详情。
func (j *JournalHandler) GetJournal(ctx context.Context, in *GetJournalInput) (*dto.HumaOut[*dto.JournalWithComment], error) {
	journals, err := j.JournalService.GetByJournalIDs(ctx, []int32{in.JournalID})
	if err != nil {
		return dto.HumaErr[*dto.JournalWithComment](err)
	}
	if len(journals) == 0 {
		return dto.HumaErr[*dto.JournalWithComment](xerr.WithStatus(nil, xerr.StatusBadRequest))
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, []*entity.Journal{journals[in.JournalID]})
	if err != nil {
		return dto.HumaErr[*dto.JournalWithComment](err)
	}
	return dto.HumaOK(journalDTOs[0])
}

// ListTopCommentInput 获取顶级评论输入。
type ListTopCommentInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	Page      int   `query:"page" doc:"页码"`
}

// ListTopComment 获取日志的顶级评论分页列表。
func (j *JournalHandler) ListTopComment(ctx context.Context, in *ListTopCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page, PageSize: pageSize},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
	}
	commentQuery.ContentID = &in.JournalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, totalCount, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, comments)
	commenVOs, err := j.JournalCommentAssembler.ConvertToWithHasChildren(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commenVOs, totalCount, commentQuery.Page))
}

// ListChildrenInput 获取子评论输入。
type ListChildrenInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	ParentID  int32 `path:"parentID" doc:"父评论ID"`
}

// ListChildren 获取指定评论的子评论列表。
func (j *JournalHandler) ListChildren(ctx context.Context, in *ListChildrenInput) (*dto.HumaOut[[]*dto.Comment], error) {
	children, err := j.JournalCommentService.GetChildren(ctx, in.ParentID, in.JournalID, consts.CommentTypeJournal)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, children)
	comments, err := j.JournalCommentAssembler.ConvertToDTOList(ctx, children)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(comments)
}

// ListCommentTreeInput 获取评论树输入。
type ListCommentTreeInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	Page      int   `query:"page" doc:"页码"`
}

// ListCommentTree 获取日志的评论树形分页列表。
func (j *JournalHandler) ListCommentTree(ctx context.Context, in *ListCommentTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page, PageSize: pageSize},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
	}
	commentQuery.ContentID = &in.JournalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.ParentID = util.Int32Ptr(0)

	allComments, err := j.JournalCommentService.GetByContentID(ctx, in.JournalID, consts.CommentTypeJournal, commentQuery.Sort)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, allComments)
	commentVOs, total, err := j.JournalCommentAssembler.PageConvertToVOs(ctx, allComments, commentQuery.Page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, total, commentQuery.Page))
}

// ListCommentInput 获取评论列表输入。
type ListCommentInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
	Page      int   `query:"page" doc:"页码"`
}

// ListComment 获取日志的评论平铺分页列表。
func (j *JournalHandler) ListComment(ctx context.Context, in *ListCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page, PageSize: pageSize},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
	}
	commentQuery.ContentID = &in.JournalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, total, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, comments)
	result, err := j.JournalCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(result, total, commentQuery.Page))
}

// CreateCommentInput 创建评论输入。
type CreateCommentInput struct {
	Body param.Comment `doc:"评论内容"`
}

// CreateComment 创建日志评论。
func (j *JournalHandler) CreateComment(ctx context.Context, in *CreateCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	p := in.Body
	if p.AuthorURL != "" {
		err := util.Validate.Var(p.AuthorURL, "http_url")
		if err != nil {
			return dto.HumaErr[*dto.Comment](xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error"))
		}
	}
	p.Author = template.HTMLEscapeString(p.Author)
	p.AuthorURL = template.HTMLEscapeString(p.AuthorURL)
	p.Content = template.HTMLEscapeString(p.Content)
	p.Email = template.HTMLEscapeString(p.Email)
	p.CommentType = consts.CommentTypeJournal
	result, err := j.JournalCommentService.CreateBy(ctx, &p)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	comment, err := j.JournalCommentAssembler.ConvertToDTO(ctx, result)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(comment)
}

// LikeJournalInput 日志点赞输入。
type LikeJournalInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
}

// Like 日志点赞。
func (j *JournalHandler) Like(ctx context.Context, in *LikeJournalInput) (*dto.HumaOut[any], error) {
	err := j.JournalService.IncreaseLike(ctx, in.JournalID)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
