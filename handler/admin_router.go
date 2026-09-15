package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

// 本文件把 admin API（/api/admin，需鉴权）的路由注册按 subgroup 拆成独立方法。
// 目的：每个 subgroup 一个函数，迁移到 huma 时可分别修改，避免多文件并行改动时冲突。
// 迁移规则：把函数内的 gin 注册（s.wrapHandler(...)）替换为 huma.Register(api, ...)，
// 并在 handler 中把方法改写为 huma 签名。api 即挂在 authRouter 上的 adminHumaAPI。

// registerAdminBase 注册 authRouter 下无 subgroup 的直接路由。
func (s *Server) registerAdminBase(rg *gin.RouterGroup, api huma.API) {
	rg.POST("/logout", s.wrapHandler(s.AdminHandler.LogOut))
	rg.POST("/password/code", s.wrapHandler(s.AdminHandler.SendResetCode))
	rg.GET("/environments", s.wrapHandler(s.AdminHandler.GetEnvironments))
	rg.GET("/airpress/logfile", s.wrapHandler(s.AdminHandler.GetLogFiles))
}

// registerAdminAttachment 注册 /attachments 相关路由。
func (s *Server) registerAdminAttachment(rg *gin.RouterGroup, api huma.API) {
	attachmentRouter := rg.Group("/attachments")
	attachmentRouter.POST("/upload", s.wrapHandler(s.AttachmentHandler.UploadAttachment))
	attachmentRouter.POST("/uploads", s.wrapHandler(s.AttachmentHandler.UploadAttachments))
	attachmentRouter.DELETE("/:id", s.wrapHandler(s.AttachmentHandler.DeleteAttachment))
	attachmentRouter.DELETE("", s.wrapHandler(s.AttachmentHandler.DeleteAttachmentInBatch))
	attachmentRouter.GET("", s.wrapHandler(s.AttachmentHandler.QueryAttachment))
	attachmentRouter.GET("/:id", s.wrapHandler(s.AttachmentHandler.GetAttachmentByID))
	attachmentRouter.PUT("/:id", s.wrapHandler(s.AttachmentHandler.UpdateAttachment))
	attachmentRouter.GET("/media_types", s.wrapHandler(s.AttachmentHandler.GetAllMediaType))
	attachmentRouter.GET("types", s.wrapHandler(s.AttachmentHandler.GetAllTypes))
}

// registerAdminBackup 注册 /backups 相关路由。
func (s *Server) registerAdminBackup(rg *gin.RouterGroup, api huma.API) {
	backupRouter := rg.Group("/backups")
	backupRouter.POST("/work-dir", s.wrapHandler(s.BackupHandler.BackupWholeSite))
	backupRouter.GET("/work-dir", s.wrapHandler(s.BackupHandler.ListBackups))
	backupRouter.GET("/work-dir/*path", s.BackupHandler.HandleWorkDir)
	backupRouter.DELETE("/work-dir", s.wrapHandler(s.BackupHandler.DeleteBackups))
	backupRouter.POST("/data", s.wrapHandler(s.BackupHandler.ExportData))
	backupRouter.DELETE("/data", s.wrapHandler(s.BackupHandler.DeleteDataFile))
	backupRouter.GET("/data/*path", s.BackupHandler.HandleData)
	backupRouter.POST("/markdown/export", s.wrapHandler(s.BackupHandler.ExportMarkdown))
	backupRouter.POST("/markdown/import", s.wrapHandler(s.BackupHandler.ImportMarkdown))
	backupRouter.GET("/markdown/fetch", s.wrapHandler(s.BackupHandler.GetMarkDownBackup))
	backupRouter.GET("/markdown/export", s.wrapHandler(s.BackupHandler.ListMarkdowns))
	backupRouter.DELETE("/markdown/export", s.wrapHandler(s.BackupHandler.DeleteMarkdowns))
	backupRouter.GET("/markdown/export/:filename", s.BackupHandler.DownloadMarkdown)
}

// registerAdminCategory 注册 /categories 相关路由。
func (s *Server) registerAdminCategory(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/categories/batch", Summary: "批量更新分类", Tags: []string{"admin/categories"}}, s.CategoryHandler.UpdateCategoryBatch)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/categories/{categoryID}", Summary: "获取分类详情", Tags: []string{"admin/categories"}}, s.CategoryHandler.GetCategoryByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/categories", Summary: "获取分类列表", Tags: []string{"admin/categories"}}, s.CategoryHandler.ListAllCategory)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/categories/tree_view", Summary: "获取分类树", Tags: []string{"admin/categories"}}, s.CategoryHandler.ListAsTree)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/categories", Summary: "创建分类", Tags: []string{"admin/categories"}}, s.CategoryHandler.CreateCategory)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/categories/{categoryID}", Summary: "更新分类", Tags: []string{"admin/categories"}}, s.CategoryHandler.UpdateCategory)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/categories/{categoryID}", Summary: "删除分类", Tags: []string{"admin/categories"}}, s.CategoryHandler.DeleteCategory)
}

// registerAdminPost 注册 /posts 及其 /comments 子路由。
func (s *Server) registerAdminPost(rg *gin.RouterGroup, api huma.API) {
	postRouter := rg.Group("/posts")

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts", Summary: "获取文章列表", Tags: []string{"admin/posts"}}, s.PostHandler.ListPosts)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/latest", Summary: "获取最新文章列表", Tags: []string{"admin/posts"}}, s.PostHandler.ListLatestPosts)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/status/{status}", Summary: "按状态获取文章列表", Tags: []string{"admin/posts"}}, s.PostHandler.ListPostsByStatus)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/{postID}", Summary: "获取文章详情", Tags: []string{"admin/posts"}}, s.PostHandler.GetByPostID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/posts", Summary: "创建文章", Tags: []string{"admin/posts"}}, s.PostHandler.CreatePost)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/{postID}", Summary: "更新文章", Tags: []string{"admin/posts"}}, s.PostHandler.UpdatePost)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/{postID}/status/{status}", Summary: "更新文章状态", Tags: []string{"admin/posts"}}, s.PostHandler.UpdatePostStatus)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/status/{status}", Summary: "批量更新文章状态", Tags: []string{"admin/posts"}}, s.PostHandler.UpdatePostStatusBatch)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/{postID}/status/draft/content", Summary: "更新文章草稿内容", Tags: []string{"admin/posts"}}, s.PostHandler.UpdatePostDraft)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/posts/{postID}", Summary: "删除文章", Tags: []string{"admin/posts"}}, s.PostHandler.DeletePost)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/posts", Summary: "批量删除文章", Tags: []string{"admin/posts"}}, s.PostHandler.DeletePostBatch)

	// PreviewPost 返回 HTML/文件流，保持 gin handler
	postRouter.GET("/:postID/preview", s.PostHandler.PreviewPost)

	// /posts/comments 子块
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/comments", Summary: "获取文章评论列表", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.ListPostComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/comments/latest", Summary: "获取最新文章评论", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.ListPostCommentLatest)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/comments/{postID}/tree_view", Summary: "获取文章评论树形结构", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.ListPostCommentAsTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/posts/comments/{postID}/list_view", Summary: "获取文章评论列表（带父评论信息）", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.ListPostCommentWithParent)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/posts/comments", Summary: "创建文章评论", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.CreatePostComment)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/comments/{commentID}", Summary: "更新文章评论", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.UpdatePostComment)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/comments/{commentID}/status/{status}", Summary: "更新文章评论状态", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.UpdatePostCommentStatus)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/posts/comments/status/{status}", Summary: "批量更新文章评论状态", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.UpdatePostCommentStatusBatch)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/posts/comments/{commentID}", Summary: "删除文章评论", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.DeletePostComment)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/posts/comments", Summary: "批量删除文章评论", Tags: []string{"admin/posts/comments"}}, s.PostCommentHandler.DeletePostCommentBatch)
}

// registerAdminOption 注册 /options 相关路由。
func (s *Server) registerAdminOption(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/options", Summary: "获取全部选项列表", Tags: []string{"admin/options"}}, s.OptionHandler.ListAllOptions)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/options/map_view", Summary: "获取选项地图视图", Tags: []string{"admin/options"}}, s.OptionHandler.ListAllOptionsAsMap)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/options/map_view/keys", Summary: "按key获取选项地图", Tags: []string{"admin/options"}}, s.OptionHandler.ListAllOptionsAsMapWithKey)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/options/saving", Summary: "保存选项", Tags: []string{"admin/options"}}, s.OptionHandler.SaveOption)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/options/map_view/saving", Summary: "以地图形式保存选项", Tags: []string{"admin/options"}}, s.OptionHandler.SaveOptionWithMap)
}

// registerAdminLog 注册 /logs 相关路由。
func (s *Server) registerAdminLog(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/logs/latest", Summary: "获取最新日志", Tags: []string{"admin/logs"}}, s.LogHandler.PageLatestLog)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/logs", Summary: "分页查询日志", Tags: []string{"admin/logs"}}, s.LogHandler.PageLog)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/logs/clear", Summary: "清空日志", Tags: []string{"admin/logs"}}, s.LogHandler.ClearLog)
}

// registerAdminStatistic 注册 /statistics 相关路由。
func (s *Server) registerAdminStatistic(rg *gin.RouterGroup, api huma.API) {
	statisticRouter := rg.Group("/statistics")
	statisticRouter.GET("", s.wrapHandler(s.StatisticHandler.Statistics))
	statisticRouter.GET("user", s.wrapHandler(s.StatisticHandler.StatisticsWithUser))
}

// registerAdminSheet 注册 /sheets 及其 /comments 子路由。
func (s *Server) registerAdminSheet(rg *gin.RouterGroup, api huma.API) {
	sheetRouter := rg.Group("/sheets")

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/{sheetID}", Summary: "获取页面详情", Tags: []string{"admin/sheets"}}, s.SheetHandler.GetSheetByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets", Summary: "获取页面列表", Tags: []string{"admin/sheets"}}, s.SheetHandler.ListSheet)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/sheets", Summary: "创建页面", Tags: []string{"admin/sheets"}}, s.SheetHandler.CreateSheet)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/sheets/{sheetID}", Summary: "更新页面", Tags: []string{"admin/sheets"}}, s.SheetHandler.UpdateSheet)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/sheets/{sheetID}/{status}", Summary: "更新页面状态", Tags: []string{"admin/sheets"}}, s.SheetHandler.UpdateSheetStatus)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/sheets/{sheetID}/status/draft/content", Summary: "更新页面草稿内容", Tags: []string{"admin/sheets"}}, s.SheetHandler.UpdateSheetDraft)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/sheets/{sheetID}", Summary: "删除页面", Tags: []string{"admin/sheets"}}, s.SheetHandler.DeleteSheet)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/independent", Summary: "获取独立页面列表", Tags: []string{"admin/sheets"}}, s.SheetHandler.IndependentSheets)

	// PreviewSheet 返回 HTML/文件流，保持 gin handler
	sheetRouter.GET("/preview/:sheetID", s.SheetHandler.PreviewSheet)

	// /sheets/comments 子块
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/comments", Summary: "获取页面评论列表", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.ListSheetComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/comments/latest", Summary: "获取最新页面评论", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.ListSheetCommentLatest)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/comments/{sheetID}/tree_view", Summary: "获取页面评论树形结构", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.ListSheetCommentAsTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/sheets/comments/{sheetID}/list_view", Summary: "获取页面评论列表（带父评论信息）", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.ListSheetCommentWithParent)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/sheets/comments", Summary: "创建页面评论", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.CreateSheetComment)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/sheets/comments/{commentID}/status/{status}", Summary: "更新页面评论状态", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.UpdateSheetCommentStatus)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/sheets/comments/status/{status}", Summary: "批量更新页面评论状态", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.UpdateSheetCommentStatusBatch)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/sheets/comments/{commentID}", Summary: "删除页面评论", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.DeleteSheetComment)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/sheets/comments", Summary: "批量删除页面评论", Tags: []string{"admin/sheets/comments"}}, s.SheetCommentHandler.DeleteSheetCommentBatch)
}

// registerAdminJournal 注册 /journals 及其 /comments 子路由。
func (s *Server) registerAdminJournal(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals", Summary: "获取日志分页列表", Tags: []string{"admin/journals"}}, s.JournalHandler.ListJournal)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/latest", Summary: "获取最新日志列表", Tags: []string{"admin/journals"}}, s.JournalHandler.ListLatestJournal)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/journals", Summary: "创建日志", Tags: []string{"admin/journals"}}, s.JournalHandler.CreateJournal)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/journals/{journalID}", Summary: "更新日志", Tags: []string{"admin/journals"}}, s.JournalHandler.UpdateJournal)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/journals/{journalID}", Summary: "删除日志", Tags: []string{"admin/journals"}}, s.JournalHandler.DeleteJournal)

	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/comments", Summary: "获取日志评论列表", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.ListJournalComment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/comments/latest", Summary: "获取最新日志评论", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.ListJournalCommentLatest)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/comments/{journalID}/tree_view", Summary: "获取日志评论树形列表", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.ListJournalCommentAsTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/journals/comments/{journalID}/list_view", Summary: "获取日志评论带父级列表", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.ListJournalCommentWithParent)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/journals/comments", Summary: "创建日志评论", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.CreateJournalComment)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/journals/comments/{commentID}/status/{status}", Summary: "更新日志评论状态", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.UpdateJournalCommentStatus)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/journals/comments/status/{status}", Summary: "批量更新日志评论状态", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.UpdateJournalStatusBatch)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/journals/comments/{commentID}", Summary: "更新日志评论", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.UpdateJournalComment)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/journals/comments/{commentID}", Summary: "删除日志评论", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.DeleteJournalComment)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/journals/comments", Summary: "批量删除日志评论", Tags: []string{"admin/journals"}}, s.JournalCommentHandler.DeleteJournalCommentBatch)
}

// registerAdminLink 注册 /links 相关路由。
func (s *Server) registerAdminLink(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/links", Summary: "获取链接列表", Tags: []string{"admin/links"}}, s.LinkHandler.ListLinks)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/links/{id}", Summary: "获取链接详情", Tags: []string{"admin/links"}}, s.LinkHandler.GetLinkByID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/links", Summary: "创建链接", Tags: []string{"admin/links"}}, s.LinkHandler.CreateLink)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/links/{id}", Summary: "更新链接", Tags: []string{"admin/links"}}, s.LinkHandler.UpdateLink)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/links/{id}", Summary: "删除链接", Tags: []string{"admin/links"}}, s.LinkHandler.DeleteLink)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/links/teams", Summary: "获取链接团队列表", Tags: []string{"admin/links"}}, s.LinkHandler.ListLinkTeams)
}

// registerAdminMenu 注册 /menus 相关路由。
func (s *Server) registerAdminMenu(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/menus", Summary: "获取菜单列表", Tags: []string{"admin/menus"}}, s.MenuHandler.ListMenus)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/menus/tree_view", Summary: "获取菜单树", Tags: []string{"admin/menus"}}, s.MenuHandler.ListMenusAsTree)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/menus/team/tree_view", Summary: "按团队获取菜单树", Tags: []string{"admin/menus"}}, s.MenuHandler.ListMenusAsTreeByTeam)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/menus/{id}", Summary: "获取菜单详情", Tags: []string{"admin/menus"}}, s.MenuHandler.GetMenuByID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/menus", Summary: "创建菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.CreateMenu)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/menus/batch", Summary: "批量创建菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.CreateMenuBatch)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/menus/{id}", Summary: "更新菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.UpdateMenu)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/menus/batch", Summary: "批量更新菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.UpdateMenuBatch)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/menus/{id}", Summary: "删除菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.DeleteMenu)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/menus/batch", Summary: "批量删除菜单", Tags: []string{"admin/menus"}}, s.MenuHandler.DeleteMenuBatch)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/menus/teams", Summary: "获取菜单团队列表", Tags: []string{"admin/menus"}}, s.MenuHandler.ListMenuTeams)
}

// registerAdminTag 注册 /tags 相关路由。
func (s *Server) registerAdminTag(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/tags", Summary: "获取标签列表", Tags: []string{"admin/tags"}}, s.TagHandler.ListTags)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/tags/{id}", Summary: "获取标签详情", Tags: []string{"admin/tags"}}, s.TagHandler.GetTagByID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/tags", Summary: "创建标签", Tags: []string{"admin/tags"}}, s.TagHandler.CreateTag)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/tags/{id}", Summary: "更新标签", Tags: []string{"admin/tags"}}, s.TagHandler.UpdateTag)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/tags/{id}", Summary: "删除标签", Tags: []string{"admin/tags"}}, s.TagHandler.DeleteTag)
}

// registerAdminPhoto 注册 /photos 相关路由。
func (s *Server) registerAdminPhoto(rg *gin.RouterGroup, api huma.API) {
	photoRouter := rg.Group("/photos")
	photoRouter.GET("/latest", s.wrapHandler(s.PhotoHandler.ListPhoto))
	photoRouter.GET("", s.wrapHandler(s.PhotoHandler.PagePhotos))
	photoRouter.GET("/:id", s.wrapHandler(s.PhotoHandler.GetPhotoByID))
	photoRouter.DELETE("/batch", s.wrapHandler(s.PhotoHandler.DeletePhotoBatch))
	photoRouter.POST("", s.wrapHandler(s.PhotoHandler.CreatePhoto))
	photoRouter.POST("/batch", s.wrapHandler(s.PhotoHandler.CreatePhotoBatch))
	photoRouter.PUT("/:id", s.wrapHandler(s.PhotoHandler.UpdatePhoto))
	photoRouter.GET("/teams", s.wrapHandler(s.PhotoHandler.ListPhotoTeams))
}

// registerAdminUser 注册 /users 相关路由。
func (s *Server) registerAdminUser(rg *gin.RouterGroup, api huma.API) {
	userRouter := rg.Group("/users")
	userRouter.GET("/profiles", s.wrapHandler(s.UserHandler.GetCurrentUserProfile))
	userRouter.PUT("/profiles", s.wrapHandler(s.UserHandler.UpdateUserProfile))
	userRouter.PUT("/profiles/password", s.wrapHandler(s.UserHandler.UpdatePassword))
	userRouter.PUT("/mfa/generate", s.wrapHandler(s.UserHandler.GenerateMFAQRCode))
	userRouter.PUT("/mfa/update", s.wrapHandler(s.UserHandler.UpdateMFA))
}

// registerAdminTheme 注册 themes 相关路由（原 group 路径为 "themes"，无前导斜杠，保持原样）。
func (s *Server) registerAdminTheme(rg *gin.RouterGroup, api huma.API) {
	themeRouter := rg.Group("themes")
	themeRouter.GET("/activation", s.wrapHandler(s.ThemeHandler.GetActivatedTheme))
	themeRouter.GET("/:themeID", s.wrapHandler(s.ThemeHandler.GetThemeByID))
	themeRouter.GET("", s.wrapHandler(s.ThemeHandler.ListAllThemes))
	themeRouter.GET("/activation/files", s.wrapHandler(s.ThemeHandler.ListActivatedThemeFile))
	themeRouter.GET("/:themeID/files", s.wrapHandler(s.ThemeHandler.ListThemeFileByID))
	themeRouter.GET("files/content", s.wrapHandler(s.ThemeHandler.GetThemeFileContent))
	themeRouter.GET("/:themeID/files/content", s.wrapHandler(s.ThemeHandler.GetThemeFileContentByID))
	themeRouter.PUT("/files/content", s.wrapHandler(s.ThemeHandler.UpdateThemeFile))
	themeRouter.PUT("/:themeID/files/content", s.wrapHandler(s.ThemeHandler.UpdateThemeFileByID))
	themeRouter.GET("activation/template/custom/sheet", s.wrapHandler(s.ThemeHandler.ListCustomSheetTemplate))
	themeRouter.GET("activation/template/custom/post", s.wrapHandler(s.ThemeHandler.ListCustomPostTemplate))
	themeRouter.POST("/:themeID/activation", s.wrapHandler(s.ThemeHandler.ActivateTheme))
	themeRouter.GET("activation/configurations", s.wrapHandler(s.ThemeHandler.GetActivatedThemeConfig))
	themeRouter.GET("/:themeID/configurations", s.wrapHandler(s.ThemeHandler.GetThemeConfigByID))
	themeRouter.GET("/:themeID/configurations/groups/:group", s.wrapHandler(s.ThemeHandler.GetThemeConfigByGroup))
	themeRouter.GET("/:themeID/configurations/groups", s.wrapHandler(s.ThemeHandler.GetThemeConfigGroupNames))
	themeRouter.GET("activation/settings", s.wrapHandler(s.ThemeHandler.GetActivatedThemeSettingMap))
	themeRouter.GET("/:themeID/settings", s.wrapHandler(s.ThemeHandler.GetThemeSettingMapByID))
	themeRouter.GET("/:themeID/groups/:group/settings", s.wrapHandler(s.ThemeHandler.GetThemeSettingMapByGroupAndThemeID))
	themeRouter.POST("activation/settings", s.wrapHandler(s.ThemeHandler.SaveActivatedThemeSetting))
	themeRouter.POST("/:themeID/settings", s.wrapHandler(s.ThemeHandler.SaveThemeSettingByID))
	themeRouter.DELETE("/:themeID", s.wrapHandler(s.ThemeHandler.DeleteThemeByID))
	themeRouter.POST("upload", s.wrapHandler(s.ThemeHandler.UploadTheme))
	themeRouter.PUT("upload/:themeID", s.wrapHandler(s.ThemeHandler.UpdateThemeByUpload))
	themeRouter.POST("fetching", s.wrapHandler(s.ThemeHandler.FetchTheme))
	themeRouter.PUT("fetching/:themeID", s.wrapHandler(s.ThemeHandler.UpdateThemeByFetching))
	themeRouter.POST("reload", s.wrapHandler(s.ThemeHandler.ReloadTheme))
	themeRouter.GET("activation/template/exists", s.wrapHandler(s.ThemeHandler.TemplateExist))
}

// registerAdminEmail 注册 /mails 相关路由。
func (s *Server) registerAdminEmail(rg *gin.RouterGroup, api huma.API) {
	emailRouter := rg.Group("/mails")
	emailRouter.POST("/test", s.wrapHandler(s.EmailHandler.Test))
}
