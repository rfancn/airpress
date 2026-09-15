package admin

import (
	"context"
	"strconv"

	"github.com/rfancn/airpress/model/dto"
	"github.com/rfancn/airpress/model/param"
	"github.com/rfancn/airpress/service"
	"github.com/rfancn/airpress/util/xerr"
)

type OptionHandler struct {
	OptionService service.OptionService
}

func NewOptionHandler(optionService service.OptionService) *OptionHandler {
	return &OptionHandler{
		OptionService: optionService,
	}
}

// ListAllOptionsInput 选项列表查询无输入参数。
type ListAllOptionsInput struct{}

// ListAllOptions 获取全部选项列表。
func (o *OptionHandler) ListAllOptions(ctx context.Context, _ *ListAllOptionsInput) (*dto.HumaOut[[]*dto.Option], error) {
	options, err := o.OptionService.ListAllOption(ctx)
	if err != nil {
		return dto.HumaErr[[]*dto.Option](err)
	}
	return dto.HumaOK(options)
}

// SaveOptionInput 保存选项输入。
type SaveOptionInput struct {
	Body []*param.Option `doc:"选项列表"`
}

// SaveOption 批量保存选项。
func (o *OptionHandler) SaveOption(ctx context.Context, in *SaveOptionInput) (*dto.HumaOut[any], error) {
	optionMap := make(map[string]string, len(in.Body))
	for _, option := range in.Body {
		optionMap[option.Key] = option.Value
	}
	err := o.OptionService.Save(ctx, optionMap)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}

// ListAllOptionsAsMapInput 选项地图视图查询无输入参数。
type ListAllOptionsAsMapInput struct{}

// ListAllOptionsAsMap 获取全部选项的键值地图。
func (o *OptionHandler) ListAllOptionsAsMap(ctx context.Context, _ *ListAllOptionsAsMapInput) (*dto.HumaOut[map[string]interface{}], error) {
	options, err := o.OptionService.ListAllOption(ctx)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	result := make(map[string]interface{})
	for _, option := range options {
		result[option.Key] = option.Value
	}
	return dto.HumaOK(result)
}

// ListAllOptionsAsMapWithKeyInput 按指定 key 列表获取选项地图。
type ListAllOptionsAsMapWithKeyInput struct {
	Body []string `doc:"选项key列表"`
}

// ListAllOptionsAsMapWithKey 获取指定 key 的选项键值地图。
func (o *OptionHandler) ListAllOptionsAsMapWithKey(ctx context.Context, in *ListAllOptionsAsMapWithKeyInput) (*dto.HumaOut[map[string]interface{}], error) {
	options, err := o.OptionService.ListAllOption(ctx)
	if err != nil {
		return dto.HumaErr[map[string]interface{}](err)
	}
	keyMap := make(map[string]struct{})
	for _, key := range in.Body {
		keyMap[key] = struct{}{}
	}
	result := make(map[string]interface{})
	for _, option := range options {
		if _, ok := keyMap[option.Key]; ok {
			result[option.Key] = option.Value
		}
	}
	return dto.HumaOK(result)
}

// SaveOptionWithMapInput 以地图形式保存选项。
type SaveOptionWithMapInput struct {
	Body map[string]interface{} `doc:"选项键值对"`
}

// SaveOptionWithMap 以地图形式批量保存选项。
func (o *OptionHandler) SaveOptionWithMap(ctx context.Context, in *SaveOptionWithMapInput) (*dto.HumaOut[any], error) {
	temp := make(map[string]string)
	for key, value := range in.Body {
		var v string
		switch value := value.(type) {
		case int32:
			v = strconv.Itoa(int(value))
		case int64:
			v = strconv.FormatInt(value, 10)
		case int:
			v = strconv.Itoa(value)
		case string:
			v = value
		case bool:
			v = strconv.FormatBool(value)
		case float64:
			v = strconv.FormatFloat(value, 'f', -1, 64)
		case float32:
			v = strconv.FormatFloat(float64(value), 'f', -1, 32)
		default:
			return dto.HumaErr[any](xerr.BadParam.New("key=%v,value=%v", key, value).WithStatus(xerr.StatusBadRequest).WithMsg("Parameter type is incorrect"))
		}
		temp[key] = v
	}
	err := o.OptionService.Save(ctx, temp)
	if err != nil {
		return dto.HumaErr[any](err)
	}
	return dto.HumaOK[any](nil)
}
