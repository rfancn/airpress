package service

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
)

type SheetService interface {
	BasePostService
	Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Post, int64, error)
	Create(ctx context.Context, sheetParam *param.Sheet) (*entity.Post, error)
	Update(ctx context.Context, sheetID int32, sheetParam *param.Sheet) (*entity.Post, error)
	Preview(ctx context.Context, sheetID int32) (string, error)
	CountVisit(ctx context.Context) (int64, error)
	CountLike(ctx context.Context) (int64, error)
	ListIndependentSheets(ctx context.Context) ([]*dto.IndependentSheet, error)
}
