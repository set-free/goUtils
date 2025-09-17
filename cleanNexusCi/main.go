package main

import (
	"fmt"
	"log"
	"os"
)
import "github.com/akamensky/argparse"

const (
	asyncDelJobs int = 1
	daysBefore   int = 5
)

type Cred struct {
	usr string // TODO: не хранить секреты в string
	pwd string // TODO: не хранить секреты в string
}

func argumentParse() (string, int, int) {
	parser := argparse.NewParser("print", "URL for Nexus CI")
	nexusUrl := parser.String("n", "nexus-url", &argparse.Options{Required: true})
	days := parser.Int("d", "days", &argparse.Options{})
	delJobs := parser.Int("j", "jobs-queue", &argparse.Options{})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Println(parser.Usage(err))
		log.Fatal(err.Error())
	}

	if *days == 0 {
		*days = daysBefore
		log.Printf("Dont't set argument '--days', set default = %d", *days)
	}

	if *delJobs == 0 {
		*delJobs = asyncDelJobs
		log.Printf("Dont't set argument '--jobs-queue', set default = %d", *delJobs)
	}

	return *nexusUrl, *days, *delJobs
}

func getEnv(envName string) string {
	env, envExists := os.LookupEnv(envName)
	if envExists != true {
		err := fmt.Errorf("не найдена переменная окружения %s", envName)
		log.Fatal(err.Error())
	}
	return env
}

func getCred() Cred {
	cred := Cred{}
	cred.usr = getEnv("USERNAME")
	cred.pwd = getEnv("PASSWORD")
	return cred
}

func main() {
	nexusUrl, deleteAfterDays, jobQueue := argumentParse()
	cred := getCred()
	xmlData := HttpGet(nexusUrl, cred)
	xmlStruct := ParseXml(xmlData)
	oldArchives := FindOldArchives(xmlStruct, getDateBefore(deleteAfterDays))
	deleteOldArchives(oldArchives, cred, jobQueue)
}
