package main

func NewCrossScheme() Scheme {
	return SchemeGen("cross", func(m *SnowMan) {
		m.Flash(0, 5)
		m.Flash(1, 4)
		m.Flash(2, 3)
		m.Flash(6, 7)
		m.Flash(0, 1, 2, 3, 4, 5, 6, 7, 8)
	})
}
