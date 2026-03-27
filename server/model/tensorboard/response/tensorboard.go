package response

// TensorBoardItem TensorBoard列表项
type TensorBoardItem struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	LogsPath          string `json:"logsPath"`
	Status            string `json:"status"`
	CreationTimestamp string `json:"creationTimestamp"`
}

// GetTensorBoardListResp 获取TensorBoard列表响应
type GetTensorBoardListResp struct {
	List  []TensorBoardItem `json:"list"`
	Total int64             `json:"total"`
}

// AddTensorBoardResp 创建TensorBoard响应
type AddTensorBoardResp struct {
	Name string `json:"name"`
}

// UpdateTensorBoardResp 更新TensorBoard响应
type UpdateTensorBoardResp struct {
	Name string `json:"name"`
}

// DeleteTensorBoardResp 删除TensorBoard响应
type DeleteTensorBoardResp struct {
	Name string `json:"name"`
}
