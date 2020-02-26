package squeeze

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const prefix = "_"

var (
	redColor = color.New(color.FgRed)
	greenColor = color.New(color.FgGreen)
	greenColorFunc = greenColor.SprintFunc()
)

func Squeeze(context *cli.Context) {
	rootPath := "./"
	if context.Args().Len() > 0 {
		rootPath = context.Args().First()
	}

	files, err := findFiles(rootPath, context)
	if err != nil {
		_, _ = redColor.Println(err)
	} else {
		if len(files) > 0 {
			squeezeFiles(files, context)
			_, _ = greenColor.Println("DONE!")
		}
	}
}

func findFiles(rootPath string, context *cli.Context) ([]string, error) {
	var files []string
	var err error

	info, errInfo := os.Stat(rootPath)

	if os.IsNotExist(errInfo) {
		err = errors.New("path error")
	} else {
		if !info.IsDir() {
			files, err = filepath.Glob(rootPath)
		} else if context.Bool("r") {
			files, err = recursionGlob(rootPath, context.String("ext"))
		} else {
			files, err = filepath.Glob(rootPath + "/*" + context.String("ext"))
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

func squeezeFiles(files []string, context *cli.Context) {
	var wg sync.WaitGroup

	squeezeFile := func(nameFile string, info os.FileInfo,async bool) {
		if async {
			defer wg.Done()
		}
		mapStat, err := GetMapStat(nameFile, context.Int("dlen"), context.String("dpat"))
		if err == nil {
			name := info.Name()
			if !strings.HasPrefix(name, prefix) {
				err = ReturnResult(filepath.Dir(nameFile)+"/" + prefix + name, mapStat)
				if err != nil {
					_, _ = redColor.Println(err)
				} else {
					fmt.Printf("%s: %s\n", greenColorFunc("OK"), nameFile)
					if context.Bool("rm") {
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
		if int(info.Size()/1048576) >= context.Int("as") {
			wg.Add(1)
			go squeezeFile(nameFile, info, true)
		} else {
			squeezeFile(nameFile, info, false)
		}
	}
	wg.Wait()
}