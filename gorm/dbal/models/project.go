package models

// Code generated by a tool. DO NOT EDIT.
// https://sqldalmaker.sourceforge.net/

type Project struct {
	PID   int64  `json:"p_id" gorm:"column:p_id;primaryKey;autoIncrement;not null"`
	PName string `json:"p_name" gorm:"column:p_name;size:256;not null"`
}
