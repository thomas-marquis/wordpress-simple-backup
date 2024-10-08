package application

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "Application", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
)
