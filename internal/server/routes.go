package server

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kgantsov/filo/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RegisterFile implements the "upsert" logic.
// If the file (directory path + filename) exists, it's updated.
// If not, it's created.
func (h *APIHandler) RegisterFile(ctx context.Context, input *RegisterFileInput) (*RegisterFileOutput, error) {
	var file model.File

	// Check if the record already exists
	result := h.DB.WithContext(ctx).Where(
		"directory_path = ? AND filename = ?", input.Body.DirectoryPath, input.Body.Filename,
	).First(&file)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// --- Create new record ---
		file = model.File{
			DirectoryPath: input.Body.DirectoryPath,
			Filename:      input.Body.Filename,
			FileType:      input.Body.FileType,
			Size:          input.Body.Size,
			Checksum:      input.Body.Checksum,
			// CreatedAt and UpdatedAt are set by GORM by default
		}
		if err := h.DB.WithContext(ctx).Create(&file).Error; err != nil {
			log.Error().Msgf("Error creating file record: %v", err)
			return nil, huma.Error500InternalServerError("Failed to create file record")
		}
	} else if result.Error == nil {
		// --- Update existing record ---
		// We only update fields that might change
		file.Size = input.Body.Size
		file.Checksum = input.Body.Checksum
		// GORM automatically updates 'UpdatedAt' on Save
		if err := h.DB.WithContext(ctx).Save(&file).Error; err != nil {
			log.Error().Msgf("Error updating file record: %v", err)
			return nil, huma.Error500InternalServerError("Failed to update file record")
		}
	} else {
		// Other database error
		log.Error().Msgf("Error finding file record: %v", result.Error)
		return nil, huma.Error500InternalServerError("Database error")
	}

	return &RegisterFileOutput{Body: file}, nil
}

// ListFiles retrieves a paginated and filtered list of files.
func (h *APIHandler) ListFiles(ctx context.Context, input *ListFilesInput) (*ListFilesOutput, error) {
	var files []model.File
	var total int64

	// Start building the query
	query := h.DB.WithContext(ctx).Model(&model.File{})

	if input.DirectoryPath != "" {
		query = query.Where("directory_path = ?", input.DirectoryPath)
	}
	if input.Filename != "" {
		query = query.Where("filename = ?", input.Filename)
	}
	if input.Checksum != "" {
		query = query.Where("checksum = ?", input.Checksum)
	}

	// Get total count for pagination
	if err := query.Count(&total).Error; err != nil {
		log.Error().Msgf("Error counting files: %v", err)
		return nil, huma.Error500InternalServerError("Failed to count files")
	}

	// Get the paginated results
	if err := query.Limit(input.Limit).Offset(input.Offset).Find(&files).Error; err != nil {
		log.Error().Msgf("Error listing files: %v", err)
		return nil, huma.Error500InternalServerError("Failed to list files")
	}

	return &ListFilesOutput{
		Body: struct {
			Files []model.File `json:"files" doc:"The list of file records."`
			Total int64        `json:"total" doc:"Total number of matching records."`
		}{
			Files: files,
			Total: total,
		},
	}, nil
}

// GetFileByID retrieves a single file record by its primary key.
func (h *APIHandler) GetFileByID(ctx context.Context, input *GetFileInput) (*GetFileOutput, error) {
	var file model.File

	result := h.DB.WithContext(ctx).First(&file, input.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, huma.Error404NotFound("file not found")
	} else if result.Error != nil {
		log.Error().Msgf("Error getting file by ID %d: %v", input.ID, result.Error)
		return nil, result.Error
	}

	return &GetFileOutput{Body: file}, nil
}

// DeleteFile deletes a file record based on its directory path and filename.
func (h *APIHandler) DeleteFile(ctx context.Context, input *DeleteFileInput) (*DeleteFileOutput, error) {

	result := h.DB.WithContext(ctx).Where(
		"directory_path = ? AND filename = ?", input.DirectoryPath, input.Filename,
	).Delete(&model.File{})

	if result.Error != nil {
		log.Error().Msgf("Error deleting file %s@%s: %v", input.DirectoryPath, input.Filename, result.Error)
		return nil, huma.Error500InternalServerError("Failed to delete file record")
	}

	if result.RowsAffected == 0 {
		// No record was found to delete
		return nil, huma.Error404NotFound("file not found at specified directory and filename")
	}

	return &DeleteFileOutput{}, nil
}
