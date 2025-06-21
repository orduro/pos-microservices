package main

type config struct {
	addr string
	env  string
}

func NewConfig() *config {
	cfg := config{
		addr: ":8080",
		env:  "dev",
	}

	return &cfg
}
