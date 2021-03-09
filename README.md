# subjectmilter

A postfix milter (mail filter) checking the subject of an E-Mail. If it contains a predefined bad word the mail will be rejected.

## Postfix has that built in dummy! It's called header_checks!

I've been using header_checks for some time now. It works great when the subject header is in plain text. If it's RFC 2047 encoded it gets nasty. From my observations every spam mail encodes their subjects now and postfix won't decode it prior.

## Why not use SpamAssassin?

Big bloated software written in perl which adds another complexity layer for such an easy task.

## Installation

1) Build the binary: `make`
2) Upload it to your server
3) Copy and adapt the systemd service file
4) Create a text file under `/etc/subjects.txt` and fill it with bad strings or subjects line by line. See `example_subjects.txt`
5) Add or adapt `smtpd_milters` in your postfix config. Example: `smtpd_milters = inet:127.0.0.1:12301`
6) Reload postfix and start subjectmilter

## What happens afterwards?

If it detects bad strings the mail will be rejected with a `550 Fuck off` message. Example log output:

```
:-$ journalctl -u subjectmilter --follow
Mar 09 18:33:47 mail.veloc1ty.de subjectmilter[27318]: time="2021-03-09T18:33:47Z" level=info msg="Loading bad strings"
Mar 09 18:33:47 mail.veloc1ty.de subjectmilter[27318]: time="2021-03-09T18:33:47Z" level=info msg="Loaded bad strings" Amount=39
Mar 09 18:33:47 mail.veloc1ty.de subjectmilter[27318]: time="2021-03-09T18:33:47Z" level=info msg="Started signal handler"
Mar 09 18:33:47 mail.veloc1ty.de subjectmilter[27318]: time="2021-03-09T18:33:47Z" level=info msg="Subjectmilter initalized"
Mar 09 18:35:05 mail.veloc1ty.de subjectmilter[27318]: time="2021-03-09T18:35:05Z" level=info msg="Nothing to nag about" Subject="Stack Overflow Newsletter - Tuesday, March 9, 2021"
```

## Hot reloading

By sending a HUP signal to the process subjectmilter will read the changes from `/etc/subjects.txt`.
