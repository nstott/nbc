include $(GOROOT)/src/Make.inc

TARG=nbc
GOFILES=\
	ngram.go\
	mongongram.go\
	class.go\
	nbc.go\

include $(GOROOT)/src/Make.cmd
