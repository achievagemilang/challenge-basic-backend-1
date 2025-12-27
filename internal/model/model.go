package model

type WebResponse[T any] struct {
	Ok   bool `json:"ok"`
	Data T    `json:"data"`
}

type ErrorResponse struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err"`
	Msg string `json:"msg"`
}

type PageResponse[T any] struct {
	Ok     bool         `json:"ok"`
	Data   []T          `json:"data,omitempty"`
	Paging PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
