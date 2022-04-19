package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/config"
	"github.com/ThelonKarrde/Kube-Image-Guardian/internal/validation"
	"github.com/gofiber/fiber/v2"
)

func main() {
	startConfig := config.StartUpConfig{}
	startConfig.InitConfig()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var vl validation.Validator
	vl.New(startConfig.ConfigPath, infoLog, errorLog)

	app := fiber.New()
	app.Post("/", vl.Validate)

	if err := app.ListenTLS(fmt.Sprintf(":%s", startConfig.Port), startConfig.TlsCertPath, startConfig.TlsKeyPath); err != nil {
		errorLog.Fatal(err)
	}
}
