package main

import (
	"fmt"
)

type LineNoCurveChart struct {
	LineChart
}

func (l *LineNoCurveChart) NewChart(name string) string {
	return fmt.Sprintf("new Chart(document.getElementById(\"%s\").getContext(\"2d\")).Line(eval('('+lineJsonStr+')'), {pointDot:false, bezierCurve:false});", name)
}

func init() {
	line := new(LineNoCurveChart)
	line.name = "line_no_curve"

	ChartHandlers["line_no_curve"] = line
}
