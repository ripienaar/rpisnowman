package main

// Scheme is a light scheme the snowman can perform
type Scheme interface {
	Run(s *SnowMan)
}

// BaseScheme is a basic scheme that the SchemeGen factory uses to generate new schemes
type BaseScheme struct {
	name   string
	runner func(*SnowMan)
}

// Run implements Scheme
func (s *BaseScheme) Run(m *SnowMan) {
	m.Log().Infof("Running scheme %s", s.name)
	s.runner(m)
}

// SchemeGen is a factory for light schemes
func SchemeGen(name string, r func(*SnowMan)) Scheme {
	return &BaseScheme{
		name:   name,
		runner: r,
	}
}
