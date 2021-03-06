package config

import (
	"flag"
	"os"
	"path"
)

var Role string //role can either be "router", "shard", "client" or "standalone"
var Leader_addr string
var Password string
var AES_key string
var Host string
var Port string
var Client_port string
var App_port string
var Save_disk bool
var Debug bool
var Auto_shard bool
var Secure_connect bool

var Max_Worker uint64
var Auto_Max_Worker bool

var Enable_https bool
var Encrypt_disk bool
var Encrypt_bit_length bool
var Enable_ML bool
var Javascript_timeout int // milli seconds
var ID uint64

var RootDir string

func init() {
	flag.StringVar(&Role, "role", "standalone", "Specify role. Default is router, option: shard, client")
	flag.StringVar(&Leader_addr, "ip", "", "Specify router ip. Default is stand alone router")
	flag.StringVar(&Password, "pwd", "a empty password", "Specify password. Default is empty")
	flag.StringVar(&AES_key, "Aes_key", "a empty password", "Specify aes key. Default is empty")
	flag.StringVar(&Host, "host", "localhost", "specify host addr")
	flag.StringVar(&Port, "port", "8080", "specify hosting port")
	flag.StringVar(&Client_port, "client_port", ":54620", "specify client port")
	flag.StringVar(&App_port, "app_port", "36730", "specify application port")

	flag.Uint64Var(&Max_Worker, "max_worker", 1200, "maximum concurrent workers")
	flag.BoolVar(&Auto_Max_Worker, "disable_auto_max_worker", true, "automatic change the max_worker base on workload")
	//flag.StringVar(&sql_file, "sqlfile", ":memory:", "specify sql file path")
	flag.BoolVar(&Save_disk, "disk", false, "save copy to disk")
	flag.BoolVar(&Debug, "debug", false, "developer debug option")
	flag.BoolVar(&Enable_ML, "ML", false, "Enable built in machine learning service")
	flag.BoolVar(&Secure_connect, "sc", false, "specify to use securite connection and the bit width")
	flag.IntVar(&Javascript_timeout, "jstimeout", 500, "javascript vm timeout duration in milliseconds")

	h, _ := os.UserConfigDir()
	flag.StringVar(&RootDir, "rootdir", path.Join(h, "dejavuDB"), "specify application root directory")
	os.Chdir(RootDir)
	flag.Parse()
}
