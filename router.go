package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"

	"github.com/DonovanDiamond/milter"
)

type RouterMilter struct {
	milter.Milter
	host              string
	ip                net.IP
	helo              string
	from              string
	to                []string
	rejectedFrom      map[string]bool
	rejectedTo        map[string]bool
	rejectedToRegex   []*regexp.Regexp
	rejectedToSha256  map[string]bool
}

func (e *RouterMilter) Connect(host string, family string, port uint16, addr net.IP, m *milter.Modifier) (milter.Response, error) {
	e.host = host
	e.ip = addr
	log.Printf("Connect from %s [%s]", host, addr)
	return milter.RespContinue, nil
}

func (e *RouterMilter) Helo(name string, m *milter.Modifier) (milter.Response, error) {
	e.helo = name
	log.Printf("HELO: %s", name)
	return milter.RespContinue, nil
}

func (e *RouterMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	// reset state for new message
	e.from = from
	e.to = nil

	fromLower := strings.ToLower(from)

	// check if sender is rejected
	if e.rejectedFrom[fromLower] {
		log.Printf("[%s] Rejected MAIL FROM: %s", e.ip, from)
		return milter.RespReject, nil
	}
	log.Printf("[%s] MAIL FROM: %s", e.ip, from)
	return milter.RespContinue, nil
}

func (e *RouterMilter) RcptTo(to string, m *milter.Modifier) (milter.Response, error) {
	toLower := strings.ToLower(to)

	// check if recipient is rejected
	if e.rejectedTo[toLower] {
		log.Printf("[%s] Rejected RCPT TO: %s", e.ip, to)
		return milter.RespReject, nil
	}

	for _, re := range e.rejectedToRegex {
		if re.MatchString(toLower) {
			log.Printf("[%s] Rejected RCPT TO (regex): %s", e.ip, to)
			return milter.RespReject, nil
		}
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(toLower)))
	if e.rejectedToSha256[hash] {
		log.Printf("[%s] Rejected RCPT TO (hash): %s", e.ip, to)
		return milter.RespReject, nil
	}

	// save recipient address for later reference
	e.to = append(e.to, to)
	log.Printf("[%s] RCPT TO: %s", e.ip, to)
	return milter.RespContinue, nil
}

func (e *RouterMilter) Header(name string, value string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *RouterMilter) Headers(h textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *RouterMilter) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *RouterMilter) Body(m *milter.Modifier) (milter.Response, error) {
	log.Printf("[%s] Message from %s to %v accepted", e.ip, e.from, e.to)
	return milter.RespContinue, nil
}

func (e *RouterMilter) Abort(m *milter.Modifier) error {
	log.Printf("[%s] Transaction aborted", e.ip)
	e.from = ""
	e.to = nil
	return nil
}
