package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

// GetCategoryByIDInput 分类详情查询输入。
type GetCategoryByIDInput struct {
	CategoryID int32 `path:"categoryID" doc:"分类ID"`
}

// GetCategoryByID 获取分类详情。
func (c *CategoryHandler) GetCategoryByID(ctx context.Context, in *GetCategoryByIDInput) (*dto.HumaOut[*dto.CategoryDTO], error) {
	category, err := c.CategoryService.GetByID(ctx, in.CategoryID)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	data, err := c.CategoryService.ConvertToCategoryDTO(ctx, category)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	return dto.HumaOK(data)
}

// ListAllCategoryInput 分类列表查询输入。
type ListAllCategoryInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
	More bool     `query:"more" doc:"是否返回带文章数的分类"`
}

// ListAllCategory 获取分类列表。
func (c *CategoryHandler) ListAllCategory(ctx context.Context, in *ListAllCategoryInput) (*dto.HumaOut[interface{}], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"priority,asc"}
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

// ListAsTreeInput 分类树查询输入。
type ListAsTreeInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// ListAsTree 获取分类树。
func (c *CategoryHandler) ListAsTree(ctx context.Context, in *ListAsTreeInput) (*dto.HumaOut[[]*vo.CategoryVO], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"priority,asc"}
	}
	data, err := c.CategoryService.ListAsTree(ctx, sort, false)
	if err != nil {
		return dto.HumaErr[[]*vo.CategoryVO](err)
	}
	return dto.HumaOK(data)
}

// CreateCategoryInput 创建分类输入。
type CreateCategoryInput struct {
	Body param.Category `doc:"分类参数"`
}

// CreateCategory 创建分类。
func (c *CategoryHandler) CreateCategory(ctx context.Context, in *CreateCategoryInput) (*dto.HumaOut[*dto.CategoryDTO], error) {
	category, err := c.CategoryService.Create(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	data, err := c.CategoryService.ConvertToCategoryDTO(ctx, category)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	return dto.HumaOK(data)
}

// UpdateCategoryInput 更新分类输入。
type UpdateCategoryInput struct {
	CategoryID int32          `path:"categoryID" doc:"分类ID"`
	Body       param.Category `doc:"分类参数"`
}

// UpdateCategory 更新分类。
func (c *CategoryHandler) UpdateCategory(ctx context.Context, in *UpdateCategoryInput) (*dto.HumaOut[*dto.CategoryDTO], error) {
	in.Body.ID = in.CategoryID
	category, err := c.CategoryService.Update(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	data, err := c.CategoryService.ConvertToCategoryDTO(ctx, category)
	if err != nil {
		return dto.HumaErr[*dto.CategoryDTO](err)
	}
	return dto.HumaOK(data)
}

// UpdateCategoryBatchInput 批量更新分类输入。
type UpdateCategoryBatchInput struct {
	Body []*param.Category `doc:"批量分类参数"`
}

// UpdateCategoryBatch 批量更新分类。
func (c *CategoryHandler) UpdateCategoryBatch(ctx context.Context, in *UpdateCategoryBatchInput) (*dto.HumaOut[[]*dto.CategoryDTO], error) {
	categories, err := c.CategoryService.UpdateBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[[]*dto.CategoryDTO](err)
	}
	data, err := c.CategoryService.ConvertToCategoryDTOs(ctx, categories)
	if err != nil {
		return dto.HumaErr[[]*dto.CategoryDTO](err)
	}
	return dto.HumaOK(data)
}

// DeleteCategoryInput 删除分类输入。
type DeleteCategoryInput struct {
	CategoryID int32 `path:"categoryID" doc:"分类ID"`
}

// DeleteCategory 删除分类。
func (c *CategoryHandler) DeleteCategory(ctx context.Context, in *DeleteCategoryInput) (*dto.HumaOut[interface{}], error) {
	err := c.CategoryService.Delete(ctx, in.CategoryID)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}
