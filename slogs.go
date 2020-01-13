package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slogs/squeeze"
	"strings"
	"sync"
)

const prefix = "_"

type ArgsAndFlags struct {
	rootPath, extFlag, datePatternFlag string
	recursionFlag, removeFlag bool
	dateLengthFlag, sizeFileAsync int
}

func main() {
	extFlag := flag.String("ext", ".txt", "Extension (default: \".txt\")")
	recursionFlag := flag.Bool("r", false, "Recursion flag (default: false)")
	removeFlag := flag.Bool("rm", false, "Remove files (default: false")
	dateLengthFlag := flag.Int("dlen", 15, "Date length (default: 15")
	datePatternFlag := flag.String("dpat", "Jan _2 15:04:05", "Date pattern (\"Jan _2 15:04:05\"")
	sizeFileAsync := flag.Int("as", 5, "File size (more or equal) for async (default: 5MB")
	flag.Parse()

	rootPath := "./"
	if len(flag.Args()) > 0 {
		rootPath = flag.Args()[0]
	}

	flags := ArgsAndFlags{
		rootPath: rootPath,
		extFlag: *extFlag,
		recursionFlag: *recursionFlag,
		removeFlag: *removeFlag,
		dateLengthFlag: *dateLengthFlag,
		datePatternFlag: *datePatternFlag,
		sizeFileAsync: *sizeFileAsync,
	}

	files, err := findFiles(flags)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) > 0 {
		squeezeFiles(files, flags)
	}

	fmt.Println("Done!")
}

func findFiles(flags ArgsAndFlags) ([]string, error) {
	var files []string
	var err error

	info, errInfo := os.Stat(flags.rootPath)

	if os.IsNotExist(errInfo) {
		err = errors.New("path error")
	} else {
		if !info.IsDir() {
			files, err = filepath.Glob(flags.rootPath)
		} else if flags.recursionFlag {
			files, err = recursionGlob(flags.rootPath, flags.extFlag)
		} else {
			files, err = filepath.Glob(flags.rootPath + "/*" + flags.extFlag)
		}
	}

	return files, err
}

func recursionGlob(filePath string, extension string) ([]string, error) {
	var files []string

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == extension {
				files = append(files, path)
			}
		}
		return nil
	})

	return files, err
}

func squeezeFiles(files []string, flags ArgsAndFlags) {
	var wg sync.WaitGroup

	squeezeFile := func(nameFile string, info os.FileInfo,async bool) {
		if async {
			defer wg.Done()
		}
		mapStat, err := squeeze.GetMapStat(nameFile, flags.dateLengthFlag, flags.datePatternFlag)
		if err == nil {
			name := info.Name()
			if !strings.HasPrefix(name, prefix) {
				err = squeeze.ReturnResult(filepath.Dir(nameFile)+"/" + prefix + name, mapStat)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("OK: " + nameFile)
					if flags.removeFlag {
						_ = os.Remove(nameFile)
					}
				}
			}
		}
	}

	for _, nameFile := range files {
		info, _ := os.Stat(nameFile)
		if int(info.Size()/1048576) >= flags.sizeFileAsync {
			wg.Add(1)
			go squeezeFile(nameFile, info, true)
		} else {
			squeezeFile(nameFile, info, false)
		}
	}
	wg.Wait()
}