package client

type Client interface {
	RequestLogin(string, string) bool
}

type ClientImpl struct {
	logged bool
}

func New() Client {
	return &ClientImpl{}
}

func (ci *ClientImpl)RequestLogin(username string, passwd string) bool {
	// request forward to server...
	// server job
	if username == "haha" && passwd == "12345" {
		ci.logged = true
		return true
	}
	// ***server job end***

	return false
}



