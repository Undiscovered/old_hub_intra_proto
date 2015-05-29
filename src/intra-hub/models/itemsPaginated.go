package models

import (
	"github.com/astaxie/beego"
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
	beego.Warn(pagesToShow)
	sort.Sort(Alpha(pagesToShow))
	beego.Warn(pagesToShow)
	currentPageIndex := 0
	firstPageIndex := 0
	lastPageIndex := 0
	for index, value := range pagesToShow {
		switch value {
		case 1:
			firstPageIndex = index
			if p.CurrentPage == 1 {
				currentPageIndex = index
			}
		case p.CurrentPage:
			currentPageIndex = index
		case p.TotalPageCount:
			lastPageIndex = index
		}
	}
	beego.Warn(firstPageIndex, currentPageIndex, lastPageIndex)
	if firstPageIndex != 0 {
        if lastPageIndex != 0 {
            pagesToShow = pagesToShow[firstPageIndex : lastPageIndex]
        } else {
            pagesToShow = pagesToShow[firstPageIndex : currentPageIndex+5]
        }
	} else if lastPageIndex != 0 {
        pagesToShow = pagesToShow[:lastPageIndex]
	} else {
		pagesToShow = pagesToShow[:currentPageIndex+5]
	}
	beego.Warn(pagesToShow)
	p.PagesToShow = pagesToShow
}
