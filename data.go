package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/v1/yaml"
	"io/ioutil"
)

// loadJSONData unmarshals JSON-encoded data from the file specified by path,
// and returns the result as an interface{}. If the path is "-", then data will
// be acquired from os.Stdin.
func loadJSONData(path string) (ret interface{}, err error) {
	file, err := openStream(path)
	if err != nil {
		return
	}
	defer closeStream(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ret)
	return // whether err==nil or not, our work is done
}

// loadYAMLData unmarshals YAML-encoded data from the file specified by path,
// and returns the result as an interface{}. If the path is "-", then data will
// be acquired from os.Stdin.
func loadYAMLData(path string) (ret interface{}, err error) {
	file, err := openStream(path)
	if err != nil {
		return
	}
	defer closeStream(file)

	rawYAML, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	ret = make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(rawYAML), &ret)
	return
}

// loadTOMLData unmarshals TOML-encoded data from the file specified by path,
// and returns the result as an interface{}. If the path is "-", then data will
// be acquired from os.Stdin.
func loadTOMLData(path string) (ret interface{}, err error) {
	file, err := openStream(path)
	if err != nil {
		return
	}
	defer closeStream(file)

	_, err = toml.DecodeReader(file, &ret)
	return
}
