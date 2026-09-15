package handler

import (
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/handler/content/api"
)

// TestContentHumaRegistration 实际执行 content API 的 huma 注册，
// 验证：① 所有 handler 输出类型能被 huma 推导（否则注册时 panic）；
// ② 产出 OpenAPI 3.1；③ 28 条路径且都带 /api/content 前缀。
// handler 的 service 字段为 nil 无妨——注册阶段不会调用 handler 本体。
func TestContentHumaRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	contentGroup := engine.Group("/api/content")
	humaAPI, prefix := newHumaAPI(engine, contentGroup, "AirPress Content API")

	s := &Server{
		ContentAPIArchiveHandler:  &api.ArchiveHandler{},
		ContentAPICategoryHandler: &api.CategoryHandler{},
		ContentAPIJournalHandler:  &api.JournalHandler{},
		ContentAPILinkHandler:     &api.LinkHandler{},
		ContentAPIPostHandler:     &api.PostHandler{},
		ContentAPISheetHandler:    &api.SheetHandler{},
		ContentAPIOptionHandler:   &api.OptionHandler{},
		ContentAPIPhotoHandler:    &api.PhotoHandler{},
		ContentAPICommentHandler:  &api.CommentHandler{},
	}

	s.registerContentHumaAPI(humaAPI) // panic 即代表某 handler 输出类型无法推导

	merged := &huma.OpenAPI{
		OpenAPI:    humaAPI.OpenAPI().OpenAPI,
		Info:       humaAPI.OpenAPI().Info,
		Paths:      map[string]*huma.PathItem{},
		Components: &huma.Components{},
	}
	mergeInto(merged, humaAPI.OpenAPI(), prefix)

	if merged.OpenAPI != "3.1.0" {
		t.Fatalf("expected 3.1.0, got %s", merged.OpenAPI)
	}
	if len(merged.Paths) != 28 {
		t.Fatalf("expected 28 content paths, got %d: %v", len(merged.Paths), keys(merged.Paths))
	}
	for p := range merged.Paths {
		if !strings.HasPrefix(p, "/api/content") {
			t.Errorf("path missing /api/content prefix: %s", p)
		}
	}
}

func keys(m map[string]*huma.PathItem) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
