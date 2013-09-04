package main

import (
	"net/http"
	"text/template"
)

const html = `{{define "T"}}
<!doctype html>
<html>
	<head>
		<title>Line Chart</title>
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

var ChartHandlers = make(map[string]ChartIf)

func handler(w http.ResponseWriter, r *http.Request) {
	datas, err := LookupCurrentDir(".")
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

		canvas := chart.Canvas("test", prop.Height, prop.Width)
		Args["Canvas"] = canvas

		newChart := chart.NewChart("test")
		Args["NewChart"] = newChart

		if json, err := chart.JsonCode(c); err != nil {
			w.Write([]byte(err.Error()))
			return
		} else {
			Args["JsonCode"] = json
		}
	}

	t, err1 := template.New("foo").Parse(html)
	if err1 != nil {
		w.Write([]byte(err1.Error()))
		return
	}

	err = t.ExecuteTemplate(w, "T", Args)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(addr, nil)
}
