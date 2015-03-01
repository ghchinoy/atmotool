package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

// Configuration
type Configuration struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	usage := `SOA Software Community Manager Helper Tool.

Usage:
  atmosphere zip --prefix <prefix> --config <config> [--dir <dir>]
  atmosphere upload less <file> --config <config>
  atmosphere upload file --path <path> --config <config> <files>...
  atmosphere upload all --config <config> [--dir <dir>]
  atmosphere -h | --help
  atmosphere --version

Options:
  -h --help  Show help message and exit.
  --version  Show version and exit.
  --dir=<dir>  Directory. [default: .]
  --path=<cms_path>  CM CMS path.
`
	arguments, _ := docopt.Parse(usage, nil, true, "1.0 cirrus", false)

	// Debug for command-line args
	/*
		var keys []string
		for k := range arguments {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// print the argument keys and values
		for _, k := range keys {
			fmt.Printf("%9s %v\n", k, arguments[k])
		}
	*/

	configLocation, _ := arguments["<config>"].(string)
	configBytes, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Printf("Error opening config file %s\n", err)
		flag.Usage()
		os.Exit(1)
	}
	var config Configuration
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Unable to parse configuration file. %s\n", err)
		os.Exit(1)
	}

	if len(config.Password) < 1 {
		fmt.Printf("Missing or blank password.")
		os.Exit(1)
	}

	if arguments["upload"] == true {
		if arguments["less"] == true {
			lessFilePath := arguments["<file>"].(string)
			uploadLessFile(lessFilePath, config)
		} else if arguments["all"] == true {
			dir, _ := arguments["--dir"].(string)
			uploadAllHelper(dir, config)
		} else if arguments["file"] == true {
			var files []string
			for _, v := range arguments["<files>"].([]string) {
				files = append(files, v)
			}
			path, _ := arguments["--path"].(string)
			upload(files, config, path)
		}
	} else if arguments["zip"] == true {
		prefix, _ := arguments["<prefix>"].(string)
		dir, _ := arguments["--dir"].(string)
		zipPredefinedPath(prefix, dir)
	}
}

func zipPredefinedPath(prefix string, dir string) {
	//fmt.Printf("Zipping '%s' with prefix '%s'\n", dir, prefix)

	// test to see if expected paths exist
	resources_dir := dir + "/resources/theme/default"
	content_dir := dir + "/landing"

	resources_src, err := os.Stat(resources_dir)
	if err != nil {
		fmt.Printf("Error with resource dir ", err)
		os.Exit(1)
	}
	if !resources_src.IsDir() {
		fmt.Printf("%s is not a directory.", resources_dir)
	}

	content_src, err := os.Stat(content_dir)
	if err != nil {
		fmt.Printf("Error with content dir ", err)
		os.Exit(1)
	}
	if !content_src.IsDir() {
		fmt.Printf("%s is not a directory.", content_dir)
		os.Exit(1)
	}

	exclusions := []string{".DS_Store", ".zip"}

	// Get the file lists
	resourcesFileList := listFilesInDir(resources_dir, "", exclusions, false)
	contentFileList := listFilesInDir(content_dir, "", exclusions, false)

	// Create zip files
	writeZipTo(createZipBuffer(resourcesFileList), prefix+"_resourcesThemeDefault.zip")
	zipTheseFiles(contentFileList, prefix+"_contentHomeLanding.zip")
}

// writes a zip buffer to a filename
func writeZipTo(zipbuffer *bytes.Buffer, filename string) {

	fout, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fout.Close(); err != nil {
			panic(err)
		}
	}()

	fw := bufio.NewWriter(fout)
	fw.Write(zipbuffer.Bytes())
	fw.Flush()

	fileinfo, _ := os.Stat(filename)
	fmt.Printf("Created %s (%v)\n", filename, fileinfo.Size())

}

// creates a zip byte buffer and returns
func createZipBuffer(filesList []string) *bytes.Buffer {

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, file := range filesList {
		f, err := w.Create(file)
		if err != nil {
			log.Fatal(err)
		}
		byts, _ := ioutil.ReadFile(file)
		_, err = f.Write(byts)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	return buf

}

// convenience method
func zipTheseFiles(filesList []string, filename string) {

	writeZipTo(createZipBuffer(filesList), filename)

}

func listFilesInDir(path string, relpath string, exclusions []string, debug bool) []string {

	if debug {
		fmt.Printf("%s\n", path)
	}

	var fileList []string

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {

		exclude := false
		for _, e := range exclusions {
			if exclude = strings.Contains(f.Name(), e); exclude {
				break
			}
		}

		if debug && exclude {
			fmt.Printf("  - %s\n", f.Name())
		}

		if !exclude {
			if debug {
				fmt.Printf("   %s\n", f.Name())
			}
			if f.IsDir() {
				dirFileList := listFilesInDir(path+"/"+f.Name(), f.Name()+"/", exclusions, debug)
				fileList = append(fileList, dirFileList...)
			}
			fileList = append(fileList, relpath+f.Name())
		}

	}

	return fileList
}

func uploadLessFile(lessFilePath string, config Configuration) {
	fmt.Printf("Uploading Less file %s to %s\n", lessFilePath, config.Url)
}

func uploadAllHelper(dir string, config Configuration) {
	fmt.Printf("Uploading all in %s to %s\n", dir, config.Url)
}

func upload(files []string, config Configuration, path string) {
	fmt.Printf("Uploading to %s cms location %s these: %s\n", config.Url, path, files)
}
