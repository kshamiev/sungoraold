package config

import (
	"database/sql"

	"github.com/kshamiev/sungora/pkg/app"
	"github.com/kshamiev/sungora/pkg/logger"
	"github.com/kshamiev/sungora/protob"
)

type Component struct {
	Db            *sql.DB
	Lg            logger.Logger
	Cfg           *Config
	Wp            *app.Scheduler
	WsBus         app.WSBus
	SungoraClient protob.SungoraClient
}