package types

import (

	//"github.com/polyAnalytica/equityConsultancy/src/db"
	"gorm.io/gorm"
)

type question struct {
	gorm.Model
	Id       uint     `gorm:"primarykey; autoIncrement:true; unique"`
	Question string   `gorm:"size:255;not null;unique" json:"question"`
	Options  []string `gorm:"type:string[];not null" json:"options"`
}
