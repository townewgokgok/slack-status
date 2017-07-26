package main

import (
	"io/ioutil"

	"flag"
	"fmt"

	"os"

	"github.com/kyokomi/emoji"
	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
)

type statusTemplate struct {
	Text  string `yaml:"text,omitempty"`
	Emoji string `yaml:"emoji,omitempty"`
}

type settings struct {
	Token     string                    `yaml:"token"`
	Templates map[string]statusTemplate `yaml:"templates"`
}

var s settings

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: slack-status <template ID>")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Templates:")
	for id, tmpl := range s.Templates {
		emoji.Fprintln(os.Stderr, "- "+id+" : :"+tmpl.Emoji+": "+tmpl.Text)
	}
	os.Exit(1)
}

func main() {
	var err error

	// Load settings
	data, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		panic("Failed to unmarshall settings: " + err.Error())
	}

	// Parse arguments
	flag.Parse()
	id := flag.Arg(0)
	if id == "" {
		usage()
	}

	tmpl, ok := s.Templates[id]
	if !ok {
		fmt.Fprintln(os.Stderr, `Template "`+id+`" is not defined in settings.yml`)
		fmt.Fprintln(os.Stderr, "")
		usage()
	}

	// Request
	api := slack.New(s.Token)
	err = api.SetUserCustomStatus(tmpl.Text, ":" + tmpl.Emoji + ":")
	if err != nil {
		panic("Failed to change status: " + err.Error())
	}

	emoji.Println(":" + tmpl.Emoji + ": " + tmpl.Text)
}
