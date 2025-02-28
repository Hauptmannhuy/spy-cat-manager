package main

type Spy struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Breed      string `json:"breed"`
	Experience int    `json:"experience"`
	Salary     int    `json:"salary"`
}

type Mission struct {
	Id            int      `json:"id"`
	SpyID         int      `json:"spy_id"`
	Targets       []Target `json:"targets"`
	CompleteState bool     `json:"completed"`
}

type Target struct {
	Id            int    `json:"id"`
	MissionId     int    `json:"mission_id"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	CompleteState bool   `json:"completed"`
}

type Note struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
