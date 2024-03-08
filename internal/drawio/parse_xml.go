package drawio

import (
	"encoding/xml"
	"fmt"
	"os"
)

// ParseXML parses a draw.io XML file and returns an MxFile struct.
func ParseXML(fileName string) (*MxFile, error) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer xmlFile.Close()

	var mxFile MxFile
	if err := xml.NewDecoder(xmlFile).Decode(&mxFile); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &mxFile, nil
}
