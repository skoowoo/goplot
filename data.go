package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	START = iota
	PROP
	DATA
)

var (
	prop_header_str = []byte("===")
	data_header_str = []byte("---")
)

type ChartPropType struct {
	Name   string
	Width  int
	Height int
}

type ChartItemType struct {
	key    string
	values []int
}

func newChartItem(d string) *ChartItemType {
	item := &ChartItemType{}
	item.values = make([]int, 0, 2)

	fields := strings.Split(d, " ")
	toValue := false

	for _, f := range fields {
		if len(f) == 0 {
			continue
		}

		if toValue {
			if i, err := strconv.Atoi(f); err != nil {
				log.Println(err)
			} else {
				item.values = append(item.values, i)
			}
		} else {
			item.key = f
			toValue = true
		}
	}

	return item
}

type ChartDataType struct {
	prop  *bytes.Buffer
	items []*ChartItemType
}

func newChartData() *ChartDataType {
	c := new(ChartDataType)
	c.prop = bytes.NewBuffer(make([]byte, 0, 128))
	c.items = make([]*ChartItemType, 0, 10)
	return c
}

func (c *ChartDataType) appendProp(p []byte) {
	c.prop.Write(p)
}

func (c *ChartDataType) appendValue(item *ChartItemType) {
	c.items = append(c.items, item)
}

func (c *ChartDataType) Prop() (p ChartPropType, err error) {
	b := c.prop.Bytes()
	b = bytes.Trim(b, " ")
	err = json.Unmarshal(b, &p)
	return
}

func (c *ChartDataType) ItemNum() int {
	return len(c.items)
}

func (c *ChartDataType) ItemName() []string {
	names := make([]string, 0, 5)
	for _, it := range c.items {
		names = append(names, it.key)
	}
	return names
}

func (c *ChartDataType) ValueNum() int {
	if len(c.items) == 0 {
		return 0
	}
	return len(c.items[0].values)
}

func (c *ChartDataType) ItemValue(i int) []int {
	values := make([]int, 0, 5)
	for _, it := range c.items {
		if i >= len(it.values) {
			return nil
		}
		values = append(values, it.values[i])
	}
	return values
}

func ParseDataFile(file string) ([]*ChartDataType, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c *ChartDataType
	charts := make([]*ChartDataType, 0, 2)
	reader := bufio.NewReader(f)
	status := START

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if line = bytes.Trim(line, " "); len(line) == 0 {
			continue
		}

		switch status {
		case START:
			if bytes.Compare(line, prop_header_str) == 0 {
				status = PROP
				c = newChartData()
			} else {
				return nil, errors.New("invalid chart file")
			}
		case PROP:
			if bytes.Compare(line, data_header_str) == 0 {
				status = DATA
			} else {
				c.appendProp(line)
			}
		case DATA:
			if bytes.Compare(line, prop_header_str) == 0 {
				status = PROP
				charts = append(charts, c)
				c = newChartData()
			} else {
				item := newChartItem(string(line))
				c.appendValue(item)
			}
		}
	}

	charts = append(charts, c)
	return charts, nil
}

func LookupCurrentDir(dir string) ([]*ChartDataType, error) {
	var data []*ChartDataType

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if ok, err := filepath.Match("*.chart", f.Name()); err != nil {
			return err
		} else if ok {
			if path == f.Name() {
				data, err = ParseDataFile(path)
				return err
			}
		}
		return nil
	})

	return data, err
}
