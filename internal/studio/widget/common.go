package widget

type buttonState = int

const (
	buttonStateUp buttonState = 1 + iota
	buttonStateOver
	buttonStateDown
)
