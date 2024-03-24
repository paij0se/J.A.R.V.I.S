#!/bin/bash
echo "Downloading"
sudo apt-get install portaudio19-dev
git clone https://github.com/paij0se/J.A.R.V.I.S.
cd "J.A.R.V.I.S."
echo "Installing"
go get .
echo "compiling"
go build -o jarvis main.go
echo "Done!, run ./J.A.R.V.I.S./jarvis"
