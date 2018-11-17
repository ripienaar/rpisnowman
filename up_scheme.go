package main

func NewUpScheme() Scheme {
	return SchemeGen("up", func(m *SnowMan) {
		m.Flash(2, 5)
		m.Flash(1, 4)
		m.Flash(0, 3)
		m.Flash(8)
		m.Flash(6, 7)
	})
}
