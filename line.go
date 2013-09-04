package main

import (
	"encoding/json"
	"fmt"
)

type lineDataSetsType struct {
	FillColor        string `json:"fillColor"`
	StrokeColor      string `json:"strokeColor"`
	PointColor       string `json:"pointColor"`
	PointStrokeColor string `json:"pointStrokeColor"`
	Data             []int  `json:"data"`
}

type lineDataType struct {
	Labels   []string            `json:"labels"`
	Datasets []*lineDataSetsType `json:"datasets"`
}

type LineChart struct {
	name string
}

func (l *LineChart) Canvas(name string, height int, width int) string {
	if height == 0 {
		height = 300
	}
	if width == 0 {
		width = 400
	}
	return fmt.Sprintf("<canvas id=\"%s\" height=\"%d\" width=\"%d\"></canvas>", name, height, width)
}

func (l *LineChart) JsonCode(c *ChartDataType) (string, error) {
	lines := new(lineDataType)

	lineNum := c.ValueNum()

	lines.Labels = c.ItemName()
	lines.Datasets = make([]*lineDataSetsType, 0, lineNum)

	for i := 0; i < lineNum; i++ {
		line := &lineDataSetsType{}
		line.FillColor = "rgba(220,220,220,0)"
		line.StrokeColor = GetColorValue(i)
		line.PointColor = GetColorValue(i)
		line.PointStrokeColor = "#fff"

		line.Data = c.ItemValue(i)
		lines.Datasets = append(lines.Datasets, line)
	}

	b, err := json.Marshal(lines)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("var lineJsonStr = '%s';", string(b)), nil
}

func (l *LineChart) NewChart(name string) string {
	return fmt.Sprintf("new Chart(document.getElementById(\"%s\").getContext(\"2d\")).Line(eval('('+lineJsonStr+')'));", name)
}

func init() {
	line := new(LineChart)
	line.name = "line"

	ChartHandlers["line"] = line
}
