package terminal

import "strings"

// Embedded 3x5 glyphs keep rendering independent of network, figlet and locale.
var glyphs = map[rune]string{
	'A': " # |# #|###|# #|# #", 'B': "## |# #|## |# #|## ", 'C': " ##|#  |#  |#  | ##", 'D': "## |# #|# #|# #|## ", 'E': "###|#  |## |#  |###", 'F': "###|#  |## |#  |#  ", 'G': " ##|#  |# #|# #| ##", 'H': "# #|# #|###|# #|# #", 'I': "###| # | # | # |###", 'J': "  #|  #|  #|# #| # ", 'K': "# #|# #|## |# #|# #", 'L': "#  |#  |#  |#  |###", 'M': "# #|###|###|# #|# #", 'N': "# #|###|###|###|# #", 'O': " # |# #|# #|# #| # ", 'P': "## |# #|## |#  |#  ", 'Q': " # |# #|# #| ##|  #", 'R': "## |# #|## |# #|# #", 'S': " ##|#  | # |  #|## ", 'T': "###| # | # | # | # ", 'U': "# #|# #|# #|# #|###", 'V': "# #|# #|# #|# #| # ", 'W': "# #|# #|###|###|# #", 'X': "# #|# #| # |# #|# #", 'Y': "# #|# #| # | # | # ", 'Z': "###|  #| # |#  |###", '0': "###|# #|# #|# #|###", '1': " # |## | # | # |###", '2': "## |  #| # |#  |###", '3': "## |  #| # |  #|## ", '4': "# #|# #|###|  #|  #", '5': "###|#  |## |  #|## ", '6': " ##|#  |###|# #|###", '7': "###|  #| # | # | # ", '8': "###|# #|###|# #|###", '9': "###|# #|###|  #|## ", '-': "   |   |###|   |   ", '_': "   |   |   |   |###", ' ': "   |   |   |   |   ",
}

func ascii(name string) string { return asciiWidth(name, 80) }

func asciiWidth(name string, width int) string {
	clean := cleanText(name)
	letters := []rune(strings.ToUpper(clean))
	for _, r := range letters {
		if _, ok := glyphs[r]; !ok {
			return fitText(clean, width) + "\n"
		}
	}
	var out strings.Builder
	columns := (width - 1) / 4
	if columns < 1 {
		columns = 1
	}
	if columns > 16 {
		columns = 16
	}
	for start := 0; start < len(letters); start += columns {
		end := start + columns
		if end > len(letters) {
			end = len(letters)
		}
		for row := 0; row < 5; row++ {
			for _, r := range letters[start:end] {
				g, ok := glyphs[r]
				if !ok {
					g = "###|  #| # |   | # "
				}
				out.WriteString(strings.Split(g, "|")[row])
				out.WriteByte(' ')
			}
			out.WriteByte('\n')
		}
	}
	out.WriteString(fitText(clean, width) + "\n")
	return out.String()
}
