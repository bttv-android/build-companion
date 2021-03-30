package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Fix() {
	args := os.Args[2:]

	if len(args) < 2 {
		fmt.Println("usage: public-fixer fix <path/to/base> [other-public.xml]")
		os.Exit(1)
	}

	base := args[0]
	other := args[1:]

	baseResources := xmlify(base + "/res/values/public.xml")

	missingMap := make(map[string]int)    // map id to index in baseResources
	changedMap := make(map[string]string) // map old name to new name

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
					oldName := baseResources.Public[baseIndex].Name
					baseResources.Public[baseIndex].Name = name
					delete(missingMap, id)
					changedMap[oldName] = name
				}
			}
		}
	}

	if val := len(missingMap); val != 0 {
		fmt.Println("Still", val, "missing names")
	}

	// replace in all xml files
	err := filepath.Walk(base+"/res",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(info.Name(), ".xml") {
				return nil
			}
			file, err := os.OpenFile(path, os.O_RDWR, 0644)
			handleErr(err)

			bytes, err := ioutil.ReadAll(file)
			handleErr(err)

			str := string(bytes)

			for old, newv := range changedMap {
				str = strings.ReplaceAll(str, old+"\"", newv+"\"")
				str = strings.ReplaceAll(str, old+"<", newv+"<")
			}

			err = file.Truncate(0)
			handleErr(err)
			_, err = file.WriteAt([]byte(str), 0)
			handleErr(err)

			err = file.Close()
			handleErr(err)
			return nil
		})
	handleErr(err)

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
