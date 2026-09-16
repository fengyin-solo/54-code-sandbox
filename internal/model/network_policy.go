package model

import (
	"strings"
	"time"
)

const (
	NetPolicyActive   = "active"
	NetPolicyDisabled = "disabled"
)

type NetworkPolicy struct {
	ID            string    `json:"id"`
	SandboxID     string    `json:"sandbox_id"`
	AllowedHosts  []string  `json:"allowed_hosts"`
	BlockedHosts  []string  `json:"blocked_hosts"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (n *NetworkPolicy) Validate() error {
	n.SandboxID = strings.TrimSpace(n.SandboxID)
	if n.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if n.Status == "" {
		n.Status = NetPolicyActive
	}
	if n.Status != NetPolicyActive && n.Status != NetPolicyDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func (n *NetworkPolicy) IsHostAllowed(host string) bool {
	host = strings.TrimSpace(host)
	for _, h := range n.BlockedHosts {
		if strings.TrimSpace(h) == host {
			return false
		}
	}
	if len(n.AllowedHosts) == 0 {
		return true
	}
	for _, h := range n.AllowedHosts {
		if strings.TrimSpace(h) == host {
			return true
		}
	}
	return false
}

type NetworkPolicyFilter struct {
	SandboxID string
	Status    string
}

func (f NetworkPolicyFilter) Match(n *NetworkPolicy) bool {
	if f.SandboxID != "" && n.SandboxID != f.SandboxID {
		return false
	}
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	return true
}
