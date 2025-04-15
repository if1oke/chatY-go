package client

import (
	"fmt"
	"strings"
)

func (m *Model) View() string {
	messages := strings.Join(m.Messages, "\n")
	return fmt.Sprintf("%s\n\n[%s]> %s", messages, m.Username, m.Input)
}
