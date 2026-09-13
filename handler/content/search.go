package content

import (
	"html"

	"github.com/gin-gonic/gin"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/handler/binding"
	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/property"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/service/assembler"
	"github.com/rfancn/airpress/template"
	"github.com/rfancn/airpress/util"
	"github.com/rfancn/airpress/util/xerr"
)

type SearchHandler struct {
	PostAssembler assembler.PostAssembler
	PostService   service.PostService
	OptionService service.OptionService
	ThemeService  service.ThemeService
}

func NewSearchHandler(
	postAssembler assembler.PostAssembler,
	postService service.PostService,
	optionService service.OptionService,
	themeService service.ThemeService,
) *SearchHandler {
	return &SearchHandler{
		PostAssembler: postAssembler,
		PostService:   postService,
		OptionService: optionService,
		ThemeService:  themeService,
	}
}

func (s *SearchHandler) Search(ctx *gin.Context, model template.Model) (string, error) {
	return s.search(ctx, 0, model)
}

func (s *SearchHandler) PageSearch(ctx *gin.Context, model template.Model) (string, error) {
	page, err := util.ParamInt32(ctx, "page")
	if err != nil {
		return "", err
	}
	return s.search(ctx, int(page)-1, model)
}

func (s *SearchHandler) search(ctx *gin.Context, pageNum int, model template.Model) (string, error) {
	keyword, err := util.MustGetQueryString(ctx, "keyword")
	if err != nil {
		return "", err
	}
	sort := param.Sort{}
	err = ctx.ShouldBindWith(&sort, binding.CustomFormBinding)
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if len(sort.Fields) == 0 {
		sort = s.OptionService.GetPostSort(ctx)
	}
	defaultPageSize := s.OptionService.GetIndexPageSize(ctx)
	page := param.Page{
		PageNum:  pageNum,
		PageSize: defaultPageSize,
	}
	postQuery := param.PostQuery{
		Page:     page,
		Sort:     &sort,
		Keyword:  &keyword,
		Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
	}
	posts, total, err := s.PostService.Page(ctx, postQuery)
	if err != nil {
		return "", err
	}
	postVOs, err := s.PostAssembler.ConvertToListVO(ctx, posts)
	if err != nil {
		return "", err
	}
	model["is_search"] = true
	model["keyword"] = html.EscapeString(keyword)
	model["posts"] = dto.NewPage(postVOs, total, page)
	model["meta_keywords"] = s.OptionService.GetOrByDefault(ctx, property.SeoKeywords)
	model["meta_description"] = s.OptionService.GetOrByDefault(ctx, property.SeoDescription)
	return s.ThemeService.Render(ctx, "search")
}
