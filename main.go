package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TSAP-Laval/acquisition-backend/api"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

func main() {

	// Récupération de la structure des
	// configurations de l'api
	var a api.AcquisitionConfiguration

	// Récupération des configurations
	// dans les variables d'environnement
	// du système d'exploitation
	err := envconfig.Process("TSAP", &a)

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
