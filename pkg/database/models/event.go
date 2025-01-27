package models

import (
	"time"

	"gorm.io/datatypes"
)

type BaseLocalPolicyEvent struct {
	EventName   string         `gorm:"column:event_name;type:varchar(63);not null" json:"eventName"`
	PolicyID    string         `gorm:"column:policy_id;type:uuid;not null" json:"policyId"`
	Message     string         `gorm:"column:message;type:text" json:"message"`
	LeafHubName string         `gorm:"size:63;not null" json:"-"`
	Reason      string         `gorm:"column:reason;type:text" json:"reason"`
	Count       int            `gorm:"column:count;type:integer;not null;default:0" json:"count"`
	Source      datatypes.JSON `gorm:"column:source;type:jsonb" json:"source"`
	CreatedAt   time.Time      `gorm:"column:created_at;default:now();not null" json:"createdAt"`
	Compliance  string         `gorm:"column:compliance" json:"compliance"`
}

type LocalClusterPolicyEvent struct {
	BaseLocalPolicyEvent
	ClusterID string `gorm:"column:cluster_id;type:uuid;not null" json:"clusterId"`
}

func (LocalClusterPolicyEvent) TableName() string {
	return "event.local_policies"
}

type LocalRootPolicyEvent struct {
	BaseLocalPolicyEvent
}

func (LocalRootPolicyEvent) TableName() string {
	return "event.local_root_policies"
}
