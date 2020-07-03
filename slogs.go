package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"squeeze/squeeze"
)

func main() {
	app := &cli.App{
		Name: "slogs",
		Usage: "to squeeze log files",
		ArgsUsage: "[FilePath or DirPath (default: ./)]",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "ext",
				Value: ".txt",
				Usage: "extension",
			},
			&cli.BoolFlag{
				Name: "r",
				Value: false,
				Usage: "recursion flag",
			},
			&cli.BoolFlag{
				Name: "rm",
				Value: false,
				Usage: "remove files",
			},
			&cli.IntFlag{
				Name: "dlen",
				Value: 15,
				Usage: "date length",
			},
			&cli.StringFlag{
				Name: "dpat",
				Value: "Jan _2 15:04:05",
				Usage: "date pattern. See https://golang.org/src/time/format.go",
			},
			&cli.IntFlag{
				Name: "as",
				Value: 5,
				Usage: "file size (more or equal) for async, MB",
			},
		},
		Action: func(context *cli.Context) error {
			rootPath := "./"
			if context.Args().Len() > 0 {
				rootPath = context.Args().First()
			}
			flags := squeeze.Flags{
				ExtFlag:         context.String("ext"),
				DatePatternFlag: context.String("dpat"),
				RecursionFlag:   context.Bool("r"),
				RemoveFlag:      context.Bool("rm"),
				DateLengthFlag:  context.Int("dlen"),
				SizeFileAsync:   context.Int("as"),
			}
			squeeze.Execute(rootPath, flags)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}