package scripts

import (
	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
	"path/filepath"
	"regexp"
	"fmt"
	"os"
	// "sync"
	"io/ioutil"
	"encoding/json"
)

// Describes license file metadata entry (couples license type with it's associated
// confidence)
type FileData struct {
	LicenseType string
	Confidence float32
}

// Retrieves a map of all projects mapping to all of it's subfolders (found recursively)
// The projects are expected to be in format of Github repos: root/author/repo
// Keys represent the relative path ("root/author/repo"), while values represent
// relative paths of subfolders (e.g. ["root/author/repo/sub", "root/author/repo/sub/sub"])
func getProjects(root string) map[string][]string {
	projectMap := make(map[string][]string)
	projectPattern := regexp.MustCompile(fmt.Sprintf("%s\\/[^\\/]+\\/[^\\/]+$", root))
	projectContentPattern := regexp.MustCompile(fmt.Sprintf("%s\\/[^\\/]+\\/[^\\/]+(.*)", root))
	currentProject := ""

	_ = filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if (projectPattern.MatchString(path)) {
			currentProject = path
		}

		if (f.IsDir() && currentProject != "" && projectContentPattern.MatchString(path)) {
			projectMap[currentProject] = append(projectMap[currentProject], path)
		}

		return nil
	})

	return projectMap
}

// Retrieves all license files from the specified directories.
// The script will check all directories (without crawling them further)
// and construct a map where the keys are the filenames, and the values
// are arrays of possible licenses (license type + confidence, e.g. MIT and 0.95)
func processLicensesFromDirs(dirs []string) map[string][]FileData {
	projectLicenses := make(map[string][]FileData)
		
	for _, dir := range dirs {
		dirLicenses := getLicenses(dir)

		// Copy directory licenses to licenses found in project
		for file := range dirLicenses {
			projectLicenses[file] = dirLicenses[file]
		}
	}

	return projectLicenses
}

// Retrieves licenses in a target directory.
// The returned map's keys correspond to filenames, which
// map to arrays of possible license types and their associated confidence
// (e.g. MIT with confidence 0.89, GPL3 with confidence 0.15)
func getLicenses(dir string) map[string][]FileData {
	target, _ := filer.FromDirectory(dir)
	licenses, _ := licensedb.Detect(target)

	fileSet := make(map[string][]FileData)

	for licenseType, match := range licenses {
		for file, confidence := range match.Files {
			fullPath := filepath.Join(dir, file)
			fileData := FileData{ LicenseType: licenseType, Confidence: confidence}
			fileSet[fullPath] = append(fileSet[fullPath], fileData) // Add to set
		}
	}

	return fileSet
}

// Retrieves license files for a target projects root path.
// The projects in the target root path are expected to be in structure: "<projectsPath>/author/repo"
// A projectCachePath is supplied to access the cache of the projects, if it exists (otherwise
// this file is created).
// The returned file maps projects to the corresponding file types. A project might have more
// than one license file (thus the value is a map), and each license file might have
// more than 1 possible value of licensing (thus an array of file metadata is supplied, with
// corresponding license types and confidence levels for each)
func GetLicensesForProjects(projectsPath string, projectCachePath string) map[string]map[string][]FileData {
	resultLicenses := make(map[string]map[string][]FileData)

	var projectSet map[string][]string
	fmt.Printf("Checking cache in %s\n", projectCachePath)

	jsonFile, err := os.Open(projectCachePath)
	if err != nil {
		projectSet = getProjects(projectsPath)

		resultJson, _ := json.MarshalIndent(projectSet, "", "\t")
		_ = ioutil.WriteFile(projectCachePath, resultJson, 0644)
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)	
		json.Unmarshal([]byte(byteValue), &projectSet)
	}
	defer jsonFile.Close()

	fmt.Printf("Found %d projects\n", len(projectSet))

	for project, dirs := range projectSet {
		projectLicenses := processLicensesFromDirs(dirs)
		resultLicenses[project] = projectLicenses
	}

	// Concurrent code (seems like it has a race condition, as it results in non-deterministic runs)
	// TODO: Fix race condition before using the parallel code!
	// var licensesSyncMap sync.Map
	// var wg sync.WaitGroup
	// wg.Add(len(projectSet))

	// for project, dirs := range projectSet {
	// 	go func(project string, dirs []string) {
	// 		projectLicenses := processLicensesFromDirs(project, dirs)
	// 		licensesSyncMap.Store(project, projectLicenses)
	// 		// fmt.Printf("Processed: %s\n", project)
	// 		defer wg.Done()
	// 	}(project, dirs)
	// }

	// fmt.Println("Waiting for goroutines to finish")
	// wg.Wait()

	// // Copy sync map to regular map
	// licensesSyncMap.Range(func(key, value interface{}) bool {
	// 	resultLicenses[key.(string)] = value.(map[string][]FileData)
	// 	return true
	// })

	return resultLicenses
}
