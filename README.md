# Nuke GO

this project is intended to be nothing more than an in-memory cache server with local persistence.

if we can also make it faster, simpler to install in the cloud and less complicated than the others it will already be a great success


This project is written in GO, because I was interested in learning this language and above all because it seems to me to be an excellent compromise between speed and simplicity. (of course writing it in rust would have been better but... who knows, maybe one day...)


It includes by default a TCP server that only accepts the following 3 commands:
1) PUSH (to enter a value)
2) POP (to extract a value)
3) READ (to read a value)

