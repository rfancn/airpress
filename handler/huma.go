package handler

import (
	"net/http"
	"path"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// schemaNamer 在 huma 默认命名基础上加包名前缀，避免跨包同名类型冲突。
// 项目里 param/dto/vo/entity 包存在大量同名类型（如 Comment、Post），
// 而 huma 默认 namer 只取简单类型名会导致 "duplicate name" panic。
// 例：param.Comment -> ParamComment，dto.Comment -> DtoComment。
func schemaNamer(t reflect.Type, hint string) string {
	base := t
	for base.Kind() == reflect.Pointer || base.Kind() == reflect.Slice || base.Kind() == reflect.Array {
		base = base.Elem()
	}
	name := huma.DefaultSchemaNamer(t, hint)
	if base.Name() == "" {
		return name
	}
	pkg := path.Base(base.PkgPath())
	if pkg == "" || pkg == "." {
		return name
	}
	return upperFirst(pkg) + name
}

// upperFirst 将字符串首字母转为大写。
func upperFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 'a' - 'A'
	}
	return string(r)
}

// newHumaAPI 在指定的 gin group 上创建 huma API。
// 由于使用 humagin.NewWithGroup，该 group 的 gin 中间件（如鉴权、日志）对 huma 路由自动生效。
// 返回值包含 group 的 basePath 前缀，用于在合并文档时修正 OpenAPI 路径（huma 的 op.Path 不含 group 前缀）。
func newHumaAPI(engine *gin.Engine, group *gin.RouterGroup, title string) (huma.API, string) {
	cfg := huma.DefaultConfig(title, "1.0.0")
	// 用带包名前缀的自定义 registry，避免跨包同名类型的 schema 命名冲突
	cfg.Components = &huma.Components{
		Schemas: huma.NewMapRegistry("#/components/schemas/", schemaNamer),
	}
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

// registerContentHumaAPI 注册 content API（/api/content/*，公开）的 huma 路由。
// 独立成方法以便测试直接调用验证（RegisterRouters 依赖完整 DI，难以单测）。
func (s *Server) registerContentHumaAPI(api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/archives/years", Summary: "按年份归档", Tags: []string{"content/archives"}}, s.ContentAPIArchiveHandler.ListYearArchives)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/archives/months", Summary: "按月归档", Tags: []string{"content/archives"}}, s.ContentAPIArchiveHandler.ListMonthArchives)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/categories", Summary: "获取分类列表", Tags: []string{"content/categories"}}, s.ContentAPICategoryHandler.ListCategories)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/categories/{slug}/posts", Summary: "获取分类下的文章列表", Tags: []string{"content/categories"}}, s.ContentAPICategoryHandler.ListPosts)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/links", Summary: "获取链接列表", Tags: []string{"content/links"}}, s.ContentAPILinkHandler.ListLinks)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/links/team_view", Summary: "获取链接团队视图", Tags: []string{"content/links"}}, s.ContentAPILinkHandler.LinkTeamVO)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/options/comment", Summary: "获取评论选项", Tags: []string{"content/options"}}, s.ContentAPIOptionHandler.Comment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/comments/{commentID}/likes", Summary: "评论点赞", Tags: []string{"content/comments"}}, s.ContentAPICommentHandler.Like)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals", Summary: "获取日志列表", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.ListJournal)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/{journalID}", Summary: "获取日志详情", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.GetJournal)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/{journalID}/comments/top_view", Summary: "获取顶级评论", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.ListTopComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/{journalID}/comments/{parentID}/children", Summary: "获取子评论", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.ListChildren)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/{journalID}/comments/tree_view", Summary: "获取评论树", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.ListCommentTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/{journalID}/comments/list_view", Summary: "获取评论列表", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.ListComment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/journals/comments", Summary: "创建评论", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.CreateComment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/journals/{journalID}/likes", Summary: "日志点赞", Tags: []string{"content/journals"}}, s.ContentAPIJournalHandler.Like)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/photos/{photoID}/likes", Summary: "图片点赞", Tags: []string{"content/photos"}}, s.ContentAPIPhotoHandler.Like)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/{postID}/comments/top_view", Summary: "获取文章顶级评论列表", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.ListTopComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/{postID}/comments/{parentID}/children", Summary: "获取文章评论子列表", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.ListChildren)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/{postID}/comments/tree_view", Summary: "获取文章评论树形列表", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.ListCommentTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/{postID}/comments/list_view", Summary: "获取文章评论列表", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.ListComment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/posts/comments", Summary: "创建文章评论", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.CreateComment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/posts/{postID}/likes", Summary: "点赞文章", Tags: []string{"content/posts"}}, s.ContentAPIPostHandler.Like)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/{sheetID}/comments/top_view", Summary: "获取页面顶级评论列表", Tags: []string{"content/sheets"}}, s.ContentAPISheetHandler.ListTopComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/{sheetID}/comments/{parentID}/children", Summary: "获取页面评论子列表", Tags: []string{"content/sheets"}}, s.ContentAPISheetHandler.ListChildren)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/{sheetID}/comments/tree_view", Summary: "获取页面评论树形列表", Tags: []string{"content/sheets"}}, s.ContentAPISheetHandler.ListCommentTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/{sheetID}/comments/list_view", Summary: "获取页面评论列表", Tags: []string{"content/sheets"}}, s.ContentAPISheetHandler.ListComment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/sheets/comments", Summary: "创建页面评论", Tags: []string{"content/sheets"}}, s.ContentAPISheetHandler.CreateComment)
}
