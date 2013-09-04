package main

var colorSets = [...]string{
	"#2980B9", // blue
	"#C0392B", // red
	"#F39C12", // yellow
	"#8E44AD", // WISTERIA
	"#16A085", // green
	"#2C3E50", // black
}

func GetColorValue(i int) string {
	if i >= len(colorSets) {
		return colorSets[0]
	}
	return colorSets[i]
}
