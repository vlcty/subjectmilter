# subjectmilter

A postfix milter (mail filter) checking the subject of an E-Mail. If it contains a bad word the mail will be rejected.

## Postfix has that built in dummy! It's called header_checks!

I've been using header_checks for some time now. It works great when the subject header is in plain text. If it's RFC 2047 encoded it gets nasty. But all spam mails encode their subjects now. And sadly postfix can't decode the subject first before doing the checks. A leightweight alternative had to be created.

## Why not use SpamAssassin

Big bloated software written in perl which adds another complexity layer for such an easy task. No thanks.

## Installation

1) Get the dependency: `go get github.com/mschneider82/milter`
2) Build the binary: `make`
3) Upload it to your server
4) Copy and adapt the systemd service file
5) Create a text file under `/etc/subjects.txt` and fill it with bad strings or subjects line by line. See `example_subjects.txt`
6) Add or adapt `smtpd_milters`. Example: `smtpd_milters = inet:127.0.0.1:12301`
7) Reload postfix and start subjectmilter

## What happens afterwards?

If it detects bad strings the mail will be rejected with `550 Fuck off`.

Example log output:

```
:-$ journalctl -u subjectmilter --follow
Feb 04 19:32:16 mail.veloc1ty.de subjectmilter[91368]: Subject to analyze: "Stack Overflow Newsletter - Tuesday, February 4, 2020"
Feb 04 19:32:16 mail.veloc1ty.de subjectmilter[91368]: Nothing to nag about. Continuing!
Feb 04 20:15:36 mail.veloc1ty.de subjectmilter[91368]: Subject to analyze: "Warten auf Fuckbuddy"
Feb 04 20:15:36 mail.veloc1ty.de subjectmilter[91368]: Bad string "Fuckbuddy" detected. Fuck off sent!
Feb 04 20:19:34 mail.veloc1ty.de subjectmilter[91368]: Subject to analyze: "Warten auf Fuckbuddy"
Feb 04 20:19:34 mail.veloc1ty.de subjectmilter[91368]: Bad string "Fuckbuddy" detected. Fuck off sent!
Feb 04 21:02:00 mail.veloc1ty.de subjectmilter[91368]: Subject to analyze: "I_Instacheat Request ausstehend"
Feb 04 21:02:00 mail.veloc1ty.de subjectmilter[91368]: Bad string "I_Instacheat" detected. Fuck off sent!
```

## Reloading bad words

After adding more subjects to `/etc/subjects.txt` you can send a HUP signal to the process. It reloads the bad strings afterwards.

## Future of this project

Maybe I'll adapt it in the future to use regex instead of fixed strings.
