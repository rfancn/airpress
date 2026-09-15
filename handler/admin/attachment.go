package admin

import (
	"context"
	"mime/multipart"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/entity"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util/xerr"
)

type AttachmentHandler struct {
	AttachmentService service.AttachmentService
}

func NewAttachmentHandler(attachmentService service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{
		AttachmentService: attachmentService,
	}
}

// QueryAttachmentInput 附件查询输入。
type QueryAttachmentInput struct {
	Page           int    `query:"page" doc:"页码"`
	Size           int    `query:"size" doc:"每页数量"`
	Keyword        string `query:"keyword" doc:"关键字"`
	MediaType      string `query:"mediaType" doc:"媒体类型"`
	AttachmentType int32  `query:"attachmentType" default:"-1" doc:"附件类型"`
}

// QueryAttachment 分页查询附件列表。
func (a *AttachmentHandler) QueryAttachment(ctx context.Context, in *QueryAttachmentInput) (*dto.HumaOut[*dto.Page], error) {
	queryParam := &param.AttachmentQuery{
		Page: param.Page{
			PageNum:  in.Page,
			PageSize: in.Size,
		},
		Keyword:   in.Keyword,
		MediaType: in.MediaType,
	}
	// huma 不支持指针 query 参数，用 default:"-1" 表示未提供
	if in.AttachmentType >= 0 {
		t := consts.AttachmentType(in.AttachmentType)
		queryParam.AttachmentType = &t
	}
	attachments, totalCount, err := a.AttachmentService.Page(ctx, queryParam)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	attachmentDTOs, err := a.AttachmentService.ConvertToDTOs(ctx, attachments)
	if err != nil {
		return dto.HumaErr[*dto.Page](err)
	}
	return dto.HumaOK(dto.NewPage(attachmentDTOs, totalCount, queryParam.Page))
}

// GetAttachmentByIDInput 按ID获取附件输入。
type GetAttachmentByIDInput struct {
	ID int32 `path:"id" doc:"附件ID"`
}

// GetAttachmentByID 按ID获取附件。
func (a *AttachmentHandler) GetAttachmentByID(ctx context.Context, in *GetAttachmentByIDInput) (*dto.HumaOut[*entity.Attachment], error) {
	if in.ID < 0 {
		return dto.HumaErr[*entity.Attachment](xerr.BadParam.New("id < 0").WithStatus(xerr.StatusBadRequest).WithMsg("param error"))
	}
	attachment, err := a.AttachmentService.GetAttachment(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*entity.Attachment](err)
	}
	return dto.HumaOK(attachment)
}

// UploadAttachmentInput 上传单个附件输入。
type UploadAttachmentInput struct {
	RawBody multipart.Form
}

// UploadAttachment 上传单个附件。
func (a *AttachmentHandler) UploadAttachment(ctx context.Context, in *UploadAttachmentInput) (*dto.HumaOut[*dto.AttachmentDTO], error) {
	files := in.RawBody.File["file"]
	if len(files) == 0 {
		return dto.HumaErr[*dto.AttachmentDTO](xerr.BadParam.New("上传文件错误").WithStatus(xerr.StatusBadRequest))
	}
	attachment, err := a.AttachmentService.Upload(ctx, files[0])
	if err != nil {
		return dto.HumaErr[*dto.AttachmentDTO](err)
	}
	return dto.HumaOK(attachment)
}

// UploadAttachmentsInput 批量上传附件输入。
type UploadAttachmentsInput struct {
	RawBody multipart.Form
}

// UploadAttachments 批量上传附件。
func (a *AttachmentHandler) UploadAttachments(ctx context.Context, in *UploadAttachmentsInput) (*dto.HumaOut[[]*dto.AttachmentDTO], error) {
	files := in.RawBody.File["files"]
	if len(files) == 0 {
		return dto.HumaErr[[]*dto.AttachmentDTO](xerr.BadParam.New("empty files").WithStatus(xerr.StatusBadRequest).WithMsg("empty files"))
	}
	attachmentDTOs := make([]*dto.AttachmentDTO, 0)
	for _, file := range files {
		attachment, err := a.AttachmentService.Upload(ctx, file)
		if err != nil {
			return dto.HumaErr[[]*dto.AttachmentDTO](err)
		}
		attachmentDTOs = append(attachmentDTOs, attachment)
	}
	return dto.HumaOK(attachmentDTOs)
}

// UpdateAttachmentInput 更新附件输入。
type UpdateAttachmentInput struct {
	ID   int32                  `path:"id" doc:"附件ID"`
	Body param.AttachmentUpdate `doc:"附件更新参数"`
}

// UpdateAttachment 更新附件。
func (a *AttachmentHandler) UpdateAttachment(ctx context.Context, in *UpdateAttachmentInput) (*dto.HumaOut[*entity.Attachment], error) {
	attachment, err := a.AttachmentService.Update(ctx, in.ID, &in.Body)
	if err != nil {
		return dto.HumaErr[*entity.Attachment](err)
	}
	return dto.HumaOK(attachment)
}

// DeleteAttachmentInput 按ID删除附件输入。
type DeleteAttachmentInput struct {
	ID int32 `path:"id" doc:"附件ID"`
}

// DeleteAttachment 按ID删除附件。
func (a *AttachmentHandler) DeleteAttachment(ctx context.Context, in *DeleteAttachmentInput) (*dto.HumaOut[*entity.Attachment], error) {
	attachment, err := a.AttachmentService.Delete(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*entity.Attachment](err)
	}
	return dto.HumaOK(attachment)
}

// DeleteAttachmentInBatchInput 批量删除附件输入。
type DeleteAttachmentInBatchInput struct {
	Body []int32 `doc:"附件ID列表"`
}

// DeleteAttachmentInBatch 批量删除附件。
func (a *AttachmentHandler) DeleteAttachmentInBatch(ctx context.Context, in *DeleteAttachmentInBatchInput) (*dto.HumaOut[[]*entity.Attachment], error) {
	attachments, err := a.AttachmentService.DeleteBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[[]*entity.Attachment](err)
	}
	return dto.HumaOK(attachments)
}

// GetAllMediaTypeInput 获取所有媒体类型无输入参数。
type GetAllMediaTypeInput struct{}

// GetAllMediaType 获取所有媒体类型。
func (a *AttachmentHandler) GetAllMediaType(ctx context.Context, _ *GetAllMediaTypeInput) (*dto.HumaOut[[]string], error) {
	mediaTypes, err := a.AttachmentService.GetAllMediaTypes(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(mediaTypes)
}

// GetAllTypesInput 获取所有附件类型无输入参数。
type GetAllTypesInput struct{}

// GetAllTypes 获取所有附件类型。
func (a *AttachmentHandler) GetAllTypes(ctx context.Context, _ *GetAllTypesInput) (*dto.HumaOut[[]consts.AttachmentType], error) {
	attachmentTypes, err := a.AttachmentService.GetAllTypes(ctx)
	if err != nil {
		return dto.HumaErr[[]consts.AttachmentType](err)
	}
	return dto.HumaOK(attachmentTypes)
}
