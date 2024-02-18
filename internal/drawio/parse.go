package drawio

import (
	"encoding/xml"
	"fmt"
	"os"
)

type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram Diagram  `xml:"diagram"`
}

type Diagram struct {
	XMLName      xml.Name     `xml:"diagram"`
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

type MxGraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Root    Root     `xml:"root"`
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	MxCells []MxCell `xml:"mxCell"`
}

type MxCell struct {
	XMLName  xml.Name `xml:"mxCell"`
	Id       string   `xml:"id,attr"`
	Value    string   `xml:"value,attr"`
	Style    string   `xml:"style,attr"`
	Parent   string   `xml:"parent,attr"`
	Vertex   bool     `xml:"vertex,attr"`
	Source   string   `xml:"source,attr"`
	Target   string   `xml:"target,attr"`
	Geometry Geometry `xml:"mxGeometry"`
}

type Geometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
}

func Parse(fileName string) (*MxFile, error) {
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
