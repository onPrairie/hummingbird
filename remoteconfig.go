package main

type RemoteCon struct {
	// Version     string `json:"Version"`
	// StationName string `json:"StationName"`
	// Ftpuname    string "Ftpuname"
	// Ftppass     string "Ftppass"
	StationCode   string `json:"StationCode"`
	StationType   int    `json:"StationType"`
	BaudRate      string `json:"BaudRate"`
	Com           int    `json:"Com"`
	Cap0          string `json:"cap0"`
	Cap1          string `json:"cap1"`
	Cap2          string `json:"cap2"`
	Cap3          string `json:"cap3"`
	Path          string `json:"Path"`
	Ftpip         string `json:"Ftpip"`
	Updbind       int    `json:"Updbind"`
	BaudRateOther string `json:"BaudRateOther"`
	ComOther      int    `json:"ComOther"`
	IsComOther    bool   `json:"IsComOther"`
	Lasertcp      string `json:"lasertcp"`
	Maxweighte    int    `json:"maxweighte"`
	Bindaddress   string `json:"bindaddress"`
	Url           string `json:"url"`
}
