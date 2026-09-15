package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
)

type LogHandler struct {
	LogService service.LogService
}

func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{
		LogService: logService,
	}
}

// PageLatestLogInput 最新日志查询输入。
type PageLatestLogInput struct {
	Top int `query:"top" doc:"获取数量"`
}

// PageLatestLog 获取最新日志列表。
func (l *LogHandler) PageLatestLog(ctx context.Context, in *PageLatestLogInput) (*dto.HumaOut[[]*dto.Log], error) {
	top := in.Top
	if top == 0 {
		top = 10
	}
	logs, _, err := l.LogService.PageLog(ctx, param.Page{PageSize: top}, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return dto.HumaErr[[]*dto.Log](err)
	}
	logDTOs := make([]*dto.Log, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, l.LogService.ConvertToDTO(log))
	}
	return dto.HumaOK(logDTOs)
}

// PageLogInput 日志分页查询输入。
type PageLogInput struct {
	Page int      `query:"page" doc:"页码"`
	Size int      `query:"size" doc:"每页数量"`
	Sort []string `query:"sort" doc:"排序字段"`
}

// PageLog 分页查询日志。
func (l *LogHandler) PageLog(ctx context.Context, in *PageLogInput) (*dto.HumaOut[*dto.Page], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort = &param.Sort{Fields: []string{"createTime,desc"}}
	}
	logs, totalCount, err := l.LogService.PageLog(ctx, param.Page{PageNum: in.Page, PageSize: in.Size}, sort)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	logDTOs := make([]*dto.Log, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, l.LogService.ConvertToDTO(log))
	}
	return dto.HumaOK(dto.NewPage(logDTOs, totalCount, param.Page{PageNum: in.Page, PageSize: in.Size}))
}

// ClearLogInput 清空日志无输入参数。
type ClearLogInput struct{}

// ClearLog 清空日志。
func (l *LogHandler) ClearLog(ctx context.Context, _ *ClearLogInput) (*dto.HumaOut[any], error) {
	err := l.LogService.Clear(ctx)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
