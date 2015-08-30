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
	"github.com/dndx/bigmail"
	"log"
)

func main() {
	senders := new([10]*bigmail.Sender)
	incoming := make(chan *bigmail.Message)

	for i := range senders {
		var err error
		if senders[i], err = bigmail.NewSender("Your SMTP Server Address", incoming, true, 0); err != nil {
			log.Panic(err)
		}
	}

	msg := bigmail.Message{}
	msg.To = []string{"example@example.com"}
	msg.From = "test@example.com"
	msg.Subject = "Subject"
	msg.Body = "Hello from Golang!"

	log.Println("Sending message")
	for i := 0; i < 1; i++ {
		incoming <- &msg
		log.Println("Message sent")
	}
	close(incoming)
	bigmail.Wait()
}
