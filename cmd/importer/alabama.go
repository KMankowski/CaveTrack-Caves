package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/KMankowski/CaveTrack-Caves/internal/models"
)

var ErrInvalidXML = errors.New("cannot parse XML")
var ErrNameTooShort = errors.New("name has less than 2 words")
var ErrInvalidCaveID = errors.New("id is not parseable int")

type AlabamaCave struct {
	Latitude    float64 `xml:"lat,attr"`
	Longitude   float64 `xml:"lon,attr"`
	Name        string  `xml:"name"`
	Comment     string  `xml:"cmt"`
	Description string  `xml:"desc"`
}

type AlabamaCaves struct {
	Caves []AlabamaCave `xml:"wpt"`
}

// <wpt lat="12.345678901" lon="-12.345678901">
//
//	<name>2  Cave Cave</name>
//	<cmt>text about the cave</cmt>
//	<desc>description of the cave</desc>
//
// </wpt>
// <wpt lat="12.345678901" lon="-12.345678901">
//
//	<name>2 E2 Cave Cave</name>
//	<cmt>text about the cave</cmt>
//	<desc>description of the cave</desc>
//
// </wpt>
func parseAlabamaCaves(rawXML io.Reader) ([]models.Cave, error) {
	decoder := xml.NewDecoder(rawXML)
	var rawCaves AlabamaCaves
	if err := decoder.Decode(&rawCaves); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidXML, err)
	}

	// If two wpt tags have the same id (first char of <name>),
	// then only keep the first as the rest are other entrances
	seenCaveIDs := make(map[int]struct{})
	parsedCaves := make([]models.Cave, 0)
	for _, rawCave := range rawCaves.Caves {
		words := strings.Fields(rawCave.Name)
		if len(words) < 2 {
			return nil, fmt.Errorf("parsing cave name %q: %w", rawCave.Name, ErrNameTooShort)
		}

		id, err := strconv.Atoi(words[0])
		if err != nil {
			return nil, fmt.Errorf("parsing id of cave name %q: %w", rawCave.Name, ErrInvalidCaveID)
		}

		if _, ok := seenCaveIDs[id]; ok {
			continue
		}
		seenCaveIDs[id] = struct{}{}

		parsedCave := models.Cave{
			Name:      strings.Join(words[1:], " "),
			State:     "Alabama",
			County:    strings.Fields(rawCave.Comment)[0],
			Latitude:  rawCave.Latitude,
			Longitude: rawCave.Longitude,
		}
		parsedCaves = append(parsedCaves, parsedCave)
	}

	return parsedCaves, nil
}
