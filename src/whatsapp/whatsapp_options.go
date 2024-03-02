package whatsapp

// Whatsapp service options, setted on start, so if want to changed then, you have to restart the entire service
type WhatsappOptions struct {

	// should handle groups messages
	Groups WhatsappOptionsBoolean `json:"groups,omitempty"`

	// should handle broadcast messages
	Broadcasts WhatsappOptionsBoolean `json:"broadcasts,omitempty"`

	// should emit read receipts
	ReadReceipts WhatsappOptionsBoolean `json:"readreceipts,omitempty"`

	// should auto reject calls
	RejectCalls WhatsappOptionsBoolean `json:"rejectcalls,omitempty"`

	// nil for no sync, 0 for all, X for specific days
	HistorySync *uint32 `json:"historysync,omitempty"`

	// default log level
	LogLevel string `json:"loglevel,omitempty"`
}

func (source WhatsappOptions) IsDefault() bool {
	return source.Groups == UnknownBooleanType &&
		source.Broadcasts == UnknownBooleanType &&
		source.ReadReceipts == UnknownBooleanType &&
		source.RejectCalls == UnknownBooleanType &&
		source.HistorySync == nil &&
		len(source.LogLevel) == 0
}

/*
<summary>
	default options from environment variables
	should be set on main
</summary>
*/
var Options WhatsappOptions

func (source WhatsappOptions) HandleRejectCalls(local *bool) bool {
	switch source.RejectCalls {
	case ForcedFalseBooleanType:
		return false
	case ForcedTrueBooleanType:
		return true
	default:
		if local != nil {
			return *local
		}

		return source.RejectCalls.ToBoolean(WhatsappRejectCalls)
	}
}

func (source WhatsappOptions) HandleReadReceipts(local *bool) bool {
	switch source.ReadReceipts {
	case ForcedFalseBooleanType:
		return false
	case ForcedTrueBooleanType:
		return true
	default:
		if local != nil {
			return *local
		}

		return source.ReadReceipts.ToBoolean(WhatsappReadReceipts)
	}
}

func (source WhatsappOptions) HandleGroups(local *bool) bool {
	switch source.Groups {
	case ForcedFalseBooleanType:
		return false
	case ForcedTrueBooleanType:
		return true
	default:
		if local != nil {
			return *local
		}

		return source.Groups.ToBoolean(WhatsappGroups)
	}
}

func (source WhatsappOptions) HandleBroadcasts(local *bool) bool {
	switch source.Broadcasts {
	case ForcedFalseBooleanType:
		return false
	case ForcedTrueBooleanType:
		return true
	default:
		if local != nil {
			return *local
		}

		return source.Broadcasts.ToBoolean(WhatsappBroadcasts)
	}
}
