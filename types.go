package main

import (
	"github.com/spyhere/gator/internal/config"
	"github.com/spyhere/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
