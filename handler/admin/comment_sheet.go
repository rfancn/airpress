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
)

type SheetCommentHandler struct {
	SheetCommentService   service.SheetCommentService
	BaseCommentService    service.BaseCommentService
	OptionService         service.OptionService
	SheetService          service.SheetService
	SheetAssembler        assembler.SheetAssembler
	SheetCommentAssembler assembler.SheetCommentAssembler
}

func NewSheetCommentHandler(
	sheetCommentService service.SheetCommentService,
	baseCommentService service.BaseCommentService,
	optionService service.OptionService,
	sheetService service.SheetService,
	sheetAssembler assembler.SheetAssembler,
	sheetCommentAssembler assembler.SheetCommentAssembler,
) *SheetCommentHandler {
	return &SheetCommentHandler{
		SheetCommentService:   sheetCommentService,
		BaseCommentService:    baseCommentService,
		OptionService:         optionService,
		SheetService:          sheetService,
		SheetAssembler:        sheetAssembler,
		SheetCommentAssembler: sheetCommentAssembler,
	}
}

// ListSheetCommentInput 页面评论列表查询输入。
type ListSheetCommentInput struct {
	Page      int      `query:"page" doc:"页码"`
	Size      int      `query:"size" doc:"每页数量"`
	Sort      []string `query:"sort" doc:"排序字段"`
	ContentID int32    `query:"contentId" doc:"内容ID"`
	Keyword   string   `query:"keyword" doc:"关键字"`
	ParentID  int32    `query:"parentID" doc:"父评论ID"`
}

// ListSheetComment 获取页面评论列表（分页）。
func (s *SheetCommentHandler) ListSheetComment(ctx context.Context, in *ListSheetCommentInput) (*dto.HumaOut[*dto.Page], error) {
	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.Size},
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
	}
	if in.ContentID != 0 {
		commentQuery.ContentID = &in.ContentID
	}
	if in.Keyword != "" {
		commentQuery.Keyword = &in.Keyword
	}
	if in.ParentID != 0 {
		commentQuery.ParentID = &in.ParentID
	}
	comments, totalCount, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentDTOs, err := s.ConvertToWithSheet(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentDTOs, totalCount, commentQuery.Page))
}

// ListSheetCommentLatestInput 最新页面评论查询输入。
type ListSheetCommentLatestInput struct {
	Top int32 `query:"top" doc:"返回数量"`
}

// ListSheetCommentLatest 获取最新页面评论列表。
func (s *SheetCommentHandler) ListSheetCommentLatest(ctx context.Context, in *ListSheetCommentLatestInput) (*dto.HumaOut[[]*vo.SheetCommentWithSheet], error) {
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(in.Top)},
	}
	comments, _, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return dto.HumaErr[[]*vo.SheetCommentWithSheet](err)
	}
	result, err := s.ConvertToWithSheet(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*vo.SheetCommentWithSheet](err)
	}
	return dto.HumaOK(result)
}

// ListSheetCommentAsTreeInput 页面评论树形列表查询输入。
type ListSheetCommentAsTreeInput struct {
	SheetID int32 `path:"sheetID" doc:"页面ID"`
	Page    int   `query:"page" doc:"页码"`
}

// ListSheetCommentAsTree 获取页面评论树形结构（分页）。
func (s *SheetCommentHandler) ListSheetCommentAsTree(ctx context.Context, in *ListSheetCommentAsTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := s.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: in.Page}
	allComments, err := s.SheetCommentService.GetByContentID(ctx, in.SheetID, consts.CommentTypeSheet, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentVOs, totalCount, err := s.SheetCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, totalCount, page))
}

// ListSheetCommentWithParentInput 页面评论列表（带父评论信息）查询输入。
type ListSheetCommentWithParentInput struct {
	SheetID int32 `path:"sheetID" doc:"页面ID"`
	Page    int   `query:"page" doc:"页码"`
}

// ListSheetCommentWithParent 获取页面评论列表（带 parentVO 信息，分页）。
func (s *SheetCommentHandler) ListSheetCommentWithParent(ctx context.Context, in *ListSheetCommentWithParentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := s.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: in.Page}
	comments, totalCount, err := s.SheetCommentService.Page(ctx, param.CommentQuery{
		ContentID: &in.SheetID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentsWithParent, err := s.SheetCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentsWithParent, totalCount, page))
}

// CreateSheetCommentInput 创建页面评论输入。
type CreateSheetCommentInput struct {
	Body param.AdminComment `doc:"评论参数"`
}

// CreateSheetComment 创建页面评论。
func (s *SheetCommentHandler) CreateSheetComment(ctx context.Context, in *CreateSheetCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	commentParam := in.Body
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	blogURL, err := s.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	commonParam := param.Comment{
		Author:            user.Username,
		Email:             user.Email,
		AuthorURL:         blogURL,
		Content:           commentParam.Content,
		PostID:            commentParam.PostID,
		ParentID:          commentParam.ParentID,
		AllowNotification: true,
		CommentType:       consts.CommentTypeSheet,
	}
	comment, err := s.BaseCommentService.CreateBy(ctx, &commonParam)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	dtoComment, err := s.SheetCommentAssembler.ConvertToDTO(ctx, comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(dtoComment)
}

// UpdateSheetCommentStatusInput 更新页面评论状态输入。
type UpdateSheetCommentStatusInput struct {
	CommentID int32  `path:"commentID" doc:"评论ID"`
	Status    string `path:"status" doc:"评论状态（PUBLISHED/AUDITING/RECYCLE）"`
}

// UpdateSheetCommentStatus 更新页面评论状态。
func (s *SheetCommentHandler) UpdateSheetCommentStatus(ctx context.Context, in *UpdateSheetCommentStatusInput) (*dto.HumaOut[*entity.Comment], error) {
	status, err := consts.CommentStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	comment, err := s.SheetCommentService.UpdateStatus(ctx, in.CommentID, status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	return dto.HumaOK(comment)
}

// UpdateSheetCommentStatusBatchInput 批量更新页面评论状态输入。
type UpdateSheetCommentStatusBatchInput struct {
	Status int32   `path:"status" doc:"评论状态（数字：0=PUBLISHED,1=AUDITING,2=RECYCLE）"`
	Body   []int32 `doc:"评论ID列表"`
}

// UpdateSheetCommentStatusBatch 批量更新页面评论状态。
func (s *SheetCommentHandler) UpdateSheetCommentStatusBatch(ctx context.Context, in *UpdateSheetCommentStatusBatchInput) (*dto.HumaOut[[]*dto.Comment], error) {
	comments, err := s.SheetCommentService.UpdateStatusBatch(ctx, in.Body, consts.CommentStatus(in.Status))
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	result, err := s.SheetCommentAssembler.ConvertToDTOList(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// DeleteSheetCommentInput 删除页面评论输入。
type DeleteSheetCommentInput struct {
	CommentID int32 `path:"commentID" doc:"评论ID"`
}

// DeleteSheetComment 删除页面评论。
func (s *SheetCommentHandler) DeleteSheetComment(ctx context.Context, in *DeleteSheetCommentInput) (*dto.HumaOut[any], error) {
	if err := s.SheetCommentService.Delete(ctx, in.CommentID); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// DeleteSheetCommentBatchInput 批量删除页面评论输入。
type DeleteSheetCommentBatchInput struct {
	Body []int32 `doc:"评论ID列表"`
}

// DeleteSheetCommentBatch 批量删除页面评论。
func (s *SheetCommentHandler) DeleteSheetCommentBatch(ctx context.Context, in *DeleteSheetCommentBatchInput) (*dto.HumaOut[any], error) {
	if err := s.SheetCommentService.DeleteBatch(ctx, in.Body); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// ConvertToWithSheet 将评论转换为带页面信息的 VO。
func (s *SheetCommentHandler) ConvertToWithSheet(ctx context.Context, comments []*entity.Comment) ([]*vo.SheetCommentWithSheet, error) {
	postIDs := make([]int32, 0, len(comments))
	for _, comment := range comments {
		postIDs = append(postIDs, comment.PostID)
	}
	posts, err := s.SheetService.GetByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.SheetCommentWithSheet, 0, len(comments))
	for _, comment := range comments {
		commentDTO, err := s.SheetCommentAssembler.ConvertToDTO(ctx, comment)
		if err != nil {
			return nil, err
		}
		commentWithSheet := &vo.SheetCommentWithSheet{
			Comment: *commentDTO,
		}
		result = append(result, commentWithSheet)
		post, ok := posts[comment.PostID]
		if ok {
			commentWithSheet.PostMinimal, err = s.SheetAssembler.ConvertToMinimalDTO(ctx, post)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
