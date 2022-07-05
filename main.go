package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	bulletinUrl = "http://www.columbia.edu/cu/bulletin/uwb/sel/COMS_Fall2022_text.html"
	pattern     = `([A-Z]{4} [A-Z]\d{4})\s{2,}<a href="(.+)">(.+)<\/a>\s{2,}(\d+)\s{2,}(\d)\s{2,}(.+)\s+(M|T|MW|TR|F|R|W)\s([0-9-apm:]+)\s{2,}[A-Za-z0-9 ]+\s{2,}([A-Za-z, ]+)`
	outputFile  = "./fall-2022-courses.csv"
)

func main() {
	resp, err := http.Get(bulletinUrl)
	if err != nil {
		log.Fatalf("Unable to get the class information from %s", bulletinUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("Unable to get the class information: %s", err.Error())
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Println(err.Error())
	}

	currMatchNum := 0
	currRow := []string{}
	contents := []byte{}
	contents = append(contents, []byte("Course Code,Link,Sec,Call Number,Credits,Title,Days,Time,Faculty\n")...)

	for matchNum, match := range re.FindAllSubmatch(body, -1) {
		if matchNum != currMatchNum {
			row := strings.Join(currRow, ",") + "\n"
			contents = append(contents, []byte(row)...)
			currRow = []string{}
			currMatchNum += 1
		}
		for groupIdx, group := range match {
			if groupIdx == 0 {
				log.Printf(`matchNum: %d, groupIdx: %d, group: %s`, matchNum, groupIdx, group)
				continue
			} else {
				var strGroup string
				if groupIdx == 2 {
					strGroup = "http://columbia.edu"
				}
				strGroup += strings.Replace(string(group), ",", " ", 1)
				currRow = append(currRow, strings.TrimSpace(strGroup))
			}
		}
	}
	if len(currRow) > 0 {
		row := strings.Join(currRow, ",") + "\n"
		contents = append(contents, []byte(row)...)
	}

	if err := os.Remove(outputFile); err != nil {
		log.Println(err.Error())
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	ioutil.WriteFile(outputFile, contents, 0644)
}
