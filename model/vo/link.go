package vo

import "github.com/rfancn/airpress/model/dto"

type LinkTeamVO struct {
	Team  string
	Links []*dto.Link
}
