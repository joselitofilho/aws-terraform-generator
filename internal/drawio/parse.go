package drawio

import (
	"encoding/xml"
	"fmt"
	"os"
)

// MxFile represents the root element of the draw.io XML file
type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram Diagram  `xml:"diagram"`
}

// Diagram represents the diagram element within the draw.io XML file
type Diagram struct {
	XMLName      xml.Name     `xml:"diagram"`
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

// MxGraphModel represents the graph model element within the draw.io XML file
type MxGraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Root    Root     `xml:"root"`
}

// Root represents the root element within the graph model of the draw.io XML file
type Root struct {
	XMLName xml.Name `xml:"root"`
	MxCells []MxCell `xml:"mxCell"`
}

// MxCell represents a cell element within the draw.io XML file
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

// Geometry represents the geometry element within a cell of the draw.io XML file
type Geometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
}

// Parse parses a draw.io XML file and returns an MxFile struct
func Parse(fileName string) (*MxFile, error) {
	// Open the XML file
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer xmlFile.Close()

	// Decode the XML file into an MxFile struct
	var mxFile MxFile
	if err := xml.NewDecoder(xmlFile).Decode(&mxFile); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &mxFile, nil
}
