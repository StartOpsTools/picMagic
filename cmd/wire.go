//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/qx66/picMagic/internal/biz"
	"github.com/qx66/picMagic/internal/conf"
	"go.uber.org/zap"
)

func initApp(*conf.Bootstrap, *zap.Logger) (*app, error) {
	panic(wire.Build(
		biz.ProviderSet,
		newApp))
}
