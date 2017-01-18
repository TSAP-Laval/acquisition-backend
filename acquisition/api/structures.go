package api

import (
	"github.com/jinzhu/gorm"
)

type Actions_type struct {
	gorm.Model
	ID               string `json:"Id,omitempty"`
	Name             string `json:"Name,omitempty"`
	Description      string `json:"Description,omitempty"`
	Id_movement_type int    `json:"Id_movement_Type,omitempty"`
}
