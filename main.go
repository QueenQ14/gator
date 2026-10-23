package main

import (
	"fmt"

	"github.com/QueenQ14/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("QueenQ")
	cfg = config.Read()
	fmt.Println(cfg.DBUrl, cfg.UserName)

}
