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
type Error ErrorType

const (
	LobbyFull ErrorType = iota
	NotEnoughPlayers
	WrongPhase
	InvalidInput
)
type Ack struct {}
type GovernmentElectionFailed struct {
	ForcedPolicy bool
}
type GovernmentElectionSuccess struct {}
type PolicyEnacted struct {
	win DidAnyoneWin
	power SpecialPowers
}
type StartVeto struct {}
type VetoResult struct {
	success bool
	force bool
	enacted PolicyEnacted
}
