package admin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type SheetHandler struct {
	SheetService   service.SheetService
	PostService    service.PostService
	SheetAssembler assembler.SheetAssembler
}

func NewSheetHandler(sheetService service.SheetService, postService service.PostService, sheetAssembler assembler.SheetAssembler) *SheetHandler {
	return &SheetHandler{
		SheetService:   sheetService,
		PostService:    postService,
		SheetAssembler: sheetAssembler,
	}
}

// GetSheetByIDInput 页面详情查询输入。
type GetSheetByIDInput struct {
	SheetID int32 `path:"sheetID" doc:"页面ID"`
}

// GetSheetByID 获取页面详情。
func (s *SheetHandler) GetSheetByID(ctx context.Context, in *GetSheetByIDInput) (*dto.HumaOut[*vo.SheetDetail], error) {
	sheet, err := s.SheetService.GetByPostID(ctx, in.SheetID)
	if err != nil {
		return dto.HumaErr[*vo.SheetDetail](err)
	}
	sheetDetailVO, err := s.SheetAssembler.ConvertToDetailVO(ctx, sheet)
	if err != nil {
		return dto.HumaErr[*vo.SheetDetail](err)
	}
	return dto.HumaOK(sheetDetailVO)
}

// ListSheetInput 页面列表查询输入。
type ListSheetInput struct {
	Page int `query:"page" doc:"页码"`
	Size int `query:"size" doc:"每页数量"`
}

// ListSheet 获取页面列表（分页）。
func (s *SheetHandler) ListSheet(ctx context.Context, in *ListSheetInput) (*dto.HumaOut[*dto.Page], error) {
	page := param.Page{PageNum: in.Page, PageSize: in.Size}
	sheets, totalCount, err := s.SheetService.Page(ctx, page, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	sheetVOs, err := s.SheetAssembler.ConvertToListVO(ctx, sheets)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(sheetVOs, totalCount, page))
}

// IndependentSheetsInput 独立页面查询无输入参数。
type IndependentSheetsInput struct{}

// IndependentSheets 获取独立页面列表。
func (s *SheetHandler) IndependentSheets(ctx context.Context, _ *IndependentSheetsInput) (*dto.HumaOut[[]*dto.IndependentSheet], error) {
	sheets, err := s.SheetService.ListIndependentSheets(ctx)
	if err != nil {
		return dto.HumaErr[[]*dto.IndependentSheet](err)
	}
	return dto.HumaOK(sheets)
}

// CreateSheetInput 创建页面输入。
type CreateSheetInput struct {
	Body param.Sheet `doc:"页面参数"`
}

// CreateSheet 创建页面。
func (s *SheetHandler) CreateSheet(ctx context.Context, in *CreateSheetInput) (*dto.HumaOut[*vo.SheetDetail], error) {
	sheetParam := in.Body
	sheet, err := s.SheetService.Create(ctx, &sheetParam)
	if err != nil {
		return dto.HumaErr[*vo.SheetDetail](err)
	}
	sheetDetailVO, err := s.SheetAssembler.ConvertToDetailVO(ctx, sheet)
	if err != nil {
		return dto.HumaErr[*vo.SheetDetail](err)
	}
	return dto.HumaOK(sheetDetailVO)
}

// UpdateSheetInput 更新页面输入。
type UpdateSheetInput struct {
	SheetID int32       `path:"sheetID" doc:"页面ID"`
	Body    param.Sheet `doc:"页面参数"`
}

// UpdateSheet 更新页面。
func (s *SheetHandler) UpdateSheet(ctx context.Context, in *UpdateSheetInput) (*dto.HumaOut[*entity.Post], error) {
	sheetParam := in.Body
	post, err := s.SheetService.Update(ctx, in.SheetID, &sheetParam)
	if err != nil {
		return dto.HumaErr[*entity.Post](err)
	}
	return dto.HumaOK(post)
}

// UpdateSheetStatusInput 更新页面状态输入。
type UpdateSheetStatusInput struct {
	SheetID int32  `path:"sheetID" doc:"页面ID"`
	Status  string `path:"status" doc:"页面状态（PUBLISHED/DRAFT/RECYCLE/INTIMATE）"`
}

// UpdateSheetStatus 更新页面状态。
func (s *SheetHandler) UpdateSheetStatus(ctx context.Context, in *UpdateSheetStatusInput) (*dto.HumaOut[*entity.Post], error) {
	status, err := consts.PostStatusFromString(in.Status)
	if err != nil {
		return dto.HumaErr[*entity.Post](err)
	}
	if status < consts.PostStatusPublished || status > consts.PostStatusIntimate {
		return dto.HumaErr[*entity.Post](xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error"))
	}
	post, err := s.SheetService.UpdateStatus(ctx, in.SheetID, status)
	if err != nil {
		return dto.HumaErr[*entity.Post](err)
	}
	return dto.HumaOK(post)
}

// UpdateSheetDraftInput 更新页面草稿内容输入。
type UpdateSheetDraftInput struct {
	SheetID int32             `path:"sheetID" doc:"页面ID"`
	Body    param.PostContent `doc:"草稿内容"`
}

// UpdateSheetDraft 更新页面草稿内容。
func (s *SheetHandler) UpdateSheetDraft(ctx context.Context, in *UpdateSheetDraftInput) (*dto.HumaOut[*dto.PostDetail], error) {
	post, err := s.SheetService.UpdateDraftContent(ctx, in.SheetID, in.Body.Content, in.Body.OriginalContent)
	if err != nil {
		return dto.HumaErr[*dto.PostDetail](err)
	}
	postDetailDTO, err := s.SheetAssembler.ConvertToDetailDTO(ctx, post)
	if err != nil {
		return dto.HumaErr[*dto.PostDetail](err)
	}
	return dto.HumaOK(postDetailDTO)
}

// DeleteSheetInput 删除页面输入。
type DeleteSheetInput struct {
	SheetID int32 `path:"sheetID" doc:"页面ID"`
}

// DeleteSheet 删除页面。
func (s *SheetHandler) DeleteSheet(ctx context.Context, in *DeleteSheetInput) (*dto.HumaOut[any], error) {
	if err := s.SheetService.Delete(ctx, in.SheetID); err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// PreviewSheet 预览页面（保持 gin handler，返回 HTML/文件流）。
func (s *SheetHandler) PreviewSheet(ctx *gin.Context) {
	sheetID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}
	previewPath, err := s.SheetService.Preview(ctx, sheetID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}
	ctx.String(http.StatusOK, previewPath)
}
