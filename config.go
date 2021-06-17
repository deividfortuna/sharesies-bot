package main

import (
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

func load(p string, v interface{}) error {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, v)
	if err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(v)
}
