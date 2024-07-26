package postgres

import (
	"janction/model"

	"gorm.io/gorm/clause"
)

func UpsertNodeRegistration(nodeInfo *model.NodeRegistration) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ping_at"}),
	}).Create(nodeInfo).Error
}

func GetNodeRegistrationByNodeID(nodeID string) (*model.NodeRegistration, error) {
	var node model.NodeRegistration
	if err := db.Where("node_id = ?", nodeID).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}
