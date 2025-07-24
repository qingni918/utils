package safe

import (
	"github.com/qingni918/utils/logger_aggregate"
)

func Goroutine(f func(), logger logger_aggregate.Logger) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()

		f()
	}()
}
