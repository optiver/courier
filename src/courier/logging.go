package main

import (
	"errors"
	"log"
)

type LogLevel int

const (
	Verbose LogLevel = iota
	Normal
	Quiet
)

var GlobalLogLevel LogLevel

func SetupLogging(quiet, verbose bool) error {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	if quiet && verbose {
		return errors.New("cannot specify both quiet and verbose")
	}
	if quiet {
		GlobalLogLevel = Quiet
	} else if verbose {
		GlobalLogLevel = Verbose
	} else {
		GlobalLogLevel = Normal
	}
	return nil
}

func LogDebug(format string, args ...interface{}) {
	if GlobalLogLevel != Verbose {
		return
	}
	log.Printf("[DEBUG] "+format+"\n", args...)
}

func LogInfo(format string, args ...interface{}) {
	if GlobalLogLevel == Quiet {
		return
	}
	if cmdLineArgs.colour {
		log.Printf("\u001B[32m[INFO ]\u001B[0m "+format+"\n", args...)
	} else {
		log.Printf("[INFO ] "+format+"\n", args...)
	}
}

func LogWarn(format string, args ...interface{}) {
	if GlobalLogLevel == Quiet {
		return
	}
	if cmdLineArgs.colour {
		log.Printf("\u001B[33m[WARN ]\u001B[0m "+format+"\n", args...)
	} else {
		log.Printf("[WARN ] "+format+"\n", args...)
	}
}

func LogError(format string, args ...interface{}) {
	if GlobalLogLevel == Quiet {
		return
	}
	if cmdLineArgs.colour {
		log.Printf("\u001B[31m[ERROR]\u001B[0m "+format+"\n", args...)
	} else {
		log.Printf("[ERROR] "+format+"\n", args...)
	}
}
