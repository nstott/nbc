include $(GOROOT)/src/Make.inc

TARG=nbc
GOFILES=\
	ngram.go\
	mongo.go\
	class.go\
	nbc.go\

include $(GOROOT)/src/Make.cmd
