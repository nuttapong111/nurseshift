package entities

import (
	"time"

	"github.com/google/uuid"
)

// PackageType represents the subscription package type
type PackageType string

const (
	PackageTypeTrial      PackageType = "trial"
	PackageTypeStandard   PackageType = "standard"
	PackageTypeEnterprise PackageType = "enterprise"
)

// Organization represents an organization in the system
type Organization struct {
	ID                    uuid.UUID   `json:"id"`
	Name                  string      `json:"name"`
	Description           *string     `json:"description,omitempty"`
	Email                 string      `json:"email"`
	Phone                 *string     `json:"phone,omitempty"`
	Address               *string     `json:"address,omitempty"`
	Website               *string     `json:"website,omitempty"`
	LicenseNumber         *string     `json:"licenseNumber,omitempty"`
	SubscriptionExpiresAt *time.Time  `json:"subscriptionExpiresAt,omitempty"`
	PackageType           PackageType `json:"packageType"`
	MaxUsers              int         `json:"maxUsers"`
	MaxDepartments        int         `json:"maxDepartments"`
	CreatedAt             time.Time   `json:"createdAt"`
	UpdatedAt             time.Time   `json:"updatedAt"`
}

// IsSubscriptionActive checks if the organization's subscription is still active
func (o *Organization) IsSubscriptionActive() bool {
	if o.SubscriptionExpiresAt == nil {
		return true // Unlimited subscription
	}
	return time.Now().Before(*o.SubscriptionExpiresAt)
}

// GetPackageDisplayName returns the display name for the package type
func (o *Organization) GetPackageDisplayName() string {
	switch o.PackageType {
	case PackageTypeTrial:
		return "แพ็คเกจทดลองใช้"
	case PackageTypeStandard:
		return "แพ็คเกจมาตรฐาน"
	case PackageTypeEnterprise:
		return "แพ็คเกจระดับองค์กร"
	default:
		return "ไม่ระบุ"
	}
}
