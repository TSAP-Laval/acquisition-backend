package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TSAP-Laval/acquisition-backend/api"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	// On récupère la configuration
	// de l'environnement & on la passe au service
	var a api.AcquisitionConfiguration

	err := envconfig.Process("tsap", &a)

	if err != nil {
		panic(err)
	}

	service := api.New(os.Stdout, &a)
	service.Start()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to stop server...")
	reader.ReadString('\n')

	service.Stop()

	if err != nil {
		panic(err)
	}
}
