package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func Yml() {
	args := os.Args[2:]

	if len(args) < 2 {
		fmt.Println("usage: public-fixer yml <path/to/base> [other-apktool.yml]")
		os.Exit(1)
	}

	base := args[0]
	other := args[1:]

	all := make(map[string]bool)

	basePath := base + "/apktool.yml"
	getDNC(basePath, all)
	for i := 0; i < len(other); i++ {
		path := other[i]
		getDNC(path, all)
	}

	pre, post := getPreAndPostDNC(basePath)
	result := pre
	for _, v := range mapToArray(all) {
		result = result + "- " + v + "\n"
	}
	result += post
	file, err := os.OpenFile(basePath, os.O_RDWR, 0644)
	handleErr(err)
	err = file.Truncate(0)
	handleErr(err)
	_, err = file.WriteAt([]byte(result), 0)
	file.Close()
}

func readFile(path string) string {
	file, err := os.Open(path)
	handleErr(err)
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	handleErr(err)
	content := string(bytes)
	return content
}

func getDNC(path string, buffer map[string]bool) {
	content := readFile(path)
	scanner := bufio.NewScanner(strings.NewReader(content))

	pre := true

	for scanner.Scan() {
		line := scanner.Text()
		if pre {
			// objective: find doNotCompress-line
			if strings.HasSuffix(line, "doNotCompress:") {
				pre = false
			}
		} else {
			// objective: collect all items in array
			if !strings.HasPrefix(line, "- ") {
				break
			} else {
				buffer[line[2:]] = true
			}
		}
	}
}

func getPreAndPostDNC(path string) (string, string) {
	content := readFile(path)
	scanner := bufio.NewScanner(strings.NewReader(content))
	state := "pre"

	preStr := ""
	postStr := ""

	for scanner.Scan() {
		line := scanner.Text()
		if state == "pre" {
			preStr += line + "\n"
			if strings.HasSuffix(line, "doNotCompress:") {
				state = "in"
			}
		} else if state == "in" {
			if strings.HasPrefix(line, "- ") {
				continue
			} else {
				state = "post"
			}
		} else if state == "post" {
			postStr += line + "\n"
		}
	}
	return preStr, postStr
}

func mapToArray(m map[string]bool) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}
