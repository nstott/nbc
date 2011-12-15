include $(GOROOT)/src/Make.inc

TARG=nbc
GOFILES=\
	ngram.go\
	mongongram.go\
	nbc.go\

include $(GOROOT)/src/Make.cmd
