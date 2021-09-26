package configs

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
)

type Conf struct {
	SheetsAdr string
}

func InitConf() *Conf {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	return envVar(local)
}

func envVar(local bool) *Conf {
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			println(err.Error())
			return &Conf{}
		}
	}
	return &Conf{
		os.Getenv("SHEETSAPI_ID"),
	}
}