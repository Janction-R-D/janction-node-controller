package model

import "time"

type NodeRegistration struct {
	NodeID           string    `json:"node_id" gorm:"primarykey;type:varchar(40)"`
	ArchitectureType string    `json:"architecture_type" gorm:"type:varchar(40)"`
	UseCPU           int       `json:"use_cpu"`
	UseGPU           int       `json:"use_gpu"`
	PingAt           time.Time `json:"ping_at"`
	Status           string    `json:"status" gorm:"-"`
}
