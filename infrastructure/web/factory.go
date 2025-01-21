package web

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	errInvalidWebServerInstance = errors.New("invalid router server instance")
)

const (
	InstanceGin int = iota
	InstanceEcho
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func NewServer(instance int, db *gorm.DB) (Server, error) {
	config := NewConfigWeb()
	switch instance {
	// case InstanceGin:
	// 	return NewGinServer(config.Host, config.Port, config.CorsAllowOrigins, db)
	case InstanceEcho:
		return NewEchoServer(config.Host, config.Port, db)
	default:
		panic(errInvalidWebServerInstance)
	}
}
