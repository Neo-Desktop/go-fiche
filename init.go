package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// GitCommit holds the git commit message on compile time
var GitCommit string

// GitBranch holds the active git branch on compile time
var GitBranch string

// GitState holds the dirty-state on compile time
var GitState string

// BuildDate date of compile
var BuildDate string

// Version version string as contained in VERSION
var Version string

// init parses the command line arguments with pflag and viper
func init() {
	pflag.BoolP("help", "h", false, "Prints this help message")
	pflag.StringP("output", "o", "./code", "Relative or absolute path to the directory where you want to store user-posted pastes.")
	pflag.StringP("domain", "d", "localhost", "This will be used as a prefix for an output received by the client. Value will be prepended with http[s].")
	pflag.IntP("port", "p", 9999, "Port in which the service should listen on.")
	pflag.BoolP("https", "S", false, fmt.Sprintf("If set, %s returns url with https prefix instead of http.", AppName))
	pflag.IntP("buffer", "B", 32768, "This parameter defines size of the buffer used for getting data from the user. Maximum size (in bytes) of all input files is defined by this value.")
	pflag.StringP("log", "l", "", "Log file. This file has to be user-writable.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("help") {
		fmt.Printf("%s! - Version %s, Built on %s from Git tag [%s:%s-%s)\n", AppName, Version, BuildDate, GitBranch, GitCommit, GitState)
		pflag.Usage()
		os.Exit(2)
	}

	if viper.GetBool("https") {
		viper.Set("uriprefix", "https")
	} else {
		viper.Set("uriprefix", "http")
	}
}
