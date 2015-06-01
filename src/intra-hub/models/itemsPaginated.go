package models

import (
	"github.com/soniah/tcgl/sort"
)

type ItemPaginated struct {
	Items          interface{} `json:"items"`
	ItemCount      int         `json:"itemCount"`
	TotalItemCount int         `json:"totalItemCount"`
	CurrentPage    int         `json:"currentPage"`
	TotalPageCount int         `json:"totalPageCount"`
	PagesToShow    []int       `json:"-"`
}

type Alpha []int

func (a Alpha) Len() int           { return len(a) }
func (a Alpha) Less(i, j int) bool { return a[i] < a[j] }
func (a Alpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (p *ItemPaginated) SetPagesToShow() {
	currentPage := p.CurrentPage
	indexArray := 0
	pagesToShow := make([]int, 1)
	pagesToShow[0] = currentPage
	for indexArray != 4 {
		currentPage--
		pagesToShow = append(pagesToShow, currentPage)
		indexArray++
	}
	currentPage = p.CurrentPage
	for indexArray != 8 {
		currentPage++
		pagesToShow = append(pagesToShow, currentPage)
		indexArray++
	}
	sort.Sort(Alpha(pagesToShow))
	firstPageIndex := 0
	lastPageIndex := 0
	for index, value := range pagesToShow {
		switch {
		case value == 1 && firstPageIndex == 0:
			firstPageIndex = index
		case value == p.TotalPageCount && lastPageIndex == 0:
			lastPageIndex = index
		case value == p.CurrentPage-4 && firstPageIndex == 0:
			firstPageIndex = index
		case value == p.CurrentPage+4 && lastPageIndex == 0:
			lastPageIndex = index
		}
	}
    if p.TotalPageCount == 1 {
        lastPageIndex = firstPageIndex
    }
	pagesToShow = pagesToShow[firstPageIndex : lastPageIndex+1]
	p.PagesToShow = pagesToShow
}
