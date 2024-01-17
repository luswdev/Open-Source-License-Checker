package main

import (
	"strings"
	"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"encoding/csv"
	"log"
	"github.com/google/licensecheck"
)

var (
	tempRow []string
)

func writeReport(records [][]string, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			log.Fatal(err)
		}
	}
}

func licenseChk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println("err?")
		return err
	}

	if strings.Index(d.Name(), "COPYING") != -1 ||
		strings.Index(d.Name(), "LICENSE") != -1 {
		buf, err := os.ReadFile(s)
		if err != nil {
			panic(err)
		}

		cov := licensecheck.Scan(buf)
		if len(cov.Match) > 0 {
			tempRow = append(tempRow, cov.Match[0].ID)
		}
	}

	return nil
}

func walkDir(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}

	csvRes := [][]string{[]string{"source", "license"}}
	for _, src := range entries {
		if src.IsDir() {
			tempRow = []string{src.Name()}

			srcFullPath := filepath.Join(root, src.Name())
			filepath.WalkDir(srcFullPath, licenseChk)

			fmt.Println(tempRow)
			csvRes = append(csvRes, tempRow)
		}
	}

	outPath := strings.ReplaceAll(root, "/", "_") + "license.csv"
	writeReport(csvRes, outPath)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: os_license_chk [buildroot_dir]")
		return
	}

	var rootDir []string = os.Args[1:len(os.Args)]
	for _, path := range rootDir {
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Println("invaild input:", path)
		}

		if fileInfo.IsDir() {
			walkDir(path)
		} else {
			fmt.Println("not a dir:", path)
		}
	}
}

