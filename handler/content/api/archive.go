package api

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
)

type ArchiveHandler struct {
	PostService   service.PostService
	PostAssembler assembler.PostAssembler
}

func NewArchiveHandler(postService service.PostService, postAssemeber assembler.PostAssembler) *ArchiveHandler {
	return &ArchiveHandler{
		PostService:   postService,
		PostAssembler: postAssemeber,
	}
}

// ListYearArchivesInput 按年份归档查询无输入参数。
type ListYearArchivesInput struct{}

// ListYearArchives 获取按年份归档的文章列表。
func (a *ArchiveHandler) ListYearArchives(ctx context.Context, _ *ListYearArchivesInput) (*dto.HumaOut[[]*vo.ArchiveYear], error) {
	posts, err := a.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
	if err != nil {
		return dto.HumaErr[[]*vo.ArchiveYear](err)
	}
	vos, err := a.PostAssembler.ConvertToArchiveYearVOs(ctx, posts)
	if err != nil {
		return dto.HumaErr[[]*vo.ArchiveYear](err)
	}
	return dto.HumaOK(vos)
}

// ListMonthArchivesInput 按月份归档查询无输入参数。
type ListMonthArchivesInput struct{}

// ListMonthArchives 获取按月份归档的文章列表。
func (a *ArchiveHandler) ListMonthArchives(ctx context.Context, _ *ListMonthArchivesInput) (*dto.HumaOut[[]*vo.ArchiveMonth], error) {
	posts, err := a.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
	if err != nil {
		return dto.HumaErr[[]*vo.ArchiveMonth](err)
	}
	vos, err := a.PostAssembler.ConvertTOArchiveMonthVOs(ctx, posts)
	if err != nil {
		return dto.HumaErr[[]*vo.ArchiveMonth](err)
	}
	return dto.HumaOK(vos)
}
