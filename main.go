package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"gopkg.in/yaml.v3"

	"github.com/logikone/gensops/config"
	"github.com/logikone/gensops/gatherer"
)

func main() {
	cfg := config.Config{}

	byt, err := ioutil.ReadFile(".gensops.yaml")

	if err != nil {
		fmt.Println("error opening config file")
		os.Exit(1)
	}

	if err := yaml.Unmarshal(byt, &cfg); err != nil {
		fmt.Printf("error decoding yaml: %s", err)
		os.Exit(1)
	}

	g := gatherer.NewAWSGatherer(&cfg.AWS)

	keys := g.Gather()

	fmt.Println(awsutil.Prettify(keys))
}
