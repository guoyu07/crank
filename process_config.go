package main

import (
	"encoding/json"
	"io"
	"os"
)

type ProcessConfig struct {
	Command      string `json:"command"`
	StartTimeout int    `json:"start_timeout"`
	StopTimeout  int    `json:"stop_timeout"`
}

func LoadProcessConfig(path string) (processConfig *ProcessConfig, err error) {
	var reader io.ReadCloser
	if reader, err = os.Open(path); err != nil {
		return
	}
	defer reader.Close()

	processConfig = new(ProcessConfig)
	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(processConfig)
	return
}

func (self *ProcessConfig) Save(path string) error {
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	jsonEncoder := json.NewEncoder(writer)
	return jsonEncoder.Encode(self)
}
