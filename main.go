package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/bearcherian/poswago/convert"
	"github.com/bearcherian/poswago/postman"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no source specified")
		os.Exit(1)
	}
	source := os.Args[1]
	var target string
	if len(os.Args) > 2 {
		target = os.Args[2]
	} else {
		ext := filepath.Ext(source)
		target = fmt.Sprintf("%s.swagger.json", source[0:len(source)-len(ext)])
	}

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	if source == "" {
		fmt.Println("no source specified for conversion")
		os.Exit(1)
	}

	fmt.Printf("source: %s\n", source)
	fmt.Printf("target: %s\n", target)

	// read file
	postmanSource, err := ioutil.ReadFile(source)
	if err != nil {
		log.WithFields(log.Fields{
			"source": source,
			"error":  err,
		}).Error("unable to read source")
		os.Exit(1)
	}

	// unmarshal to struct
	postmanSpec := postman.Spec{}
	if err := json.Unmarshal(postmanSource, &postmanSpec); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("could not unmarshal JSON form postman spec")
		os.Exit(1)
	}

	swaggerSpec, err := convert.PostmanToSwagger(postmanSpec)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to convert postmand spec to swagger")
		os.Exit(1)
	}

	swaggerData, err := json.Marshal(swaggerSpec)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to marshall swagger spec")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(target, swaggerData, 0644); err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"target": target,
		}).Error("could not write to target file")
		os.Exit(1)
	}

}
