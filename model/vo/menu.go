package vo

import "github.com/rfancn/airpress/model/dto"

type Menu struct {
	dto.Menu
	Children []*Menu `json:"children"`
}
