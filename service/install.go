package service

import (
	"context"

	"github.com/rfancn/airpress/model/param"
)

type InstallService interface {
	InstallBlog(ctx context.Context, installParam param.Install) error
}
