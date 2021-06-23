package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Strings() {
	args := os.Args
	if len(args) <= 3 {
		fmt.Println("usage: build-companion strings <path/to/our/res> <path/to/their/res>")
		os.Exit(1)
	}
	ourResDirPath := args[2]
	theirResDirPath := args[3]

	matches, err := filepath.Glob(ourResDirPath + "/values*/strings.xml")
	handleErr(err)
	for _, ourStringsXmlPath := range matches {
		relStringsXmlPath, err := filepath.Rel(ourResDirPath, ourStringsXmlPath)
		handleErr(err)
		theirStringsXmlPath := filepath.Join(theirResDirPath, relStringsXmlPath)

		_, err = os.Stat(theirStringsXmlPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println(theirStringsXmlPath, "does not exist")
				createCopy(ourStringsXmlPath, theirStringsXmlPath)
			} else {
				handleErr(err)
			}
		}
		handleFile(ourStringsXmlPath, theirStringsXmlPath)
	}
}

func createCopy(from, to string) {
	source, err := os.Open(from)
	defer source.Close()

	newdir, _ := filepath.Split(to)

	err = os.MkdirAll(newdir, os.ModePerm)
	if err != nil {
		if !os.IsNotExist(err) {
			handleErr(err)
		}
	}

	dest, err := os.Create(to)
	defer dest.Close()
	_, err = io.Copy(dest, source)
	handleErr(err)
}

type stringsxml struct {
	XMLName              xml.Name             `xml:"resources"`
	StringXMLElements    []stringXMLElement   `xml:"string"`
	StringArraysElements []stringArrayElement `xml:"string-array"`
	PluralElements 		 []pluralsElement 	  `xml:"plurals"`
}

type stringXMLElement struct {
	XMLName   xml.Name `xml:"string"`
	Name      string   `xml:"name,attr"`
	Value     string   `xml:",innerxml"`
	Formatted string   `xml:"formatted,attr,omitempty"`
}

type pluralsElement struct {
	XMLName   xml.Name `xml:"string"`
	Name      string   `xml:"name,attr"`
	Value     string   `xml:",innerxml"`
}

type stringArrayElement struct {
	XMLName xml.Name `xml:"string-array"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}

func handleFile(ourStringsXmlPath, theirStringsXmlPath string) {
	oursxml := getxml(ourStringsXmlPath)
	theirxml := getxml(theirStringsXmlPath)

	theirxmlmap := make(map[string]int)

	for i, el := range theirxml.StringXMLElements {
		theirxmlmap[el.Name] = i
	}
	for i, el := range theirxml.PluralElements {
		theirxmlmap[el.Name] = i
	}
	for i, el := range theirxml.StringArraysElements {
		theirxmlmap[el.Name] = i
	}

	for _, el := range oursxml.StringXMLElements {
		i, got := theirxmlmap[el.Name]

		if !got {
			fmt.Println(ourStringsXmlPath, "has", el.Name, theirStringsXmlPath, "does not, creating new element...")
			theirxml.StringXMLElements = append(theirxml.StringXMLElements, stringXMLElement{
				XMLName:   theirxml.StringXMLElements[0].XMLName,
				Name:      el.Name,
				Value:     el.Value,
				Formatted: "",
			})
		} else if theirxml.StringXMLElements[i].Value != el.Value {
			fmt.Println("overwriting", el.Name, "in", theirStringsXmlPath)
			theirxml.StringXMLElements[i].Value = el.Value
		}
	}

	for _, el := range oursxml.StringArraysElements {
		i, got := theirxmlmap[el.Name]

		if !got {
			fmt.Println(ourStringsXmlPath, "has", el.Name, theirStringsXmlPath, "does not, creating new element...")
			theirxml.StringArraysElements = append(theirxml.StringArraysElements, stringArrayElement{
				XMLName:   theirxml.StringArraysElements[0].XMLName,
				Name:      el.Name,
				Value:     el.Value,
			})
		} else if theirxml.StringArraysElements[i].Value != el.Value {
			fmt.Println("overwriting", el.Name, "in", theirStringsXmlPath)
			theirxml.StringArraysElements[i].Value = el.Value
		}
	}

	for _, el := range oursxml.PluralElements {
		i, got := theirxmlmap[el.Name]

		if !got {
			fmt.Println(ourStringsXmlPath, "has", el.Name, theirStringsXmlPath, "does not, creating new element...")
			theirxml.PluralElements = append(theirxml.PluralElements, pluralsElement{
				XMLName:   theirxml.PluralElements[0].XMLName,
				Name:      el.Name,
				Value:     el.Value,
			})
		} else if theirxml.PluralElements[i].Value != el.Value {
			fmt.Println("overwriting", el.Name, "in", theirStringsXmlPath)
			theirxml.PluralElements[i].Value = el.Value
		}
	}

	str := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"
	res, err := xml.MarshalIndent(theirxml, "", "    ")
	handleErr(err)
	str = str + string(res)

	file, err := os.OpenFile(theirStringsXmlPath, os.O_RDWR, 0644)
	handleErr(err)
	err = file.Truncate(0)
	handleErr(err)
	_, err = file.WriteAt([]byte(str), 0)
	handleErr(err)
	file.Close()
}

func getxml(path string) stringsxml {
	file, err := os.Open(path)
	handleErr(err)

	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var resources stringsxml
	err = xml.Unmarshal(byteValue, &resources)
	handleErr(err)

	return resources
}
