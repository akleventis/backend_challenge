package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalCSV(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description string
		content     string
		expRes      []*EncodeData
		expErr      bool
	}{
		{
			description: "success parsing csv file",
			content:     "skip,first,line\nhttps://adventuretime.com/,bit.ly,prismo\nhttps://taco.com/,bit.ly,cat",
			expRes: []*EncodeData{
				{longURL: "https://adventuretime.com/", bitlink: "bit.ly/prismo", clicks: 0},
				{longURL: "https://taco.com/", bitlink: "bit.ly/cat", clicks: 0},
			},
		},
		{
			description: "error csv format",
			content:     "skip,first,line\nhttps://adventuretime.com/,bit.ly,prismo,oops",
			expRes:      nil,
			expErr:      true,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			var buffer bytes.Buffer
			buffer.WriteString(tc.content)

			c, err := UnmarshalCSV(&buffer)
			if tc.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			assert.Equal(tc.expRes, c)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description string
		content     string
		expRes      []*DecodeData
		expErr      bool
	}{
		{
			description: "success parsing of JSON file",
			content: `[{"bitlink": "http://bit.ly/3hxENM5", "user_agent": "ice king", "timestamp": "2021-05-13T00:00:00Z", "referrer": "pepperment butler", "remote_ip": "123"},
			{"bitlink": "http://amzn.to/3C5IIJm", "user_agent": "princess bubblegum", "timestamp": "2020-06-01T00:00:00Z", "referrer": "gunther", "remote_ip": "123"}]`,
			expRes: []*DecodeData{
				{
					Bitlink:   "http://bit.ly/3hxENM5",
					Timestamp: time.Date(2021, 05, 13, 0, 0, 0, 0, time.UTC),
				},
				{
					Bitlink:   "http://amzn.to/3C5IIJm",
					Timestamp: time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC),
				},
			},
			expErr: false,
		},
		{
			description: "invalid JSON file",
			content:     `[{"bitlink": "http://bit.ly/3hxENM5", "user_agent": "ice king", "timestamp": "2021-05-13T00:00:00Z", "referrer": "pepperment butler", "remote_ip": "123"}`,
			expRes: []*DecodeData(
				nil,
			),
			expErr: true,
		},
	}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			var buffer bytes.Buffer
			buffer.WriteString(tc.content)

			c, err := UnmarshalJSON(&buffer)
			if tc.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			assert.Equal(tc.expRes, c)

		})
	}
}

// func TestUpdateClickData(t *testing.T) {
// 	assert := assert.New(t)

// 	tests := []struct {
// 		description string
// 		decodeData  []*DecodeData
// 		encodeData  map[string]*EncodeData
// 		expRes      map[string]*EncodeData
// 		expErr      bool
// 	}{
// 		{
// 			description: "success updating click data",
// 			decodeData: []*DecodeData{
// 				{
// 					Bitlink:   "http://bit.ly/3hxENM5", // +1
// 					Timestamp: time.Date(2021, 05, 13, 0, 0, 0, 0, time.UTC),
// 				},
// 				{
// 					Bitlink:   "http://bit.ly/jake", // skip
// 					Timestamp: time.Date(2021, 06, 01, 0, 0, 0, 0, time.UTC),
// 				},
// 				{
// 					Bitlink:   "http://bit.ly/3hxENM5", // skip
// 					Timestamp: time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC),
// 				},
// 				{
// 					Bitlink:   "http://bit.ly/3hxENM5", // +1
// 					Timestamp: time.Date(2021, 06, 01, 0, 0, 0, 0, time.UTC),
// 				},
// 			},
// 			encodeData: map[string]*EncodeData{
// 				"bit.ly/3hxENM5": {
// 					longURL: "https://adventuretime.com/",
// 					clicks:  0,
// 				},
// 			},
// 			expRes: map[string]*EncodeData{
// 				"bit.ly/3hxENM5": {
// 					longURL: "https://adventuretime.com/",
// 					clicks:  2,
// 				},
// 			},
// 			expErr: false,
// 		},
// 		{
// 			description: "url parse error",
// 			decodeData: []*DecodeData{
// 				{
// 					Bitlink:   "%foo.html",
// 					Timestamp: time.Date(2021, 05, 13, 0, 0, 0, 0, time.UTC),
// 				},
// 			},
// 			encodeData: map[string]*EncodeData{
// 				"bit.ly/3hxENM5": {
// 					longURL: "https://adventuretime.com/",
// 					clicks:  0,
// 				},
// 			},
// 			expRes: map[string]*EncodeData(nil),
// 			expErr: true,
// 		},
// 	}
// 	for i, tc := range tests {
// 		tc := tc
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			t.Parallel()
// 			t.Log(tc.description)

// 			ed, err := updateClickData(tc.decodeData, tc.encodeData)
// 			if tc.expErr {
// 				assert.Error(err)
// 			} else {
// 				assert.NoError(err)
// 			}

// 			assert.Equal(tc.expRes, ed)

// 		})
// 	}
// }

// func TestArrMapSort(t *testing.T) {
// 	assert := assert.New(t)

// 	tests := []struct {
// 		description string
// 		encodeData  map[string]*EncodeData
// 		expRes      []map[string]int
// 	}{
// 		{
// 			description: "success sort and map conversion",
// 			encodeData: map[string]*EncodeData{
// 				"bit.ly/11111": {longURL: "https://jakethedog.com/", clicks: 300},
// 				"bit.ly/22222": {longURL: "https://marceline.com/", clicks: 100},
// 				"bit.ly/33333": {longURL: "https://adventuretime.com/", clicks: 500},
// 			},
// 			expRes: []map[string]int{
// 				{"https://adventuretime.com/": 500}, {"https://jakethedog.com/": 300}, {"https://marceline.com/": 100},
// 			},
// 		},
// 	}
// 	for i, tc := range tests {
// 		tc := tc
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			t.Parallel()
// 			t.Log(tc.description)

// 			res := arrMapSort(tc.encodeData)

// 			assert.Equal(tc.expRes, res)
// 		})
// 	}
// }
