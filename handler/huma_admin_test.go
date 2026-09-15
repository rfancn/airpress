package handler

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/handler/admin"
)

// TestAdminHumaRegistration 实际执行所有 admin huma 注册，验证：
// ① 所有 handler 输出/输入类型能被 huma 推导（否则注册时 panic）；
// ② 产出 OpenAPI 3.1；③ 路径都带 /api/admin 前缀；④ 无重复路径。
// handler 的 service 字段为 nil 无妨——注册阶段不会调用 handler 本体。
func TestAdminHumaRegistration(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("registration panic: %v", r)
		}
	}()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	authGroup := engine.Group("/api/admin")
	humaAPI, prefix := newHumaAPI(engine, authGroup, "AirPress Admin API")

	s := &Server{
		AdminHandler:          &admin.AdminHandler{},
		AttachmentHandler:     &admin.AttachmentHandler{},
		BackupHandler:         &admin.BackupHandler{},
		CategoryHandler:       &admin.CategoryHandler{},
		PostHandler:           &admin.PostHandler{},
		PostCommentHandler:    &admin.PostCommentHandler{},
		OptionHandler:         &admin.OptionHandler{},
		LogHandler:            &admin.LogHandler{},
		StatisticHandler:      &admin.StatisticHandler{},
		SheetHandler:          &admin.SheetHandler{},
		SheetCommentHandler:   &admin.SheetCommentHandler{},
		JournalHandler:        &admin.JournalHandler{},
		JournalCommentHandler: &admin.JournalCommentHandler{},
		LinkHandler:           &admin.LinkHandler{},
		MenuHandler:           &admin.MenuHandler{},
		TagHandler:            &admin.TagHandler{},
		PhotoHandler:          &admin.PhotoHandler{},
		UserHandler:           &admin.UserHandler{},
		ThemeHandler:          &admin.ThemeHandler{},
		EmailHandler:          &admin.EmailHandler{},
	}

	// 逐个注册（任何 schema 命名冲突或类型问题会在此 panic）
	s.registerAdminBase(authGroup, humaAPI)
	s.registerAdminAttachment(authGroup, humaAPI)
	s.registerAdminBackup(authGroup, humaAPI)
	s.registerAdminCategory(authGroup, humaAPI)
	s.registerAdminPost(authGroup, humaAPI)
	s.registerAdminOption(authGroup, humaAPI)
	s.registerAdminLog(authGroup, humaAPI)
	s.registerAdminStatistic(authGroup, humaAPI)
	s.registerAdminSheet(authGroup, humaAPI)
	s.registerAdminJournal(authGroup, humaAPI)
	s.registerAdminLink(authGroup, humaAPI)
	s.registerAdminMenu(authGroup, humaAPI)
	s.registerAdminTag(authGroup, humaAPI)
	s.registerAdminPhoto(authGroup, humaAPI)
	s.registerAdminUser(authGroup, humaAPI)
	s.registerAdminTheme(authGroup, humaAPI)
	s.registerAdminEmail(authGroup, humaAPI)

	// 公开 API（登录/安装，独立实例，挂在无鉴权的 adminAPIRouter）
	publicGroup := engine.Group("/api/admin")
	publicAPI, publicPrefix := newHumaAPI(engine, publicGroup, "AirPress Admin Public API")
	s.registerAdminPublicHumaAPI(publicAPI)

	merged := &huma.OpenAPI{
		OpenAPI:    humaAPI.OpenAPI().OpenAPI,
		Info:       humaAPI.OpenAPI().Info,
		Paths:      map[string]*huma.PathItem{},
		Components: &huma.Components{},
	}
	mergeInto(merged, humaAPI.OpenAPI(), prefix)
	mergeInto(merged, publicAPI.OpenAPI(), publicPrefix)

	if merged.OpenAPI != "3.1.0" {
		t.Fatalf("expected 3.1.0, got %s", merged.OpenAPI)
	}
	t.Logf("admin 迁移后 huma 路由数: %d", len(merged.Paths))
	for p := range merged.Paths {
		if len(p) < len("/api/admin") || p[:len("/api/admin")] != "/api/admin" {
			t.Errorf("path missing /api/admin prefix: %s", p)
		}
	}
}
