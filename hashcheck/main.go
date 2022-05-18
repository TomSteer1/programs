package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

var hashes map[string]int
var fileNames map[string][]string
var emptyFiles []string
var clearFuncs = make(map[string]func()) //create a map for storing clear funcs

func init() {
	clearFuncs["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearFuncs["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func listFiles() {
	clear()
	hashes = make(map[string]int)
	fileNames = make(map[string][]string)
	emptyFiles = make([]string, 0)
	fmt.Println("Scanning directory...")
	content, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range content {
		if file.IsDir() {
			continue
		}
		file, err := os.Open(file.Name())
		fStat, err2 := file.Stat()
		if err2 == nil && fStat.Size() == 0 {
			emptyFiles = append(emptyFiles, file.Name())
			continue
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		hash := md5.New()
		io.Copy(hash, file)
		hashes[fmt.Sprintf("%x", hash.Sum(nil))]++
		fileNames[fmt.Sprintf("%x", hash.Sum(nil))] = append(fileNames[fmt.Sprintf("%x", hash.Sum(nil))], file.Name())
	}
}

func cleanUp() {
	dupe := false
	for hash, count := range hashes {
		if count > 1 {
			dupe = true
			options := fileNames[hash]
			clear()
			fmt.Println(("Duplicate file found\nSelect a file to keep: (Enter 0 to skip)"))
			for count, option := range options {
				fmt.Printf("%x. %s \n", count+1, option)
			}
			input := -1
			for input < 0 || input > len(options) {
				fmt.Scanf("%d\n", &input)
				if input < 0 || input > len(options) {
					fmt.Println("Invalid input")
				}
			}
			if input == 0 {
				continue
			}
			for id, file := range options {
				if id+1 != input {
					fmt.Println("Deleting", file)
					os.Remove(file)
				}
			}
			fmt.Println("Press enter to continue")
			fmt.Scanf("\n")
		}
	}
	if !dupe {
		fmt.Println("No duplicates found")
	}
}

func brokenFiles() {
	if len(emptyFiles) == 0 {
		fmt.Println("No empty files found")
		return
	}
	for _, file := range emptyFiles {
		input := -1
		clear()
		fmt.Printf("Empty file found : %s \n", file)
		fmt.Println("Delete file? (1. Yes, 2. No)")
		for input < 1 || input > 2 {
			fmt.Scanf("%d\n", &input)
			if input < 1 || input > 2 {
				fmt.Println("Invalid input")
			}
		}
		if input == 1 {
			fmt.Println("Deleting", file)
			os.Remove(file)
		}
	}
}

func clear() {
	clearFuncs[runtime.GOOS]()
}

func main() {
	clear()
	fmt.Println("Starting...")
	for {
		fmt.Println("1. Find Duplicate Files \n2. Find Empty Files \n3. Exit")
		input := -1
		for input < 1 || input > 3 {
			fmt.Scanf("%d\n", &input)
			if input < 1 || input > 3 {
				fmt.Println("Invalid input")
			}
		}
		switch input {
		case 1:
			listFiles()
			cleanUp()
		case 2:
			listFiles()
			brokenFiles()
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		fmt.Printf("Task Complete \nPress enter to continue")
		fmt.Scanf("\n")
		clear()
	}
}
