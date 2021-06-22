package main

type SimpleData struct {
	Filename           string
	EstimatedPrintTime float64
	Completion         float64
	State              string

	ToolActual, ToolTarget float64
	BedActual, BedTarget   float64
	Printing               bool
}

type JobData struct {
	Job struct {
		AveragePrintTime   float64 `json:"averagePrintTime"`
		EstimatedPrintTime float64 `json:"estimatedPrintTime"`
		Filament           struct {
			Tool0 struct {
				Length float64 `json:"length"`
				Volume float64 `json:"volume"`
			} `json:"tool0"`
		} `json:"filament"`
		File struct {
			Date    int    `json:"date"`
			Display string `json:"display"`
			Name    string `json:"name"`
			Origin  string `json:"origin"`
			Path    string `json:"path"`
			Size    int    `json:"size"`
		} `json:"file"`
		LastPrintTime float64 `json:"lastPrintTime"`
		User          string  `json:"user"`
	} `json:"job"`
	Progress struct {
		Completion          float64     `json:"completion"`
		Filepos             int         `json:"filepos"`
		PrintTime           int         `json:"printTime"`
		PrintTimeLeft       int         `json:"printTimeLeft"`
		PrintTimeLeftOrigin interface{} `json:"printTimeLeftOrigin"`
	} `json:"progress"`
	State string `json:"state"`
}
