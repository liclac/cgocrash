package main

import (
	"gopkg.in/olebedev/go-duktape.v2"
	"log"
	"os"
	"os/signal"
)

func Func(ctx *duktape.Context) int {
	num := ctx.RequireNumber(-1)
	ctx.Pop()
	ctx.PushNumber(num)
	return 1
}

func main() {
	for i := 0; i < 2000; i++ {
		go func() {
			ctx := duktape.New()
			ctx.PushGlobalGoFunction("func", Func)

			for {
				if err := ctx.PevalString(`func();`); err != nil {
					log.Fatalf("Couldn't call: %s", err)
				}
				ctx.Pop() // Return value leaks without this
			}
		}()
	}

	// Block until killed by Ctrl+C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
