package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Resources struct {
	XMLName xml.Name `xml:"resources"`
	Public  []Public `xml:"public"`
}

type Public struct {
	XMLName xml.Name `xml:"public"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Id      string   `xml:"id,attr"`
}

func main() {

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "fix":
		Fix()
	case "yml":
		Yml()
	default:
		fmt.Println("public-fixer")
		fmt.Println("available commands:")
		fmt.Println("* fix - find and replace all placeholders in xml files")
		fmt.Println("* yml - add DNC of splits to base apktool.yml")
		fmt.Println("")
		os.Exit(1)
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
