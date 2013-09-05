package main

import (
	"encoding/json"
	"fmt"
)

type pieDataSetsType struct {
	Color string `json:"color"`
	Value int    `json:"value"`
}

type PieChart struct {
	name string
}

func (b *PieChart) Canvas(name string, height int, width int) string {
	if height == 0 {
		height = 300
	}
	if width == 0 {
		width = 400
	}
	return fmt.Sprintf("<canvas id=\"%s\" height=\"%d\" width=\"%d\"></canvas>", name, height, width)
}

func (l *PieChart) JsonCode(c *ChartDataType) (string, error) {
	pieNum := c.ItemNum()
	items := c.ItemValue(0)

	datasets := make([]*pieDataSetsType, 0, pieNum)

	for i := 0; i < pieNum; i++ {
		pie := &pieDataSetsType{}
		pie.Color = GetColorValue(i)
		pie.Value = items[i]

		datasets = append(datasets, pie)
	}

	b, err := json.Marshal(datasets)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("var pieJsonStr = '%s';", string(b)), nil
}

func (l *PieChart) NewChart(name string) string {
	return fmt.Sprintf("new Chart(document.getElementById(\"%s\").getContext(\"2d\")).Pie(eval('('+pieJsonStr+')'));", name)
}

func init() {
	pie := new(PieChart)
	pie.name = "pie"

	ChartHandlers["pie"] = pie
}
