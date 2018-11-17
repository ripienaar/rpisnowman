package main

func NewSimpleScheme() Scheme {
	return SchemeGen("simple", func(m *SnowMan) {
		seq := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1, 0}
		for _, pin := range seq {
			m.Toggle(pin)
			sleep()
		}

		m.ToggleAll()
	})
}
