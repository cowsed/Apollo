package main

type ServerData struct {
	Temperature struct {
		Tool0 struct {
			Actual float64 `json:"actual"`
			Target float64 `json:"target"`
			Offset int     `json:"offset"`
		} `json:"tool0"`
		Tool1 struct {
			Actual float64     `json:"actual"`
			Target interface{} `json:"target"`
			Offset int         `json:"offset"`
		} `json:"tool1"`
		Bed struct {
			Actual float64 `json:"actual"`
			Target float64 `json:"target"`
			Offset int     `json:"offset"`
		} `json:"bed"`
		History []struct {
			Time  int `json:"time"`
			Tool0 struct {
				Actual float64 `json:"actual"`
				Target float64 `json:"target"`
			} `json:"tool0"`
			Tool1 struct {
				Actual float64     `json:"actual"`
				Target interface{} `json:"target"`
			} `json:"tool1"`
			Bed struct {
				Actual float64 `json:"actual"`
				Target float64 `json:"target"`
			} `json:"bed"`
		} `json:"history"`
	} `json:"temperature"`
	Sd struct {
		Ready bool `json:"ready"`
	} `json:"sd"`
	State struct {
		Text  string `json:"text"`
		Flags struct {
			Operational   bool `json:"operational"`
			Paused        bool `json:"paused"`
			Printing      bool `json:"printing"`
			Cancelling    bool `json:"cancelling"`
			Pausing       bool `json:"pausing"`
			SdReady       bool `json:"sdReady"`
			Error         bool `json:"error"`
			Ready         bool `json:"ready"`
			ClosedOrError bool `json:"closedOrError"`
		} `json:"flags"`
	} `json:"state"`
}
