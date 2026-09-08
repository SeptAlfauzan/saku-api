package domain

type OCRRequest struct {
	ImageBase64 string `json:"image"`
	MimeType    string `json:"mime_type"`
}
