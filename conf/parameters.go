package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type _Parameters struct {
	HOSTNAME string
	CONTEXT  string
	PORT     string
}

var (
	Parameters = _Parameters{}
)

func init() {
	//secret := beego.AppConfig.String("secret")
	data, err := ioutil.ReadFile("parameters.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(data), &Parameters)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//log.Printf("parameters:\n%v\n\n", p)
}
