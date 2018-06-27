package db

import (
	"time"
)

type Employee struct {
	ID               uint   `gorm:"primary_key"`
	FullName         string `json:"full_name" xlsx:"1"`
	Post             string `json:"post" xlsx:"0"`
	Email            string `json:"email" xlsx:"2"`
	ContactInfo      string `xlsx:"3"`
	Phone            string `xlsx:"4"`
	OrganizationName string `json:"organization_name"`
	OrganizationID   uint
	RealPhone        string    `json:"real_phone"`
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
