package sortable

import (
	"fmt"
)

type Direction uint8

const (
	ASC = 0 + iota
	DESC
)

var types = [...]string {
	"ASC",
	"DESC",
}

func (s Direction) String() string {
	if (s < ASC || s > DESC) {
		return  "ASC"
	}
	return types[s]
}

func ParseDirection(d string) Direction {
	if d == "desc" {
		return DESC
	}
	return ASC
}


type Sortable struct {
	Name string `json:"name"`
	Direction Direction `json:"direction"`
}

func (s *Sortable) String() string {
	return fmt.Sprintf("ORDER BY %s %s ", s.Name, s.Direction.String())
}