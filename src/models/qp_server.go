package models

import (
	"time"

	"github.com/nocodeleaks/quepasa/whatsapp"
)

/*
<summary>

	Database representation for whatsapp controller service

</summary>
*/
type QpServer struct {
	// Optional whatsapp options
	// ------------------------
	whatsapp.WhatsappOptions

	// Public token
	Token string `db:"token" json:"token" validate:"max=100"`

	// Whatsapp session id
	Wid      string `db:"wid" json:"wid" validate:"max=255"`
	Verified bool   `db:"verified" json:"verified"`
	Devel    bool   `db:"devel" json:"devel"`

	User      string    `db:"user" json:"user,omitempty" validate:"max=36"`
	Timestamp time.Time `db:"timestamp" json:"timestamp,omitempty"`
}

func (source QpServer) GetWId() string {
	return source.Wid
}

//#region VIEW TRICKS

// used for view
func (source QpServer) IsSetCalls() bool {
	return source.Calls != whatsapp.UnSetBooleanType
}

// used for view
func (source QpServer) GetCalls() bool {
	return source.Calls.Boolean()
}

// used for view
func (source QpServer) IsSetReadReceipts() bool {
	return source.ReadReceipts != whatsapp.UnSetBooleanType
}

// used for view
func (source QpServer) GetReadReceipts() bool {
	return source.ReadReceipts.Boolean()
}

// used for view
func (source QpServer) IsSetBroadcasts() bool {
	return source.Broadcasts != whatsapp.UnSetBooleanType
}

// used for view
func (source QpServer) GetBroadcasts() bool {
	return source.Broadcasts.Boolean()
}

// used for view
func (source QpServer) IsSetGroups() bool {
	return source.Groups != whatsapp.UnSetBooleanType
}

// used for view
func (source QpServer) GetGroups() bool {
	return source.Groups.Boolean()
}

//#endregion
