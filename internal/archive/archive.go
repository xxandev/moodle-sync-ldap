package archive

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const MAIN_DIR string = "archive"

var (
	//go:embed archive
	archive embed.FS
	once    sync.Once
)

// Install - unpacking of resources
func Unpack(path string) (err error) {
	if _, er := os.Stat(path); er == nil {
		return fmt.Errorf("archive don`t unpacked, main dir already exists: %v", path)
	}
	once.Do(func() { err = unpack(path, MAIN_DIR) })
	return
}

func unpack(outpath, inpath string) error {
	list, err := archive.ReadDir(inpath)
	if err != nil {
		return fmt.Errorf("unpack resources, error read %s: %v", inpath, err)
	}
	for n := range list {
		inFP := inpath + "/" + list[n].Name()
		if list[n].IsDir() {
			unpack(outpath, inFP)
			continue
		}
		outFP, err := filepath.Abs(outpath + "/" + strings.Replace(inFP, MAIN_DIR, "", 1))
		if err != nil {
			// log.Printf("unpack resources, warning filepath.abs(%v/%v): %v\n", outpath, inFP, err)
			continue
		}
		createDir(filepath.Dir(outFP))
		// if err := createDir(filepath.Dir(outFP)); err != nil {
		// 	log.Printf("unpack resources, warning create dir: %v\n", err)
		// }
		content, err := archive.ReadFile(inFP)
		if err != nil {
			// log.Printf("unpack resources, error read file %v: %v\n", inFP, err)
			continue
		}
		createFile(outFP, content)
		// if err := createFile(outFP, content); err != nil {
		// 	log.Printf("unpack resources, warning create file %v: %v\n", outFP, err)
		// }
	}
	return nil
}

func createDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("dir already exists %v", path)
	}
	return os.MkdirAll(path, os.ModePerm)
}

func createFile(path string, content []byte) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists %v", path)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return file.Close()
}
