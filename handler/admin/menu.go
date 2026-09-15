package admin

import (
	"context"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/model/vo"
	"github.com/rfancn/airpress/service"
)

type MenuHandler struct {
	MenuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		MenuService: menuService,
	}
}

// ListMenusInput 菜单列表查询输入。
type ListMenusInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// ListMenus 获取菜单列表。
func (m *MenuHandler) ListMenus(ctx context.Context, in *ListMenusInput) (*dto.HumaOut[[]*dto.Menu], error) {
	fields := in.Sort
	if len(fields) == 0 {
		fields = []string{"team,desc", "priority,asc"}
	} else {
		fields = append(fields, "priority,asc")
	}
	sort := &param.Sort{Fields: fields}
	menus, err := m.MenuService.List(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTOs(ctx, menus))
}

// ListMenusAsTreeInput 菜单树查询输入。
type ListMenusAsTreeInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
}

// ListMenusAsTree 获取菜单树。
func (m *MenuHandler) ListMenusAsTree(ctx context.Context, in *ListMenusAsTreeInput) (*dto.HumaOut[[]*vo.Menu], error) {
	fields := in.Sort
	if len(fields) == 0 {
		fields = []string{"team,desc", "priority,asc"}
	} else {
		fields = append(fields, "priority,asc")
	}
	sort := &param.Sort{Fields: fields}
	data, err := m.MenuService.ListAsTree(ctx, sort)
	if err != nil {
		return dto.HumaErr[[]*vo.Menu](err)
	}
	return dto.HumaOK(data)
}

// ListMenusAsTreeByTeamInput 按团队获取菜单树查询输入。
type ListMenusAsTreeByTeamInput struct {
	Sort []string `query:"sort" doc:"排序字段"`
	Team string   `query:"team" doc:"团队名称"`
}

// ListMenusAsTreeByTeam 按团队获取菜单树。
func (m *MenuHandler) ListMenusAsTreeByTeam(ctx context.Context, in *ListMenusAsTreeByTeamInput) (*dto.HumaOut[[]*vo.Menu], error) {
	fields := in.Sort
	if len(fields) == 0 {
		fields = []string{"priority,asc"}
	}
	sort := &param.Sort{Fields: fields}
	if in.Team == "" {
		data, err := m.MenuService.ListAsTree(ctx, sort)
		if err != nil {
			return dto.HumaErr[[]*vo.Menu](err)
		}
		return dto.HumaOK(data)
	}
	data, err := m.MenuService.ListAsTreeByTeam(ctx, in.Team, sort)
	if err != nil {
		return dto.HumaErr[[]*vo.Menu](err)
	}
	return dto.HumaOK(data)
}

// GetMenuByIDInput 菜单详情查询输入。
type GetMenuByIDInput struct {
	ID int32 `path:"id" doc:"菜单ID"`
}

// GetMenuByID 获取菜单详情。
func (m *MenuHandler) GetMenuByID(ctx context.Context, in *GetMenuByIDInput) (*dto.HumaOut[*dto.Menu], error) {
	menu, err := m.MenuService.GetByID(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTO(ctx, menu))
}

// CreateMenuInput 创建菜单输入。
type CreateMenuInput struct {
	Body param.Menu `doc:"菜单参数"`
}

// CreateMenu 创建菜单。
func (m *MenuHandler) CreateMenu(ctx context.Context, in *CreateMenuInput) (*dto.HumaOut[*dto.Menu], error) {
	menu, err := m.MenuService.Create(ctx, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTO(ctx, menu))
}

// CreateMenuBatchInput 批量创建菜单输入。
type CreateMenuBatchInput struct {
	Body []*param.Menu `doc:"批量菜单参数"`
}

// CreateMenuBatch 批量创建菜单。
func (m *MenuHandler) CreateMenuBatch(ctx context.Context, in *CreateMenuBatchInput) (*dto.HumaOut[[]*dto.Menu], error) {
	menus, err := m.MenuService.CreateBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[[]*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTOs(ctx, menus))
}

// UpdateMenuInput 更新菜单输入。
type UpdateMenuInput struct {
	ID   int32      `path:"id" doc:"菜单ID"`
	Body param.Menu `doc:"菜单参数"`
}

// UpdateMenu 更新菜单。
func (m *MenuHandler) UpdateMenu(ctx context.Context, in *UpdateMenuInput) (*dto.HumaOut[*dto.Menu], error) {
	menu, err := m.MenuService.Update(ctx, in.ID, &in.Body)
	if err != nil {
		return dto.HumaErr[*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTO(ctx, menu))
}

// UpdateMenuBatchInput 批量更新菜单输入。
type UpdateMenuBatchInput struct {
	Body []*param.Menu `doc:"批量菜单参数"`
}

// UpdateMenuBatch 批量更新菜单。
func (m *MenuHandler) UpdateMenuBatch(ctx context.Context, in *UpdateMenuBatchInput) (*dto.HumaOut[[]*dto.Menu], error) {
	menus, err := m.MenuService.UpdateBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[[]*dto.Menu](err)
	}
	return dto.HumaOK(m.MenuService.ConvertToDTOs(ctx, menus))
}

// DeleteMenuInput 删除菜单输入。
type DeleteMenuInput struct {
	ID int32 `path:"id" doc:"菜单ID"`
}

// DeleteMenu 删除菜单。
func (m *MenuHandler) DeleteMenu(ctx context.Context, in *DeleteMenuInput) (*dto.HumaOut[interface{}], error) {
	err := m.MenuService.Delete(ctx, in.ID)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// DeleteMenuBatchInput 批量删除菜单输入。
type DeleteMenuBatchInput struct {
	Body []int32 `doc:"菜单ID列表"`
}

// DeleteMenuBatch 批量删除菜单。
func (m *MenuHandler) DeleteMenuBatch(ctx context.Context, in *DeleteMenuBatchInput) (*dto.HumaOut[interface{}], error) {
	err := m.MenuService.DeleteBatch(ctx, in.Body)
	if err != nil {
		return dto.HumaErr[interface{}](err)
	}
	return dto.HumaOK[interface{}](nil)
}

// ListMenuTeamsInput 菜单团队列表查询无输入参数。
type ListMenuTeamsInput struct{}

// ListMenuTeams 获取菜单团队列表。
func (m *MenuHandler) ListMenuTeams(ctx context.Context, _ *ListMenuTeamsInput) (*dto.HumaOut[[]string], error) {
	data, err := m.MenuService.ListTeams(ctx)
	if err != nil {
		return dto.HumaErr[[]string](err)
	}
	return dto.HumaOK(data)
}
