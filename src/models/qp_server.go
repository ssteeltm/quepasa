package models

import "time"

/*
<summary>
	Database representation for whatsapp controller service
</summary>
*/
type QpServer struct {
	// Public token
	Token string `db:"token" json:"token" validate:"max=100"`

	// Whatsapp session id
	Wid      string `db:"wid" json:"wid" validate:"max=255"`
	Verified bool   `db:"verified" json:"verified"`
	Devel    bool   `db:"devel" json:"devel"`

	// Optional whatsapp options
	// ------------------------

	// handle groups
	Groups *bool `db:"groups" json:"groups,omitempty"`

	// handle broadcast messages
	Broadcasts *bool `db:"broadcasts" json:"broadcasts,omitempty"`

	// handle read receipt messages
	ReadReceipts *bool `db:"readreceipts" json:"readreceipts,omitempty"`

	// auto reject calls
	RejectCalls *bool `db:"rejectcalls" json:"rejectcalls,omitempty"`

	User      string    `db:"user" json:"user,omitempty" validate:"max=36"`
	Timestamp time.Time `db:"timestamp" json:"timestamp,omitempty"`
}

func (source QpServer) GetWId() string {
	return source.Wid
}

//#region VIEW TRICKS

// used for view
func (source QpServer) IsSetRejectCalls() bool {
	return source.RejectCalls != nil
}

// used for view
func (source QpServer) GetRejectCalls() bool {
	return *source.RejectCalls
}

// used for view
func (source QpServer) IsSetReadReceipts() bool {
	return source.ReadReceipts != nil
}

// used for view
func (source QpServer) GetReadReceipts() bool {
	return *source.ReadReceipts
}

// used for view
func (source QpServer) IsSetBroadcasts() bool {
	return source.Broadcasts != nil
}

// used for view
func (source QpServer) GetBroadcasts() bool {
	return *source.Broadcasts
}

// used for view
func (source QpServer) IsSetGroups() bool {
	return source.Groups != nil
}

// used for view
func (source QpServer) GetGroups() bool {
	return *source.Groups
}

//#endregion
