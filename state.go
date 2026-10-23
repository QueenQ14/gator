package main

import (
	"github.com/QueenQ14/gator/internal/config"
	"github.com/QueenQ14/gator/internal/database"
)

type state struct {
	db     *database.Queries
	Config *config.Config
}
