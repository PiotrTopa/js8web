package model

type db interface {
	Insert() error
	Update() error
}
