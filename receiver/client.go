package receiver

import (
	"github.com/quickfixgo/quickfix"
	"fmt"
)

type Client struct {
	*quickfix.MessageRouter
}


func NewClient() *Client {

	r := &Client{
		MessageRouter: quickfix.NewMessageRouter(),
	}
	return r
}


//OnCreate implemented as part of Application interface
func (e Client) OnCreate(sessionID quickfix.SessionID) {
	return
}

//OnLogon implemented as part of Application interface
func (e Client) OnLogon(sessionID quickfix.SessionID) {
	return
}

//OnLogout implemented as part of Application interface
func (e Client) OnLogout(sessionID quickfix.SessionID) {
	return
}

//FromAdmin implemented as part of Application interface
func (e Client) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

//ToAdmin implemented as part of Application interface
func (e Client) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	return
}

//ToApp implemented as part of Application interface
func (e Client) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", &msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e Client) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

