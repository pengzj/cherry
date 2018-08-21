package pagination

import (
	"fmt"
)

type Pagination struct {
	TotalCount int `json:"totalCount"`
	TotalPage int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
	PageNum int `json:"pageNum"`
}

func (p *Pagination) GetSkip() int {
	if p.CurrentPage <= 1 {
		p.CurrentPage = 1
	}
	return (p.CurrentPage -1) * p.PageNum
}

func (p *Pagination) GetLimit() int {
	return p.PageNum
}

func (p *Pagination) SetTotalCount(totalCount int)  {
	if p.PageNum == 0 {
		p.PageNum = 20;
	}
	p.TotalCount = totalCount
	p.TotalPage = totalCount / p.PageNum
	if totalCount % p.PageNum != 0 {
		p.TotalPage++
	}
}

func (p *Pagination) String() string {
	return fmt.Sprintf("Limit %d, %d ", p.GetSkip(), p.GetLimit())
}