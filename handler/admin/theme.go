package admin

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type ThemeHandler struct {
	ThemeService  service.ThemeService
	OptionService service.OptionService
}

func NewThemeHandler(l service.ThemeService, o service.OptionService) *ThemeHandler {
	return &ThemeHandler{
		ThemeService:  l,
		OptionService: o,
	}
}

// --- huma handlers (JSON API) ---

// GetThemeByIDInput 按ID获取主题输入。
type GetThemeByIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// GetThemeByID 按ID获取主题。
func (t *ThemeHandler) GetThemeByID(ctx context.Context, in *GetThemeByIDInput) (*dto.HumaOut[*dto.ThemeProperty], error) {
	theme, err := t.ThemeService.GetThemeByID(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[*dto.ThemeProperty](err)
	}
	return dto.HumaOK(theme)
}

// ListAllThemesInput 获取所有主题列表无输入参数。
type ListAllThemesInput struct{}

// ListAllThemes 获取所有主题列表。
func (t *ThemeHandler) ListAllThemes(ctx context.Context, _ *ListAllThemesInput) (*dto.HumaOut[[]*dto.ThemeProperty], error) {
	themes, err := t.ThemeService.ListAllTheme(ctx)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeProperty](err)
	}
	return dto.HumaOK(themes)
}

// ListActivatedThemeFileInput 获取已激活主题文件列表无输入参数。
type ListActivatedThemeFileInput struct{}

// ListActivatedThemeFile 获取已激活主题文件列表。
func (t *ThemeHandler) ListActivatedThemeFile(ctx context.Context, _ *ListActivatedThemeFileInput) (*dto.HumaOut[[]*dto.ThemeFile], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeFile](err)
	}
	files, err := t.ThemeService.ListThemeFiles(ctx, activatedThemeID)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeFile](err)
	}
	return dto.HumaOK(files)
}

// ListThemeFileByIDInput 按ID获取主题文件列表输入。
type ListThemeFileByIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// ListThemeFileByID 按ID获取主题文件列表。
func (t *ThemeHandler) ListThemeFileByID(ctx context.Context, in *ListThemeFileByIDInput) (*dto.HumaOut[[]*dto.ThemeFile], error) {
	files, err := t.ThemeService.ListThemeFiles(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeFile](err)
	}
	return dto.HumaOK(files)
}

// GetThemeFileContentInput 获取已激活主题文件内容输入。
type GetThemeFileContentInput struct {
	Path string `query:"path" required:"true" doc:"文件路径"`
}

// GetThemeFileContent 获取已激活主题文件内容。
func (t *ThemeHandler) GetThemeFileContent(ctx context.Context, in *GetThemeFileContentInput) (*dto.HumaOut[string], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[string](err)
	}
	content, err := t.ThemeService.GetThemeFileContent(ctx, activatedThemeID, in.Path)
	if err != nil {
		return dto.HumaErr[string](err)
	}
	return dto.HumaOK(content)
}

// GetThemeFileContentByIDInput 按ID获取主题文件内容输入。
type GetThemeFileContentByIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
	Path    string `query:"path" required:"true" doc:"文件路径"`
}

// GetThemeFileContentByID 按ID获取主题文件内容。
func (t *ThemeHandler) GetThemeFileContentByID(ctx context.Context, in *GetThemeFileContentByIDInput) (*dto.HumaOut[string], error) {
	content, err := t.ThemeService.GetThemeFileContent(ctx, in.ThemeID, in.Path)
	if err != nil {
		return dto.HumaErr[string](err)
	}
	return dto.HumaOK(content)
}

// UpdateThemeFileInput 更新已激活主题文件内容输入。
type UpdateThemeFileInput struct {
	Body param.ThemeContent `doc:"主题文件内容"`
}

// UpdateThemeFile 更新已激活主题文件内容。
func (t *ThemeHandler) UpdateThemeFile(ctx context.Context, in *UpdateThemeFileInput) (*dto.HumaOut[any], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	err = t.ThemeService.UpdateThemeFile(ctx, activatedThemeID, in.Body.Path, in.Body.Content)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// UpdateThemeFileByIDInput 按ID更新主题文件内容输入。
type UpdateThemeFileByIDInput struct {
	ThemeID string             `path:"themeID" doc:"主题ID"`
	Body    param.ThemeContent `doc:"主题文件内容"`
}

// UpdateThemeFileByID 按ID更新主题文件内容。
func (t *ThemeHandler) UpdateThemeFileByID(ctx context.Context, in *UpdateThemeFileByIDInput) (*dto.HumaOut[any], error) {
	err := t.ThemeService.UpdateThemeFile(ctx, in.ThemeID, in.Body.Path, in.Body.Content)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// ListCustomSheetTemplateInput 获取自定义页面模板列表无输入参数。
type ListCustomSheetTemplateInput struct{}

// ListCustomSheetTemplate 获取自定义页面模板列表。
func (t *ThemeHandler) ListCustomSheetTemplate(ctx context.Context, _ *ListCustomSheetTemplateInput) (*dto.HumaOut[[]string], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	templates, err := t.ThemeService.ListCustomTemplates(ctx, activatedThemeID, consts.ThemeCustomSheetPrefix)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(templates)
}

// ListCustomPostTemplateInput 获取自定义文章模板列表无输入参数。
type ListCustomPostTemplateInput struct{}

// ListCustomPostTemplate 获取自定义文章模板列表。
func (t *ThemeHandler) ListCustomPostTemplate(ctx context.Context, _ *ListCustomPostTemplateInput) (*dto.HumaOut[[]string], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	templates, err := t.ThemeService.ListCustomTemplates(ctx, activatedThemeID, consts.ThemeCustomPostPrefix)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(templates)
}

// ActivateThemeInput 激活主题输入。
type ActivateThemeInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// ActivateTheme 激活主题。
func (t *ThemeHandler) ActivateTheme(ctx context.Context, in *ActivateThemeInput) (*dto.HumaOut[*dto.ThemeProperty], error) {
	theme, err := t.ThemeService.ActivateTheme(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[*dto.ThemeProperty](err)
	}
	return dto.HumaOK(theme)
}

// GetActivatedThemeInput 获取已激活主题无输入参数。
type GetActivatedThemeInput struct{}

// GetActivatedTheme 获取已激活主题。
func (t *ThemeHandler) GetActivatedTheme(ctx context.Context, _ *GetActivatedThemeInput) (*dto.HumaOut[*dto.ThemeProperty], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[*dto.ThemeProperty](err)
	}
	theme, err := t.ThemeService.GetThemeByID(ctx, activatedThemeID)
	if err != nil {
		return dto.HumaErr[*dto.ThemeProperty](err)
	}
	return dto.HumaOK(theme)
}

// GetActivatedThemeConfigInput 获取已激活主题配置无输入参数。
type GetActivatedThemeConfigInput struct{}

// GetActivatedThemeConfig 获取已激活主题配置。
func (t *ThemeHandler) GetActivatedThemeConfig(ctx context.Context, _ *GetActivatedThemeConfigInput) (*dto.HumaOut[[]*dto.ThemeConfigGroup], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeConfigGroup](err)
	}
	config, err := t.ThemeService.GetThemeConfig(ctx, activatedThemeID)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeConfigGroup](err)
	}
	return dto.HumaOK(config)
}

// GetThemeConfigByIDInput 按ID获取主题配置输入。
type GetThemeConfigByIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// GetThemeConfigByID 按ID获取主题配置。
func (t *ThemeHandler) GetThemeConfigByID(ctx context.Context, in *GetThemeConfigByIDInput) (*dto.HumaOut[[]*dto.ThemeConfigGroup], error) {
	config, err := t.ThemeService.GetThemeConfig(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeConfigGroup](err)
	}
	return dto.HumaOK(config)
}

// GetThemeConfigByGroupInput 按分组获取主题配置输入。
type GetThemeConfigByGroupInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
	Group   string `path:"group" doc:"配置分组名称"`
}

// GetThemeConfigByGroup 按分组获取主题配置。
func (t *ThemeHandler) GetThemeConfigByGroup(ctx context.Context, in *GetThemeConfigByGroupInput) (*dto.HumaOut[[]*dto.ThemeConfigItem], error) {
	themeSettings, err := t.ThemeService.GetThemeConfig(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[[]*dto.ThemeConfigItem](err)
	}
	for _, setting := range themeSettings {
		if setting.Name == in.Group {
			return dto.HumaOK(setting.Items)
		}
	}
	return dto.HumaOK[[]*dto.ThemeConfigItem](nil)
}

// GetThemeConfigGroupNamesInput 获取主题配置分组名称列表输入。
type GetThemeConfigGroupNamesInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// GetThemeConfigGroupNames 获取主题配置分组名称列表。
func (t *ThemeHandler) GetThemeConfigGroupNames(ctx context.Context, in *GetThemeConfigGroupNamesInput) (*dto.HumaOut[[]string], error) {
	themeSettings, err := t.ThemeService.GetThemeConfig(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	groupNames := make([]string, len(themeSettings))
	for index, setting := range themeSettings {
		groupNames[index] = setting.Name
	}
	return dto.HumaOK(groupNames)
}

// GetActivatedThemeSettingMapInput 获取已激活主题设置地图无输入参数。
type GetActivatedThemeSettingMapInput struct{}

// GetActivatedThemeSettingMap 获取已激活主题设置地图。
func (t *ThemeHandler) GetActivatedThemeSettingMap(ctx context.Context, _ *GetActivatedThemeSettingMapInput) (*dto.HumaOut[map[string]interface{}], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	settings, err := t.ThemeService.GetThemeSettingMap(ctx, activatedThemeID)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	return dto.HumaOK(settings)
}

// GetThemeSettingMapByIDInput 按ID获取主题设置地图输入。
type GetThemeSettingMapByIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// GetThemeSettingMapByID 按ID获取主题设置地图。
func (t *ThemeHandler) GetThemeSettingMapByID(ctx context.Context, in *GetThemeSettingMapByIDInput) (*dto.HumaOut[map[string]interface{}], error) {
	settings, err := t.ThemeService.GetThemeSettingMap(ctx, in.ThemeID)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	return dto.HumaOK(settings)
}

// GetThemeSettingMapByGroupAndThemeIDInput 按分组获取主题设置地图输入。
type GetThemeSettingMapByGroupAndThemeIDInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
	Group   string `path:"group" doc:"配置分组名称"`
}

// GetThemeSettingMapByGroupAndThemeID 按分组获取主题设置地图。
func (t *ThemeHandler) GetThemeSettingMapByGroupAndThemeID(ctx context.Context, in *GetThemeSettingMapByGroupAndThemeIDInput) (*dto.HumaOut[map[string]interface{}], error) {
	settings, err := t.ThemeService.GetThemeGroupSettingMap(ctx, in.ThemeID, in.Group)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	return dto.HumaOK(settings)
}

// SaveActivatedThemeSettingInput 保存已激活主题设置输入。
type SaveActivatedThemeSettingInput struct {
	Body map[string]interface{} `doc:"主题设置键值对"`
}

// SaveActivatedThemeSetting 保存已激活主题设置。
func (t *ThemeHandler) SaveActivatedThemeSetting(ctx context.Context, in *SaveActivatedThemeSettingInput) (*dto.HumaOut[any], error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	err = t.ThemeService.SaveThemeSettings(ctx, activatedThemeID, in.Body)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// SaveThemeSettingByIDInput 按ID保存主题设置输入。
type SaveThemeSettingByIDInput struct {
	ThemeID string                 `path:"themeID" doc:"主题ID"`
	Body    map[string]interface{} `doc:"主题设置键值对"`
}

// SaveThemeSettingByID 按ID保存主题设置。
func (t *ThemeHandler) SaveThemeSettingByID(ctx context.Context, in *SaveThemeSettingByIDInput) (*dto.HumaOut[any], error) {
	err := t.ThemeService.SaveThemeSettings(ctx, in.ThemeID, in.Body)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// DeleteThemeByIDInput 按ID删除主题输入。
type DeleteThemeByIDInput struct {
	ThemeID        string `path:"themeID" doc:"主题ID"`
	DeleteSettings bool   `query:"deleteSettings" default:"false" doc:"是否同时删除设置"`
}

// DeleteThemeByID 按ID删除主题。
func (t *ThemeHandler) DeleteThemeByID(ctx context.Context, in *DeleteThemeByIDInput) (*dto.HumaOut[any], error) {
	err := t.ThemeService.DeleteTheme(ctx, in.ThemeID, in.DeleteSettings)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// FetchThemeInput 远程拉取主题输入。
type FetchThemeInput struct {
	URI string `query:"uri" doc:"主题下载地址"`
}

// FetchTheme 远程拉取主题。
func (t *ThemeHandler) FetchTheme(ctx context.Context, in *FetchThemeInput) (*dto.HumaOut[*dto.ThemeProperty], error) {
	theme, err := t.ThemeService.Fetch(ctx, in.URI)
	if err != nil {
		return dto.HumaErr[*dto.ThemeProperty](err)
	}
	return dto.HumaOK(theme)
}

// UpdateThemeByFetchingInput 远程更新主题输入。
type UpdateThemeByFetchingInput struct {
	ThemeID string `path:"themeID" doc:"主题ID"`
}

// UpdateThemeByFetching 远程更新主题（未实现）。
func (t *ThemeHandler) UpdateThemeByFetching(ctx context.Context, in *UpdateThemeByFetchingInput) (*dto.HumaOut[any], error) {
	return dto.HumaErr[any](xerr.WithMsg(nil, "not support").WithStatus(xerr.StatusInternalServerError))
}

// ReloadThemeInput 重载主题无输入参数。
type ReloadThemeInput struct{}

// ReloadTheme 重载主题。
func (t *ThemeHandler) ReloadTheme(ctx context.Context, _ *ReloadThemeInput) (*dto.HumaOut[any], error) {
	err := t.ThemeService.ReloadTheme(ctx)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// TemplateExistInput 检查模板是否存在输入。
type TemplateExistInput struct {
	Template string `query:"template" required:"true" doc:"模板名称"`
}

// TemplateExist 检查模板是否存在。
func (t *ThemeHandler) TemplateExist(ctx context.Context, in *TemplateExistInput) (*dto.HumaOut[bool], error) {
	exist, err := t.ThemeService.TemplateExist(ctx, in.Template)
	if err != nil {
		return dto.HumaErr[bool](err)
	}
	return dto.HumaOK(exist)
}

// --- gin handlers (文件上传，保持 gin 不变) ---

func (t *ThemeHandler) UploadTheme(ctx *gin.Context) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UploadTheme(ctx, fileHeader)
}

func (t *ThemeHandler) UpdateThemeByUpload(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UpdateThemeByUpload(ctx, themeID, fileHeader)
}
