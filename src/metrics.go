package main

import (
	"errors"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	metricTemplate = "letter.count."
)

// populateInventory This integration does not have inventory, so I simply return nil.
func populateMetrics(metricSet *metric.Set, args argumentList) error {
	metrics, err := fetchMetrics(args.FolderToScan, args.ExcludeLetters)
	if err != nil {
		return err
	}

	for n, m := range metrics {
		err = metricSet.SetMetric(n, m, metric.GAUGE)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchMetrics(folderToScan string, excludeLetters string) (metrics map[string]float64, errWalking error){
	metrics = make(map[string]float64)

	errWalking = filepath.WalkDir(folderToScan, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !errors.Is(err, fs.ErrPermission) { // Go 1.13
		//if err != nil && err != fs.ErrPermission {
			return err
		} else if errors.Is(err, fs.ErrPermission) {
			return nil
		} else if d.IsDir() || isExcludedFile(d.Name(), excludeLetters) {
			return nil
		}

		letter := strings.ToLower(string(d.Name()[0]))

		if _, ok := metrics[metricTemplate+letter]; !ok {
			metrics[metricTemplate+letter] = 0
		}

		metrics[metricTemplate+letter]++

		return nil
	})

	return
}

func isExcludedFile(filename string, excludeLetters string) bool {
	el := strings.ToLower(excludeLetters)

	for _, filter := range el {
		if []rune(filename)[0] == filter {
			return true
		}
	}

	return false
}
