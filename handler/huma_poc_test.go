package handler

import (
	"encoding/json"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

// TestHumaPoC 验证 PoC 关键机制（不依赖数据库）：
// 1. huma 能从泛型 handler 推导 schema 且不 panic
// 2. 产出 OpenAPI 3.1
// 3. 合并后的文档路径是否包含 gin group 前缀
func TestHumaPoC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	adminGroup := engine.Group("/api/admin")
	contentGroup := engine.Group("/api/content")

	adminAPI, adminPrefix := newHumaAPI(engine, adminGroup, "Admin")
	contentAPI, contentPrefix := newHumaAPI(engine, contentGroup, "Content")

	// handler 方法值允许在 nil 接收者上获取（仅注册阶段不调用，schema 由类型静态推导）
	var s Server
	huma.Register(adminAPI, huma.Operation{Method: "GET", Path: "/posts/{postID}", Tags: []string{"t"}}, s.PostHandler.GetByPostID)
	huma.Register(contentAPI, huma.Operation{Method: "GET", Path: "/archives/years", Tags: []string{"t"}}, s.ContentAPIArchiveHandler.ListYearArchives)

	merged := &huma.OpenAPI{
		OpenAPI:    adminAPI.OpenAPI().OpenAPI,
		Info:       adminAPI.OpenAPI().Info,
		Paths:      map[string]*huma.PathItem{},
		Components: &huma.Components{},
	}
	mergeInto(merged, adminAPI.OpenAPI(), adminPrefix)
	mergeInto(merged, contentAPI.OpenAPI(), contentPrefix)

	if merged.OpenAPI != "3.1.0" {
		t.Fatalf("expected 3.1.0, got %s", merged.OpenAPI)
	}
	t.Logf("openapi version: %s", merged.OpenAPI)
	for k := range merged.Paths {
		t.Logf("path: %s", k)
	}
	// 校验路径已带上 group 前缀
	if _, ok := merged.Paths["/api/admin/posts/{postID}"]; !ok {
		t.Fatalf("missing prefixed admin path")
	}
	if _, ok := merged.Paths["/api/content/archives/years"]; !ok {
		t.Fatalf("missing prefixed content path")
	}
	// 打印生成的 schema，确认 BaseDTO[T] 的 data 字段有具体结构（而非 any）
	if merged.Components != nil && merged.Components.Schemas != nil {
		for k := range merged.Components.Schemas.Map() {
			t.Logf("schema: %s", k)
		}
	}
	b, err := json.Marshal(merged)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("openapi json:\n%s", string(b))
}
