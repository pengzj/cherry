package pagination

import (
	"fmt"
)

type Pagination struct {
	TotalCount int `json:"totalCount"`
	TotalPage int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
	PageSize int `json:"pageSize"`
}

func (p *Pagination) GetSkip() int {
	if p.CurrentPage <= 1 {
		p.CurrentPage = 1
	}
	return (p.CurrentPage -1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	return p.PageSize
}

func (p *Pagination) SetTotalCount(totalCount int)  {
	if p.PageSize == 0 {
		p.PageSize = 20;
	}
	p.TotalCount = totalCount
	p.TotalPage = totalCount / p.PageSize
	if totalCount % p.PageSize != 0 {
		p.TotalPage++
	}
}

func (p *Pagination) String() string {
	return fmt.Sprintf("Limit %d, %d ", p.GetSkip(), p.GetLimit())
}