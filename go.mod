module github.com/malumar/templater

go 1.17

replace (
	github.com/malumar/domaintools => ..\domaintools
	github.com/malumar/strutils => ..\strutils
)

require (
	github.com/malumar/domaintools v0.0.0-00010101000000-000000000000
	github.com/malumar/strutils v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20211118161319-6a13c67c3ce4
)
