package main

import (
	"bufio"
	"fmt"
	"mime"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mschneider82/milter"
)

const SUBJECTS_TXT_FILE string = "/etc/subjects.txt"

var (
	badstrings      []string
	fuckOffResponse milter.Response
)

type MyFilter struct {
}

func (e *MyFilter) Init(sid, mid string) {
	return
}

func (e *MyFilter) Disconnect() {
	return
}

func (e *MyFilter) Connect(name, value string, port uint16, ip net.IP, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) Helo(h string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) MailFrom(name string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) RcptTo(name string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	if name == "Subject" {
		dec := mime.WordDecoder{}
		decoded, decodeError := dec.DecodeHeader(value)

		if decodeError != nil {
			fmt.Printf("Passing because of decode error: %s\n", decodeError.Error())
		} else {

			fmt.Printf("Subject to analyze: \"%s\"\n", decoded)

			for _, badString := range badstrings {
				if strings.Contains(decoded, badString) {

					fmt.Printf("Bad string \"%s\" detected. Fuck off sent!\n", badString)

					return fuckOffResponse, nil
				}
			}

			fmt.Println("Nothing to nag about. Continuing!")
		}
	}

	return milter.RespContinue, nil
}

func (e *MyFilter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *MyFilter) Body(m *milter.Modifier) (milter.Response, error) {
	return milter.RespAccept, nil
}

func main() {
	badstrings = LoadBadStrings()
	fuckOffResponse = milter.NewResponseStr('y', "550 Fuck off")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP)

	go HandleSignals(signals)

	socket, socketErr := net.Listen("tcp", "127.0.0.1:1339")
	if socketErr != nil {
		fmt.Printf("Error creating socket: %s\n", socketErr.Error())
		os.Exit(1)
	} else {
		defer socket.Close()

		init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
			return &MyFilter{},
				milter.OptNone,
				milter.OptNoBody
		}

		errhandler := func(e error) {
			fmt.Printf("Panic happend: %s\n", e.Error())
		}

		server := milter.Server{
			Listener:      socket,
			MilterFactory: init,
			ErrHandlers:   []func(error){errhandler},
			Logger:        nil,
		}
		defer server.Close()

		fmt.Println("Subjectmilter initalized")

		server.RunServer()
	}
}

func LoadBadStrings() []string {
	fmt.Println("Loading badstrings")

	strings := make([]string, 0)

	file, err := os.Open(SUBJECTS_TXT_FILE)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", SUBJECTS_TXT_FILE, err.Error())
		return strings
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	fmt.Printf("Read %d subjects from %s\n", len(strings), SUBJECTS_TXT_FILE)

	return strings
}

func HandleSignals(signals chan os.Signal) {
	fmt.Println("Signal handler started")

	for {
		<-signals

		badstrings = LoadBadStrings()
	}
}
