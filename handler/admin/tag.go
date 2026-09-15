package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type TagHandler struct {
	PostTagService service.PostTagService
	TagService     service.TagService
}

func NewTagHandler(postTagService service.PostTagService, tagService service.TagService) *TagHandler {
	return &TagHandler{
		PostTagService: postTagService,
		TagService:     tagService,
	}
}

// ListTagsInput 标签列表查询输入。
type ListTagsInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
	More bool     `query:"more" doc:"是否返回带文章数的标签"`
}

// ListTags 获取标签列表。
func (t *TagHandler) ListTags(ctx context.Context, in *ListTagsInput) (*dto.HumaOut[interface{}], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"createTime,desc"}
	}
	if in.More {
		data, err := t.PostTagService.ListAllTagWithPostCount(ctx, sort)
		if err != nil {
			return dto.HumaErr[interface{}](err)
		}
		return dto.HumaOK[interface{}](data)
	}
	tags, err := t.TagService.ListAll(ctx, sort)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	data, err := t.TagService.ConvertToDTOs(ctx, tags)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](data)
}

// GetTagByIDInput 标签详情查询输入。
type GetTagByIDInput struct {
	ID int32 `path:"id" doc:"标签ID"`
}

// GetTagByID 获取标签详情。
func (t *TagHandler) GetTagByID(ctx context.Context, in *GetTagByIDInput) (*dto.HumaOut[*dto.Tag], error) {
	tag, err := t.TagService.GetByID(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	data, err := t.TagService.ConvertToDTO(ctx, tag)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	return dto.HumaOK(data)
}

// CreateTagInput 创建标签输入。
type CreateTagInput struct {
	Body param.Tag `doc:"标签参数"`
}

// CreateTag 创建标签。
func (t *TagHandler) CreateTag(ctx context.Context, in *CreateTagInput) (*dto.HumaOut[*dto.Tag], error) {
	tag, err := t.TagService.Create(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	data, err := t.TagService.ConvertToDTO(ctx, tag)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	return dto.HumaOK(data)
}

// UpdateTagInput 更新标签输入。
type UpdateTagInput struct {
	ID   int32     `path:"id" doc:"标签ID"`
	Body param.Tag `doc:"标签参数"`
}

// UpdateTag 更新标签。
func (t *TagHandler) UpdateTag(ctx context.Context, in *UpdateTagInput) (*dto.HumaOut[*dto.Tag], error) {
	tag, err := t.TagService.Update(ctx, in.ID, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	data, err := t.TagService.ConvertToDTO(ctx, tag)
	if err != nil {
		return dto.HumaErr[*dto.Tag](err)
	}
	return dto.HumaOK(data)
}

// DeleteTagInput 删除标签输入。
type DeleteTagInput struct {
	ID int32 `path:"id" doc:"标签ID"`
}

// DeleteTag 删除标签。
func (t *TagHandler) DeleteTag(ctx context.Context, in *DeleteTagInput) (*dto.HumaOut[interface{}], error) {
	err := t.TagService.Delete(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}
