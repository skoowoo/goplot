package main

import (
	"fmt"
)

type LineNoDotChart struct {
	LineChart
}

func (l *LineNoDotChart) NewChart(name string) string {
	return fmt.Sprintf("new Chart(document.getElementById(\"%s\").getContext(\"2d\")).Line(eval('('+lineJsonStr+')'), {pointDot : false});", name)
}

func init() {
	line := new(LineNoDotChart)
	line.name = "line_no_dot"

	ChartHandlers["line_no_dot"] = line
}
