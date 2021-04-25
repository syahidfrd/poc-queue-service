package server

type Service interface {
	Run() (err error)
}
