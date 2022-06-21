package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const BASE_URL = "https://www.gitignore.io/api/"

func main() {
	app := &cli.App{
		Name:  "tignore",
		Usage: "generate .gitignore file with cli",
		Action: func(c *cli.Context) error {
			var tools string

			for i := 0; i < c.NArg(); i++ {
				tools += c.Args().Get(i) + ","
			}
			tools = strings.TrimRight(tools, ",")

			URL := BASE_URL + tools
			resp, err := http.Get(URL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			byteArray, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			// 既存の.gitignoreがある場合上書きしない
			_, err = os.Stat(".gitignore")
			if errors.Is(err, os.ErrNotExist) {
				file, err := os.Create(".gitignore")
				defer file.Close()
				if err != nil {
					fmt.Println(err)
					return err
				}
				file.Write(byteArray)
				fmt.Println(".gitignore has been generated->", tools)
			} else {
				fmt.Println(".gitignore already exists.")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
