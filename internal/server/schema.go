package server

import "github.com/kgantsov/filo/internal/model"

// RegisterFileInput defines the request body for registering a file.
type RegisterFileInput struct {
	Body struct {
		DirectoryPath string `json:"directory_path" doc:"Directory path where the file is located." required:"true"`
		Filename      string `json:"filename" doc:"Filename of the file." required:"true"`
		FileType      string `json:"file_type" doc:"File type (e.g., file, directory)." required:"true"`
		Size          int64  `json:"size" doc:"File size in bytes."`
		Checksum      string `json:"checksum" doc:"SHA256 checksum." required:"true"`
	}
}

// RegisterFileOutput defines the response body for a successful registration.
type RegisterFileOutput struct {
	Body model.File `json:"body" doc:"The created or updated file record."`
}

// ListFilesInput defines query parameters for filtering files.
type ListFilesInput struct {
	DirectoryPath string `query:"directory_path" doc:"Directory path where the file is located."`
	Filename      string `query:"filename" doc:"Filename of the file."`
	Checksum      string `query:"checksum" doc:"Filter by checksum (optional)."`
	Limit         int    `query:"limit" doc:"Max number of results." default:"20"`
	Offset        int    `query:"offset" doc:"Offset for pagination." default:"0"`
}

// ListFilesOutput defines the response for listing files.
type ListFilesOutput struct {
	Body struct {
		Files []model.File `json:"files" doc:"The list of file records."`
		Total int64        `json:"total" doc:"Total number of matching records."`
	}
}

// GetFileInput defines the path parameter for getting a single file.
type GetFileInput struct {
	ID uint `path:"id" doc:"The unique ID of the file record."`
}

// GetFileOutput defines the response for getting a single file.
type GetFileOutput struct {
	Body model.File `json:"body" doc:"The requested file record."`
}

// DeleteFileInput defines query parameters for deleting a file.
// We use DeleteFileInput + Path as this is what an agent knows, not the DB ID.
type DeleteFileInput struct {
	DirectoryPath string `query:"directory_path" doc:"Directory path of the file to delete." required:"true"`
	Filename      string `query:"filename" doc:"Filename of the file to delete." required:"true"`
}

// DeleteFileOutput defines the empty response for a successful deletion.
type DeleteFileOutput struct {
	// No body needed for a 204 No Content response
}
