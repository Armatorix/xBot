package main

type Config struct {
	Email    string   `env:"EMAIL"`
	Password string   `env:"PASSWORD"`
	User     string   `env:"USERNAME"`
	Tags     []string `env:"TAGS"`
	Localdev bool     `env:"LOCALDEV" envDefault:"false"`
}
