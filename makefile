BUILD_DIR=build
LDFLAGS="-s -w"
GIT_REPO=github.com/felfele/mobileswarm
LIB=$(GIT_REPO)/lib

all: clean android ios

android:
	-gomobile bind -target=android -o $(BUILD_DIR)/mobileswarm.aar -ldflags=$(LDFLAGS) $(LIB)

ios:
	-gomobile bind -target=ios -o $(BUILD_DIR)/Mobileswarm.framework -ldflags=$(LDFLAGS) $(LIB)

clean:
	-rm -rf ./$(BUILD_DIR)/*
