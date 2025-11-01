package model

import "time"

type File struct {
	ID            uint      `gorm:"primarykey" json:"id" doc:"Unique identifier for the file."`
	DirectoryPath string    `gorm:"uniqueIndex:idx_directory_path_filename" json:"directory_path" doc:"Directory path where the file is located."`
	Filename      string    `gorm:"uniqueIndex:idx_directory_path_filename" json:"filename" doc:"Name of the file."`
	FileType      string    `gorm:"index" json:"file_type" doc:"Type of the file."`
	Size          int64     `json:"size" doc:"File size in bytes."`
	Checksum      string    `json:"checksum" gorm:"index" doc:"SHA256 checksum of the file content."`
	CreatedAt     time.Time `json:"createdAt" doc:"Timestamp when the record was created (first discovery)."`
	UpdatedAt     time.Time `json:"updatedAt" doc:"Timestamp when the record was last updated (last seen)."`
}
