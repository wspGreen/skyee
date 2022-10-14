package codec

import "github.com/wspGreen/skyee/slog"

func Marshal(funcName *string, params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			slog.Fatal(err)
		}
	}()
}
