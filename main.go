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
	bitlink string
	clicks  int
}

type DecodeData struct {
	Bitlink   string    `json:"bitlink"`
	Timestamp time.Time `json:"timestamp"`
}

type Data struct {
	EncodeData []*EncodeData
	DecodeData []*DecodeData
}

type ClickMapStruct struct {
	longURL string
	clicks  int
}

// bitlink => longURL mapping
type LinkMap map[string]string

// longURL => click mapping
type ClickMap map[string]int

// return data structure
type SortedSlice []map[string]int

func UnmarshalCSV(f io.Reader) ([]*EncodeData, error) {
	l, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var ed []*EncodeData
	for i, v := range l {
		if i == 0 {
			continue
		}
		e := &EncodeData{
			bitlink: fmt.Sprintf("%s/%s", v[1], v[2]),
			longURL: v[0],
			clicks:  0,
		}
		ed = append(ed, e)
	}
	return ed, nil
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

func (d *Data) updateClickData(linkMap LinkMap) ([]*ClickMapStruct, error) {
	cm := make(ClickMap)

	for _, v := range d.EncodeData {
		cm[v.longURL] = v.clicks
	}

	for _, v := range d.DecodeData {
		if v.Timestamp.Year() != 2021 {
			continue
		}

		u, err := url.Parse(v.Bitlink)
		if err != nil {
			return nil, err
		}

		id := fmt.Sprintf("%s%s", u.Host, u.Path)
		if val, ok := linkMap[id]; ok {
			cm[val] += 1
		}
	}

	cms := []*ClickMapStruct{}
	for k, v := range cm {
		cms = append(cms, &ClickMapStruct{k, v})
	}
	return cms, nil
}

func (d *Data) fillLinkMap() LinkMap {
	m := make(LinkMap)
	for _, v := range d.EncodeData {
		m[v.bitlink] = v.longURL
	}
	return m
}

func sortData(cms []*ClickMapStruct) {
	sort.Slice(cms, func(i, j int) bool {
		return cms[i].clicks > cms[j].clicks
	})
}

func convertToMap(cm []*ClickMapStruct) SortedSlice {
	var res SortedSlice
	for _, v := range cm {
		res = append(res, map[string]int{v.longURL: v.clicks})
	}
	return res
}

func main() {
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

	ed, err := UnmarshalCSV(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	dd, err := UnmarshalJSON(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	data := Data{
		EncodeData: ed,
		DecodeData: dd,
	}

	lm := data.fillLinkMap()

	cm, err := data.updateClickData(lm)
	if err != nil {
		log.Fatal(err)
	}

	sortData(cm)

	resMapping := convertToMap(cm)

	fmt.Println(resMapping)
}
