package main

import "github.com/aeilang/urlshortener/app"

func main() {
	applcation, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	applcation.Run()
}
