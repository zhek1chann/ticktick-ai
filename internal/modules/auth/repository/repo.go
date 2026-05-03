package repository

import (
	"ticktick-ai/pkg/db"
)

type Repo struct {
	db db.Client
}

func NewAuthRepo(db db.Client) *Repo {
	return &Repo{db: db}
}
