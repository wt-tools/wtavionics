package config

import "fmt"

type local struct{}

func New() *local {
	return &local{}
}

func (l *local) CurrentPlayer() string {
	return "ZenAviator"
}

func (l *local) GamePoint(path string) string {
	return fmt.Sprintf("http://localhost:8111/%s", path)
}
