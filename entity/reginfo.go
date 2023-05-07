package entity

type RegInfo struct {
	Username string `json:"username"`
	Passwd string `json:"passwd"`
}
type RegStatus struct {
	Success bool `json:"success"`
}

type LogInfo RegInfo
type LogStatus RegStatus
