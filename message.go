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

// A Message represents an email message that needs to be sent
// Bigmail will automatically add appropriate envelope headers to make the message valid before sending
type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}
