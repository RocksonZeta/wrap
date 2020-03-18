package iriswrap

type JqueryUploadResult struct {
	Files []JqueryUploadFile `json:"files"`
}

func NewJqueryUploadResult(name string, size int64, err string) JqueryUploadResult {
	return JqueryUploadResult{Files: []JqueryUploadFile{{Name: name, Size: size, Error: err}}}
}

type JqueryUploadFile struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Error        string `json:"error"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	DeleteUrl    string `json:"deleteUrl"`
	DeleteType   string `json:"deleteType"`
}
