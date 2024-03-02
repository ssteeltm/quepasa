package whatsapp

type WhatsappOptionsBoolean uint

const (
	// Value not setted
	UnknownBooleanType WhatsappOptionsBoolean = iota

	// False
	FalseBooleanType

	// Forced False
	ForcedFalseBooleanType

	// True
	TrueBooleanType

	// Forced True
	ForcedTrueBooleanType
)

// converts to boolean passing default value for unknown option
func (source WhatsappOptionsBoolean) ToBoolean(v bool) bool {
	switch source {
	case FalseBooleanType, ForcedFalseBooleanType:
		return false
	case TrueBooleanType, ForcedTrueBooleanType:
		return true
	default:
		return v
	}
}

func (source WhatsappOptionsBoolean) ToString() string {
	switch source {
	case FalseBooleanType:
		return "false"
	case ForcedFalseBooleanType:
		return "forcedfalse"
	case TrueBooleanType:
		return "true"
	case ForcedTrueBooleanType:
		return "forcedtrue"
	default:
		return "unknown"
	}
}
