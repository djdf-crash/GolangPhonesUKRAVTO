package db

import (
	"time"
)

type Employee struct {
	ID               uint      `gorm:"primary_key"`
	FullName         string    `json:"full_name"`
	Post             string    `json:"post"`
	Email            string    `json:"email"`
	PhoneMobile      string    `json:"phone_mobile"`
	Phone            string    `json:"phone"`
	Department       string    `json:"department"`
	Section          string    `json:"section"`
	OrganizationName string    `json:"organization_name"`
	Category         string    `json:"category"`
	OrganizationID   uint      `json:"organization_id"`
	Address          string    `json:"address"`
	LastUpdate       time.Time `json:"last_update"`
	IsDelete         bool      `json:"delete"`
}

// set User's table name to be `profiles`
func (Employee) TableName() string {
	return "employees"
}

type User struct {
	Token      string `gorm:"primary_key"`
	DeviceID   string
	Email      string
	LastUpdate time.Time
}

// set User's table name to be `profiles`
func (User) TableName() string {
	return "users"
}

type Organization struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	LastUpdatePhones time.Time
	IsDelete         bool
}

func (Organization) TableName() string {
	return "organizations"
}
