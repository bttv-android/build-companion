package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("usage: public-fixer <base-public.xml> [other-public.xml]")
		os.Exit(1)
	}

	base := args[0]
	other := args[1:]

	baseResources := xmlify(base)

	missingMap := make(map[string]int) // map id to index in baseResources

	for i := 0; i < len(baseResources.Public); i++ {
		public := baseResources.Public[i]
		name := public.Name
		id := public.Id
		if strings.HasPrefix(name, "APKTOOL_DUMMY_") {
			missingMap[id] = i
		}
	}

	fmt.Println("Found", len(missingMap), "missing values")

	for i := 0; i < len(other); i++ {
		filename := other[i]
		otherRes := xmlify(filename)

		for j := 0; j < len(otherRes.Public); j++ {
			public := otherRes.Public[j]
			name := public.Name
			id := public.Id
			if baseIndex, ok := missingMap[id]; ok {
				if strings.HasPrefix(name, "APKTOOL_DUMMY_") {
					fmt.Println("Found id, but it's also a dummy:", id)
				} else {
					fmt.Println(id, "->", name)
					baseResources.Public[baseIndex].Name = name
					delete(missingMap, id)
				}
			}
		}
	}

	if val := len(missingMap); val != 0 {
		fmt.Println("Still", val, "missing names")
	}

	writexml(base, baseResources)
}

func xmlify(path string) Resources {
	file, err := os.Open(path)
	handleErr(err)

	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var resources Resources
	err = xml.Unmarshal(byteValue, &resources)
	handleErr(err)

	return resources
}

func writexml(path string, res Resources) {
	bytes, err := xml.MarshalIndent(res, "", "    ")
	handleErr(err)
	str := string(bytes)
	str = strings.ReplaceAll(str, "></public>", " />")
	str = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" + str

	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	handleErr(err)

	_, err = file.WriteString(str)
	handleErr(err)
	file.Close()
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
