package main

import (
	"bufio"
	"mime"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mschneider82/milter"
	log "github.com/sirupsen/logrus"
)

const (
	SUBJECTS_TXT_FILE string = "/etc/subjects.txt"
	FUCK_OFF_MESSAGE  string = "550 Fuck off"
)

type MyFilter struct {
	badstrings []string

	decodedSubject    string
	detectedBadString string
}

func (filter *MyFilter) ContainsBadString(value string) bool {
	decoder := mime.WordDecoder{}
	decoded, decodeError := decoder.DecodeHeader(value)

	if decodeError != nil {
		return false
	}

	decoded = strings.ReplaceAll(decoded, "\n", "")
	filter.decodedSubject = decoded

	for _, badString := range filter.badstrings {
		if strings.Contains(decoded, badString) {
			filter.detectedBadString = badString
			return true
		}
	}

	return false
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
		containsBadString := e.ContainsBadString(value)
		logger := log.WithField("Subject", e.decodedSubject)

		if containsBadString {
			logger.WithField("Bad string", e.detectedBadString).Info("Detected bad word. Fuck off!")
			return milter.NewResponseStr(milter.SMFIR_REPLYCODE, FUCK_OFF_MESSAGE), nil
		} else {
			logger.Info("Nothing to nag about")
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
	badstrings := LoadBadStrings()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP)

	go func() {
		log.Info("Started signal handler")

		for {
			<-signals

			badstrings = LoadBadStrings()

			log.WithField("Amount", len(badstrings)).Info("Loaded badstrings")
		}
	}()

	socket, socketErr := net.Listen("tcp", "127.0.0.1:1339")
	if socketErr != nil {
		log.WithField("Error", socketErr.Error()).Fatal("Was not able to bind port")
	} else {
		defer socket.Close()

		init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
			return &MyFilter{
					badstrings: badstrings},
				milter.OptNone,
				milter.OptNoBody
		}

		errhandler := func(e error) {
			log.WithField("Error", e.Error()).Error("Error while parsing message")
		}

		server := milter.Server{
			Listener:      socket,
			MilterFactory: init,
			ErrHandlers:   []func(error){errhandler},
			Logger:        nil,
		}
		defer server.Close()

		log.Info("Subjectmilter initalized")

		server.RunServer()
	}
}

func LoadBadStrings() []string {
	log.Info("Loading bad strings")

	strings := make([]string, 0)

	file, err := os.Open(SUBJECTS_TXT_FILE)
	if err != nil {
		log.WithField("Error", err.Error()).Fatal("Error reading bad strings")
		return strings
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	log.WithField("Amount", len(strings)).Info("Loaded bad strings")

	return strings
}
