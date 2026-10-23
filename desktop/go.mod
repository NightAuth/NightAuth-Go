module main

go 1.20

require (
	nightauth v0.0.0
)

require golang.org/x/crypto v0.21.0 // indirect

replace nightauth => ../sdk
