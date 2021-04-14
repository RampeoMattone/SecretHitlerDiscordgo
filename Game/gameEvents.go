package Game

// Input Calls
type Join *string
type Start struct{}
type ElectChancellor struct {
	caller   *string
	proposal *string
}
type VoteGovernment struct {
	caller *string
	vote   bool
}
type PolicyChoice struct {
	caller   *string
	selected uint8
}
type PolicyVeto struct {
	caller *string
	choice bool
}

// Output Calls

type ErrorType int8

const (
	LobbyFull ErrorType = iota
	NotEnoughPlayers
	WrongPhase
	InvalidInput
)

type Error ErrorType

type AckType int8

const (
	General AckType = iota
)

type Ack AckType
