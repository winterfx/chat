package app

import (
	"chat/internal/transport/rest"
)

func Run() {
	//1. init config
	//2. init logger
	//3. init dependencies, e.g. db, cache, etc.
	rest.LaunchApiServer()
}
