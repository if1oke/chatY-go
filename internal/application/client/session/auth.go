package session

import (
	"fmt"
	"strings"
)

func (s *ClientSession) Authenticate() error {
	authMsg := fmt.Sprintf("/auth %s %s\n", s.username, s.password)
	_, err := s.writer.WriteString(authMsg)
	if err != nil {
		return err
	}
	err = s.writer.Flush()
	if err != nil {
		return err
	}

	response, err := s.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if !strings.Contains(response, "Successfully") {
		return fmt.Errorf("authentication failed: %s", response)
	}

	return nil
}
