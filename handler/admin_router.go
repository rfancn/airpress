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
	categoryRouter := rg.Group("/categories")
	categoryRouter.PUT("/batch", s.wrapHandler(s.CategoryHandler.UpdateCategoryBatch))
	categoryRouter.GET("/:categoryID", s.wrapHandler(s.CategoryHandler.GetCategoryByID))
	categoryRouter.GET("", s.wrapHandler(s.CategoryHandler.ListAllCategory))
	categoryRouter.GET("/tree_view", s.wrapHandler(s.CategoryHandler.ListAsTree))
	categoryRouter.POST("", s.wrapHandler(s.CategoryHandler.CreateCategory))
	categoryRouter.PUT("/:categoryID", s.wrapHandler(s.CategoryHandler.UpdateCategory))
	categoryRouter.DELETE("/:categoryID", s.wrapHandler(s.CategoryHandler.DeleteCategory))
}

// registerAdminPost 注册 /posts 及其 /comments 子路由。
func (s *Server) registerAdminPost(rg *gin.RouterGroup, api huma.API) {
	postRouter := rg.Group("/posts")
	postRouter.GET("", s.wrapHandler(s.PostHandler.ListPosts))
	postRouter.GET("/latest", s.wrapHandler(s.PostHandler.ListLatestPosts))
	postRouter.GET("/status/:status", s.wrapHandler(s.PostHandler.ListPostsByStatus))
	// 已迁移到 huma：huma.Register(api, ... "/posts/{postID}")
	// postRouter.GET("/:postID", s.wrapHandler(s.PostHandler.GetByPostID))
	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/posts/{postID}",
		Summary: "获取文章详情",
		Tags:    []string{"admin/posts"},
	}, s.PostHandler.GetByPostID)
	postRouter.POST("", s.wrapHandler(s.PostHandler.CreatePost))
	postRouter.PUT("/:postID", s.wrapHandler(s.PostHandler.UpdatePost))
	postRouter.PUT("/:postID/status/:status", s.wrapHandler(s.PostHandler.UpdatePostStatus))
	postRouter.PUT("/status/:status", s.wrapHandler(s.PostHandler.UpdatePostStatusBatch))
	postRouter.PUT("/:postID/status/draft/content", s.wrapHandler(s.PostHandler.UpdatePostDraft))
	postRouter.DELETE("/:postID", s.wrapHandler(s.PostHandler.DeletePost))
	postRouter.DELETE("", s.wrapHandler(s.PostHandler.DeletePostBatch))
	postRouter.GET("/:postID/preview", s.PostHandler.PreviewPost)
	{
		postCommentRouter := postRouter.Group("/comments")
		postCommentRouter.GET("", s.wrapHandler(s.PostCommentHandler.ListPostComment))
		postCommentRouter.GET("/latest", s.wrapHandler(s.PostCommentHandler.ListPostCommentLatest))
		postCommentRouter.GET("/:postID/tree_view", s.wrapHandler(s.PostCommentHandler.ListPostCommentAsTree))
		postCommentRouter.GET("/:postID/list_view", s.wrapHandler(s.PostCommentHandler.ListPostCommentWithParent))
		postCommentRouter.POST("", s.wrapHandler(s.PostCommentHandler.CreatePostComment))
		postCommentRouter.PUT("/:commentID", s.wrapHandler(s.PostCommentHandler.UpdatePostComment))
		postCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.PostCommentHandler.UpdatePostCommentStatus))
		postCommentRouter.PUT("/status/:status", s.wrapHandler(s.PostCommentHandler.UpdatePostCommentStatusBatch))
		postCommentRouter.DELETE("/:commentID", s.wrapHandler(s.PostCommentHandler.DeletePostComment))
		postCommentRouter.DELETE("", s.wrapHandler(s.PostCommentHandler.DeletePostCommentBatch))
	}
}

// registerAdminOption 注册 /options 相关路由。
func (s *Server) registerAdminOption(rg *gin.RouterGroup, api huma.API) {
	optionRouter := rg.Group("/options")
	optionRouter.GET("", s.wrapHandler(s.OptionHandler.ListAllOptions))
	optionRouter.GET("/map_view", s.wrapHandler(s.OptionHandler.ListAllOptionsAsMap))
	optionRouter.POST("/map_view/keys", s.wrapHandler(s.OptionHandler.ListAllOptionsAsMapWithKey))
	optionRouter.POST("/saving", s.wrapHandler(s.OptionHandler.SaveOption))
	optionRouter.POST("/map_view/saving", s.wrapHandler(s.OptionHandler.SaveOptionWithMap))
}

// registerAdminLog 注册 /logs 相关路由。
func (s *Server) registerAdminLog(rg *gin.RouterGroup, api huma.API) {
	logRouter := rg.Group("/logs")
	logRouter.GET("/latest", s.wrapHandler(s.LogHandler.PageLatestLog))
	logRouter.GET("", s.wrapHandler(s.LogHandler.PageLog))
	logRouter.GET("/clear", s.wrapHandler(s.LogHandler.ClearLog))
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
	sheetRouter.GET("/:sheetID", s.wrapHandler(s.SheetHandler.GetSheetByID))
	sheetRouter.GET("", s.wrapHandler(s.SheetHandler.ListSheet))
	sheetRouter.POST("", s.wrapHandler(s.SheetHandler.CreateSheet))
	sheetRouter.PUT("/:sheetID", s.wrapHandler(s.SheetHandler.UpdateSheet))
	sheetRouter.PUT("/:sheetID/:status", s.wrapHandler(s.SheetHandler.UpdateSheetStatus))
	sheetRouter.PUT("/:sheetID/status/draft/content", s.wrapHandler(s.SheetHandler.UpdateSheetDraft))
	sheetRouter.DELETE("/:sheetID", s.wrapHandler(s.SheetHandler.DeleteSheet))
	sheetRouter.GET("/preview/:sheetID", s.SheetHandler.PreviewSheet)
	sheetRouter.GET("/independent", s.wrapHandler(s.SheetHandler.IndependentSheets))
	{
		sheetCommentRouter := sheetRouter.Group("/comments")
		sheetCommentRouter.GET("", s.wrapHandler(s.SheetCommentHandler.ListSheetComment))
		sheetCommentRouter.GET("/latest", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentLatest))
		sheetCommentRouter.GET("/:sheetID/tree_view", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentAsTree))
		sheetCommentRouter.GET("/:sheetID/list_view", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentWithParent))
		sheetCommentRouter.POST("/", s.wrapHandler(s.SheetCommentHandler.CreateSheetComment))
		sheetCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.SheetCommentHandler.UpdateSheetCommentStatus))
		sheetCommentRouter.PUT("/status/:status", s.wrapHandler(s.SheetCommentHandler.UpdateSheetCommentStatusBatch))
		sheetCommentRouter.DELETE("/:commentID", s.wrapHandler(s.SheetCommentHandler.DeleteSheetComment))
		sheetCommentRouter.DELETE("", s.wrapHandler(s.SheetCommentHandler.DeleteSheetCommentBatch))
	}
}

// registerAdminJournal 注册 /journals 及其 /comments 子路由。
func (s *Server) registerAdminJournal(rg *gin.RouterGroup, api huma.API) {
	journalRouter := rg.Group("/journals")
	journalRouter.GET("", s.wrapHandler(s.JournalHandler.ListJournal))
	journalRouter.GET("/latest", s.wrapHandler(s.JournalHandler.ListLatestJournal))
	journalRouter.POST("", s.wrapHandler(s.JournalHandler.CreateJournal))
	journalRouter.PUT("/:journalID", s.wrapHandler(s.JournalHandler.UpdateJournal))
	journalRouter.DELETE("/:journalID", s.wrapHandler(s.JournalHandler.DeleteJournal))
	{
		journalCommentRouter := journalRouter.Group("/comments")
		journalCommentRouter.GET("", s.wrapHandler(s.JournalCommentHandler.ListJournalComment))
		journalCommentRouter.GET("/latest", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentLatest))
		journalCommentRouter.GET("/:journalID/tree_view", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentAsTree))
		journalCommentRouter.GET("/:journalID/list_view", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentWithParent))
		journalCommentRouter.POST("/", s.wrapHandler(s.JournalCommentHandler.CreateJournalComment))
		journalCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.JournalCommentHandler.UpdateJournalCommentStatus))
		journalCommentRouter.PUT("/status/:status", s.wrapHandler(s.JournalCommentHandler.UpdateJournalStatusBatch))
		journalCommentRouter.PUT("/:commentID", s.wrapHandler(s.JournalCommentHandler.UpdateJournalComment))
		journalCommentRouter.DELETE("/:commentID", s.wrapHandler(s.JournalCommentHandler.DeleteJournalComment))
		journalCommentRouter.DELETE("", s.wrapHandler(s.JournalCommentHandler.DeleteJournalCommentBatch))
	}
}

// registerAdminLink 注册 /links 相关路由。
func (s *Server) registerAdminLink(rg *gin.RouterGroup, api huma.API) {
	linkRouter := rg.Group("/links")
	linkRouter.GET("", s.wrapHandler(s.LinkHandler.ListLinks))
	linkRouter.GET("/:id", s.wrapHandler(s.LinkHandler.GetLinkByID))
	linkRouter.POST("", s.wrapHandler(s.LinkHandler.CreateLink))
	linkRouter.PUT("/:id", s.wrapHandler(s.LinkHandler.UpdateLink))
	linkRouter.DELETE("/:id", s.wrapHandler(s.LinkHandler.DeleteLink))
	linkRouter.GET("/teams", s.wrapHandler(s.LinkHandler.ListLinkTeams))
}

// registerAdminMenu 注册 /menus 相关路由。
func (s *Server) registerAdminMenu(rg *gin.RouterGroup, api huma.API) {
	menuRouter := rg.Group("/menus")
	menuRouter.GET("", s.wrapHandler(s.MenuHandler.ListMenus))
	menuRouter.GET("/tree_view", s.wrapHandler(s.MenuHandler.ListMenusAsTree))
	menuRouter.GET("/team/tree_view", s.wrapHandler(s.MenuHandler.ListMenusAsTreeByTeam))
	menuRouter.GET("/:id", s.wrapHandler(s.MenuHandler.GetMenuByID))
	menuRouter.POST("", s.wrapHandler(s.MenuHandler.CreateMenu))
	menuRouter.POST("/batch", s.wrapHandler(s.MenuHandler.CreateMenuBatch))
	menuRouter.PUT("/:id", s.wrapHandler(s.MenuHandler.UpdateMenu))
	menuRouter.PUT("/batch", s.wrapHandler(s.MenuHandler.UpdateMenuBatch))
	menuRouter.DELETE("/:id", s.wrapHandler(s.MenuHandler.DeleteMenu))
	menuRouter.DELETE("/batch", s.wrapHandler(s.MenuHandler.DeleteMenuBatch))
	menuRouter.GET("/teams", s.wrapHandler(s.MenuHandler.ListMenuTeams))
}

// registerAdminTag 注册 /tags 相关路由。
func (s *Server) registerAdminTag(rg *gin.RouterGroup, api huma.API) {
	tagRouter := rg.Group("/tags")
	tagRouter.GET("", s.wrapHandler(s.TagHandler.ListTags))
	tagRouter.GET("/:id", s.wrapHandler(s.TagHandler.GetTagByID))
	tagRouter.POST("", s.wrapHandler(s.TagHandler.CreateTag))
	tagRouter.PUT("/:id", s.wrapHandler(s.TagHandler.UpdateTag))
	tagRouter.DELETE("/:id", s.wrapHandler(s.TagHandler.DeleteTag))
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
