package query

type Request struct {
	Skip  int    `json:"skip"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Order int    `json:"order"`
}
