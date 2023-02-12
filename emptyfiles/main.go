package main

import (
	"fmt"
	"os"
)

var emptyFiles []string

func main() {
	checkDir(os.Args[1])
}

func checkDir(root string) {
	_, err := os.Stat(root)
	if err != nil {
		return
	}
	os.Chdir(root)
	dir, err := os.ReadDir(".")
	for _, file := range dir {
		if file.IsDir() {
			checkDir(file.Name())
			os.Chdir("..")
		} else {
			fStat, err := os.Stat(file.Name())
			if err != nil {
				continue
			}
			if fStat.Size() == 0 {
				emptyFiles = append(emptyFiles, file.Name())
				fmt.Println(file.Name())
			}
		}
	}
}
