package squeeze

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	prefix = "_"
	dividerSize = 1048576
)

var (
	redColor = color.New(color.FgRed)
	greenColor = color.New(color.FgGreen)
	greenColorFunc = greenColor.SprintFunc()
)

type Flags struct {
	ExtFlag, DatePatternFlag string
	RecursionFlag, RemoveFlag bool
	DateLengthFlag, SizeFileAsync int
}

func Squeeze(rootPath string, flags Flags) {
	files, err := findFiles(rootPath, flags)
	if err != nil {
		_, _ = redColor.Println(err)
	} else {
		if len(files) > 0 {
			squeezeFiles(files, flags)
			_, _ = greenColor.Println("DONE!")
		}
	}
}

func findFiles(rootPath string, flags Flags) ([]string, error) {
	var files []string
	var err error

	info, errInfo := os.Stat(rootPath)

	if os.IsNotExist(errInfo) {
		err = errors.New("path error")
	} else {
		if !info.IsDir() {
			files, err = filepath.Glob(rootPath)
		} else if flags.RecursionFlag {
			files, err = recursionGlob(rootPath, flags.ExtFlag)
		} else {
			files, err = filepath.Glob(rootPath + "/*" + flags.ExtFlag)
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

func squeezeFiles(files []string, flags Flags) {
	var wg sync.WaitGroup

	squeezeFile := func(nameFile string, info os.FileInfo, async bool) {
		if async {
			defer wg.Done()
		}
		mapStat, err := GetMapStat(nameFile, flags.DateLengthFlag, flags.DatePatternFlag)
		if err == nil {
			name := info.Name()
			if !strings.HasPrefix(name, prefix) {
				err = ReturnResult(filepath.Dir(nameFile) + "/" + prefix + name, mapStat)
				if err != nil {
					_, _ = redColor.Println(err)
				} else {
					fmt.Printf("%s: %s\n", greenColorFunc("OK"), nameFile)
					if flags.RemoveFlag {
						_ = os.Remove(nameFile)
					}
				}
			}
		} else {
			_, _ = redColor.Println(err)
		}
	}

	for _, nameFile := range files {
		info, _ := os.Stat(nameFile)
		if int(info.Size() / dividerSize) >= flags.SizeFileAsync {
			wg.Add(1)
			go squeezeFile(nameFile, info, true)
		} else {
			squeezeFile(nameFile, info, false)
		}
	}
	wg.Wait()
}