package authentication

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"bytes"

	"github.com/authelia/authelia/internal/configuration/schema"
	"github.com/authelia/authelia/internal/logging"
)

// ExternalUserProvider is a provider reading details from a file.
type ExternalUserProvider struct {
	configuration *schema.ExternalAuthenticationBackendConfiguration
}

// NewExternalUserProvider creates a new instance of ExternalUserProvider.
func NewExternalUserProvider(configuration *schema.ExternalAuthenticationBackendConfiguration) *ExternalUserProvider {

	errs := checkExternalExecutable(configuration.Executable)
	if errs != nil {
		for _, err := range errs {
			logging.Logger().Error(err)
		}

		os.Exit(1)
	}

	return &ExternalUserProvider{
		configuration: configuration,
	}
}

func checkExternalExecutable(path string) []error {

	pathExecutable, _ := exec.LookPath(path)
	
	_, err := os.Stat(pathExecutable)
	if err != nil {
		errs := []error{
			fmt.Errorf("Unable to find external executable: %v", path),
		}

		return errs
	}
	
	return nil
}

// CheckUserPassword checks if provided password matches for the given user.
func (p *ExternalUserProvider) CheckUserPassword(username string, password string) (bool, error) {

	pathExecutable, _ := exec.LookPath(p.configuration.Executable)
	fmt.Printf("pathExecutable: %s\n", pathExecutable);
	
	buffer := bytes.Buffer{}
	buffer.Write([]byte(username + "\n" + password + "\n"))
	
	execExternal := &exec.Cmd {
		Path: pathExecutable,
		Stdin: &buffer,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	
	if err := execExternal.Run(); err != nil {
		return false, errors.New("Wrong username or password!\n")
	}
	
	return true, nil
}

// GetDetails retrieve the groups a user belongs to.
func (p *ExternalUserProvider) GetDetails(username string) (*UserDetails, error) {

	return &UserDetails{
		Username:    username,
		DisplayName: username,
		Groups:      []string{ username },
		Emails:      []string{ username + "@sample.com" },
	}, nil	
}

// UpdatePassword update the password of the given user.
func (p *ExternalUserProvider) UpdatePassword(username string, newPassword string) error {

	return errors.New("Not implemented, yet!")
}
