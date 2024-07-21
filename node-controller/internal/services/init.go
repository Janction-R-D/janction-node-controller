package services

import (
	"node-controller/config"
	"node-controller/internal/services/iam"
	"node-controller/internal/services/repo"
	"node-controller/internal/services/schedule"
)

func Init() {
	iam.InitAuthService(&repo.UserRepo{}, config.CacheInstance, config.CasbinModelPath, config.CasbinPolicyPath)
	schedule.InitScheduleService(&repo.UserRepo{})
}
