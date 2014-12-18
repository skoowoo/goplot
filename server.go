package main

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"
)

const html = `{{define "T"}}
<!doctype html>
<html>
    <head>
        <script>
            {{.Chartjs}}
        </script>
        <meta name = "viewport" content = "initial-scale = 1, user-scalable = no">
        <style>
            canvas{
            }
        </style>
    </head>
    <body>
        <div style="padding-top:30px;">
            By <a href="http://www.bigendian123.com/skoo.html" target="_blank">skoo</a>
        </div>
        <div style="padding-top:30px;">
        	{{.ItemNames}}
        </div>
        <div style="padding-top:30px;"></div>
        {{.Canvas}}
        <script>
            {{.JsonCode}}
            {{.NewChart}}
        </script>
    </body>
</html>
{{end}}
`

type ChartIf interface {
	Canvas(string, int, int) string
	JsonCode(*ChartDataType) (string, error)
	NewChart(string) string
}

var (
	ChartHandlers = make(map[string]ChartIf)
	ChartFiles    []string
	Index         int
)

func ItemNames(prop ChartPropType) string {
	//<font color="red">xxx</font></br>
	var s string
	for i, n := range prop.Item {
		s += fmt.Sprintf("<b><font color=\"%s\">%s</font></b></br>", GetColorValue(i), n)
	}

	return s
}

func handler(w http.ResponseWriter, r *http.Request) {
	file := ChartFiles[Index]
	Index++
	if Index >= len(ChartFiles) {
		Index = 0
	}

	datas, err := ParseDataFile(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if len(datas) == 0 {
		return
	}
	c := datas[0]

	var chart ChartIf
	var Args = map[string]string{
		"Chartjs": Chartjs,
	}

	if prop, err := c.Prop(); err != nil {
		w.Write([]byte(err.Error()))
		return
	} else {
		chart = ChartHandlers[prop.Name]

		canvas := chart.Canvas("line", prop.Height, prop.Width)
		Args["Canvas"] = canvas

		newChart := chart.NewChart("line")
		Args["NewChart"] = newChart

		if json, err := chart.JsonCode(c); err != nil {
			w.Write([]byte(err.Error()))
			return
		} else {
			Args["JsonCode"] = json
		}

		Args["ItemNames"] = ItemNames(prop)
	}

	if t, err := template.New("foo").Parse(html); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		if err = t.ExecuteTemplate(w, "T", Args); err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	var err error
	ChartFiles, err = LookupCurrentDir(".")
	if err != nil {
		return err
	}

	if len(ChartFiles) == 0 {
		return errors.New("No chart data.")
	}

	return http.ListenAndServe(addr, nil)
}
