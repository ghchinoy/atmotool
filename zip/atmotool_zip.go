package zip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	singleFileByteLimit = 107374182400 // 1 GB
	chunkSize           = 4096         // 4 KB
)

var (
	exclusions = []string{".DS_Store", ".zip", ".conf"}
)

// filepath WalkFunc doesn't allow custom params
// This struct will help
type zipper struct {
	srcFolder string
	destFile  string
	writer    *zip.Writer
}

func contains(slice []string, item string) bool {
	ok := false
	for _, s := range slice {
		if strings.Contains(item, s) {
			ok = true
		}
	}
	return ok

	/*
		set := make(map[string]struct{}, len(slice))
		for _, s := range slice {
			set[s] = struct{}{}
		}
		_, ok := set[item]
		return ok
	*/
}

func copyContents(r io.Reader, w io.Writer) error {
	var size int64
	b := make([]byte, chunkSize)
	for {
		// check for large file size
		size += chunkSize
		if size > singleFileByteLimit {
			return errors.New("File too large to zip in this tool.")
		}
		// read into memory
		length, err := r.Read(b[:cap(b)])
		if err != nil {
			if err != io.EOF {
				return err
			}
			if length == 0 {
				break
			}
		}
		// write chunk to zip
		_, err = w.Write(b[:length])
		if err != nil {
			return err
		}
	}
	return nil
}

// internal zip file, called by filepath.Walk on each file
func (z *zipper) zipFile(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	// only zip files, since dirs are created by files inside them
	if !f.Mode().IsRegular() || f.Size() == 0 {
		return nil
	}

	// Exclusions
	if contains(exclusions, f.Name()) {
		//if strings.HasSuffix(f.Name(), ".conf") {
		//log.Println("Skipping", f.Name())
		return nil
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// new file in zip
	fileName := strings.TrimPrefix(path, z.srcFolder+"/")
	w, err := z.writer.Create(fileName)
	if err != nil {
		return err
	}
	// copy contents to zip writer
	err = copyContents(file, w)
	if err != nil {
		return err
	}
	return nil
}

func (z *zipper) zipFolder() error {
	// create zip file
	zipFile, err := os.Create(z.destFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	// zip writer
	z.writer = zip.NewWriter(zipFile)
	err = filepath.Walk(z.srcFolder, z.zipFile)
	if err != nil {
		return nil
	}
	// close zip file
	err = z.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

// zips given folder to file named
func ZipFolder(srcFolder string, destFile string) error {
	z := &zipper{
		srcFolder: srcFolder,
		destFile:  destFile,
	}
	return z.zipFolder()
}
