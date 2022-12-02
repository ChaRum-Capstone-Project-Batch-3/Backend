package helper

type BaseResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Page        `json:"pagination"`
}

type Page struct {
	Size        int `json:"size,omitempty"`
	TotalData   int `json:"totalData,omitempty"`
	CurrentPage int `json:"currentPage,omitempty"`
	TotalPage   int `json:"totalPage,omitempty"`
}
