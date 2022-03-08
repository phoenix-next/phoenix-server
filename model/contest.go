package model

type ContestT struct {
	ID      uint64 `json:"ID"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type CreateContestQ struct {
	OrgID      uint64   `json:"orgID"`
	Name       string   `json:"name"`
	Profile    string   `json:"profile"`
	Readable   int      `json:"readable"`
	ProblemIDs []uint64 `json:"problemIDs"`
}

type GetContestA struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	ID         uint64   `json:"ID"`
	Name       string   `json:"name"`
	Profile    string   `json:"profile"`
	ProblemIDs []uint64 `json:"problemIDs"`
}

type UpdateContestQ struct {
	Name       string   `json:"name"`
	Profile    string   `json:"profile"`
	ProblemIDs []uint64 `json:"problemIDs"`
}

type GetContestListA struct {
	Success     bool       `json:"success"`
	Message     string     `json:"message"`
	ContestList []ContestT `json:"contestList"`
}
