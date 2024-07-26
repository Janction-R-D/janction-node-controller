package logic

import (
	"janction/dao/postgres"
	"janction/model"
	"time"
)

func RegisterNode(p *model.FormRegisterNode) (err error) {
	err = postgres.UpsertNodeRegistration(&model.NodeRegistration{
		NodeID:           p.NodeID,
		ArchitectureType: p.ArchitectureType,
		UseCPU:           p.UseCPU,
		UseGPU:           p.UseGPU,
	})
	return
}

func Ping(nodeID string) (err error) {
	err = postgres.UpsertNodeRegistration(&model.NodeRegistration{
		NodeID: nodeID,
		PingAt: time.Now(),
	})
	return
}
