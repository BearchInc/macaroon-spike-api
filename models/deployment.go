package models

import "github.com/drborges/appx"

type DeploymentStatus string

var (
	DeploymentPending          = DeploymentStatus("Pending")
	DeploymentAwaitingApproval = DeploymentStatus("Awaiting Approval")
	DeploymentApproved         = DeploymentStatus("Approved")
	DeploymentRejected         = DeploymentStatus("Rejected")
)

type Deployment struct {
	appx.Model
	Status       DeploymentStatus `json:"status"`
	Commit       string           `json:"commit"`
	Requester    string           `json:"requester"`
	RequesterGCM string           `json:"gcm"`
	Approver     string           `json:"approver"`
}

func (d *Deployment) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "Deployments",
		Incomplete: true,
	}
}
