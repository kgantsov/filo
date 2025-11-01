package server

import "gorm.io/gorm"

type APIHandler struct {
	DB *gorm.DB
}
