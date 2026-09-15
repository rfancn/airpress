package admin

import (
	"context"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type JournalHandler struct {
	JournalService service.JournalService
}

func NewJournalHandler(journalService service.JournalService) *JournalHandler {
	return &JournalHandler{
		JournalService: journalService,
	}
}

// ListJournalInput 日志列表查询输入。
type ListJournalInput struct {
	Page        int    `query:"page" doc:"页码"`
	Size        int    `query:"size" doc:"每页数量"`
	Keyword     string `query:"keyword" doc:"关键词"`
	JournalType int32  `query:"journalType" doc:"日志类型"`
}

// ListJournal 获取日志分页列表。
func (j *JournalHandler) ListJournal(ctx context.Context, in *ListJournalInput) (*dto.HumaOut[*dto.Page], error) {
	journalQuery := param.JournalQuery{
		Page: param.Page{PageNum: in.Page, PageSize: in.Size},
	}
	// huma 不支持指针 query 参数，空串/0 视为未提供
	if in.Keyword != "" {
		journalQuery.Keyword = &in.Keyword
	}
	if in.JournalType != 0 {
		jt := consts.JournalType(in.JournalType)
		journalQuery.JournalType = &jt
	}
	journalQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	journals, totalCount, err := j.JournalService.ListJournal(ctx, journalQuery)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, journals)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(journalDTOs, totalCount, journalQuery.Page))
}

// ListLatestJournalInput 最新日志查询输入。
type ListLatestJournalInput struct {
	Top int `query:"top" doc:"获取数量"`
}

// ListLatestJournal 获取最新日志列表。
func (j *JournalHandler) ListLatestJournal(ctx context.Context, in *ListLatestJournalInput) (*dto.HumaOut[[]*dto.JournalWithComment], error) {
	top := in.Top
	if top == 0 {
		top = 10
	}
	journalQuery := param.JournalQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: top},
	}
	journals, _, err := j.JournalService.ListJournal(ctx, journalQuery)
	if err != nil {
		return dto.HumaErr[[]*dto.JournalWithComment](err)
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, journals)
	if err != nil {
		return dto.HumaErr[[]*dto.JournalWithComment](err)
	}
	return dto.HumaOK(journalDTOs)
}

// CreateJournalInput 创建日志输入。
type CreateJournalInput struct {
	Body param.Journal `doc:"日志内容"`
}

// CreateJournal 创建日志。
func (j *JournalHandler) CreateJournal(ctx context.Context, in *CreateJournalInput) (*dto.HumaOut[*dto.Journal], error) {
	journalParam := in.Body
	if journalParam.Content == "" {
		journalParam.Content = journalParam.SourceContent
	}
	journal, err := j.JournalService.Create(ctx, &journalParam)
	if err != nil {
		return dto.HumaErr[*dto.Journal](err)
	}
	return dto.HumaOK(j.JournalService.ConvertToDTO(journal))
}

// UpdateJournalInput 更新日志输入。
type UpdateJournalInput struct {
	JournalID int32         `path:"journalID" doc:"日志ID"`
	Body      param.Journal `doc:"日志内容"`
}

// UpdateJournal 更新日志。
func (j *JournalHandler) UpdateJournal(ctx context.Context, in *UpdateJournalInput) (*dto.HumaOut[*entity.Journal], error) {
	journalParam := in.Body
	journal, err := j.JournalService.Update(ctx, in.JournalID, &journalParam)
	if err != nil {
		return dto.HumaErr[*entity.Journal](err)
	}
	return dto.HumaOK(journal)
}

// DeleteJournalInput 删除日志输入。
type DeleteJournalInput struct {
	JournalID int32 `path:"journalID" doc:"日志ID"`
}

// DeleteJournal 删除日志。
func (j *JournalHandler) DeleteJournal(ctx context.Context, in *DeleteJournalInput) (*dto.HumaOut[any], error) {
	err := j.JournalService.Delete(ctx, in.JournalID)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
