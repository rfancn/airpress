package api

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/service"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		PhotoService: photoService,
	}
}

// LikePhotoInput 图片点赞输入。
type LikePhotoInput struct {
	PhotoID int32 `path:"photoID" doc:"图片ID"`
}

// Like 图片点赞。
func (p *PhotoHandler) Like(ctx context.Context, in *LikePhotoInput) (*dto.HumaOut[any], error) {
	err := p.PhotoService.IncreaseLike(ctx, in.PhotoID)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
