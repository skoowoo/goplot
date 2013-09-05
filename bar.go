package main

import (
	"encoding/json"
	"fmt"
)

type barDataSetsType struct {
	FillColor   string `json:"fillColor"`
	StrokeColor string `json:"strokeColor"`
	Data        []int  `json:"data"`
}

type barDataType struct {
	Labels   []string           `json:"labels"`
	Datasets []*barDataSetsType `json:"datasets"`
}

type BarChart struct {
	name string
}

func (b *BarChart) Canvas(name string, height int, width int) string {
	if height == 0 {
		height = 300
	}
	if width == 0 {
		width = 400
	}
	return fmt.Sprintf("<canvas id=\"%s\" height=\"%d\" width=\"%d\"></canvas>", name, height, width)
}

func (l *BarChart) JsonCode(c *ChartDataType) (string, error) {
	bars := new(barDataType)

	barNum := c.ValueNum()

	bars.Labels = c.ItemName()
	bars.Datasets = make([]*barDataSetsType, 0, barNum)

	for i := 0; i < barNum; i++ {
		bar := &barDataSetsType{}
		bar.FillColor = GetColorValue(i)
		bar.StrokeColor = GetColorValue(i)

		bar.Data = c.ItemValue(i)
		bars.Datasets = append(bars.Datasets, bar)
	}

	b, err := json.Marshal(bars)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("var barJsonStr = '%s';", string(b)), nil
}

func (l *BarChart) NewChart(name string) string {
	return fmt.Sprintf("new Chart(document.getElementById(\"%s\").getContext(\"2d\")).Bar(eval('('+barJsonStr+')'));", name)
}

func init() {
	bar := new(BarChart)
	bar.name = "bar"

	ChartHandlers["bar"] = bar
}
