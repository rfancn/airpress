package admin

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		PhotoService: photoService,
	}
}

// ListPhotoInput 照片列表查询输入。
type ListPhotoInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// ListPhoto 获取照片列表。
func (p *PhotoHandler) ListPhoto(ctx context.Context, in *ListPhotoInput) (*dto.HumaOut[[]*dto.Photo], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"createTime,desc"}
	}
	photos, err := p.PhotoService.List(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*dto.Photo](err)
	}
	return dto.HumaOK(p.PhotoService.ConvertToDTOs(ctx, photos))
}

// PagePhotosInput 照片分页查询输入。
type PagePhotosInput struct {
	Page int      `query:"page" doc:"页码"`
	Size int      `query:"size" doc:"每页数量"`
	Sort []string `query:"sort" doc:"排序字段"`
}

// PagePhotos 分页获取照片列表。
func (p *PhotoHandler) PagePhotos(ctx context.Context, in *PagePhotosInput) (*dto.HumaOut[*dto.Page], error) {
	sort := &param.Sort{Fields: in.Sort}
	if len(sort.Fields) == 0 {
		sort.Fields = []string{"createTime,desc"}
	}
	page := param.Page{PageNum: in.Page, PageSize: in.Size}
	photos, totalCount, err := p.PhotoService.Page(ctx, page, sort)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(p.PhotoService.ConvertToDTOs(ctx, photos), totalCount, page))
}

// GetPhotoByIDInput 照片详情查询输入。
type GetPhotoByIDInput struct {
	ID int32 `path:"id" doc:"照片ID"`
}

// GetPhotoByID 获取照片详情。
func (p *PhotoHandler) GetPhotoByID(ctx context.Context, in *GetPhotoByIDInput) (*dto.HumaOut[*dto.Photo], error) {
	photo, err := p.PhotoService.GetByID(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*dto.Photo](err)
	}
	return dto.HumaOK(p.PhotoService.ConvertToDTO(ctx, photo))
}

// CreatePhotoInput 创建照片输入。
type CreatePhotoInput struct {
	Body param.Photo `doc:"照片参数"`
}

// CreatePhoto 创建照片。
func (p *PhotoHandler) CreatePhoto(ctx context.Context, in *CreatePhotoInput) (*dto.HumaOut[*dto.Photo], error) {
	photo, err := p.PhotoService.Create(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Photo](err)
	}
	return dto.HumaOK(p.PhotoService.ConvertToDTO(ctx, photo))
}

// CreatePhotoBatchInput 批量创建照片输入。
type CreatePhotoBatchInput struct {
	Body []*param.Photo `doc:"照片参数列表"`
}

// CreatePhotoBatch 批量创建照片。
func (p *PhotoHandler) CreatePhotoBatch(ctx context.Context, in *CreatePhotoBatchInput) (*dto.HumaOut[[]*dto.Photo], error) {
	photos, err := p.PhotoService.CreateBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[[]*dto.Photo](err)
	}
	return dto.HumaOK(p.PhotoService.ConvertToDTOs(ctx, photos))
}

// UpdatePhotoInput 更新照片输入。
type UpdatePhotoInput struct {
	ID   int32       `path:"id" doc:"照片ID"`
	Body param.Photo `doc:"照片参数"`
}

// UpdatePhoto 更新照片。
func (p *PhotoHandler) UpdatePhoto(ctx context.Context, in *UpdatePhotoInput) (*dto.HumaOut[*dto.Photo], error) {
	photo, err := p.PhotoService.Update(ctx, in.ID, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Photo](err)
	}
	return dto.HumaOK(p.PhotoService.ConvertToDTO(ctx, photo))
}

// DeletePhoto 删除照片（未注册路由，保持 gin 签名）。
func (p *PhotoHandler) DeletePhoto(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, p.PhotoService.Delete(ctx, id)
}

// DeletePhotoBatchInput 批量删除照片输入。
type DeletePhotoBatchInput struct {
	Body []int32 `doc:"照片ID列表"`
}

// DeletePhotoBatch 批量删除照片。
func (p *PhotoHandler) DeletePhotoBatch(ctx context.Context, in *DeletePhotoBatchInput) (*dto.HumaOut[any], error) {
	for _, id := range in.Body {
		err := p.PhotoService.Delete(ctx, id)
		if err != nil {
			return dto.HumaErr[any](err)
		}
	}
	return dto.HumaOK[any](nil)
}

// ListPhotoTeamsInput 照片团队查询无输入参数。
type ListPhotoTeamsInput struct{}

// ListPhotoTeams 获取照片团队列表。
func (p *PhotoHandler) ListPhotoTeams(ctx context.Context, _ *ListPhotoTeamsInput) (*dto.HumaOut[[]string], error) {
	teams, err := p.PhotoService.ListTeams(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(teams)
}
