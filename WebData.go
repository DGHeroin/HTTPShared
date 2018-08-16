package HTTPShared

type ErrorData struct {
	Code int   `json:"code"`
	Msg string `json:"msg"`
}

type GetActionData struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version uint64 `json:"version"`
}

type PutActionData struct {
	Version uint64 `json:"version"`
}