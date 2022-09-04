package panicfiles

type Config struct {
	Program    string `json:"program"`
	Version    string `json:"version"`
	Panicfile  string `json:"panicfile"`
	Fileformat string `json:"fileformat"`
}
const version = "2.1.5"
const Panicconfigname =  "PanicCfg.json"
var panicFile = "./Exception/"