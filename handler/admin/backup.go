package admin

import (
	"context"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/config"
	"github.com/rfancn/airpress/log"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type BackupHandler struct {
	BackupService service.BackupService
}

func NewBackupHandler(backupService service.BackupService) *BackupHandler {
	return &BackupHandler{
		BackupService: backupService,
	}
}

// --- huma handlers (JSON API) ---

// BackupWholeSiteInput 全站备份输入。
type BackupWholeSiteInput struct {
	Body []string `doc:"待备份项列表"`
}

// BackupWholeSite 全站备份。
func (b *BackupHandler) BackupWholeSite(ctx context.Context, in *BackupWholeSiteInput) (*dto.HumaOut[*dto.BackupDTO], error) {
	backupDTO, err := b.BackupService.BackupWholeSite(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[*dto.BackupDTO](err)
	}
	return dto.HumaOK(backupDTO)
}

// ListBackupsInput 列出全站备份无输入参数。
type ListBackupsInput struct{}

// ListBackups 列出全站备份。
func (b *BackupHandler) ListBackups(ctx context.Context, _ *ListBackupsInput) (*dto.HumaOut[[]*dto.BackupDTO], error) {
	backups, err := b.BackupService.ListFiles(ctx, config.BackupDir, service.WholeSite)
	if err != nil {
		return dto.HumaErr[[]*dto.BackupDTO](err)
	}
	return dto.HumaOK(backups)
}

// DeleteBackupsInput 删除全站备份输入。
type DeleteBackupsInput struct {
	Filename string `query:"filename" required:"true" doc:"文件名"`
}

// DeleteBackups 删除全站备份。
func (b *BackupHandler) DeleteBackups(ctx context.Context, in *DeleteBackupsInput) (*dto.HumaOut[interface{}], error) {
	err := b.BackupService.DeleteFile(ctx, config.BackupDir, in.Filename)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// ExportDataInput 导出数据无输入参数。
type ExportDataInput struct{}

// ExportData 导出数据。
func (b *BackupHandler) ExportData(ctx context.Context, _ *ExportDataInput) (*dto.HumaOut[*dto.BackupDTO], error) {
	backupDTO, err := b.BackupService.ExportData(ctx)
	if err != nil {
		return dto.HumaErr[*dto.BackupDTO](err)
	}
	return dto.HumaOK(backupDTO)
}

// DeleteDataFileInput 删除数据文件输入。
type DeleteDataFileInput struct {
	Filename string `query:"filename" required:"true" doc:"文件名"`
}

// DeleteDataFile 删除数据文件。
func (b *BackupHandler) DeleteDataFile(ctx context.Context, in *DeleteDataFileInput) (*dto.HumaOut[interface{}], error) {
	err := b.BackupService.DeleteFile(ctx, config.DataExportDir, in.Filename)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// ExportMarkdownInput 导出 Markdown 输入。
type ExportMarkdownInput struct {
	Body param.ExportMarkdown `doc:"Markdown 导出参数"`
}

// ExportMarkdown 导出 Markdown。
func (b *BackupHandler) ExportMarkdown(ctx context.Context, in *ExportMarkdownInput) (*dto.HumaOut[*dto.BackupDTO], error) {
	backupDTO, err := b.BackupService.ExportMarkdown(ctx, in.Body.NeedFrontMatter)
	if err != nil {
		return dto.HumaErr[*dto.BackupDTO](err)
	}
	return dto.HumaOK(backupDTO)
}

// ImportMarkdownInput 导入 Markdown 输入。
type ImportMarkdownInput struct {
	RawBody multipart.Form
}

// ImportMarkdown 导入 Markdown 文件。
func (b *BackupHandler) ImportMarkdown(ctx context.Context, in *ImportMarkdownInput) (*dto.HumaOut[interface{}], error) {
	files := in.RawBody.File["file"]
	if len(files) == 0 {
		return dto.HumaErr[interface{}](xerr.BadParam.New("上传文件错误").WithStatus(xerr.StatusBadRequest))
	}
	fileHeader := files[0]
	filenameExt := path.Ext(fileHeader.Filename)
	if filenameExt != ".md" && filenameExt != ".markdown" && filenameExt != ".mdown" {
		return dto.HumaErr[interface{}](xerr.BadParam.New("Unsupported format").WithStatus(xerr.StatusBadRequest))
	}
	err := b.BackupService.ImportMarkdown(ctx, fileHeader)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// GetMarkDownBackupInput 获取 Markdown 备份输入。
type GetMarkDownBackupInput struct {
	Filename string `query:"filename" required:"true" doc:"文件名"`
}

// GetMarkDownBackup 获取 Markdown 备份。
func (b *BackupHandler) GetMarkDownBackup(ctx context.Context, in *GetMarkDownBackupInput) (*dto.HumaOut[*dto.BackupDTO], error) {
	backupDTO, err := b.BackupService.GetBackup(ctx, filepath.Join(config.BackupMarkdownDir, in.Filename), service.Markdown)
	if err != nil {
		return dto.HumaErr[*dto.BackupDTO](err)
	}
	return dto.HumaOK(backupDTO)
}

// ListMarkdownsInput 列出 Markdown 备份无输入参数。
type ListMarkdownsInput struct{}

// ListMarkdowns 列出 Markdown 备份。
func (b *BackupHandler) ListMarkdowns(ctx context.Context, _ *ListMarkdownsInput) (*dto.HumaOut[[]*dto.BackupDTO], error) {
	backups, err := b.BackupService.ListFiles(ctx, config.BackupMarkdownDir, service.Markdown)
	if err != nil {
		return dto.HumaErr[[]*dto.BackupDTO](err)
	}
	return dto.HumaOK(backups)
}

// DeleteMarkdownsInput 删除 Markdown 备份输入。
type DeleteMarkdownsInput struct {
	Filename string `query:"filename" required:"true" doc:"文件名"`
}

// DeleteMarkdowns 删除 Markdown 备份。
func (b *BackupHandler) DeleteMarkdowns(ctx context.Context, in *DeleteMarkdownsInput) (*dto.HumaOut[interface{}], error) {
	err := b.BackupService.DeleteFile(ctx, config.BackupMarkdownDir, in.Filename)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// --- gin handlers (文件流，保持 gin 不变) ---

func (b *BackupHandler) GetWorkDirBackup(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(ctx, filepath.Join(config.BackupDir, filename), service.WholeSite)
}

func (b *BackupHandler) GetDataBackup(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(ctx, filepath.Join(config.DataExportDir, filename), service.JSONData)
}

func (b *BackupHandler) ListToBackupItems(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListToBackupItems(ctx)
}

func (b *BackupHandler) ListExportData(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListFiles(ctx, config.DataExportDir, service.JSONData)
}

func (b *BackupHandler) HandleWorkDir(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if path == "/api/admin/backups/work-dir/fetch" {
		wrapHandler(b.GetWorkDirBackup)(ctx)
		return
	}
	if path == "/api/admin/backups/work-dir/options" || path == "/api/admin/backups/work-dir/options/" {
		wrapHandler(b.ListToBackupItems)(ctx)
		return
	}
	b.DownloadBackups(ctx)
}

func (b *BackupHandler) DownloadBackups(ctx *gin.Context) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO[any]{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.BackupDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO[any]{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) HandleData(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if path == "/api/admin/backups/data/fetch" {
		wrapHandler(b.GetDataBackup)(ctx)
		return
	}
	if path == "/api/admin/backups/data" || path == "/api/admin/backups/data/" {
		wrapHandler(b.ListExportData)(ctx)
		return
	}
	b.DownloadData(ctx)
}

func (b *BackupHandler) DownloadData(ctx *gin.Context) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO[any]{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.DataExportDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO[any]{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) DownloadMarkdown(ctx *gin.Context) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO[any]{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.BackupMarkdownDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO[any]{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

type wrapperHandler func(ctx *gin.Context) (interface{}, error)

func wrapHandler(handler wrapperHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			log.CtxErrorf(ctx, "err=%+v", err)
			status := xerr.GetHTTPStatus(err)
			ctx.JSON(status, &dto.BaseDTO[any]{Status: status, Message: xerr.GetMessage(err)})
			return
		}

		ctx.JSON(http.StatusOK, &dto.BaseDTO[any]{
			Status:  http.StatusOK,
			Data:    data,
			Message: "OK",
		})
	}
}
