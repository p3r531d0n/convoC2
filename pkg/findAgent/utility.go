package findAgent

import (
	"os/user"
	"strings"
)

func getCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Split the username in case it's in the form HOST\username
	usernameParts := strings.Split(currentUser.Username, `\`)
	username := usernameParts[len(usernameParts)-1]

	return username, nil
}
