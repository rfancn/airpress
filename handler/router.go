package handler

import (
	"context"
	"path/filepath"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/rfancn/airpress/config"
	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/dal"
	"github.com/rfancn/airpress/handler/middleware"
)

func (s *Server) RegisterRouters() {
	router := s.Router
	// 收集所有 huma API 实例及其 group 前缀，最后合并成一份 OpenAPI 3.1 文档
	var humaAPIs []huma.API
	var humaPrefixes []string
	if config.IsDev() {
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowOrigins:     []string{},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Admin-Authorization", "Content-Type"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length"},
		}))
	}

	{
		router.GET("/ping", func(ctx *gin.Context) {
			_, _ = ctx.Writer.Write([]byte("pong"))
		})
		{
			staticRouter := router.Group("/")
			staticRouter.StaticFS(s.Config.AirPress.AdminURLPath, gin.Dir(s.Config.AirPress.AdminResourcesDir, false))
			staticRouter.StaticFS("/css", gin.Dir(filepath.Join(s.Config.AirPress.AdminResourcesDir, "css"), false))
			staticRouter.StaticFS("/js", gin.Dir(filepath.Join(s.Config.AirPress.AdminResourcesDir, "js"), false))
			staticRouter.StaticFS("/images", gin.Dir(filepath.Join(s.Config.AirPress.AdminResourcesDir, "images"), false))
			staticRouter.Use(middleware.NewCacheControlMiddleware(middleware.WithMaxAge(time.Hour*24*7)).CacheControl()).
				StaticFS(consts.AirPressUploadDir, gin.Dir(s.Config.AirPress.UploadDir, false))
			staticRouter.StaticFS("/themes/", gin.Dir(s.Config.AirPress.ThemeDir, false))
		}
		{
			adminAPIRouter := router.Group("/api/admin")
			adminAPIRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger(), s.InstallRedirectMiddleware.InstallRedirect())
			// 公开 huma API（登录/安装，无鉴权，不能挂到带鉴权的 authRouter 上）
			adminPublicHumaAPI, adminPublicPrefix := newHumaAPI(s.Router, adminAPIRouter, "AirPress Admin Public API")
			humaAPIs = append(humaAPIs, adminPublicHumaAPI)
			humaPrefixes = append(humaPrefixes, adminPublicPrefix)
			s.registerAdminPublicHumaAPI(adminPublicHumaAPI)
			// 已迁移到 huma：is_installed / login/precheck / login / refresh/{refreshToken} / installations
			// adminAPIRouter.GET("/is_installed", s.wrapHandler(s.AdminHandler.IsInstalled))
			// adminAPIRouter.POST("/login/precheck", s.wrapHandler(s.AdminHandler.AuthPreCheck))
			// adminAPIRouter.POST("/login", s.wrapHandler(s.AdminHandler.Auth))
			// adminAPIRouter.POST("/refresh/:refreshToken", s.wrapHandler(s.AdminHandler.RefreshToken))
			// adminAPIRouter.POST("/installations", s.wrapHandler(s.InstallHandler.InstallBlog))
			{
				authRouter := adminAPIRouter.Group("")
				authRouter.Use(s.AuthMiddleware.GetWrapHandler())
				// huma 挂在 authRouter 上，鉴权中间件自动生效
				adminHumaAPI, adminPrefix := newHumaAPI(s.Router, authRouter, "AirPress Admin API")
				// 全局中间件：把 gin 鉴权中间件注入的用户复制到 huma context
				adminHumaAPI.UseMiddleware(s.adminAuthUserMiddleware)
				humaAPIs = append(humaAPIs, adminHumaAPI)
				humaPrefixes = append(humaPrefixes, adminPrefix)
				s.registerAdminBase(authRouter, adminHumaAPI)
				s.registerAdminAttachment(authRouter, adminHumaAPI)
				s.registerAdminBackup(authRouter, adminHumaAPI)
				s.registerAdminCategory(authRouter, adminHumaAPI)
				s.registerAdminPost(authRouter, adminHumaAPI)
				s.registerAdminOption(authRouter, adminHumaAPI)
				s.registerAdminLog(authRouter, adminHumaAPI)
				s.registerAdminStatistic(authRouter, adminHumaAPI)
				s.registerAdminSheet(authRouter, adminHumaAPI)
				s.registerAdminJournal(authRouter, adminHumaAPI)
				s.registerAdminLink(authRouter, adminHumaAPI)
				s.registerAdminMenu(authRouter, adminHumaAPI)
				s.registerAdminTag(authRouter, adminHumaAPI)
				s.registerAdminPhoto(authRouter, adminHumaAPI)
				s.registerAdminUser(authRouter, adminHumaAPI)
				s.registerAdminTheme(authRouter, adminHumaAPI)
				s.registerAdminEmail(authRouter, adminHumaAPI)
			}
		}
		{
			contentRouter := router.Group("")
			contentRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger(), s.InstallRedirectMiddleware.InstallRedirect())

			contentRouter.POST("/content/:type/:slug/authentication", s.wrapHTMLHandler(s.ViewHandler.Authenticate))

			contentRouter.GET("", s.wrapHTMLHandler(s.IndexHandler.Index))
			contentRouter.GET("/page/:page", s.wrapHTMLHandler(s.IndexHandler.IndexPage))
			contentRouter.GET("/robots.txt", s.wrapTextHandler(s.FeedHandler.Robots))
			contentRouter.GET("/atom", s.wrapTextHandler(s.FeedHandler.Atom))
			contentRouter.GET("/atom.xml", s.wrapTextHandler(s.FeedHandler.Atom))
			contentRouter.GET("/rss", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/rss.xml", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed.xml", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed/categories/:slug", s.wrapTextHandler(s.FeedHandler.CategoryFeed))
			contentRouter.GET("/atom/categories/:slug", s.wrapTextHandler(s.FeedHandler.CategoryAtom))
			contentRouter.GET("/sitemap.xml", s.wrapTextHandler(s.FeedHandler.SitemapXML))
			contentRouter.GET("/sitemap.html", s.wrapHTMLHandler(s.FeedHandler.SitemapHTML))

			contentRouter.GET("/version", s.wrapHandler(s.ViewHandler.Version))
			contentRouter.GET("/install", s.ViewHandler.Install)
			contentRouter.GET("/logo", s.wrapHandler(s.ViewHandler.Logo))
			contentRouter.GET("/favicon", s.wrapHandler(s.ViewHandler.Favicon))
			contentRouter.GET("/search", s.wrapHTMLHandler(s.ContentSearchHandler.Search))
			contentRouter.GET("/search/page/:page", s.wrapHTMLHandler(s.ContentSearchHandler.PageSearch))
			err := s.registerDynamicRouters(contentRouter)
			if err != nil {
				s.logger.DPanic("regiterDynamicRouters err", zap.Error(err))
			}
		}
		{
			contentAPIRouter := router.Group("/api/content")
			contentAPIRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger())

			// huma 挂在 contentAPIRouter 上，content api 无需鉴权
			contentHumaAPI, contentPrefix := newHumaAPI(s.Router, contentAPIRouter, "AirPress Content API")
			humaAPIs = append(humaAPIs, contentHumaAPI)
			humaPrefixes = append(humaPrefixes, contentPrefix)
			s.registerContentHumaAPI(contentHumaAPI)
		}
	}

	// 合并所有 huma API 的 OpenAPI 文档，注册 /openapi.json
	s.registerHumaDocs(humaAPIs, humaPrefixes)
}

func (s *Server) registerDynamicRouters(contentRouter *gin.RouterGroup) error {
	ctx := context.Background()
	ctx = dal.SetCtxQuery(ctx, dal.GetQueryByCtx(ctx).ReplaceDB(dal.GetDB().Session(
		&gorm.Session{Logger: dal.DB.Logger.LogMode(logger.Warn)},
	)))

	archivePath, err := s.OptionService.GetArchivePrefix(ctx)
	if err != nil {
		return err
	}
	categoryPath, err := s.OptionService.GetCategoryPrefix(ctx)
	if err != nil {
		return err
	}
	sheetPermaLinkType, err := s.OptionService.GetSheetPermalinkType(ctx)
	if err != nil {
		return err
	}
	sheetPath, err := s.OptionService.GetSheetPrefix(ctx)
	if err != nil {
		return err
	}
	tagPath, err := s.OptionService.GetTagPrefix(ctx)
	if err != nil {
		return err
	}
	journalPath, err := s.OptionService.GetJournalPrefix(ctx)
	if err != nil {
		return err
	}

	photoPath, err := s.OptionService.GetPhotoPrefix(ctx)
	if err != nil {
		return err
	}
	linkPath, err := s.OptionService.GetLinkPrefix(ctx)
	if err != nil {
		return err
	}
	contentRouter.GET(archivePath, s.wrapHTMLHandler(s.ArchiveHandler.Archives))
	contentRouter.GET(archivePath+"/page/:page", s.wrapHTMLHandler(s.ArchiveHandler.ArchivesPage))
	contentRouter.GET(archivePath+"/:slug", s.wrapHTMLHandler(s.ArchiveHandler.ArchivesBySlug))

	contentRouter.GET(tagPath, s.wrapHTMLHandler(s.ContentTagHandler.Tags))
	contentRouter.GET(tagPath+"/:slug/page/:page", s.wrapHTMLHandler(s.ContentTagHandler.TagPostPage))
	contentRouter.GET(tagPath+"/:slug", s.wrapHTMLHandler(s.ContentTagHandler.TagPost))

	contentRouter.GET(categoryPath, s.wrapHTMLHandler(s.ContentCategoryHandler.Categories))
	contentRouter.GET(categoryPath+"/:slug", s.wrapHTMLHandler(s.ContentCategoryHandler.CategoryDetail))
	contentRouter.GET(categoryPath+"/:slug/page/:page", s.wrapHTMLHandler(s.ContentCategoryHandler.CategoryDetailPage))

	contentRouter.GET(linkPath, s.wrapHTMLHandler(s.ContentLinkHandler.Link))

	contentRouter.GET(photoPath, s.wrapHTMLHandler(s.ContentPhotoHandler.Phtotos))
	contentRouter.GET(photoPath+"/page/:page", s.wrapHTMLHandler(s.ContentPhotoHandler.PhotosPage))

	contentRouter.GET(journalPath, s.wrapHTMLHandler(s.ContentJournalHandler.Journals))
	contentRouter.GET(journalPath+"/page/:page", s.wrapHTMLHandler(s.ContentJournalHandler.JournalsPage))
	contentRouter.GET("admin_preview/"+archivePath+"/:slug", s.wrapHTMLHandler(s.ArchiveHandler.AdminArchivesBySlug))
	if sheetPermaLinkType == consts.SheetPermaLinkTypeRoot {
		contentRouter.GET("/:slug")
	} else {
		contentRouter.GET(sheetPath+"/:slug", s.wrapHTMLHandler(s.ContentSheetHandler.SheetBySlug))
	}
	contentRouter.GET("admin_preview/"+sheetPath+"/:slug", s.wrapHTMLHandler(s.ContentSheetHandler.AdminSheetBySlug))
	return nil
}
