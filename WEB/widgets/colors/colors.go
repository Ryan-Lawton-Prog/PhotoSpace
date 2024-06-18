package colors

import "image/color"

var Red = color.NRGBA{R: 0x80, A: 0xFF}

type Theme struct {
	Gray0       color.NRGBA
	Gray1       color.NRGBA
	Gray2       color.NRGBA
	Gray3       color.NRGBA
	Gray4       color.NRGBA
	Gray5       color.NRGBA
	Gray6       color.NRGBA
	Gray7       color.NRGBA
	Gray8       color.NRGBA
	Gray9       color.NRGBA
	Gray10      color.NRGBA
	Primary0    color.NRGBA
	Primary1    color.NRGBA
	Primary2    color.NRGBA
	Primary3    color.NRGBA
	Primary4    color.NRGBA
	Primary5    color.NRGBA
	Primary6    color.NRGBA
	Primary7    color.NRGBA
	Primary8    color.NRGBA
	Primary9    color.NRGBA
	Primary10   color.NRGBA
	Secondary0  color.NRGBA
	Secondary1  color.NRGBA
	Secondary2  color.NRGBA
	Secondary3  color.NRGBA
	Secondary4  color.NRGBA
	Secondary5  color.NRGBA
	Secondary6  color.NRGBA
	Secondary7  color.NRGBA
	Secondary8  color.NRGBA
	Secondary9  color.NRGBA
	Secondary10 color.NRGBA
	Error       color.NRGBA
	Warning     color.NRGBA
	Success     color.NRGBA
	Info        color.NRGBA
}

var White = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

var MainTheme = Theme{
	Gray0:     color.NRGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF},
	Gray1:     color.NRGBA{R: 0xB4, G: 0xB4, B: 0xB4, A: 0xFF},
	Gray2:     color.NRGBA{R: 0xA6, G: 0xA6, B: 0xA6, A: 0xFF},
	Gray3:     color.NRGBA{R: 0x96, G: 0x96, B: 0x96, A: 0xFF},
	Gray4:     color.NRGBA{R: 0x87, G: 0x87, B: 0x87, A: 0xFF},
	Gray5:     color.NRGBA{R: 0x78, G: 0x78, B: 0x78, A: 0xFF},
	Gray6:     color.NRGBA{R: 0x69, G: 0x69, B: 0x69, A: 0xFF},
	Gray7:     color.NRGBA{R: 0x59, G: 0x59, B: 0x59, A: 0xFF},
	Gray8:     color.NRGBA{R: 0x4A, G: 0x4A, B: 0x4A, A: 0xFF},
	Gray9:     color.NRGBA{R: 0x3B, G: 0x3B, B: 0x3B, A: 0xFF},
	Gray10:    color.NRGBA{R: 0x2C, G: 0x2C, B: 0x2C, A: 0xFF},
	Primary0:  color.NRGBA{R: 0xFF, G: 0xE5, B: 0xE5, A: 0xFF},
	Primary1:  color.NRGBA{R: 0xFF, G: 0xCE, B: 0xCE, A: 0xFF},
	Primary2:  color.NRGBA{R: 0xFF, G: 0xB7, B: 0xB7, A: 0xFF},
	Primary3:  color.NRGBA{R: 0xFF, G: 0xA0, B: 0xA0, A: 0xFF},
	Primary4:  color.NRGBA{R: 0xFF, G: 0x89, B: 0x89, A: 0xFF},
	Primary5:  color.NRGBA{R: 0xFF, G: 0x73, B: 0x73, A: 0xFF},
	Primary6:  color.NRGBA{R: 0xFF, G: 0x5D, B: 0x5D, A: 0xFF},
	Primary7:  color.NRGBA{R: 0xFF, G: 0x47, B: 0x47, A: 0xFF},
	Primary8:  color.NRGBA{R: 0xFF, G: 0x30, B: 0x30, A: 0xFF},
	Primary9:  color.NRGBA{R: 0xFF, G: 0x19, B: 0x19, A: 0xFF},
	Primary10: color.NRGBA{R: 0xFF, G: 0x02, B: 0x02, A: 0xFF},
	Error:     color.NRGBA{R: 0xF2, G: 0xA4, B: 0x78, A: 0xFF},
	Success:   color.NRGBA{R: 0x99, G: 0xD0, B: 0x78, A: 0xFF},
	Warning:   color.NRGBA{R: 0xF2, G: 0xD0, B: 0x78, A: 0xFF},
	Info:      color.NRGBA{R: 0x83, G: 0xB3, B: 0xE4, A: 0xFF},
}
