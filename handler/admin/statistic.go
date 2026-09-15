package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/service"
)

type StatisticHandler struct {
	StatisticService service.StatisticService
}

func NewStatisticHandler(l service.StatisticService) *StatisticHandler {
	return &StatisticHandler{
		StatisticService: l,
	}
}

// StatisticsInput 统计查询无输入参数。
type StatisticsInput struct{}

// Statistics 获取统计数据。
func (s *StatisticHandler) Statistics(ctx context.Context, _ *StatisticsInput) (*dto.HumaOut[*dto.Statistic], error) {
	data, err := s.StatisticService.Statistic(ctx)
	if err != nil {
		return dto.HumaErr[*dto.Statistic](err)
	}
	return dto.HumaOK(data)
}

// StatisticsWithUserInput 统计查询无输入参数。
type StatisticsWithUserInput struct{}

// StatisticsWithUser 获取带用户的统计数据。
func (s *StatisticHandler) StatisticsWithUser(ctx context.Context, _ *StatisticsWithUserInput) (*dto.HumaOut[*dto.StatisticWithUser], error) {
	data, err := s.StatisticService.StatisticWithUser(ctx)
	if err != nil {
		return dto.HumaErr[*dto.StatisticWithUser](err)
	}
	return dto.HumaOK(data)
}
