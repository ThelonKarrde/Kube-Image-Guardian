package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/validation"
)

func main() {
	startConfig := config.StartUpConfig{}
	startConfig.InitConfig()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var vl validation.Validator
	vl.New(startConfig.ConfigPath, infoLog, errorLog)
	mux := http.NewServeMux()
	mux.HandleFunc("/", vl.ImageValidation)

	if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", startConfig.Port), startConfig.TlsCertPath, startConfig.TlsKeyPath, mux); err != nil {
		errorLog.Fatal(err)
	}
}
