package main

import (
	"fmt"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

func logtype(a interface{}, b ...interface{}) {
	msg := fmt.Sprintf("1: %T", a)
	for i, b := range b {
		msg += fmt.Sprintf(", %d: %T", i+2, b)
	}
	logger.Warn(msg)

}
