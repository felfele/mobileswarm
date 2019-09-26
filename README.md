## Mobileswarm

Build [Ethereum Swarm](https://swarm.ethereum.org/) on mobile

### Dependencies

In order to build the library you will need to have `gomobile` installed.

https://godoc.org/golang.org/x/mobile/cmd/gomobile

### Building the library

After this you can build Mobileswarm:

` $ cd $GOPATH/src/github.com/felfele/mobileswarm`

#### Make android version:

` $ gomobile bind -target=android  -ldflags="-s -w" -o build/mobileswarm.aar github.com/felfele/mobileswarm/lib`

This will build an android archive (`.aar`) file called `mobileswarm.aar` in the `build/` directory. You can copy this file to your android project.

Then you can use it in the Java code like this:

```java
import mobileswarm.Mobileswarm;

// ...

    String applicationPath = getAbsolutePath();
    String listenAddress = ":0"; // bind to all interface, any available port
    String bootnodeURL = "enode://4c113504601930bf2000c29bcd98d1716b6167749f58bad703bae338332fe93cc9d9204f08afb44100dc7bea479205f5d162df579f9a8f76f8b402d339709023@3.122.203.99:30301";
    String logLevel = "debug"; // can be info, trace etc.
    String startResult = Mobileswarm.startNode(applicationPath, listenAddress, bootnodeURL, logLevel);

    // ...

    String stopResult = Mobileswarm.stopNode();
```

#### Make iOS version:

` $ gomobile bind -target=ios -ldflags="-s -w" -o build/Mobileswarm.framework github.com/felfele/mobileswarm/lib`

This will build an iOS framework, called `Mobileswarm.framework` in the `build/` directory. You can copy this directory to your iOS project.

Then you can use it in the Objective-C code like this:

```objc

#import <Mobileswarm/Mobileswarm.h>

// ...

    NSString *appFolderPath = [self getPathForDirectory:NSDocumentDirectory];
    NSString *listenAddress = @":0"; // bind to all interface, any available port
    NSString *bootnodeURL = @"";
    NSString *logLevel = @"debug"; // can be info, trace etc.
    NSString *startResult = MobileswarmStartNode(appFolderPath, listenAddress, bootnodeURL, logLevel);
    NSLog(@"startResult: %@", startResult);

    // ...

    NSString *stopResult = MobileswarmStopNode();
    NSLog(@"stopResult: %@", stopResult);


```
