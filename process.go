package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type attr struct {
	Scope       string
	Description string
}

type snippetItem struct {
	Scope       string   `json:"scope"`
	Description string   `json:"description"`
	Body        []string `json:"body"`
	Prefix      string   `json:"prefix"`
}

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// ParseSnpFiles read current work directory's .snp file and parse it to code snippet
func ParseSnpFiles() []byte {
	cwd, _ := os.Getwd()

	snippet := make(map[string]snippetItem)
	err := filepath.WalkDir(cwd, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(filePath, ".txt") {
			stripCwd := filePath[len(cwd)+1:]

			// scope
			scopeAttr := ""
			idx := strings.Index(stripCwd, string(os.PathSeparator))
			if idx != -1 {
				scopeAttr = stripCwd[:idx]
			}
			if scopeAttr == "global" { //file in global folder is use for all language
				scopeAttr = ""
			}

			// prefix
			prefixWithTxt := strings.ReplaceAll(stripCwd, string(os.PathSeparator), " ")
			prefixWithScope := prefixWithTxt[:len(prefixWithTxt)-4]
			prefix := prefixWithScope[len(scopeAttr)+1:]

			// description
			description := prefix

			// body
			rawContent, readFileErr := ioutil.ReadFile(filePath)
			if readFileErr != nil {
				log.Fatal(err)
			}
			stringContent := strings.TrimRight(strings.TrimLeft(string(rawContent), "\n"), "\n")
			body := strings.Split(stringContent, "\n")

			// item
			item := snippetItem{scopeAttr, description, body, prefix}

			snippet[prefix] = item
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	json, jsonEncodingErr := jsonMarshal(snippet)

	if jsonEncodingErr != nil {
		log.Fatal(jsonEncodingErr)
	}

	return json
}
