package main

var colorSets = [...]string{
	// see : http://www.rapidtables.com/web/color/RGB_Color.htm
	"#2980B9", // Blue
	"#C0392B", // Red
	"#F39C12", // Yellow
	"#16A085", // Green
	"#2C3E50", // Black
	"#808080", // Gray
	"#00FF00", // Lime
	"#800080", // Purple
	"#808000", // Olive
	"#000080", // Navy
	"#FF00FF", // Magenta / Fuchsia
	"#8E44AD", // WISTERIA
}

func GetColorValue(i int) string {
	if i >= len(colorSets) {
		return colorSets[0]
	}
	return colorSets[i]
}
