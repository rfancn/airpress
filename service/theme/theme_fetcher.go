package theme

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
)

type ThemeFetcher interface {
	FetchTheme(ctx context.Context, file interface{}) (*dto.ThemeProperty, error)
}
