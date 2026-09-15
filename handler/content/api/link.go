package api

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
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
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"createTime,desc"}
	}
	links, err := l.LinkService.List(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*dto.Link](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToDTOs(ctx, links))
}

// LinkTeamVOInput 链接团队视图查询输入。
type LinkTeamVOInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// LinkTeamVO 获取链接团队视图。
func (l *LinkHandler) LinkTeamVO(ctx context.Context, in *LinkTeamVOInput) (*dto.HumaOut[[]*vo.LinkTeamVO], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"createTime,desc"}
	}
	links, err := l.LinkService.List(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*vo.LinkTeamVO](err)
	}
	return dto.HumaOK(l.LinkService.ConvertToLinkTeamVO(ctx, links))
}
