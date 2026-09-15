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

type PostHandler struct {
	OptionService        service.OptionService
	PostService          service.PostService
	PostCommentService   service.PostCommentService
	PostCommentAssembler assembler.PostCommentAssembler
}

func NewPostHandler(
	optionService service.OptionService,
	postService service.PostService,
	postCommentService service.PostCommentService,
	postCommentAssembler assembler.PostCommentAssembler,
) *PostHandler {
	return &PostHandler{
		OptionService:        optionService,
		PostService:          postService,
		PostCommentService:   postCommentService,
		PostCommentAssembler: postCommentAssembler,
	}
}

// PostListTopCommentInput 获取文章顶级评论列表。
type PostListTopCommentInput struct {
	PostID int32    `path:"postID" doc:"文章ID"`
	Page   int      `query:"page" doc:"页码"`
	Sort   []string `query:"sort" doc:"排序字段"`
}

// ListTopComment 获取文章顶级评论列表（带 hasChildren 标记）。
func (p *PostHandler) ListTopComment(ctx context.Context, in *PostListTopCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := p.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.PostID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, totalCount, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(ctx, comments)
	commenVOs, err := p.PostCommentAssembler.ConvertToWithHasChildren(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commenVOs, totalCount, commentQuery.Page))
}

// PostListChildrenInput 获取文章评论子列表。
type PostListChildrenInput struct {
	PostID   int32 `path:"postID" doc:"文章ID"`
	ParentID int32 `path:"parentID" doc:"父评论ID"`
}

// ListChildren 获取文章评论的子评论列表。
func (p *PostHandler) ListChildren(ctx context.Context, in *PostListChildrenInput) (*dto.HumaOut[[]*dto.Comment], error) {
	children, err := p.PostCommentService.GetChildren(ctx, in.ParentID, in.PostID, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(ctx, children)
	result, err := p.PostCommentAssembler.ConvertToDTOList(ctx, children)
	if err != nil {
		return dto.HumaErr[[]*dto.Comment](err)
	}
	return dto.HumaOK(result)
}

// PostListCommentTreeInput 获取文章评论树形列表。
type PostListCommentTreeInput struct {
	PostID int32    `path:"postID" doc:"文章ID"`
	Page   int      `query:"page" doc:"页码"`
	Sort   []string `query:"sort" doc:"排序字段"`
}

// ListCommentTree 获取文章评论的树形结构（分页）。
func (p *PostHandler) ListCommentTree(ctx context.Context, in *PostListCommentTreeInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := p.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.PostID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	allComments, err := p.PostCommentService.GetByContentID(ctx, in.PostID, consts.CommentTypePost, commentQuery.Sort)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(ctx, allComments)
	commentVOs, total, err := p.PostCommentAssembler.PageConvertToVOs(ctx, allComments, commentQuery.Page)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(commentVOs, total, commentQuery.Page))
}

// PostListCommentInput 获取文章评论列表。
type PostListCommentInput struct {
	PostID int32    `path:"postID" doc:"文章ID"`
	Page   int      `query:"page" doc:"页码"`
	Sort   []string `query:"sort" doc:"排序字段"`
}

// ListComment 获取文章评论列表（带 parentVO 信息，分页）。
func (p *PostHandler) ListComment(ctx context.Context, in *PostListCommentInput) (*dto.HumaOut[*dto.Page], error) {
	pageSize := p.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{
		Page: param.Page{PageNum: in.Page},
	}
	if len(in.Sort) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &in.PostID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, total, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(ctx, comments)
	result, err := p.PostCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(result, total, commentQuery.Page))
}

// PostCreateCommentInput 创建文章评论。
type PostCreateCommentInput struct {
	Body param.Comment `doc:"评论内容"`
}

// CreateComment 创建文章评论。
func (p *PostHandler) CreateComment(ctx context.Context, in *PostCreateCommentInput) (*dto.HumaOut[*dto.Comment], error) {
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
	comment.CommentType = consts.CommentTypePost
	result, err := p.PostCommentService.CreateBy(ctx, &comment)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	dtoComment, err := p.PostCommentAssembler.ConvertToDTO(ctx, result)
	if err != nil {
		return dto.HumaErr[*dto.Comment](err)
	}
	return dto.HumaOK(dtoComment)
}

// PostLikeInput 点赞文章。
type PostLikeInput struct {
	PostID int32 `path:"postID" doc:"文章ID"`
}

// Like 点赞文章。
func (p *PostHandler) Like(ctx context.Context, in *PostLikeInput) (*dto.HumaOut[any], error) {
	err := p.PostService.IncreaseLike(ctx, in.PostID)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
