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
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/logout", Summary: "管理员登出", Tags: []string{"admin/base"}}, s.AdminHandler.LogOut)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/password/code", Summary: "发送重置密码验证码", Tags: []string{"admin/base"}}, s.AdminHandler.SendResetCode)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/environments", Summary: "获取环境信息", Tags: []string{"admin/base"}}, s.AdminHandler.GetEnvironments)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/airpress/logfile", Summary: "获取日志文件", Tags: []string{"admin/base"}}, s.AdminHandler.GetLogFiles)
}

// registerAdminAttachment 注册 /attachments 相关路由。
func (s *Server) registerAdminAttachment(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/attachments/upload", Summary: "上传单个附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.UploadAttachment)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/attachments/uploads", Summary: "批量上传附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.UploadAttachments)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/attachments/{id}", Summary: "按ID删除附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.DeleteAttachment)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/attachments", Summary: "批量删除附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.DeleteAttachmentInBatch)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/attachments", Summary: "分页查询附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.QueryAttachment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/attachments/{id}", Summary: "按ID获取附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.GetAttachmentByID)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/attachments/{id}", Summary: "更新附件", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.UpdateAttachment)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/attachments/media_types", Summary: "获取所有媒体类型", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.GetAllMediaType)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/attachments/types", Summary: "获取所有附件类型", Tags: []string{"admin/attachments"}}, s.AttachmentHandler.GetAllTypes)
}

// registerAdminBackup 注册 /backups 相关路由。
func (s *Server) registerAdminBackup(rg *gin.RouterGroup, api huma.API) {
	// huma JSON API
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/backups/work-dir", Summary: "全站备份", Tags: []string{"admin/backups"}}, s.BackupHandler.BackupWholeSite)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/backups/work-dir", Summary: "列出全站备份", Tags: []string{"admin/backups"}}, s.BackupHandler.ListBackups)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/backups/work-dir", Summary: "删除全站备份", Tags: []string{"admin/backups"}}, s.BackupHandler.DeleteBackups)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/backups/data", Summary: "导出数据", Tags: []string{"admin/backups"}}, s.BackupHandler.ExportData)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/backups/data", Summary: "删除数据文件", Tags: []string{"admin/backups"}}, s.BackupHandler.DeleteDataFile)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/backups/markdown/export", Summary: "导出 Markdown", Tags: []string{"admin/backups"}}, s.BackupHandler.ExportMarkdown)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/backups/markdown/import", Summary: "导入 Markdown", Tags: []string{"admin/backups"}}, s.BackupHandler.ImportMarkdown)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/backups/markdown/fetch", Summary: "获取 Markdown 备份", Tags: []string{"admin/backups"}}, s.BackupHandler.GetMarkDownBackup)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/backups/markdown/export", Summary: "列出 Markdown 备份", Tags: []string{"admin/backups"}}, s.BackupHandler.ListMarkdowns)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/backups/markdown/export", Summary: "删除 Markdown 备份", Tags: []string{"admin/backups"}}, s.BackupHandler.DeleteMarkdowns)
	// gin 文件流路由（保持 gin 不变）
	backupRouter := rg.Group("/backups")
	backupRouter.GET("/work-dir/*path", s.BackupHandler.HandleWorkDir)
	backupRouter.GET("/data/*path", s.BackupHandler.HandleData)
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
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/statistics", Summary: "获取统计数据", Tags: []string{"admin/statistics"}}, s.StatisticHandler.Statistics)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/statistics/user", Summary: "获取带用户的统计数据", Tags: []string{"admin/statistics"}}, s.StatisticHandler.StatisticsWithUser)
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
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/photos/latest", Summary: "获取照片列表", Tags: []string{"admin/photos"}}, s.PhotoHandler.ListPhoto)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/photos", Summary: "分页获取照片列表", Tags: []string{"admin/photos"}}, s.PhotoHandler.PagePhotos)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/photos/{id}", Summary: "获取照片详情", Tags: []string{"admin/photos"}}, s.PhotoHandler.GetPhotoByID)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/photos/batch", Summary: "批量删除照片", Tags: []string{"admin/photos"}}, s.PhotoHandler.DeletePhotoBatch)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/photos", Summary: "创建照片", Tags: []string{"admin/photos"}}, s.PhotoHandler.CreatePhoto)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/photos/batch", Summary: "批量创建照片", Tags: []string{"admin/photos"}}, s.PhotoHandler.CreatePhotoBatch)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/photos/{id}", Summary: "更新照片", Tags: []string{"admin/photos"}}, s.PhotoHandler.UpdatePhoto)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/photos/teams", Summary: "获取照片团队列表", Tags: []string{"admin/photos"}}, s.PhotoHandler.ListPhotoTeams)
}

// registerAdminUser 注册 /users 相关路由。
func (s *Server) registerAdminUser(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/users/profiles", Summary: "获取当前用户资料", Tags: []string{"admin/users"}}, s.UserHandler.GetCurrentUserProfile)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/users/profiles", Summary: "更新用户资料", Tags: []string{"admin/users"}}, s.UserHandler.UpdateUserProfile)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/users/profiles/password", Summary: "更新密码", Tags: []string{"admin/users"}}, s.UserHandler.UpdatePassword)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/users/mfa/generate", Summary: "生成MFA二维码", Tags: []string{"admin/users"}}, s.UserHandler.GenerateMFAQRCode)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/users/mfa/update", Summary: "更新MFA", Tags: []string{"admin/users"}}, s.UserHandler.UpdateMFA)
}

// registerAdminTheme 注册 themes 相关路由。
func (s *Server) registerAdminTheme(rg *gin.RouterGroup, api huma.API) {
	// huma JSON API
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation", Summary: "获取已激活主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetActivatedTheme)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}", Summary: "按ID获取主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes", Summary: "获取所有主题列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.ListAllThemes)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/files", Summary: "获取已激活主题文件列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.ListActivatedThemeFile)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/files", Summary: "按ID获取主题文件列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.ListThemeFileByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/files/content", Summary: "获取已激活主题文件内容", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeFileContent)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/files/content", Summary: "按ID获取主题文件内容", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeFileContentByID)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/themes/files/content", Summary: "更新已激活主题文件内容", Tags: []string{"admin/themes"}}, s.ThemeHandler.UpdateThemeFile)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/themes/{themeID}/files/content", Summary: "按ID更新主题文件内容", Tags: []string{"admin/themes"}}, s.ThemeHandler.UpdateThemeFileByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/template/custom/sheet", Summary: "获取自定义页面模板列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.ListCustomSheetTemplate)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/template/custom/post", Summary: "获取自定义文章模板列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.ListCustomPostTemplate)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/themes/{themeID}/activation", Summary: "激活主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.ActivateTheme)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/configurations", Summary: "获取已激活主题配置", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetActivatedThemeConfig)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/configurations", Summary: "按ID获取主题配置", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeConfigByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/configurations/groups/{group}", Summary: "按分组获取主题配置", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeConfigByGroup)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/configurations/groups", Summary: "获取主题配置分组名称列表", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeConfigGroupNames)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/settings", Summary: "获取已激活主题设置地图", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetActivatedThemeSettingMap)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/settings", Summary: "按ID获取主题设置地图", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeSettingMapByID)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/{themeID}/groups/{group}/settings", Summary: "按分组获取主题设置地图", Tags: []string{"admin/themes"}}, s.ThemeHandler.GetThemeSettingMapByGroupAndThemeID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/themes/activation/settings", Summary: "保存已激活主题设置", Tags: []string{"admin/themes"}}, s.ThemeHandler.SaveActivatedThemeSetting)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/themes/{themeID}/settings", Summary: "按ID保存主题设置", Tags: []string{"admin/themes"}}, s.ThemeHandler.SaveThemeSettingByID)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/themes/{themeID}", Summary: "按ID删除主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.DeleteThemeByID)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/themes/fetching", Summary: "远程拉取主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.FetchTheme)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/themes/fetching/{themeID}", Summary: "远程更新主题（未实现）", Tags: []string{"admin/themes"}}, s.ThemeHandler.UpdateThemeByFetching)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/themes/reload", Summary: "重载主题", Tags: []string{"admin/themes"}}, s.ThemeHandler.ReloadTheme)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/themes/activation/template/exists", Summary: "检查模板是否存在", Tags: []string{"admin/themes"}}, s.ThemeHandler.TemplateExist)

	// gin 文件流路由（上传主题包，保持 gin 不变）
	themeRouter := rg.Group("themes")
	themeRouter.POST("upload", s.wrapHandler(s.ThemeHandler.UploadTheme))
	themeRouter.PUT("upload/:themeID", s.wrapHandler(s.ThemeHandler.UpdateThemeByUpload))
}

// registerAdminEmail 注册 /mails 相关路由。
func (s *Server) registerAdminEmail(rg *gin.RouterGroup, api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/mails/test", Summary: "发送测试邮件", Tags: []string{"admin/mails"}}, s.EmailHandler.Test)
}
