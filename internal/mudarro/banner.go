package mudarro

import "strings"

// Embedded 3x5 glyphs keep rendering independent of network, figlet and locale.
var glyphs = map[rune]string{
	'A': " # |# #|###|# #|# #", 'B': "## |# #|## |# #|## ", 'C': " ##|#  |#  |#  | ##", 'D': "## |# #|# #|# #|## ", 'E': "###|#  |## |#  |###", 'F': "###|#  |## |#  |#  ", 'G': " ##|#  |# #|# #| ##", 'H': "# #|# #|###|# #|# #", 'I': "###| # | # | # |###", 'J': "  #|  #|  #|# #| # ", 'K': "# #|# #|## |# #|# #", 'L': "#  |#  |#  |#  |###", 'M': "# #|###|###|# #|# #", 'N': "# #|###|###|###|# #", 'O': " # |# #|# #|# #| # ", 'P': "## |# #|## |#  |#  ", 'Q': " # |# #|# #| ##|  #", 'R': "## |# #|## |# #|# #", 'S': " ##|#  | # |  #|## ", 'T': "###| # | # | # | # ", 'U': "# #|# #|# #|# #|###", 'V': "# #|# #|# #|# #| # ", 'W': "# #|# #|###|###|# #", 'X': "# #|# #| # |# #|# #", 'Y': "# #|# #| # | # | # ", 'Z': "###|  #| # |#  |###", '0': "###|# #|# #|# #|###", '1': " # |## | # | # |###", '2': "## |  #| # |#  |###", '3': "## |  #| # |  #|## ", '4': "# #|# #|###|  #|  #", '5': "###|#  |## |  #|## ", '6': " ##|#  |###|# #|###", '7': "###|  #| # | # | # ", '8': "###|# #|###|# #|###", '9': "###|# #|###|  #|## ", '-': "   |   |###|   |   ", '_': "   |   |   |   |###", ' ': "   |   |   |   |   ",
}

func ascii(name string) string {
	clean := strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, name)
	letters := []rune(strings.ToUpper(clean))
	var out strings.Builder
	for start := 0; start < len(letters); start += 16 {
		end := start + 16
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
	out.WriteString(clean + "\n")
	return out.String()
}
