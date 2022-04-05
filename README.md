# frameup

[![Go](https://github.com/royge/frameup/actions/workflows/go.yml/badge.svg)](https://github.com/royge/frameup/actions/workflows/go.yml)
[![Report](https://goreportcard.com/badge/github.com/royge/frameup)](https://goreportcard.com/report/github.com/royge/frameup)

Create custom frame or overlay PNG image to selected pictures

## Getting Started

1. Install using `go get`

		$ go get github.com/royge/frameup

1. Configure

		$ cp frameup.json.dist frameup.json

1. Help Usage

		$ frameup -h

## Usage:

	frameup [flags]
	frameup [command]

**Available Commands:**

	crop        Crop selected pictures inside the source directory.
	frame       Create frame on selected pictures inside source directory.
	help        Help about any command

**Flags:**

	-b, --bg string        Background image. (default "./assets/bg.jpg")
	-d, --delay int        Delay. (default 100)
	-e, --ext string       Picture files allowed extensions. (default ".jpg")
	-h, --help             help for frameup
	-o, --output string    Output directory.
	-l, --overlay string   Overlay or frame image. (default "./assets/overlay.jpg")
	-s, --source string    Source directory.

## Example

	$ frameup crop -s ./demo/raw/ -o ./demo/cropped

  Go to the `./demo/cropped` directory and remove the bad cuts. Only keep one
  for every dimension.

	$ frameup frame -s ./demo/cropped -o ./demo/final
