package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func addAnsibleRole(roleName string, roleLocation string) {
	ensureAnsiblefile()
	fo, err := os.OpenFile(philoteAnsiblefile, os.O_RDWR|os.O_APPEND, 0755)
	check(err)
	defer fo.Close()
	_, err = fo.WriteString("role \"" + roleName + "\"\n")
	check(err)
	fo.Sync()
}

func ensureAnsiblefile() {
	if _, err := os.Stat(philoteAnsiblefile); err != nil {
		createAnsiblefile()
	}
}

func createAnsiblefile() {
	ensurePhiloteConfigDir()
	ansiblefileTemplate := `
!/usr/bin/env ruby
#^syntax detection

site "https://galaxy.ansible.com/api/v1"
`
	fo, err := os.Create(philoteAnsiblefile)
	check(err)
	defer fo.Close()
	fo.WriteString(ansiblefileTemplate + "\n")
	fo.Sync()
}

func ensurePhiloteConfigDir() {
	if _, err := os.Stat(philoteConfig); os.IsNotExist(err) {
		os.MkdirAll(philoteConfig, 0755)
	}
}

func removeAnsiblefileRole(roleName string) {
	ensureAnsiblefile()

	fi, err := ioutil.ReadFile(philoteAnsiblefile)
	check(err)

	ansiblefileLines := strings.Split(string(fi), "\n")

	var matchStart int
	var matchEnd int
	for i, ansiblefileLine := range ansiblefileLines {
		didMatch, err := regexp.MatchString(roleName, ansiblefileLine)
		check(err)
		if didMatch {
			println("found match on line:", i+1)
			println(ansiblefileLine)
			matchStart = i
			continue
		}
		if matchStart != 0 {
			potentialNextRole, err := regexp.MatchString("role", ansiblefileLine)
			check(err)
			if potentialNextRole {
				println("found next role:", ansiblefileLine)
				matchEnd = i - 1
				break
			}
		}
	}

	if matchEnd == 0 {
		matchEnd = matchStart
	}

	ansiblefileLines = append(ansiblefileLines[:matchStart], ansiblefileLines[matchEnd+1:]...)
	updatedAnsiblefile := strings.Join(ansiblefileLines, "\n")

	fo, err := os.OpenFile(philoteAnsiblefile, os.O_RDWR|os.O_TRUNC, 0755)
	check(err)
	defer fo.Close()
	_, err = fo.WriteString(updatedAnsiblefile)
	check(err)
	fo.Sync()
}
