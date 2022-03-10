package model

import "time"

type ContestT struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type ProblemT struct {
	ProblemID   uint64 `json:"problemID"`
	ProblemName string `json:"problemName"`
	Difficulty  int    `json:"difficulty"`
}

type CreateContestQ struct {
	OrgID      uint64    `json:"orgID"`
	Name       string    `json:"name"`
	Profile    string    `json:"profile"`
	Readable   int       `json:"readable"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	ProblemIDs []uint64  `json:"problemIDs"`
}

type GetContestA struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	ID      uint64     `json:"id"`
	Name    string     `json:"name"`
	Profile string     `json:"profile"`
	Problem []ProblemT `json:"problem"`
}

type UpdateContestQ struct {
	Name       string   `json:"name"`
	Profile    string   `json:"profile"`
	ProblemIDs []uint64 `json:"problemIDs"`
}

type GetContestListA struct {
	Success     bool       `json:"success"`
	Message     string     `json:"message"`
	Total       int        `json:"total"`
	ContestList []ContestT `json:"contestList"`
}

type GetOrganizationProblemA struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	ProblemList []Problem `json:"problemList"`
}
