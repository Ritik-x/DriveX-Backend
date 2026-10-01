package workers

type FileProcessingJob struct {
	FileID     string `json:"file_id"`
	StorageKey string `json:"storage_key"`
	MimeType   string `json:"mime_type"`
}