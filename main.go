package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"runtime/debug"
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
	case "ids":
		ids()
	case "strings":
		Strings()
	default:
		fmt.Println("build-companion 5.2.0")
		fmt.Println("available commands:")
		fmt.Println("* fix - find and replace all placeholders in xml files")
		fmt.Println("* yml - add DNC of splits to base apktool.yml")
		fmt.Println("* ids - recalculate the ids of bttv_ prefixed values in public.xml")
		fmt.Println("* strings - find and add strings from mod directory")
		fmt.Println("")
		os.Exit(1)
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		os.Exit(1)
	}
}
