/* This file is part of bigmail.

bigmail is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

bigmail is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with bigmail.  If not, see <http://www.gnu.org/licenses/>. */

package bigmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

type Sender struct {
	incoming <-chan *Message // channel where messages to send will be received from
	client   *smtp.Client
	interval time.Duration // time between each send operation
}

var waiting sync.WaitGroup

// Creates a new sender, it will connect to the remote server and start the
// goroutine that does the actual sending operation.
//
// addr specifies the address of the SMTP server that the worker will connects to, as specified in smtp/Dial's document.
// incoming specifies the channel the sender will obtain jobs from.
// interval specifies the time sender sleeps between sending each message. 0 means no waiting.
func NewSender(addr string, incoming <-chan *Message, useTls bool, interval time.Duration) (*Sender, error) {
	var err error

	sender := Sender{incoming: incoming, interval: interval}
	if sender.client, err = smtp.Dial(addr); err != nil {
		return nil, err
	}

	if useTls {
		host, _, _ := net.SplitHostPort(addr)
		if err := sender.client.StartTLS(&tls.Config{ServerName: host}); err != nil {
			sender.client.Close()
			return nil, err
		}
	}

	waiting.Add(1)
	go sender.work()

	return &sender, nil
}

// The Sender's main work loop, will run in a dedicated goroutine.
func (s *Sender) work() {
	defer waiting.Done()

	for {
		var in <-chan *Message = s.incoming
		var wait <-chan time.Time

		select {
		case m, ok := <-in:
			if !ok { // time to quit
				s.client.Quit()
				return
			}

			if err := s.client.Mail(m.From); err != nil {
				log.Panic(err)
			}

			for _, to := range m.To {
				if err := s.client.Rcpt(to); err != nil {
					log.Panic(err)
				}
			}

			// Send the email body.
			wc, err := s.client.Data()
			if err != nil {
				log.Panic(err)
			}
			_, err = fmt.Fprintf(wc, `To: %s
From: %s
Subject: %s

%s`, strings.Join(m.To, ","), m.From, m.Subject, m.Body)
			if err != nil {
				log.Panic(err)
			}

			if err = wc.Close(); err != nil {
				log.Panic(err)
			}

			if s.interval > 0 {
				wait = time.After(s.interval)
				in = nil // disables sender for interval
			}
		case <-wait: // wait time is up
			in = s.incoming
		}
	}
}

// Blocks until all workers quit
func Wait() {
	waiting.Wait()
}
