package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type setting struct {
	portNum    string
	isDebug    *bool
	dotEnvPath *string
	dbPath     string
}

func getSetting() setting {
	var st setting
	st.isDebug = flag.Bool("debug", false, "enable debug mode")
	st.dotEnvPath = flag.String("dotenv", ".env", "path to dotenv file")
	flag.Parse()

	err := godotenv.Load(*st.dotEnvPath)
	if err != nil {
		log.Printf("failed to load dotenv file: %v", err)
		os.Exit(1)
	}

	*st.isDebug = *st.isDebug || os.Getenv("DEBUG") == "true"
	if *st.isDebug {
		st.portNum = os.Getenv("DEBUG_PORT")
	} else {
		st.portNum = os.Getenv("PORT")
	}
	if st.portNum == "" {
		st.portNum = "80"
	}

	st.dbPath = os.Getenv("DB_PATH")

	return st
}
