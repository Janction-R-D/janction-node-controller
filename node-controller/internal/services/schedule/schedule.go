package schedule

import (
	"node-controller/common/supports"
	"node-controller/common/wrapper"
	"node-controller/internal/proto"
)

var (
	s *ScheduleService
)

func InitScheduleService(i ScheduleInterface) {
	s = &ScheduleService{
		ISchedule: i,
	}
}

type ScheduleService struct {
	ISchedule ScheduleInterface
}

type ScheduleInterface interface {
}

func GetScheduleService() *ScheduleService {
	return s
}

func (s ScheduleService) GetTaskList(ctx *wrapper.Context, reqBody interface{}) {
	params := reqBody.(*proto.GetTaskListReq)
	supports.SendApiResponse(ctx, resp)
	return

}

func (s ScheduleService) AllocateTask(ctx *wrapper.Context, reqBody interface{}) {
	supports.SendApiResponse(ctx, nil)
	return
}
