package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	ShowVersion    bool   `default:"false" help:"Print build information and exit"`
	ExcludeLetters string `default:"" help:"Which letters to exclude from file count"`
	FolderToScan   string `default:"/home/vagrant/Source" help:""`
}

const (
	integrationName = "es.kangmak.letter-count"
)

var (
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""
)

func main() {
	log.Debug("Starting letter counter")
	defer log.Debug("Letter counter exited")

	args := argumentList{}
	logger := log.NewStdErr(args.Verbose)

	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args), integration.Logger(logger))
	if err != nil {
		log.Fatal(err)
	}

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			strings.Title(strings.Split(integrationName, ".")[2]),
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
	}

	log.SetupLogging(args.Verbose)

	//e := i.LocalEntity()
	e, err := i.Entity("test-entity", "test-namespace")
	if err != nil {
		log.Fatal(err)
	}

	if args.HasInventory() {
		log.Debug("Fetching inventory for '%s' integration", integrationName)
		if err != nil {
			log.Fatal(err)
		}

		err = populateInventory(e.Inventory, args)
		if err != nil {
			log.Fatal(err)
		}
	}

	if args.HasMetrics() {
		log.Debug("Fetching metrics for '%s' integration", integrationName)

		ms := e.NewMetricSet("LetterCounter")
		if err != nil {
			log.Fatal(err)
		}

		err = populateMetrics(ms, args)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = i.Publish()
	if err != nil {
		log.Fatal(err)
	}
}
