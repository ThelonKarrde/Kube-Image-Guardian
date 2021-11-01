package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/validation"
)

func main() {
	startConfig := config.StartUpConfig{}
	startConfig.InitConfig()

	var vl validation.Validator
	vl.ReadConfig(startConfig.ConfigPath)

	http.HandleFunc("/", vl.ImageValidation)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", startConfig.Port), startConfig.TlsCertPath, startConfig.TlsKeyPath, nil); err != nil {
		log.Fatal(err)
	}
}
