package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func InitUsers() (map[string]string, error) {
	var users map[string]string
	jsonFile, err := os.Open("users.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		return map[string]string{}, err
	}
	defer jsonFile.Close()
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		return map[string]string{}, err
	}

	return users, nil
}
