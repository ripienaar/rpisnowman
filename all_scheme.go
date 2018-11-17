package main

func NewAllScheme() Scheme {
	return SchemeGen("all", func(m *SnowMan) {
		m.ToggleAll()
		sleep()
		m.ToggleAll()
	})
}
