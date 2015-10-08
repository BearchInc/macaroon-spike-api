package models

type DeploymentStatus string

var (
	DeploymentPending          = DeploymentStatus("Pending")
	DeploymentAwaitingApproval = DeploymentStatus("Awaiting Approval")
	DeploymentApproved         = DeploymentStatus("Approved")
	DeploymentRejected         = DeploymentStatus("Rejected")
)

type Deployment struct {
	Status      DeploymentStatus
	UserID      string
	Commit      string
	ApproversID []string
}
