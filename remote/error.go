package remote

import "fmt"

type ConnectionError struct {
	Name       string // name of the remote server
	Address    string // Address of the remote server
	Username   string // Username used to log to the remote server
	AuthMethod string // Auth method used to log to the remote server
	Err        error  // Err returned
}

func (c ConnectionError) Error() string {
	return fmt.Sprintf("Couldn't authenticate to %s (%s) with user %s and authMethod %s :%v.",
		c.Name, c.Address, c.Username, c.AuthMethod, c.Err)
}
