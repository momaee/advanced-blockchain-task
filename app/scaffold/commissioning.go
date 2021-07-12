package scaffold

import (
	"backend_task/clients/redis"
	"backend_task/conf"
	"fmt"
	"html"
)

func (s *skeleton) commissioning() (err error) {

	// err = postgres.Storage.Connect(conf.GetAppConfig())
	// if err != nil {
	// 	return fmt.Errorf("%v starting postgres failed: why? %v", html.UnescapeString("&#x274C;"), err)
	// }

	err = redis.Storage.Connect(conf.GetAppConfig())
	if err != nil {
		return fmt.Errorf("%v starting redis failed: why? %v", html.UnescapeString("&#x274C;"), err)
	}

	go s.eth()
	return
}
