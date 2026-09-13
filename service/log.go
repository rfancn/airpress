package service

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
)

type LogService interface {
	PageLog(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Log, int64, error)
	ConvertToDTO(log *entity.Log) *dto.Log
	Clear(ctx context.Context) error
}
