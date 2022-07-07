package template

import (
	"log"
)

type Config struct {
	SSLEMAIL    string
	CLIENTLOG   *log.Logger
	HOST        string
	PORT        string
	LOGNAME     string
	MAXLOGSTORE int
}

type User struct {
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	PASSWORD string `json:"pwd"`
}
