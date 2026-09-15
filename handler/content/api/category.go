package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/handler/content/authentication"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/util"
)

type CategoryHandler struct {
	PostService            service.PostService
	CategoryService        service.CategoryService
	CategoryAuthentication authentication.CategoryAuthentication
	PostAssembler          assembler.PostAssembler
}

func NewCategoryHandler(postService service.PostService, categoryService service.CategoryService, categoryAuthentication *authentication.CategoryAuthentication, postAssembler assembler.PostAssembler) *CategoryHandler {
	return &CategoryHandler{
		PostService:            postService,
		CategoryService:        categoryService,
		CategoryAuthentication: *categoryAuthentication,
		PostAssembler:          postAssembler,
	}
}

// ListCategoriesInput 分类列表查询输入。
type ListCategoriesInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
	More bool     `query:"more" doc:"是否返回带文章数的分类"`
}

// ListCategories 获取分类列表。
func (c *CategoryHandler) ListCategories(ctx context.Context, in *ListCategoriesInput) (*dto.HumaOut[interface{}], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"updateTime,desc"}
	}
	if in.More {
		data, err := c.CategoryService.ListCategoryWithPostCountDTO(ctx, sort)
		if err != nil {
			return dto.HumaErr[interface{}](err)
		}
		return dto.HumaOK[interface{}](data)
	}
	categories, err := c.CategoryService.ListAll(ctx, sort)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	data, err := c.CategoryService.ConvertToCategoryDTOs(ctx, categories)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](data)
}

// ListPostsInput 分类下文章列表查询输入。
type ListPostsInput struct {
	Slug           string   `path:"slug" doc:"分类别名"`
	Page           int      `query:"page" doc:"页码"`
	Size           int      `query:"size" doc:"每页数量"`
	Sort           []string `query:"sort" doc:"排序字段"`
	Keyword        string   `query:"keyword" doc:"关键字"`
	Password       string   `query:"password" doc:"分类密码"`
	Authentication string   `cookie:"authentication" doc:"认证 token"`
}

// ListPosts 获取分类下的文章列表。
func (c *CategoryHandler) ListPosts(ctx context.Context, in *ListPostsInput) (*dto.HumaOut[*dto.Page], error) {
	category, err := c.CategoryService.GetBySlug(ctx, in.Slug)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}

	postQuery := param.PostQuery{
		Page: param.Page{
			PageNum:  in.Page,
			PageSize: in.Size,
		},
		Sort: &param.Sort{Fields: in.Sort},
	}
	// huma 不支持指针 query 参数，空关键字视为未提供（等价于原 nil）
	if in.Keyword != "" {
		postQuery.Keyword = &in.Keyword
	}
	if len(postQuery.Sort.Fields) == 0 {
		postQuery.Sort.Fields = []string{"topPriority,desc", "updateTime,desc"}
	}

	if category.Type == consts.CategoryTypeIntimate {
		token := in.Authentication
		if authenticated, _ := c.CategoryAuthentication.IsAuthenticated(ctx, token, category.ID); !authenticated {
			newToken, err := c.CategoryAuthentication.Authenticate(ctx, token, category.ID, in.Password)
			if err != nil {
				return dto.HumaErr[*dto.Page](err)
			}
			// 恢复私密分类认证成功后的 Cookie 设置（等价于原 gin 的 ctx.SetCookie）
			if hctx, ok := ctx.(huma.Context); ok {
				humagin.Unwrap(hctx).SetCookie("authentication", newToken, 1800, "/", "", false, true)
			}
		}
	}

	postQuery.WithPassword = util.BoolPtr(false)
	postQuery.Statuses = []*consts.PostStatus{consts.PostStatusPublished.Ptr(), consts.PostStatusIntimate.Ptr()}

	posts, totalCount, err := c.PostService.Page(ctx, postQuery)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	postVOs, err := c.PostAssembler.ConvertToListVO(ctx, posts)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(postVOs, totalCount, postQuery.Page))
}
