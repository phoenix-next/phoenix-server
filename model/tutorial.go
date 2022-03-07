package model

type CreateTutorialQ struct {
	OrgID    uint64 `json:"orgID"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	Readable int    `json:"readable"`
	Writable int    `json:"writable"`
}

type GetTutorialA struct {
	OrgID        uint64 `json:"orgID"`
	CreatorID    uint64 `json:"creatorID"`
	CreatorName  uint64 `json:"creatorName"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	Version      int    `json:"version"`
	TutorialPath string `json:"tutorialPath"` // 教程的下载路径
}

type UpdateTutorialQ struct {
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	Readable int    `json:"readable"`
	Writable int    `json:"writable"`
}

type GetTutorialVersionA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version int    `json:"version"`
}

type GetTutorialListA struct {
	Success      bool       `json:"success"`
	Message      string     `json:"message"`
	Total        int        `json:"total"`
	TutorialList []Tutorial `json:"tutorialList"`
}
