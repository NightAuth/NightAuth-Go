module main

go 1.20

require (
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	nightauth v0.0.0
)

require golang.org/x/crypto v0.21.0 // indirect

replace nightauth => ../sdk
