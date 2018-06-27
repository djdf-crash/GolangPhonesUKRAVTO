package handlers

import "db"

type ResponseModelEmployee struct {
	Result    bool          `json:"result"`
	Error     string        `json:"error"`
	Employees []db.Employee `json:"body"`
}

type ResponseModelOrganization struct {
	Result       bool              `json:"result"`
	Error        string            `json:"error"`
	Organization []db.Organization `json:"body"`
}

type ResponseModelToken struct {
	Result bool   `json:"result"`
	Error  string `json:"error"`
	Token  string `json:"body"`
}
