package codec

import "github.com/wspGreen/skyee/log"

func Marshal(funcName *string, params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
}
