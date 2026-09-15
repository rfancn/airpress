package api

import (
	"context"
	"html/template"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type SheetHandler struct {
	OptionService         service.OptionService
	SheetService          service.SheetService
	SheetCommentService   service.SheetCommentService
	SheetCommentAssembler assembler.SheetCommentAssembler
}

func NewSheetHandler(
	optionService service.OptionService,
	sheetService service.SheetService,
	sheetCommentService service.SheetCommentService,
	sheetCommentAssembler assembler.SheetCommentAssembler,
) *SheetHandler {
	return &SheetHandler{
		OptionService:         optionService,
		SheetService:          sheetService,
		SheetCommentService:   sheetCommentService,
		SheetCommentAssembler: sheetCommentAssembler,
	}
}

// SheetListTopCommentInput 获取页面顶级评论列表。
type SheetListTopCommentInput struct {
	SheetID int32    `path:"sheetID" doc:"页面ID"`
	Page    int      `query:"page" doc:"页码"`
	Sort    []string `query:"sort" doc:"排序字段"`
}

// ListTopComment 获取页面顶级评论列表（带 hasChildren 标记）。
func (s *SheetHandler) ListTopComment(ctx context.Context, in *SheetListTopCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := s.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.SheetID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, totalCount, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = s.SheetCommentAssembler.ClearSensitiveField(ctx, comments)
	commenVOs, err := s.SheetCommentAssembler.ConvertToWithHasChildren(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commenVOs, totalCount, commentQuery.Page))
}

// SheetListChildrenInput 获取页面评论子列表。
type SheetListChildrenInput struct {
	SheetID  int32 `path:"sheetID" doc:"页面ID"`
	ParentID int32 `path:"parentID" doc:"父评论ID"`
}

// ListChildren 获取页面评论的子评论列表。
func (s *SheetHandler) ListChildren(ctx context.Context, in *SheetListChildrenInput) (*dto.HumaOut[[]*dto.Comment], error) {
	children, err := s.SheetCommentService.GetChildren(ctx, in.ParentID, in.SheetID, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	_ = s.SheetCommentAssembler.ClearSensitiveField(ctx, children)
	result, err := s.SheetCommentAssembler.ConvertToDTOList(ctx, children)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// SheetListCommentTreeInput 获取页面评论树形列表。
type SheetListCommentTreeInput struct {
	SheetID int32    `path:"sheetID" doc:"页面ID"`
	Page    int      `query:"page" doc:"页码"`
	Sort    []string `query:"sort" doc:"排序字段"`
}

// ListCommentTree 获取页面评论的树形结构（分页）。
func (s *SheetHandler) ListCommentTree(ctx context.Context, in *SheetListCommentTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := s.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.SheetID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	allComments, err := s.SheetCommentService.GetByContentID(ctx, in.SheetID, consts.CommentTypeSheet, commentQuery.Sort)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = s.SheetCommentAssembler.ClearSensitiveField(ctx, allComments)
	commentVOs, total, err := s.SheetCommentAssembler.PageConvertToVOs(ctx, allComments, commentQuery.Page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, total, commentQuery.Page))
}

// SheetListCommentInput 获取页面评论列表。
type SheetListCommentInput struct {
	SheetID int32    `path:"sheetID" doc:"页面ID"`
	Page    int      `query:"page" doc:"页码"`
	Sort    []string `query:"sort" doc:"排序字段"`
}

// ListComment 获取页面评论列表（带 parentVO 信息，分页）。
func (s *SheetHandler) ListComment(ctx context.Context, in *SheetListCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := s.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.SheetID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, total, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = s.SheetCommentAssembler.ClearSensitiveField(ctx, comments)
	result, err := s.SheetCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(result, total, commentQuery.Page))
}

// SheetCreateCommentInput 创建页面评论。
type SheetCreateCommentInput struct {
	Body param.Comment `doc:"评论内容"`
}

// CreateComment 创建页面评论。
func (s *SheetHandler) CreateComment(ctx context.Context, in *SheetCreateCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	comment := in.Body
	if comment.AuthorURL != "" {
		err := util.Validate.Var(comment.AuthorURL, "http_url")
		if err != nil {
			return dto.HumaErr[*dto.Comment](xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error"))
		}
	}
	comment.Author = template.HTMLEscapeString(comment.Author)
	comment.AuthorURL = template.HTMLEscapeString(comment.AuthorURL)
	comment.Content = template.HTMLEscapeString(comment.Content)
	comment.Email = template.HTMLEscapeString(comment.Email)
	comment.CommentType = consts.CommentTypeSheet
	result, err := s.SheetCommentService.CreateBy(ctx, &comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	dtoComment, err := s.SheetCommentAssembler.ConvertToDTO(ctx, result)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(dtoComment)
}
