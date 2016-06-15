package ducklib

type Plugin interface {
	Intialize()
	Shutdown()
}
