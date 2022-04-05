package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

type EncodeData struct {
	longURL string
	clicks  int
}

type DecodeData struct {
	Bitlink   string    `json:"bitlink"`
	Timestamp time.Time `json:"timestamp"`
}

func UnmarshalCSV(f io.Reader) (map[string]*EncodeData, error) {
	l, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var em = make(map[string]*EncodeData)
	for i, v := range l {
		// ignore column names
		if i == 0 {
			continue
		}

		e := &EncodeData{
			longURL: v[0],
			clicks:  0,
		}
		// in case of custom domain with same hash (bit.ly/hi1234 != es.pn/hi1234)
		id := fmt.Sprintf("%s/%s", v[1], v[2])

		em[id] = e
	}

	return em, nil
}

func UnmarshalJSON(f io.Reader) ([]*DecodeData, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var dd []*DecodeData
	err = json.Unmarshal(b, &dd)
	if err != nil {
		return nil, err
	}
	return dd, nil
}

func updateClickData(dd []*DecodeData, encodedMap map[string]*EncodeData) (map[string]*EncodeData, error) {
	// main login O(n) where n is number of objects in json file
	for _, v := range dd {
		// short circuit if year is not 2021
		if v.Timestamp.Year() != 2021 {
			continue
		}

		// 0(1): parse id, replace full url with domain/hash: http://bit.ly/2kJdsg8 => bit.ly/2kJdsg8
		// Match encoded id for constant time lookup
		u, err := url.Parse(v.Bitlink)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprintf("%s%s", u.Host, u.Path)

		// O(1): increase click count on bitlink id
		if _, ok := encodedMap[id]; ok {
			encodedMap[id].clicks += 1
		}
	}

	return encodedMap, nil
}

func arrMapSort(em map[string]*EncodeData) []map[string]int {
	// Sorting
	type kv struct {
		k string
		v int
	}
	var ss []kv

	// [{https://linkedin.com/ 529} {https://youtube.com/ 469} {https://google.com/ 502} {https://github.com/ 476} {https://twitter.com/ 525} {https://reddit.com/ 542}]
	for _, v := range em {
		ss = append(ss, kv{v.longURL, v.clicks})
	}

	// [{https://reddit.com/ 542} {https://linkedin.com/ 529} {https://twitter.com/ 525} {https://google.com/ 502} {https://github.com/ 476} {https://youtube.com/ 469}]
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].v > ss[j].v
	})

	// [map[https://reddit.com/:542] map[https://linkedin.com/:529] map[https://twitter.com/:525] map[https://google.com/:502] map[https://github.com/:476] map[https://youtube.com/:469]]
	var res []map[string]int
	for _, value := range ss {
		res = append(res, map[string]int{value.k: value.v})
	}

	return res
}

func main() {
	// open both files first in case there is an error
	csvFile, err := os.Open("encodes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	jsonFile, err := os.Open("decodes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Unmarshal csv => map[string]*EncodeData for O(1) lookup, key = bitlinkID => "bit.ly/2kJO0qs: {"https://twitter.com"/", 0}
	encodedMap, err := UnmarshalCSV(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal JSON => []*DecodeData
	dd, err := UnmarshalJSON(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	// sort => convert to []map[string]int
	udpatedMap, err := updateClickData(dd, encodedMap)
	if err != nil {
		log.Fatal(err)
	}

	res := arrMapSort(udpatedMap)

	fmt.Print(res)
}
