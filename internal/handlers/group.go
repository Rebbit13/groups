package handlers

import "go.uber.org/zap"

type Group struct {
	service GroupService
	logger  *zap.Logger
}
