Experimental script that uses [go-license-detector](https://github.com/go-enry/go-license-detector) to gather license
files for all github projects in a directory.

Should be ran on a root file which contains projects in structure: `<root>/author/repo`

**How to use**
1. Intall: `go install`
2. Build: `go build`
3. Run: `./license_detection repos license_out.json cache.json`