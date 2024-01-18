package graph

type Colors struct {
	CntColors int
	colors    []string
}

func NewColors(cntColors int) *Colors {
	c := &Colors{
		CntColors: cntColors,
		colors:    make([]string, cntColors),
	}
	return c
}

func (c *Colors) SetColorName(colorId int, color string) {
	c.colors[colorId] = color
}

func (c *Colors) GetNodesToColors(codes []int) map[int]string {
	mapColors := make(map[int]string)
	for index := range codes {
		color := c.colors[codes[index]]
		mapColors[index] = color
	}
	return mapColors
}
