package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// newHumaAPI 在指定的 gin group 上创建 huma API。
// 由于使用 humagin.NewWithGroup，该 group 的 gin 中间件（如鉴权、日志）对 huma 路由自动生效。
// 返回值包含 group 的 basePath 前缀，用于在合并文档时修正 OpenAPI 路径（huma 的 op.Path 不含 group 前缀）。
func newHumaAPI(engine *gin.Engine, group *gin.RouterGroup, title string) (huma.API, string) {
	cfg := huma.DefaultConfig(title, "1.0.0")
	// 关闭 huma 自带的 /openapi 与 /docs 路由，改由合并后统一注册，避免多个 API 实例冲突。
	cfg.OpenAPIPath = ""
	cfg.DocsPath = ""
	return humagin.NewWithGroup(engine, group, cfg), group.BasePath()
}

// mergeInto 把 src 的路径（加上 prefix）与组件 schema 合并进 dest。
// dest 必须是独立对象，严禁复用某个 api.OpenAPI()，否则会原地自拷贝导致前缀重复拼接。
func mergeInto(dest, src *huma.OpenAPI, prefix string) {
	if dest.Paths == nil {
		dest.Paths = map[string]*huma.PathItem{}
	}
	for path, item := range src.Paths {
		dest.Paths[prefix+path] = item
	}
	if src.Components != nil && src.Components.Schemas != nil {
		if dest.Components == nil {
			dest.Components = &huma.Components{}
		}
		if dest.Components.Schemas != nil {
			for k, v := range src.Components.Schemas.Map() {
				dest.Components.Schemas.Map()[k] = v
			}
		} else {
			// schema registry 只读共享即可，无需深拷贝
			dest.Components.Schemas = src.Components.Schemas
		}
	}
}

// registerHumaDocs 将多个 huma API 的 OpenAPI 3.1 文档合并（修正路径前缀后），注册 /openapi.json 输出。
func (s *Server) registerHumaDocs(apis []huma.API, prefixes []string) {
	if len(apis) == 0 {
		return
	}
	first := apis[0].OpenAPI()
	merged := &huma.OpenAPI{
		OpenAPI:    first.OpenAPI,
		Info:       first.Info,
		Paths:      map[string]*huma.PathItem{},
		Components: &huma.Components{},
	}
	for i, api := range apis {
		mergeInto(merged, api.OpenAPI(), prefixes[i])
	}
	s.Router.GET("/openapi.json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, merged)
	})
}
