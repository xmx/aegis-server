package response

type REPLKind string

const (
	REPLErrorKind   REPLKind = "error"
	REPLConsoleKind REPLKind = "console"
)

type REPLMessage struct {
	Kind REPLKind `json:"kind"`
	Data any      `json:"data"`
}
