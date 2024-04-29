package clientAgentConf

import "time"

type DeployStatus string

const (
	PendingDeploy       DeployStatus = "DIRTY"
	Deploying           DeployStatus = "DEPLOYING"
	Deployed            DeployStatus = "DEPLOYED"
	DeployInitiated     DeployStatus = "IN_PROGRESS"
	DeployFailed        DeployStatus = "FAILED"
	DeployStatusUnknown DeployStatus = "UNKNOWN"
)

type ClientConfmap struct {
	ID           string       `json:"id" db:"id"`
	AgentId      string       `json:"agent_id" db:"agent_id"`
	DeployStatus DeployStatus `json:"deployStatus" db:"deploy_status"`
	DeployResult string       `json:"deployResult" db:"deploy_result"`
	LastHash     string       `json:"lastHash" db:"last_hash"`
	LastConf     string       `json:"lastConf" db:"last_config"`
	CreatedAt    time.Time    `json:"createdAt" db:"created_at"`
}
