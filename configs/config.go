package configs

import (
	"encoding/json"
	"flag"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

type Conf struct {
	SheetsAdr string
	Token     string
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
		os.Getenv("TOKEN_A"),
	}
}

func AddUser(users map[string]string, user string) error  {
	filePath:= "users.json"
	users[user] = "0"
	jsonUsers,err:= json.Marshal(users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, jsonUsers, 0644)
	if err != nil {
		return err
	}
	return nil
}
