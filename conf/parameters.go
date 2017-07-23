package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type _Parameters struct {
	HOSTNAME string
	CONTEXT  string
	PORT     string
}

var (
	// Parameters exporta _Parameters para hacer accesible lo que lee del archivo parameters.yml
	Parameters = _Parameters{}
	// File puede definir la ruta del archivo
	file = "parameters.yml"
)

func init() {
	//secret := beego.AppConfig.String("secret")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(data), &Parameters)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//log.Printf("parameters:\n%v\n\n", p)
}
