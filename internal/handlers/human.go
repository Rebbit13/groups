package handlers

import "go.uber.org/zap"

type Human struct {
	service HumanService
	logger  *zap.Logger
}
