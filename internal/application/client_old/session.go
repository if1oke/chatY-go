package client_old

import (
	"bufio"
	"chatY-go/internal/domain/auth"
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type IClientSession interface {
	Start()
}

type ClientSession struct {
	conn         net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	logger       logger.ILogger
	clientConfig config.IClientConfig
	connected    bool
	username     string
}

func NewClientSession(l logger.ILogger, c config.IClientConfig) *ClientSession {
	return &ClientSession{logger: l, clientConfig: c}
}

func (session *ClientSession) Start() {
	conn, err := session.Connect()
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	session.conn = conn
	session.reader = bufio.NewReader(conn)
	session.writer = bufio.NewWriter(conn)

	if !session.connected {
		session.waitAuth()
	}

	go func() {
		for {
			msg, err := session.reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading from server: %v", err)
				return
			}
			if strings.HasPrefix(msg, fmt.Sprintf("[%s]", session.username)) {
				continue
			}
			fmt.Print("\r\033[K")
			fmt.Print(msg)
			fmt.Print(fmt.Sprintf("[%s]>", session.username))
		}
	}()

	// TODO: Авторизацию по команде\конфиг по запросу
	// BUG: при печати если приходит сообщение - рвет строку ввода

	scanner := bufio.NewScanner(os.Stdin)
	_, _ = fmt.Fprintf(os.Stdin, fmt.Sprintf("[%s]>", session.username))
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/exit" {
			break
		}
		err = session.SendMessage(strings.TrimLeft(text, "@>"))
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
			break
		}
		_, _ = fmt.Fprintf(os.Stdin, fmt.Sprintf("[%s]>", session.username))
	}
}

func (session *ClientSession) waitAuth() {
	for !session.connected {
		fmt.Println("Authentication required!")
		cred := session.askCredentials()

		err := session.sendAuthRequest(cred)
		if err != nil {
			session.logger.Errorf("Send auth: %v", err)
			log.Fatalf("Send auth: %v", err)
		}

		buf := make([]byte, 1024)
		n, err := session.reader.Read(buf)
		if err != nil {
			session.logger.Errorf("Read conn: %v", err)
			log.Fatalf("Read conn: %v", err)
		}

		response := strings.TrimSpace(string(buf[:n]))
		session.logger.Infof("[AUTH] Server response: %s", response)

		if strings.Contains(response, "Successfully logged in") {
			session.username = cred.Username
			session.connected = true
		}
	}
}

func (session *ClientSession) sendAuthRequest(cred auth.Credentials) error {
	data := fmt.Sprintf("/auth %s %s\n", cred.Username, cred.Password)
	_, err := session.writer.WriteString(data)
	err = session.writer.Flush()
	if err != nil {
		return err
	}
	return err
}

func (session *ClientSession) askCredentials() auth.Credentials {
	username := readInput("Username: ")
	password := readInput("Password: ")
	return auth.Credentials{Username: username, Password: password}
}

func (session *ClientSession) Listen() {
	go func() {
		for {
			msg, err := session.reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading from session: %v", err)
				return
			}
			fmt.Print(msg)
		}
	}()
}

func (session *ClientSession) SendMessage(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	_, err := session.writer.WriteString(input + "\n")
	if err != nil {
		session.logger.Errorf("WriteString: %v", err)
		log.Fatalf("WriteString (sendMessage): %v", err)
	}
	return session.writer.Flush()
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func (session *ClientSession) Connect() (net.Conn, error) {
	address := fmt.Sprintf("%s:%s", session.clientConfig.ServerAddress(), session.clientConfig.ServerPort())
	return net.Dial("tcp", address)
}
