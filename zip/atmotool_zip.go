package zip

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ZipPredefinedPath(prefix string, dir string) {
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

// creates a list of file paths suitable for zipping
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
