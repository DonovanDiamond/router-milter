package session

import (
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/DonovanDiamond/milter"
)

type Session struct {
	milter.Milter

	host string
	ip   net.IP

	commands Commands
	script   string
}

func (s *Session) Log(msg string, args ...any) {
	log.Printf("[%s] %s", s.ip, fmt.Sprintf(msg, args...))
}

const OptAction = milter.OptAddHeader | milter.OptChangeHeader
const OptProtocol = milter.OptNoBody | milter.OptNoHeaders | milter.OptNoEOH

func NewSession(script string) *Session {
	return &Session{
		script: script,
	}
}

type Commands struct {
	HELO string   `json:"helo"`
	FROM string   `json:"mail"`
	RCPT []string `json:"rcpt"`
}

func (s *Session) Connect(host string, family string, port uint16, addr net.IP, m *milter.Modifier) (milter.Response, error) {
	s.host = host
	s.ip = addr
	s.Log("Connect from %s", host)
	return milter.RespContinue, nil
}

func (s *Session) Helo(name string, m *milter.Modifier) (milter.Response, error) {
	s.commands.HELO = name
	s.Log("HELO: %s", name)
	return milter.RespContinue, nil
}

func (s *Session) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	s.commands.FROM = strings.ToLower(from)
	// reset email for new from address
	s.commands.RCPT = []string{}
	s.Log("MAIL FROM: %s", from)
	return milter.RespContinue, nil
}

func (s *Session) RcptTo(to string, m *milter.Modifier) (milter.Response, error) {
	s.commands.RCPT = append(s.commands.RCPT, strings.ToLower(to))
	s.Log("RCPT TO: %s", to)
	return milter.RespContinue, nil
}

func (s *Session) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	s.Log("HEADER: %s -> %s", name, value)
	return milter.RespContinue, nil
}

func (s *Session) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	s.Log("HEADERS: %d headers processed", len(headers))
	return milter.RespContinue, nil
}

func (s *Session) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	s.Log("BODY CHUNK: %d bytes", len(chunk))
	return milter.RespContinue, nil
}

func (s *Session) Body(m *milter.Modifier) (milter.Response, error) {
	s.Log("BODY: Processing script...")
	// execute the script with the collected commands
	result, err := executeScript(s.script, s.commands)
	if err != nil {
		s.Log("Error executing script: %v", err)
		return milter.RespTempFail, nil
	}
	s.Log("Script result: %s", result)
	if result == "reject" {
		return milter.RespReject, nil
	}
	return milter.RespContinue, nil
}

func (s *Session) Abort(m *milter.Modifier) error {
	s.Log("ABORT: Transaction aborted")
	s.commands.FROM = ""
	s.commands.RCPT = []string{}
	return nil
}
