package main

import (
	"github.com/joho/godotenv"
	controllers "github.com/nocodeleaks/quepasa/controllers"
	models "github.com/nocodeleaks/quepasa/models"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	whatsmeow "github.com/nocodeleaks/quepasa/whatsmeow"

	log "github.com/sirupsen/logrus"
)

// @title chi-swagger example APIs
// @version 1.0
// @description chi-swagger example APIs
// @BasePath /
func main() {

	// Carregando variaveis de ambiente apartir de arquivo .env
	godotenv.Load()

	loglevel := models.ENV.LogLevel()

	logruslevel, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Errorf("trying parse an invalid loglevel: %s, current: %v", loglevel, log.GetLevel())
	} else {
		log.SetLevel(logruslevel)
	}

	// checks for pending database migrations
	err = models.MigrateToLatest()
	if err != nil {
		log.Fatalf("database migration error: %s", err.Error())
	}

	// should became before whatsmeow start
	title := models.ENV.AppTitle()
	if len(title) > 0 {
		whatsapp.WhatsappWebAppSystem = title
	}

	whatsapp.Options = whatsapp.WhatsappOptionsExtended{
		Groups:       models.ENV.Groups(),
		Broadcasts:   models.ENV.Broadcasts(),
		ReadReceipts: models.ENV.ReadReceipts(),
		Calls:        models.ENV.Calls(),
		ReadUpdate:   models.ENV.ReadUpdate(),
		HistorySync:  models.ENV.HistorySync(),
		LogLevel:     loglevel,
	}

	options := whatsmeow.WhatsmeowOptions{
		WhatsappOptionsExtended: whatsapp.Options,
		WMLogLevel:              models.ENV.WhatsmeowLogLevel(),
		DBLogLevel:              models.ENV.WhatsmeowDBLogLevel(),
	}

	whatsmeow.Start(options)

	// must execute after whatsmeow started
	for _, element := range models.Running {
		if handler, ok := models.MigrationHandlers[element]; ok {
			handler(element)
		}
	}

	// Inicializando serviço de controle do whatsapp
	// De forma assíncrona
	err = models.QPWhatsappStart()
	if err != nil {
		log.Fatalf("whatsapp service starting error: %s", err.Error())
	}

	err = controllers.QPWebServerStart()
	if err != nil {
		log.Info("end with errors")
	} else {
		log.Info("end")
	}
}
