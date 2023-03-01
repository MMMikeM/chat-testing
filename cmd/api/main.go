package main

import (
	"context"
	"messanger/internal"
)

func main() {
	ctx := context.Background()
	a := internal.NewApp(ctx)
	a.Run()
}
