package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/KMankowski/CaveTrack-Caves/internal/models"
)

func TestParseAlabamaCaves(t *testing.T) {
	tests := []struct {
		name          string
		inpAlabamaXML string
		expCaves      []models.Cave
		expErr        error
	}{
		{
			"Successfully parse 2 caves",
			`<?xml version="1.0" encoding="UTF-8"?>
			<gpx version="1.0" creator="GPSBabel - https://www.gpsbabel.org" xmlns="http://www.topografix.com/GPX/1/0">
				<wpt lat="12.345678901" lon="-12.345678901">
					<name>1  First Cave</name>
					<cmt>AJK about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
				<wpt lat="12.345678901" lon="-12.345678901">
					<name>2 Second Cave</name>
					<cmt>AMS about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
			</gpx>`,
			[]models.Cave{
				{
					Name:      "First Cave",
					State:     "Alabama",
					County:    "AJK",
					Latitude:  12.345678901,
					Longitude: -12.345678901,
				},
				{
					Name:      "Second Cave",
					State:     "Alabama",
					County:    "AMS",
					Latitude:  12.345678901,
					Longitude: -12.345678901,
				},
			},
			nil,
		},
		{
			"Successfully parse 1 cave with 2 entrances",
			`<?xml version="1.0" encoding="UTF-8"?>
			<gpx version="1.0" creator="GPSBabel - https://www.gpsbabel.org" xmlns="http://www.topografix.com/GPX/1/0">
				<wpt lat="12.345678901" lon="-12.345678901">
					<name>2  First Cave</name>
					<cmt>AJK about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
				<wpt lat="12.345678901" lon="-12.345678901">
					<name>2 E2 First Cave Ent. 2</name>
					<cmt>AJK about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
			</gpx>`,
			[]models.Cave{
				{
					Name:      "First Cave",
					State:     "Alabama",
					County:    "AJK",
					Latitude:  12.345678901,
					Longitude: -12.345678901,
				},
			},
			nil,
		},
		{
			"Unparseable XML",
			`<wpt>`,
			nil,
			ErrInvalidXML,
		},
		{
			"Unable to parse cave without name",
			`<?xml version="1.0" encoding="UTF-8"?>
			<gpx version="1.0" creator="GPSBabel - https://www.gpsbabel.org" xmlns="http://www.topografix.com/GPX/1/0">
				<wpt lat="12.345678901" lon="-12.345678901">
					<cmt>AJK about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
			</gpx>`,
			nil,
			ErrNameTooShort,
		},
		{
			"Unable to parse cave ID",
			`<?xml version="1.0" encoding="UTF-8"?>
			<gpx version="1.0" creator="GPSBabel - https://www.gpsbabel.org" xmlns="http://www.topografix.com/GPX/1/0">
				<wpt lat="12.345678901" lon="-12.345678901">
					<name>  First Cave</name>
					<cmt>AJK about the cave</cmt>
					<desc>description of the cave</desc>
				</wpt>
			</gpx>`,
			nil,
			ErrInvalidCaveID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outCaves, outErr := parseAlabamaCaves(strings.NewReader(test.inpAlabamaXML))

			if outErr != nil {
				if test.expErr == nil {
					t.Fatalf("unexpected error %v", outErr)
				} else if !errors.Is(outErr, test.expErr) {
					t.Fatalf("expected error %v but got error %v", test.expErr, outErr)
				}
			}

			if outCaves == nil || test.expCaves == nil {
				if outCaves != nil {
					t.Fatalf("expected nil but got %v", outCaves)
				}
				if test.expCaves != nil {
					t.Fatalf("expected %v but got nil", test.expCaves)
				}
			} else if !reflect.DeepEqual(test.expCaves, outCaves) {
				t.Fatalf("Expected %v got %v", test.expCaves, outCaves)
			}
		})
	}
}
