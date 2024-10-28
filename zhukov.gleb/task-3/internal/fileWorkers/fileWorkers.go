package fileWorkers

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"task-3/internal/structs"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "", "Path to cfg file")
	flag.Parse()
}

type FilesWorker struct{}

func NewParser() *FilesWorker {
	return &FilesWorker{}
}

func (p *FilesWorker) CfgParse() (inOutFilePaths structs.Cfg, err error) {
	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return structs.Cfg{}, err
	}

	cfg := structs.Cfg{}
	if err = yaml.Unmarshal(cfgFile, &cfg); err != nil {
		return structs.Cfg{}, err
	}

	return cfg, nil
}

func (p *FilesWorker) InputFileParse(inFilePath string) ([]structs.Currency, error) {
	inputFile, err := os.ReadFile(inFilePath)
	if err != nil {
		return nil, err
	}

	data := []byte(strings.ReplaceAll(string(inputFile), ", ", "."))
	data = []byte(strings.ReplaceAll(string(data), ",", "."))

	var valCurs structs.ValCurs
	if err = xml.Unmarshal(data, &valCurs); err != nil {
		return nil, err
	}

	sort.Slice(valCurs.Currencies, func(i, j int) bool {
		return valCurs.Currencies[i].Value > valCurs.Currencies[j].Value
	})

	return valCurs.Currencies, nil
}

func (p *FilesWorker) OutputFileParse(outFilePath string, data []structs.Currency) error {
	dir := filepath.Dir(outFilePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
