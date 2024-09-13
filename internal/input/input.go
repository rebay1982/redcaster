package input

import (
	rp "github.com/rebay1982/redpix"
)

type InputHandler struct {
	input InputVector
}

type InputVector struct {
	PlayerForward  bool
	PlayerBackward bool
	PlayerLeft     bool
	PlayerRight    bool
}

// NewInputHandler creates a new InputHandler.
func NewInputHandler() *InputHandler {
	i := &InputHandler{
		input: InputVector{},
	}

	return i
}

// HandleInputEvent handles an input event from the redpix library.
func (i *InputHandler) HandleInputEvent(e rp.InputEvent) {
	// Ignore repeated key presses.
	if e.Action == rp.IN_ACT_REPEATED {
		return
	}

	setInput := false
	// Released handling is implicit.
	if e.Action == rp.IN_ACT_PRESSED {
		setInput = true
	}

	switch e.Key {
	case rp.IN_PLAYER_FORWARD:
		i.input.PlayerForward = setInput
	case rp.IN_PLAYER_BACKWARD:
		i.input.PlayerBackward = setInput
	case rp.IN_PLAYER_LEFT:
		i.input.PlayerLeft = setInput
	case rp.IN_PLAYER_RIGHT:
		i.input.PlayerRight = setInput
	}
}

// GetInputVector returns the latest input vector.
func (i InputHandler) GetInputVector() InputVector {
	return i.input
}
