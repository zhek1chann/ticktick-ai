package service

import (
	"ticktick-ai/pkg/db"
	"ticktick-ai/pkg/jwt"
)

type AuthRepo interface {
	userRepo
}

type Service struct {
	txManager  db.TxManager
	repo       AuthRepo
	jwtManager *jwt.Manager
}

func NewService(txManager db.TxManager, repo AuthRepo, jwtManager *jwt.Manager) *Service {
	return &Service{
		txManager:  txManager,
		repo:       repo,
		jwtManager: jwtManager,
	}
}
