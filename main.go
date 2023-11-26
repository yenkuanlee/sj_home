// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var fmap = map[string]int{
	"一":   1,
	"二":   2,
	"三":   3,
	"四":   4,
	"五":   5,
	"六":   6,
	"七":   7,
	"八":   8,
	"九":   9,
	"十":   10,
	"十一":  11,
	"十二":  12,
	"十三":  13,
	"十四":  14,
	"十五":  15,
	"十六":  16,
	"十七":  17,
	"十八":  18,
	"十九":  19,
	"二十":  20,
	"二十一": 21,
	"二十二": 22,
	"二十三": 23,
	"二十四": 24,
	"二十五": 25,
	"二十六": 26,
	"二十七": 27,
	"二十八": 28,
	"二十九": 29,
	"三十":  30,
	"四十一": 41,
	"四十二": 42,
	"四十三": 43,
	"四十四": 44,
	"四十五": 45,
	"四十六": 46,
	"四十七": 47,
	"四十八": 48,
	"四十九": 49,
	"五十":  50,
	"五十一": 51,
	"五十二": 52,
	"五十三": 53,
	"五十四": 54,
	"五十五": 55,
	"五十六": 56,
	"五十七": 57,
	"五十八": 58,
	"五十九": 59,
	"六十":  60,
}

const (
	msg0 = "(總價-車位總價)/(總面積-車位總面積)"
	msg1 = "總價/總面積"
)

type obj struct {
	BuildType   string `json:"AA11"`
	Address     string `json:"a"`
	Date        string `json:"e"`
	Floor       string `json:"f"`
	Community   string `json:"bn"`
	Price       string `json:"tp"`
	CarPrice    string `json:"cp"`
	SinglePrice string `json:"p"`
	Size        string `json:"s"`
	Msg         string `json:"msg"`
	CarNumber   string `json:"l"`
}

type fobj struct {
	BuildType   string    `json:"BuildType"`
	Address     string    `json:"Address"`
	Date        time.Time `json:"Date"`
	Floor       int       `json:"Floor"`
	Community   string    `json:"Community"`
	Price       int       `json:"Price"`
	CarPrice    int       `json:"CarPrice"`
	SinglePrice int       `json:"SinglePrice"`
	Size        float64   `json:"Size"`
	GoodMsg     bool      `json:"GoodMsg"`
	CarNumber   int       `json:"CarNumber"`
}

func (f *fobj) getCarSize() float64 {
	if !f.GoodMsg || f.CarPrice == 0 {
		return 0
	}
	return (f.Size - ((float64(f.Price)-float64(f.CarPrice))/float64(f.SinglePrice))/float64(f.CarNumber))
}

func (o *obj) toFobj() *fobj {
	carNumber, _ := strconv.Atoi(o.CarNumber)
	return &fobj{
		BuildType:   o.BuildType,
		Address:     o.Address,
		Date:        o.getDate(),
		Floor:       o.getFloor(),
		Community:   o.Community,
		Price:       o.getPrice(),
		CarPrice:    o.getCarPrice(),
		SinglePrice: o.getSinglePrice(),
		Size:        o.getSize(),
		GoodMsg:     o.Msg == msg0,
		CarNumber:   carNumber,
	}
}

func (o *obj) getPrice() int {
	p := strings.ReplaceAll(o.Price, ",", "")
	pi, _ := strconv.Atoi(p)
	return pi
}

func (o *obj) getCarPrice() int {
	p := strings.ReplaceAll(o.CarPrice, ",", "")
	pi, _ := strconv.Atoi(p)
	if pi < 1000000 {
		pi = pi * 10000
	}
	return pi
}

func (o *obj) getSinglePrice() int {
	p := strings.ReplaceAll(o.SinglePrice, ",", "")
	pi, _ := strconv.Atoi(p)
	return pi
}

func (o *obj) getFloor() int {
	f, ok := fmap[strings.Split(o.Floor, "層")[0]]
	if !ok {
		return -1
	}
	return f
}

func (o *obj) getSize() float64 {
	s, _ := strconv.ParseFloat(o.Size, 32)
	return s
}

func (o *obj) getDate() time.Time {
	d := strings.Split(o.Date, "/")
	if len(d) != 3 {
		return time.Unix(0, 0)
	}
	year, err := strconv.Atoi(d[0])
	if err != nil {
		return time.Unix(0, 0)
	}
	month, err := strconv.Atoi(d[1])
	if err != nil {
		return time.Unix(0, 0)
	}
	day, err := strconv.Atoi(d[2])
	if err != nil {
		return time.Unix(0, 0)
	}
	return time.Date(year+1911, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func main() {
	dat, err := os.ReadFile("./test.json")
	if err != nil {
		panic(err)
	}

	a := []obj{}
	err = json.Unmarshal(dat, &a)
	if err != nil {
		panic(err)
	}

	res := []fobj{}

	for _, v := range a {
		if strings.Contains(v.Community, "遠雄新宿") {
			res = append(res, *v.toFobj())
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return int(res[i].Date.Unix()) > int(res[j].Date.Unix())
	})

	// log.Printf("KEVIN987: %v, %v\n", getCommunityCarSize("遠雄新宿"), getCommunityLatestCarPrice("遠雄新宿"))

	community := "遠雄新宿"
	carSize := getCommunityCarSize(community)
	latestCarPrice := getCommunityLatestCarPrice(community)
	for i := 0; i < len(res); i++ {
		if res[i].CarPrice == 0 {
			res[i].CarPrice = latestCarPrice
		}
		if !res[i].GoodMsg {
			res[i].SinglePrice = int((float64(res[i].Price) - float64(res[i].CarPrice)) / (res[i].Size - carSize))
		}
	}
	b, _ := json.Marshal(res)
	log.Printf("KEVIN123: %v\n", string(b))
}

func getCommunityCarSize(community string) float64 {
	dat, err := os.ReadFile("./test.json")
	if err != nil {
		panic(err)
	}
	a := []obj{}
	err = json.Unmarshal(dat, &a)
	if err != nil {
		panic(err)
	}
	for _, v := range a {
		if v.Community == community {
			cs := v.toFobj().getCarSize()
			if cs > 1 {
				return cs
			}
		}
	}
	return 0
}

func getCommunityLatestCarPrice(community string) int {
	dat, err := os.ReadFile("./test.json")
	if err != nil {
		panic(err)
	}
	a := []obj{}
	err = json.Unmarshal(dat, &a)
	if err != nil {
		panic(err)
	}
	res := []fobj{}
	for _, v := range a {
		if v.Community == community {
			res = append(res, *v.toFobj())
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return int(res[i].Date.Unix()) > int(res[j].Date.Unix())
	})
	for _, v := range res {
		if v.CarPrice > 0 {
			return v.CarPrice
		}
	}
	return 0
}
