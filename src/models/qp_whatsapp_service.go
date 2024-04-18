package models

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	library "github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappService struct {
	Servers     map[string]*QpWhatsappServer `json:"-"`
	DB          *QpDatabase                  `json:"-"`
	Initialized bool                         `json:"-"`

	Logger     *log.Entry  `json:"-"`
	initlock   *sync.Mutex `json:"-"`
	appendlock *sync.Mutex `json:"-"`
}

// get default log entry, never nil
func (source *QPWhatsappService) GetLogger() *log.Entry {
	if source.Logger == nil {
		logger := log.New()
		logger.SetLevel(log.ErrorLevel)

		source.Logger = logger.WithContext(context.Background())
	}

	return source.Logger
}

var WhatsappService *QPWhatsappService

func QPWhatsappStart() error {
	if WhatsappService == nil {
		log.Trace("whatsapp service starting ...")

		db := GetDatabase()
		WhatsappService = &QPWhatsappService{
			Servers:    make(map[string]*QpWhatsappServer),
			DB:         db,
			initlock:   &sync.Mutex{},
			appendlock: &sync.Mutex{},
		}

		// seeding database
		err := InitialSeed()
		if err != nil {
			return err
		}
		// iniciando servidores e cada bot individualmente
		return WhatsappService.Initialize()
	} else {
		log.Debug("attempt to start whatsapp service, already started ...")
	}
	return nil
}

// Inclui um novo servidor em um serviço já em andamento
// *Usado quando se passa pela verificação do QRCode
// *Usado quando se inicializa o sistema
func (source *QPWhatsappService) AppendNewServer(info *QpServer) (server *QpWhatsappServer, err error) {
	logger := source.GetLogger()

	// checking if it is cached already
	server, ok := source.Servers[info.Token]
	if !ok {
		// adding to cache
		logger.Infof("adding new server on cache: %s, wid: %s", info.Token, info.Wid)

		// Creating a new instance
		server, err = source.NewQpWhatsappServer(info)
		if err != nil {
			logger.Errorf("error on append new server: %s, :: %s", info.Wid, err.Error())
			return
		}

		source.Servers[info.Token] = server
	} else {
		// updating cached item
		logger.Infof("updating new server on cache: %s, wid: %s", info.Token, info.Wid)

		server.QpServer = info
	}
	return
}

func (source *QPWhatsappService) AppendPaired(paired *QpWhatsappPairing) (server *QpWhatsappServer, err error) {
	logger := source.GetLogger()

	// checking if it is cached already
	server, ok := source.Servers[paired.Token]
	if !ok {
		// adding to cache
		logger.Infof("adding paired server on cache: %s, wid: %s", paired.Token, paired.Wid)

		info := &QpServer{Token: paired.Token, Wid: paired.Wid}

		// Creating a new instance
		server, err = source.NewQpWhatsappServer(info)
		if err != nil {
			logger.Errorf("error on append new server: %s, :: %s", info.Wid, err.Error())
			return
		}

		source.Servers[info.Token] = server
	} else {
		server.Token = paired.Token
		server.Wid = paired.Wid

		// updating cached item
		logger.Infof("updating paired server on cache: %s, old wid: %s, new wid: %s", server.Token, server.Wid, paired.Wid)
	}

	server.connection = paired.conn
	server.Verified = true

	// checking user
	if paired.User != nil {
		server.User = paired.User.Username
	}

	err = server.Save()
	return
}

//region CONSTRUCTORS

// Instance a new quepasa whatsapp server control
func (service *QPWhatsappService) NewQpWhatsappServer(info *QpServer) (server *QpWhatsappServer, err error) {

	if info == nil {
		err = fmt.Errorf("missing server information")
		return
	}

	var serverLogLevel log.Level
	if info.Devel {
		serverLogLevel = log.DebugLevel
	} else {
		serverLogLevel = log.InfoLevel
	}

	logger := log.New()
	logger.SetLevel(serverLogLevel)
	logentry := logger.WithField("token", info.Token)

	if len(info.Wid) > 0 {
		logentry = logentry.WithField("wid", info.Wid)
	}

	server = &QpWhatsappServer{
		QpServer:       info,
		Reconnect:      true,
		syncConnection: &sync.Mutex{},
		syncMessages:   &sync.Mutex{},
		StartTime:      time.Now().UTC(),

		Logger:        logentry,
		StopRequested: false, // setting initial state
		db:            service.DB.Servers,
	}

	logentry.Info("server created")

	server.HandlerEnsure()
	server.WebHookEnsure()
	server.WebhookFill(info.Token, service.DB.Webhooks)
	return
}

func (source *QPWhatsappService) GetOrCreateServerFromToken(token string) (server *QpWhatsappServer, err error) {
	logger := source.GetLogger()
	logger.Debugf("locating server: %s", token)

	server, ok := source.Servers[token]
	if !ok {
		logger.Debugf("server: %s, not in cache, looking up database", token)
		exists, err := source.DB.Servers.Exists(token)
		if err != nil {
			return nil, err
		}

		var info *QpServer
		if exists {
			info, err = source.DB.Servers.FindByToken(token)
			if err != nil {
				return nil, err
			}
			logger.Debugf("server: %s, found", token)
		} else {
			info = &QpServer{
				Token: token,
			}
		}

		server, err = source.AppendNewServer(info)
		return server, err
	}

	return
}

/*
<summary>

	Get or Create a server for scanned qrcode from forms with current user informations and a whatsapp section id
	* use same token if already exists

</summary>
*/
func (service *QPWhatsappService) GetOrCreateServer(user string, wid string) (result *QpWhatsappServer, err error) {
	log.Debugf("locating server with section id: %s", wid)

	phone := library.GetPhoneByWId(wid)
	log.Infof("wid to phone: %s", phone)

	var server *QpWhatsappServer
	servers := service.GetServersForUser(user)
	for _, item := range servers {
		if item.GetNumber() == phone {
			server = item
			server.Wid = wid
			break
		}
	}

	if server == nil {
		token := uuid.New().String()
		log.Infof("creating new server with token: %s", token)
		info := &QpServer{
			Token: token,
			User:  user,
			Wid:   wid,
		}

		server, err = service.AppendNewServer(info)
		if err != nil {
			return
		}
	}

	result = server
	return
}

// delete whatsapp server and remove from cache
func (service *QPWhatsappService) Delete(server *QpWhatsappServer) (err error) {
	err = server.Delete()
	if err != nil {
		return
	}

	delete(service.Servers, server.Token)
	return
}

// method that will initiate all servers from database
func (source *QPWhatsappService) Initialize() (err error) {

	if !source.Initialized {

		servers := source.DB.Servers.FindAll()
		for _, info := range servers {

			// appending server to cache
			server, err := source.AppendNewServer(info)
			if err != nil {
				return err
			}

			logger := source.GetLogger()

			state := server.GetStatus()
			if state == whatsapp.UnPrepared || IsValidToStart(state) {

				// initialize individual server
				logger.Debugf("starting whatsapp server ... on %s state", state)
				go server.Initialize()
			} else {
				logger.Debugf("not auto starting cause state: %s", state)
			}
		}

		source.Initialized = true
	}

	return
}

// Função privada que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) GetServersForUser(username string) (servers map[string]*QpWhatsappServer) {
	servers = make(map[string]*QpWhatsappServer)
	for _, server := range service.Servers {
		if server.GetOwnerID() == username {
			servers[strings.ToLower(server.Token)] = server
		}
	}
	return
}

// Case insensitive
func (service *QPWhatsappService) FindByToken(token string) (*QpWhatsappServer, error) {
	for _, server := range service.Servers {
		if strings.ToLower(server.Token) == strings.ToLower(token) {
			return server, nil
		}
	}

	err := fmt.Errorf("server not found for token: %s", token)
	return nil, err
}

func (source *QPWhatsappService) GetUser(username string, password string) (user *QpUser, err error) {
	logger := source.GetLogger()
	logger.Debugf("finding user: %s", username)
	return source.DB.Users.Check(username, password)
}
