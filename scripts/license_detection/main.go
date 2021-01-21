package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"license_detection/scripts"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	args := os.Args[1:]
	projectsPath := args[0]
	outputFilePath := args[1]
	projectCachePath := args[2]

	startTime := time.Now()

	fmt.Printf("Processing projects path: %s\n", projectsPath)
	resultLicenses := scripts.GetLicensesForProjects(projectsPath, projectCachePath)

	// Write to JSON
	resultJson, _ := json.MarshalIndent(resultLicenses, "", "\t")
	_ = ioutil.WriteFile(outputFilePath, resultJson, 0644)

	fmt.Printf("Script elapsed: %s\n", time.Since(startTime))
}