package goutils

import "time"

func RetryOperation(operation func() error, maxRetries int) (err error) {
	for i := 0; i < maxRetries; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}

		if operation() == nil {
			return nil
		}
	}
	return err
}
