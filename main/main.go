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

package main

import (
	"bufio"
	"flag"
	"github.com/dndx/bigmail"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	var (
		server      string // address of the SMTP server that will be used for sending.
		workers     int    // how much concurrent connection to the server should we use?
		contentPath string // path to the email to send
		content     []byte // path to the email to send
		list        string // list of email addresses to receive those message
		from        string
		useTLS      bool
		sleepTime   time.Duration
		sent        int
		failed      int
	)

	flag.IntVar(&workers, "workers", 1, "Number of workers to spawn.")
	flag.StringVar(&server, "server", "", "Address of the SMTP server that will be used for sending, including the port. eg. smtp.example.com:587")
	flag.StringVar(&contentPath, "content", "", "Path to the file containing the email content.")
	flag.StringVar(&list, "list", "", "Path to the file containing the recipient addresses.")
	flag.StringVar(&from, "from", "", "From field of the email.")
	flag.BoolVar(&useTLS, "usetls", true, "Use TLS when connecting to SMTP server?")
	flag.DurationVar(&sleepTime, "sleep", 0, "Time to wait before sending two emails with the same connection.")
	flag.Parse()

	if server == "" || contentPath == "" || list == "" || from == "" {
		flag.Usage()
		return
	}

	content, err := ioutil.ReadFile(contentPath)
	if err != nil {
		log.Fatalf("Could not read content file: %v", err)
	}

	listFile, err := os.Open(list)
	if err != nil {
		log.Fatalf("Could not read list file: %v", err)
	}
	defer listFile.Close()
	scanner := bufio.NewScanner(listFile)

	incoming := make(chan *bigmail.Message)
	errors := make(chan *bigmail.Message)

	for i := 0; i < workers; i++ {
		log.Println("Worker spawned")
		if _, err := bigmail.NewSender(server, incoming, errors, useTLS, sleepTime); err != nil {
			log.Panic(err)
		}
	}

	log.Println("Begin sending message")

	var done chan bool
	if !scanner.Scan() { // done!
		close(incoming)
		incoming = nil // disable sending channel since we are done sending
		// at this point there is still a pissibility that some pending jobs might fail. We need to be able to receive
		// those jobs from errors channel before calling Wait().
		done = make(chan bool)
		go func() {
			bigmail.Wait()
			done <- true
		}()
	}

	msg := &bigmail.Message{}
	msg.To = []string{scanner.Text()}
	msg.From = from
	msg.Body = content

loop:
	for {
		select {
		case m := <-errors: // respawn sender
			if _, err := bigmail.NewSender(server, incoming, errors, true, sleepTime); err != nil {
				log.Panic(err)
			}
			log.Println("Worker spawned")
			failed++
			log.Printf("Failed to send: %v", m.To)
		case incoming <- msg:
			log.Printf("Message to %v queued", msg.To)
			sent++
			if !scanner.Scan() { // done!
				close(incoming)
				incoming = nil // disable sending channel since we are done sending
				// at this point there is still a pissibility that some pending jobs might fail. We need to be able to receive
				// those jobs from errors channel before calling Wait().
				done = make(chan bool)
				go func() {
					bigmail.Wait()
					done <- true
				}()
			}

			msg = &bigmail.Message{}
			msg.To = []string{scanner.Text()}
			msg.From = from
			msg.Body = content
		case <-done: // done!
			break loop
		}
	}

	log.Printf("Message sending finished. Sent: %d, Failed: %d", sent, failed)
}
