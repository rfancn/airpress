package admin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type PostHandler struct {
	PostService   service.PostService
	PostAssembler assembler.PostAssembler
}

func NewPostHandler(postService service.PostService, postAssembler assembler.PostAssembler) *PostHandler {
	return &PostHandler{
		PostService:   postService,
		PostAssembler: postAssembler,
	}
}

// ListPostsInput 文章列表查询输入。
type ListPostsInput struct {
	Page       int      `query:"page" doc:"页码"`
	Size       int      `query:"size" doc:"每页数量"`
	Sort       []string `query:"sort" doc:"排序字段"`
	Keyword    string   `query:"keyword" doc:"关键字"`
	CategoryID int32    `query:"categoryId" doc:"分类ID"`
	// 原 gin 逻辑为 More==nil 时视为 true，改用 default:"true" 保持语义不变
	More  bool  `query:"more" default:"true" doc:"是否返回详情列表"`
	TagID int32 `query:"tagId" doc:"标签ID"`
}

// ListPosts 获取文章列表（分页）。
func (p *PostHandler) ListPosts(ctx context.Context, in *ListPostsInput) (*dto.HumaOut[*dto.Page], error) {
	postQuery := param.PostQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.Size},
	}
	if len(in.Sort) > 0 {
		postQuery.Sort = &param.Sort{Fields: in.Sort}
	} else {
		postQuery.Sort = &param.Sort{Fields: []string{"topPriority,desc", "createTime,desc"}}
	}
	// huma 不支持指针 query 参数，空串/0 视为未提供
	if in.Keyword != "" {
		postQuery.Keyword = &in.Keyword
	}
	if in.CategoryID != 0 {
		postQuery.CategoryID = &in.CategoryID
	}
	if in.TagID != 0 {
		postQuery.TagID = &in.TagID
	}

	posts, totalCount, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	if in.More {
		postVOs, err := p.PostAssembler.ConvertToListVO(ctx, posts)
		if err != nil {
			return dto.HumaErr[*dto.Page](err)
		}
		return dto.HumaOK(dto.NewPage(postVOs, totalCount, postQuery.Page))
	}
	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(ctx, post)
		if err != nil {
			return dto.HumaErr[*dto.Page](err)
		}
		postDTOs = append(postDTOs, postDTO)
	}
	return dto.HumaOK(dto.NewPage(postDTOs, totalCount, postQuery.Page))
}

// ListLatestPostsInput 最新文章查询输入。
type ListLatestPostsInput struct {
	Top int32 `query:"top" doc:"返回数量"`
}

// ListLatestPosts 获取最新文章列表。
func (p *PostHandler) ListLatestPosts(ctx context.Context, in *ListLatestPostsInput) (*dto.HumaOut[[]*dto.PostMinimal], error) {
	top := in.Top
	if top == 0 {
		top = 10
	}
	postQuery := param.PostQuery{
		Page: param.Page{
			PageSize: int(top),
			PageNum:  0,
		},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
		Keyword:    nil,
		CategoryID: nil,
		More:       util.BoolPtr(false),
	}
	posts, _, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return dto.HumaErr[[]*dto.PostMinimal](err)
	}
	postMinimals := make([]*dto.PostMinimal, 0, len(posts))
	for _, post := range posts {
		postMinimal, err := p.PostAssembler.ConvertToMinimalDTO(ctx, post)
		if err != nil {
			return dto.HumaErr[[]*dto.PostMinimal](err)
		}
		postMinimals = append(postMinimals, postMinimal)
	}
	return dto.HumaOK(postMinimals)
}

// ListPostsByStatusInput 按状态查询文章列表输入。
type ListPostsByStatusInput struct {
	Status     int32    `path:"status" doc:"文章状态"`
	Page       int      `query:"page" doc:"页码"`
	Size       int      `query:"size" doc:"每页数量"`
	Sort       []string `query:"sort" doc:"排序字段"`
	Keyword    string   `query:"keyword" doc:"关键字"`
	CategoryID int32    `query:"categoryId" doc:"分类ID"`
	More       bool     `query:"more" doc:"是否返回详情列表"`
	TagID      int32    `query:"tagId" doc:"标签ID"`
}

// ListPostsByStatus 按状态获取文章列表（分页）。
func (p *PostHandler) ListPostsByStatus(ctx context.Context, in *ListPostsByStatusInput) (*dto.HumaOut[*dto.Page], error) {
	postQuery := param.PostQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.Size},
	}
	if len(in.Sort) > 0 {
		postQuery.Sort = &param.Sort{Fields: in.Sort}
	} else {
		postQuery.Sort = &param.Sort{Fields: []string{"createTime,desc"}}
	}
	if in.Keyword != "" {
		postQuery.Keyword = &in.Keyword
	}
	if in.CategoryID != 0 {
		postQuery.CategoryID = &in.CategoryID
	}
	if in.TagID != 0 {
		postQuery.TagID = &in.TagID
	}

	statusType := consts.PostStatus(in.Status)
	postQuery.Statuses = []*consts.PostStatus{&statusType}

	posts, totalCount, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	if in.More {
		postVOs, err := p.PostAssembler.ConvertToListVO(ctx, posts)
		if err != nil {
			return dto.HumaErr[*dto.Page](err)
		}
		return dto.HumaOK(dto.NewPage(postVOs, totalCount, postQuery.Page))
	}
	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(ctx, post)
		if err != nil {
			return dto.HumaErr[*dto.Page](err)
		}
		postDTOs = append(postDTOs, postDTO)
	}
	return dto.HumaOK(dto.NewPage(postDTOs, totalCount, postQuery.Page))
}

// GetByPostIDInput 文章详情查询输入。
type GetByPostIDInput struct {
	PostID int32 `path:"postID" doc:"文章ID"`
}

// GetByPostID 获取文章详情（huma 风格，挂在需要鉴权的 admin group 下）。
func (p *PostHandler) GetByPostID(ctx context.Context, in *GetByPostIDInput) (*dto.HumaOut[*vo.PostDetailVO], error) {
	post, err := p.PostService.GetByPostID(ctx, in.PostID)
	if err != nil {
		return dto.HumaErr[*vo.PostDetailVO](err)
	}
	postDetailVO, err := p.PostAssembler.ConvertToDetailVO(ctx, post)
	if err != nil {
		return dto.HumaErr[*vo.PostDetailVO](err)
	}
	return dto.HumaOK(postDetailVO)
}

// CreatePostInput 创建文章输入。
type CreatePostInput struct {
	Body param.Post `doc:"文章参数"`
}

// CreatePost 创建文章。
func (p *PostHandler) CreatePost(ctx context.Context, in *CreatePostInput) (*dto.HumaOut[*vo.PostDetailVO], error) {
	postParam := in.Body
	post, err := p.PostService.Create(ctx, &postParam)
	if err != nil {
		return dto.HumaErr[*vo.PostDetailVO](err)
	}
	postDetailVO, err := p.PostAssembler.ConvertToDetailVO(ctx, post)
	if err != nil {
		return dto.HumaErr[*vo.PostDetailVO](err)
	}
	return dto.HumaOK(postDetailVO)
}

// UpdatePostInput 更新文章输入。
type UpdatePostInput struct {
	PostID int32      `path:"postID" doc:"文章ID"`
	Body   param.Post `doc:"文章参数"`
}

// UpdatePost 更新文章。
func (p *PostHandler) UpdatePost(ctx context.Context, in *UpdatePostInput) (*dto.HumaOut[*entity.Post], error) {
	postParam := in.Body
	post, err := p.PostService.Update(ctx, in.PostID, &postParam)
	if err != nil {
		return dto.HumaErr[*entity.Post](err)
	}
	return dto.HumaOK(post)
}

// UpdatePostStatusInput 更新文章状态输入。
type UpdatePostStatusInput struct {
	PostID int32  `path:"postID" doc:"文章ID"`
	Status string `path:"status" doc:"文章状态（PUBLISHED/DRAFT/RECYCLE/INTIMATE）"`
}

// UpdatePostStatus 更新文章状态。
func (p *PostHandler) UpdatePostStatus(ctx context.Context, in *UpdatePostStatusInput) (*dto.HumaOut[*dto.PostMinimal], error) {
	status, err := consts.PostStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[*dto.PostMinimal](xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error"))
	}
	if int32(status) < int32(consts.PostStatusPublished) || int32(status) > int32(consts.PostStatusIntimate) {
		return dto.HumaErr[*dto.PostMinimal](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error"))
	}
	post, err := p.PostService.UpdateStatus(ctx, in.PostID, status)
	if err != nil {
		return dto.HumaErr[*dto.PostMinimal](err)
	}
	postMinimal, err := p.PostAssembler.ConvertToMinimalDTO(ctx, post)
	if err != nil {
		return dto.HumaErr[*dto.PostMinimal](err)
	}
	return dto.HumaOK(postMinimal)
}

// UpdatePostStatusBatchInput 批量更新文章状态输入。
type UpdatePostStatusBatchInput struct {
	Status string  `path:"status" doc:"文章状态（PUBLISHED/DRAFT/RECYCLE/INTIMATE）"`
	Body   []int32 `doc:"文章ID列表"`
}

// UpdatePostStatusBatch 批量更新文章状态。
func (p *PostHandler) UpdatePostStatusBatch(ctx context.Context, in *UpdatePostStatusBatchInput) (*dto.HumaOut[[]*entity.Post], error) {
	status, err := consts.PostStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[[]*entity.Post](err)
	}
	if int32(status) < int32(consts.PostStatusPublished) || int32(status) > int32(consts.PostStatusIntimate) {
		return dto.HumaErr[[]*entity.Post](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error"))
	}
	posts, err := p.PostService.UpdateStatusBatch(ctx, status, in.Body)
	if err != nil {
		return dto.HumaErr[[]*entity.Post](err)
	}
	return dto.HumaOK(posts)
}

// UpdatePostDraftInput 更新文章草稿内容输入。
type UpdatePostDraftInput struct {
	PostID int32             `path:"postID" doc:"文章ID"`
	Body   param.PostContent `doc:"草稿内容"`
}

// UpdatePostDraft 更新文章草稿内容。
func (p *PostHandler) UpdatePostDraft(ctx context.Context, in *UpdatePostDraftInput) (*dto.HumaOut[*dto.PostDetail], error) {
	post, err := p.PostService.UpdateDraftContent(ctx, in.PostID, in.Body.Content, in.Body.OriginalContent)
	if err != nil {
		return dto.HumaErr[*dto.PostDetail](err)
	}
	postDetailDTO, err := p.PostAssembler.ConvertToDetailDTO(ctx, post)
	if err != nil {
		return dto.HumaErr[*dto.PostDetail](err)
	}
	return dto.HumaOK(postDetailDTO)
}

// DeletePostInput 删除文章输入。
type DeletePostInput struct {
	PostID int32 `path:"postID" doc:"文章ID"`
}

// DeletePost 删除文章。
func (p *PostHandler) DeletePost(ctx context.Context, in *DeletePostInput) (*dto.HumaOut[any], error) {
	if err := p.PostService.Delete(ctx, in.PostID); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// DeletePostBatchInput 批量删除文章输入。
type DeletePostBatchInput struct {
	Body []int32 `doc:"文章ID列表"`
}

// DeletePostBatch 批量删除文章。
func (p *PostHandler) DeletePostBatch(ctx context.Context, in *DeletePostBatchInput) (*dto.HumaOut[any], error) {
	if err := p.PostService.DeleteBatch(ctx, in.Body); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// PreviewPost 预览文章（保持 gin handler，返回 HTML/文件流）。
func (p *PostHandler) PreviewPost(ctx *gin.Context) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		_ = ctx.Error(err)
		return
	}
	previewPath, err := p.PostService.Preview(ctx, postID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}
	ctx.String(http.StatusOK, previewPath)
}
