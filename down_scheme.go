package main

func NewDownScheme() Scheme {
	return SchemeGen("down", func(m *SnowMan) {
		m.Flash(6, 7)
		m.Flash(8)
		m.Flash(0, 3)
		m.Flash(1, 4)
		m.Flash(2, 5)
	})
}
