package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ids() {
	args := os.Args
	if len(args) <= 2 {
		println("usage: public-fixer ids <path/to/public.xml>")
		os.Exit(1)
	}
	publicxml := args[2]
	xml := xmlify(publicxml)

	maxMap := make(map[string]uint64) // map the type to the max id in use

	// 1. collect max ids for non-bttv values
	for i := 0; i < len(xml.Public); i++ {
		public := xml.Public[i]
		ttype := public.Type
		name := public.Name
		id := public.Id
		if strings.HasPrefix(name, "bttv_") {
			continue
		}
		if _, succ := maxMap[ttype]; !succ {
			maxMap[ttype] = 0
		}
		idInt := toInt(id)
		if idInt > maxMap[ttype] {
			maxMap[ttype] = idInt
		}
	}

	file, err := os.OpenFile(publicxml, os.O_RDWR, 0644)
	handleErr(err)
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	handleErr(err)
	content := string(bytes)

	// 2. replace bttv_ values
	for i := 0; i < len(xml.Public); i++ {
		public := xml.Public[i]
		ttype := public.Type
		name := public.Name
		id := public.Id
		if !strings.HasPrefix(name, "bttv_") {
			continue
		}
		maxMap[ttype] += 1
		newId := maxMap[ttype]
		find := fmt.Sprintf("<public type=\"%s\" name=\"%s\" id=\"%s\" />", ttype, name, id)
		replace := fmt.Sprintf("<public type=\"%s\" name=\"%s\" id=\"0x%x\" />", ttype, name, newId)
		content = strings.Replace(content, find, replace, 1)
	}

	err = file.Truncate(0)
	handleErr(err)

	_, err = file.WriteAt([]byte(content), 0)
	handleErr(err)
}

func toInt(id string) uint64 {
	id = strings.Replace(id, "0x", "", 1)
	i, err := strconv.ParseUint(id, 16, 64)
	handleErr(err)
	return i
}
