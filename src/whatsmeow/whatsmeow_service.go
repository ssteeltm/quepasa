package whatsmeow

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	library "github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	whatsmeow "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsmeowServiceModel struct {
	Container *sqlstore.Container
	Options   WhatsmeowOptions

	LogEntry *log.Entry `json:"-"` // log entry
}

var WhatsmeowService *WhatsmeowServiceModel

func Start(options WhatsmeowOptions) {
	if WhatsmeowService != nil {
		err := fmt.Errorf("whatsmeow service is already started, if you wanna change options, restart the service")
		panic(err)
	}

	logger := log.New()
	loglevel, err := log.ParseLevel(options.LogLevel)
	if err == nil {
		logger.SetLevel(loglevel)
	} else {
		logger.SetLevel(WhatsmeowLogLevel)
	}

	logentry := logger.WithContext(context.Background())
	logentry.Info("Starting Whatsmeow Service ....")

	dbloglevel := WhatsmeowDBLogLevel
	if len(options.DBLogLevel) > 0 {
		dbloglevel = options.DBLogLevel
	}
	dbLog := waLog.Stdout("whatsmeow/database", dbloglevel, true)

	// check if exists old whatsmeow.db
	var cs string
	if _, err := os.Stat("whatsmeow.db"); err == nil {
		cs = "file:whatsmeow.db?_foreign_keys=on"
	} else {
		// using new quepasa.sqlite
		cs = "file:whatsmeow.sqlite?_foreign_keys=on"
	}

	container, err := sqlstore.New("sqlite3", cs, dbLog)
	if err != nil {
		panic(err)
	}

	WhatsmeowService = &WhatsmeowServiceModel{
		Container: container,
		Options:   options,

		LogEntry: logentry,
	}

	showing := whatsapp.WhatsappWebAppName

	// trim spaces from app name previous setted, if exists
	previousShowing := strings.TrimSpace(whatsapp.WhatsappWebAppSystem)
	if len(previousShowing) > 0 {
		// using previous setted name
		showing = previousShowing
	}

	var version [3]uint32
	version[0] = 0
	version[1] = 9
	version[2] = 0
	store.SetOSInfo(showing, version)

	// this section is broken, history sync configurations, do nothing ......
	// --------------------------------

	historysync := WhatsmeowService.GetHistorySync()
	if historysync != nil {
		HistorySyncValue := *historysync
		logentry.Infof("setting history sync to %v days", HistorySyncValue)

		if HistorySyncValue == 0 {
			store.DeviceProps.RequireFullSync = proto.Bool(true)
			store.DeviceProps.HistorySyncConfig = &waCompanionReg.DeviceProps_HistorySyncConfig{
				FullSyncDaysLimit:   proto.Uint32(3650),
				FullSyncSizeMbLimit: proto.Uint32(102400),
			}
		} else {
			store.DeviceProps.RequireFullSync = proto.Bool(false)
			store.DeviceProps.HistorySyncConfig = &waCompanionReg.DeviceProps_HistorySyncConfig{
				RecentSyncDaysLimit: proto.Uint32(HistorySyncValue * 10),
			}
		}

		store.DeviceProps.HistorySyncConfig.StorageQuotaMb = proto.Uint32(102400)
	}
}

func (source WhatsmeowServiceModel) GetServiceOptions() whatsapp.WhatsappOptionsExtended {
	return source.Options.WhatsappOptionsExtended
}

func (source *WhatsmeowServiceModel) GetHistorySync() *uint32 {
	return source.Options.HistorySync
}

// Used for scan QR Codes
// Dont forget to attach handlers after success login
func (source *WhatsmeowServiceModel) CreateEmptyConnection() (conn *WhatsmeowConnection, err error) {
	options := &whatsapp.WhatsappConnectionOptions{
		Reconnect: false,
	}
	return source.CreateConnection(options)
}

func (source *WhatsmeowServiceModel) CreateConnection(options *whatsapp.WhatsappConnectionOptions) (conn *WhatsmeowConnection, err error) {
	client, err := source.GetWhatsAppClient(options)
	if err != nil {
		return
	}

	logentry := options.GetLogger()
	client.EnableAutoReconnect = options.GetReconnect()

	handlers := &WhatsmeowHandlers{
		WhatsappOptions:  options.WhatsappOptions,
		WhatsmeowOptions: source.Options,
		Client:           client,
		service:          source,

		LogEntry: logentry,
	}

	err = handlers.Register()
	if err != nil {
		return
	}

	conn = &WhatsmeowConnection{
		Client:   client,
		Handlers: handlers,

		LogEntry: logentry,
	}

	client.PrePairCallback = conn.PairedCallBack
	return
}

// Gets an existing store or create a new one for empty wid
func (service *WhatsmeowServiceModel) GetOrCreateStore(wid string) (str *store.Device, err error) {
	if wid == "" {
		str = service.Container.NewDevice()
	} else {
		devices, err := service.Container.GetAllDevices()
		if err != nil {
			err = fmt.Errorf("{Whatsmeow}{EX001} error on getting store from wid {%s}: %v", wid, err)
			return str, err
		} else {
			for _, element := range devices {
				if element.ID.String() == wid {
					str = element
					break
				}
			}

			if str == nil {
				err = &WhatsmeowStoreNotFoundException{Wid: wid}
				return str, err
			}
		}
	}

	return
}

// Temporary
func (service *WhatsmeowServiceModel) GetStoreForMigrated(phone string) (str *store.Device, err error) {

	devices, err := service.Container.GetAllDevices()
	if err != nil {
		err = fmt.Errorf("{Whatsmeow}{EX001} error on getting store from phone {%s}: %v", phone, err)
		return str, err
	} else {
		for _, element := range devices {
			if library.GetPhoneByWId(element.ID.String()) == phone {
				str = element
				break
			}
		}

		if str == nil {
			err = &WhatsmeowStoreNotFoundException{Wid: phone}
			return str, err
		}
	}

	return
}

func (source *WhatsmeowServiceModel) GetWhatsAppClient(options whatsapp.IWhatsappConnectionOptions) (client *whatsmeow.Client, err error) {
	loglevel := WhatsmeowClientLogLevel
	_, logerr := log.ParseLevel(source.Options.WMLogLevel)
	if logerr == nil {
		loglevel = source.Options.WMLogLevel
	}

	wid := options.GetWid()
	clientLog := waLog.Stdout("whatsmeow/client", loglevel, true)
	if len(wid) > 0 {
		clientLog = clientLog.Sub(wid)
	}

	deviceStore, err := source.GetOrCreateStore(wid)
	if deviceStore != nil {
		client = whatsmeow.NewClient(deviceStore, clientLog)
		client.AutoTrustIdentity = true
		client.EnableAutoReconnect = options.GetReconnect()
	}
	return
}

// Flush entire Whatsmeow Database
// Use with wisdom !
func (service *WhatsmeowServiceModel) FlushDatabase() (err error) {
	devices, err := service.Container.GetAllDevices()
	if err != nil {
		return
	}

	for _, element := range devices {
		err = element.Delete()
		if err != nil {
			return
		}
	}

	return
}
