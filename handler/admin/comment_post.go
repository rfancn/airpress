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

type PostCommentHandler struct {
	PostCommentService   service.PostCommentService
	OptionService        service.OptionService
	PostService          service.PostService
	PostAssembler        assembler.PostAssembler
	PostCommentAssembler assembler.PostCommentAssembler
}

func NewPostCommentHandler(
	postCommentHandler service.PostCommentService,
	optionService service.OptionService,
	postService service.PostService,
	postAssembler assembler.PostAssembler,
	postCommentAssembler assembler.PostCommentAssembler,
) *PostCommentHandler {
	return &PostCommentHandler{
		PostCommentService:   postCommentHandler,
		OptionService:        optionService,
		PostService:          postService,
		PostAssembler:        postAssembler,
		PostCommentAssembler: postCommentAssembler,
	}
}

// ListPostCommentInput 文章评论列表查询输入。
type ListPostCommentInput struct {
	Page      int      `query:"page" doc:"页码"`
	Size      int      `query:"size" doc:"每页数量"`
	Sort      []string `query:"sort" doc:"排序字段"`
	ContentID int32    `query:"contentId" doc:"内容ID"`
	Keyword   string   `query:"keyword" doc:"关键字"`
	ParentID  int32    `query:"parentID" doc:"父评论ID"`
}

// ListPostComment 获取文章评论列表（分页）。
func (p *PostCommentHandler) ListPostComment(ctx context.Context, in *ListPostCommentInput) (*dto.HumaOut[*dto.Page], error) {
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
	comments, totalCount, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentDTOs, err := p.PostCommentAssembler.ConvertToWithPost(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentDTOs, totalCount, commentQuery.Page))
}

// ListPostCommentLatestInput 最新文章评论查询输入。
type ListPostCommentLatestInput struct {
	Top int32 `query:"top" doc:"返回数量"`
}

// ListPostCommentLatest 获取最新文章评论列表。
func (p *PostCommentHandler) ListPostCommentLatest(ctx context.Context, in *ListPostCommentLatestInput) (*dto.HumaOut[[]*vo.PostCommentWithPost], error) {
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(in.Top)},
	}
	comments, _, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[[]*vo.PostCommentWithPost](err)
	}
	result, err := p.PostCommentAssembler.ConvertToWithPost(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*vo.PostCommentWithPost](err)
	}
	return dto.HumaOK(result)
}

// ListPostCommentAsTreeInput 文章评论树形列表查询输入。
type ListPostCommentAsTreeInput struct {
	PostID int32 `path:"postID" doc:"文章ID"`
	Page   int   `query:"page" doc:"页码"`
}

// ListPostCommentAsTree 获取文章评论树形结构（分页）。
func (p *PostCommentHandler) ListPostCommentAsTree(ctx context.Context, in *ListPostCommentAsTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: in.Page}
	allComments, err := p.PostCommentService.GetByContentID(ctx, in.PostID, consts.CommentTypePost, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentVOs, totalCount, err := p.PostCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, totalCount, page))
}

// ListPostCommentWithParentInput 文章评论列表（带父评论信息）查询输入。
type ListPostCommentWithParentInput struct {
	PostID int32 `path:"postID" doc:"文章ID"`
	Page   int   `query:"page" doc:"页码"`
}

// ListPostCommentWithParent 获取文章评论列表（带 parentVO 信息，分页）。
func (p *PostCommentHandler) ListPostCommentWithParent(ctx context.Context, in *ListPostCommentWithParentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	page := param.Page{PageNum: in.Page, PageSize: pageSize.(int)}
	comments, totalCount, err := p.PostCommentService.Page(ctx, param.CommentQuery{
		ContentID: &in.PostID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	commentsWithParent, err := p.PostCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentsWithParent, totalCount, page))
}

// CreatePostCommentInput 创建文章评论输入。
type CreatePostCommentInput struct {
	Body param.AdminComment `doc:"评论参数"`
}

// CreatePostComment 创建文章评论。
func (p *PostCommentHandler) CreatePostComment(ctx context.Context, in *CreatePostCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	commentParam := in.Body
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	blogURL, err := p.OptionService.GetBlogBaseURL(ctx)
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
		CommentType:       consts.CommentTypePost,
	}
	comment, err := p.PostCommentService.CreateBy(ctx, &commonParam)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	dtoComment, err := p.PostCommentAssembler.ConvertToDTO(ctx, comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(dtoComment)
}

// UpdatePostCommentInput 更新文章评论输入。
type UpdatePostCommentInput struct {
	CommentID int32         `path:"commentID" doc:"评论ID"`
	Body      param.Comment `doc:"评论参数"`
}

// UpdatePostComment 更新文章评论。
func (p *PostCommentHandler) UpdatePostComment(ctx context.Context, in *UpdatePostCommentInput) (*dto.HumaOut[*dto.Comment], error) {
	commentParam := in.Body
	if commentParam.AuthorURL != "" {
		err := util.Validate.Var(commentParam.AuthorURL, "url")
		if err != nil {
			return dto.HumaErr[*dto.Comment](xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("url is not available"))
		}
	}
	comment, err := p.PostCommentService.UpdateBy(ctx, in.CommentID, &commentParam)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	dtoComment, err := p.PostCommentAssembler.ConvertToDTO(ctx, comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(dtoComment)
}

// UpdatePostCommentStatusInput 更新文章评论状态输入。
type UpdatePostCommentStatusInput struct {
	CommentID int32  `path:"commentID" doc:"评论ID"`
	Status    string `path:"status" doc:"评论状态（PUBLISHED/AUDITING/RECYCLE）"`
}

// UpdatePostCommentStatus 更新文章评论状态。
func (p *PostCommentHandler) UpdatePostCommentStatus(ctx context.Context, in *UpdatePostCommentStatusInput) (*dto.HumaOut[*entity.Comment], error) {
	status, err := consts.CommentStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	comment, err := p.PostCommentService.UpdateStatus(ctx, in.CommentID, status)
	if err != nil {
		return dto.HumaErr[*entity.Comment](err)
	}
	return dto.HumaOK(comment)
}

// UpdatePostCommentStatusBatchInput 批量更新文章评论状态输入。
type UpdatePostCommentStatusBatchInput struct {
	Status string  `path:"status" doc:"评论状态（PUBLISHED/AUDITING/RECYCLE）"`
	Body   []int32 `doc:"评论ID列表"`
}

// UpdatePostCommentStatusBatch 批量更新文章评论状态。
func (p *PostCommentHandler) UpdatePostCommentStatusBatch(ctx context.Context, in *UpdatePostCommentStatusBatchInput) (*dto.HumaOut[[]*dto.Comment], error) {
	status, err := consts.CommentStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	comments, err := p.PostCommentService.UpdateStatusBatch(ctx, in.Body, status)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	result, err := p.PostCommentAssembler.ConvertToDTOList(ctx, comments)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// DeletePostCommentInput 删除文章评论输入。
type DeletePostCommentInput struct {
	CommentID int32 `path:"commentID" doc:"评论ID"`
}

// DeletePostComment 删除文章评论。
func (p *PostCommentHandler) DeletePostComment(ctx context.Context, in *DeletePostCommentInput) (*dto.HumaOut[any], error) {
	if err := p.PostCommentService.Delete(ctx, in.CommentID); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// DeletePostCommentBatchInput 批量删除文章评论输入。
type DeletePostCommentBatchInput struct {
	Body []int32 `doc:"评论ID列表"`
}

// DeletePostCommentBatch 批量删除文章评论。
func (p *PostCommentHandler) DeletePostCommentBatch(ctx context.Context, in *DeletePostCommentBatchInput) (*dto.HumaOut[any], error) {
	if err := p.PostCommentService.DeleteBatch(ctx, in.Body); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
