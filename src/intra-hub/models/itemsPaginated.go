package models

type ItemPaginated struct {
	Items          interface{} `json:"items"`
	ItemCount      int         `json:"itemCount"`
	TotalItemCount int         `json:"totalItemCount"`
	CurrentPage    int         `json:"currentPage"`
	TotalPageCount int         `json:"totalPageCount"`
}
