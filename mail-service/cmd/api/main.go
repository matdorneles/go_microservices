package main

type Config struct{}

func main() {
	app := Config{}

	// run web server
	app.routes()
}
