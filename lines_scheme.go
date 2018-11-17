package main

func NewLinesScheme() Scheme {
	return SchemeGen("lines", func(m *SnowMan) {
		m.Flash(0, 1, 2)
		m.Flash(3, 4, 5)
		m.Flash(6, 7, 8)
	})
}
