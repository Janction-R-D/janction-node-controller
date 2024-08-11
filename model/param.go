package model

type FormRegisterNode struct {
	NodeID           string `json:"node_id" binding:"required"`
	ArchitectureType string `json:"architecture_type" binding:"required"`
	UseCPU           int    `json:"use_cpu" binding:"required"`
	UseGPU           int    `json:"use_gpu" binding:"required"`
}

type FormPing struct {
	NodeID string `json:"node_id" binding:"required"`
}

type FormGetJob struct {
	NodeID string `json:"node_id" binding:"required"`
}

/********* External Janction Backend Interaction */
type FormGetJobType struct {
	Architecture string `json:"architecture"` // amd64 arm
	UseGPU       int    `json:"use_gpu"`      // 0 , 1 gpu available
	UseCPU       int    `json:"use_cpu"`      // 0 , 1 cpu available
}

type RespGetNonce struct {
	Nonce string `json:"nonce"`
}

type ReqLogin struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	IsNode    bool   `json:"is_node"`
}

type RespLogin struct {
	Token string `json:"token"`
}

/********* */
