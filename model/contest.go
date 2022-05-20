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
	Result      int    `json:"result"`
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
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Name      string     `json:"name"`
	Profile   string     `json:"profile"`
	StartTime string     `json:"startTime"`
	EndTime   string     `json:"endTime"`
	Problem   []ProblemT `json:"problem"`
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

type ProblemPenalty struct {
	ProblemID uint64 `json:"problemID"`
	UserID    uint64 `json:"userID"`
	TryCount  int    `json:"tryCount"`
	// 0 AC , 1 WA , 2 TLE, 3 RE, -1 未做
	Status     int `json:"status"`
	PenaltySum int `json:"penaltySum"`
}

type RankT struct {
	Rank      int              `json:"rank"`
	Name      string           `json:"name"`
	UserID    uint64           `json:"userID"`
	PassCount int              `json:"passCount"`
	Penalty   int              `json:"penalty"`
	Problems  []ProblemPenalty `json:"problems"`
}

type GetRankingListA struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Rank    []RankT `json:"rank"`
}
