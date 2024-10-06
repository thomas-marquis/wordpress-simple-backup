package infrastructure

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "infrastrcture", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
)
