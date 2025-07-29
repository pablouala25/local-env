package main

import "github.com/Bancar/reauth-bff-aws-lambda/cmd/bootstrap"

func main() {
	l := bootstrap.SetupLambda()
	defer l.Shutdown()

	l.MustStart()
}
