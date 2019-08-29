### Mobileswarm

Build [Ethereum Swarm](https://swarm.ethereum.org/) on mobile

#### Dependencies

In order to build the keydrop-go library you will need to have `gomobile` installed.

https://godoc.org/golang.org/x/mobile/cmd/gomobile

#### Building the library

After this you can build Mobileswarm:

` $ cd $GOPATH/src/github.com/felfele/mobileswarm

Make android version:

` $ gomobile bind -target=android  -ldflags="-s -w" -o build/mobileswarm.aar github.com/felfele/mobileswarm/lib`

This will build an android archive (`.aar`) file called `mobileswarm.aar` in the `build/` directory. You can copy this file to your android project.

Make iOS version:

` $ gomobile bind -target=ios -ldflags="-s -w" -o build/bin/Mobileswarm.framework github.com/felfele/mobileswarm/lib`

This will build an iOS framework, called `Mobileswarm.framework` in the `build/` directory. You can copy this directory to your iOS project.
