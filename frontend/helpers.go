package frontend

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

// function to combine the javascript files
func combineJS(scriptsDir string) {

	// combined file name
	outputFile := "combined.js"

	// the output file
	outputFilePath := scriptsDir + outputFile

	// read the directory
	files, err := os.ReadDir(scriptsDir)

	// if there was an error reading the directory, exit the program
	if err != nil {
		log.Panic("Error reading directory:", err)
	}

	// create a string to hold the combined content
	var combinedContent string

	// loop through the files
	for _, file := range files {
		// if file is scripts_combined.js, skip it
		if file.Name() == outputFile {
			continue
		}

		// if file is a .js file, combine it
		if filepath.Ext(file.Name()) == ".js" {

			// read the file
			content, err := os.ReadFile(filepath.Join(scriptsDir, file.Name()))

			// if there was an error reading the file, exit the program
			if err != nil {
				log.Panic("Error reading file:", err)
			}

			// add the content to the combined content
			combinedContent += string(content) + "\n"
		}
	}

	// setup the tool to minify the content
	minify := minify.New()

	minify.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	// minify the content
	combinedContent, err = minify.String("text/javascript", combinedContent)
	if err != nil {
		log.Panic("Error minifying content:", err)
	}

	// write the combined content to the output file
	err = os.WriteFile(outputFilePath, []byte(combinedContent), 0644)

	// if there was an error writing to the output file, exit the program
	if err != nil {
		log.Panic("Error writing to combined file:", err)
	}
}

// normalizes the the array of passed js files to /scripts/
func normalizePath(scripts []string) []string {
	var normalizedScripts []string
	for _, script := range scripts {
		normalizedScripts = append(normalizedScripts, "/scripts/"+filepath.Base(script))
	}
	return normalizedScripts
}

// remove a file from the passed string array
func removeFile(scripts []string, file string) []string {
	var newScripts []string
	for _, script := range scripts {
		if script != file {
			newScripts = append(newScripts, script)
		}
	}
	return newScripts
}

// function to get all the scripts from the scripts directory
func getAllScripts() []string {

	// get all the scripts
	scripts, err := filepath.Glob("frontend/static/scripts/*.js")

	// remove the combined file from the scripts
	scripts = removeFile(scripts, "frontend/static/scripts/combined.js")

	// normalize the paths
	scripts = normalizePath(scripts)

	// if there was an error listing the scripts, panic
	if err != nil {
		log.Panic("Error listing scripts:", err)
	}

	// if there was an error listing the scripts, exit the program
	if err != nil {
		log.Panic("Error listing scripts:", err)
	}

	// return the scripts
	return scripts

}
