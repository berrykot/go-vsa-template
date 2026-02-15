package health

// MyRequest — тело
type MyRequest struct {
	Data string `json:"data"`
}

// MyResponse — JSON-ответ
type MyResponse struct {
	Data string `json:"data"`
}
