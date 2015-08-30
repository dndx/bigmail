# bigmail
A Golang project demonstrating concurrent programming by writing a simple mass email sender.

To use, edit main/main.go and specify the following information:

* SMTP server address and port, in the format of smtp.example.com:489
* From address
* To address
* Content

If you have multiple email to send, send all of them into the channel before closing it.

# Benchmark
Using a 10 worker setup, bigmail can send 20 emails/second over an Internet environemnt. Thus it is very efficient on sending large volume of emails.

# To Do
* Nicer CLI interface
* Better error detection. Currently it just panics if error was detected.

# License
GPL v2, see LICENSE for more details.
