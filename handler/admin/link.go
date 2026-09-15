package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type LinkHandler struct {
	LinkService service.LinkService
}

func NewLinkHandler(linkService service.LinkService) *LinkHandler {
	return &LinkHandler{
		LinkService: linkService,
	}
}

// ListLinksInput 链接列表查询输入。
type ListLinksInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// ListLinks 获取链接列表。
func (l *LinkHandler) ListLinks(ctx context.Context, in *ListLinksInput) (*dto.HumaOut[[]*dto.Link], error) {
	fields := in.Sort
	if len(fields) == 0 {
		fields = []string{"team,desc", "priority,asc"}
	} else {
		fields = append(fields, "priority,asc")
	}
	sort := &param.Sort{Fields: fields}
	links, err := l.LinkService.List(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*dto.Link](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToDTOs(ctx, links))
}

// GetLinkByIDInput 链接详情查询输入。
type GetLinkByIDInput struct {
	ID int32 `path:"id" doc:"链接ID"`
}

// GetLinkByID 获取链接详情。
func (l *LinkHandler) GetLinkByID(ctx context.Context, in *GetLinkByIDInput) (*dto.HumaOut[*dto.Link], error) {
	link, err := l.LinkService.GetByID(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*dto.Link](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToDTO(ctx, link))
}

// CreateLinkInput 创建链接输入。
type CreateLinkInput struct {
	Body param.Link `doc:"链接参数"`
}

// CreateLink 创建链接。
func (l *LinkHandler) CreateLink(ctx context.Context, in *CreateLinkInput) (*dto.HumaOut[*dto.Link], error) {
	link, err := l.LinkService.Create(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Link](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToDTO(ctx, link))
}

// UpdateLinkInput 更新链接输入。
type UpdateLinkInput struct {
	ID   int32      `path:"id" doc:"链接ID"`
	Body param.Link `doc:"链接参数"`
}

// UpdateLink 更新链接。
func (l *LinkHandler) UpdateLink(ctx context.Context, in *UpdateLinkInput) (*dto.HumaOut[*dto.Link], error) {
	link, err := l.LinkService.Update(ctx, in.ID, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Link](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToDTO(ctx, link))
}

// DeleteLinkInput 删除链接输入。
type DeleteLinkInput struct {
	ID int32 `path:"id" doc:"链接ID"`
}

// DeleteLink 删除链接。
func (l *LinkHandler) DeleteLink(ctx context.Context, in *DeleteLinkInput) (*dto.HumaOut[interface{}], error) {
	err := l.LinkService.Delete(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// ListLinkTeamsInput 链接团队列表查询无输入参数。
type ListLinkTeamsInput struct{}

// ListLinkTeams 获取链接团队列表。
func (l *LinkHandler) ListLinkTeams(ctx context.Context, _ *ListLinkTeamsInput) (*dto.HumaOut[[]string], error) {
	data, err := l.LinkService.ListTeams(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(data)
}
