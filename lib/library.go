package mobileswarm

import (
	"C"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethersphere/swarm"
	config "github.com/ethersphere/swarm/api"
)

var swarmNode *node.Node

// DefaultBootnodeURL used for bootstrapping
const defaultBootnodeURL = "enode://4c113504601930bf2000c29bcd98d1716b6167749f58bad703bae338332fe93cc9d9204f08afb44100dc7bea479205f5d162df579f9a8f76f8b402d339709023@3.122.203.99:30301"
const defaultPassphrase = "test"
const keyStoreDirName = "keystore"

func getBootnodeURL(bootnodeURL string) string {
	if bootnodeURL == "" {
		return defaultBootnodeURL
	}
	return bootnodeURL
}

func newNodeWithKeystore(datadir string, ks *keystore.KeyStore, account accounts.Account) (stack *node.Node, _ error) {

	resultNode := &node.Node{}
	// Create the empty networking stack
	clientIdentifier := "SwarmMobile"
	// maxPeers := 10
	bootstrapNodes := []*enode.Node{}

	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.Version,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, keyStoreDirName),
		WSHost:      "localhost",
		WSPort:      8546,
		WSOrigins:   []string{"*"},
		WSModules:   []string{"pss"},
		// P2P: p2p.Config{
		// 	NoDiscovery:    true,
		// 	DiscoveryV5:    true,
		// 	ListenAddr:     ":0",
		// 	NAT:            nat.Any(),
		// 	MaxPeers:       maxPeers,
		// 	BootstrapNodes: append(bootstrapNodes, enode.MustParse(DefaultBootnodeURL)),
		// },
		P2P: p2p.Config{
			ListenAddr:     ":30303",
			MaxPeers:       50,
			NAT:            nat.Any(),
			BootstrapNodes: append(bootstrapNodes, enode.MustParse(defaultBootnodeURL)),
		},
	}

	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	pssAccount := account.Address.Hex()
	pssPassword := defaultPassphrase

	log.Info(fmt.Sprintf("pssAccount: %v, pssPassword %v", pssAccount, pssPassword))

	bzzSvc := func(ctx *node.ServiceContext) (node.Service, error) {
		//ks := rawStack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
		log.Info(fmt.Sprintf("ks: %v", ks))
		var a accounts.Account
		var err error
		if common.IsHexAddress(pssAccount) {
			//a, err = ks.Find(accounts.Account{Address: common.HexToAddress(config.PssAccount)})
			a = ks.Accounts()[0]
		} else if ix, ixerr := strconv.Atoi(pssAccount); ixerr == nil && ix > 0 {
			if accounts := ks.Accounts(); len(accounts) > ix {
				a = accounts[ix]
			} else {
				err = fmt.Errorf("index %d higher than number of accounts %d", ix, len(accounts))
			}
		} else {
			return nil, fmt.Errorf("Can't find swarm account key: %s", pssAccount)
		}
		if err != nil {
			return nil, fmt.Errorf("Can't find swarm account key: %v - Is the provided bzzaccount(%s) from the right datadir/Path?", err, pssAccount)
		}
		keyjson, err := ioutil.ReadFile(a.URL.Path)
		if err != nil {
			return nil, fmt.Errorf("Can't load swarm account key: %v", err)
		}

		log.Info(fmt.Sprintf("keyjson %v", keyjson))
		var bzzkey *ecdsa.PrivateKey
		//for i := 0; i < 3; i++ {
		//	password := getPassPhrase(fmt.Sprintf("Unlocking swarm account %s [%d/3]", a.Address.Hex(), i+1), i, passwords)
		//key, err := keystore.DecryptKey(keyjson, password)
		key, err := keystore.DecryptKey(keyjson, pssPassword)
		if err == nil {
			bzzkey = key.PrivateKey
		}
		//}
		if bzzkey == nil {
			return nil, fmt.Errorf("Can't decrypt swarm account key")
		}
		bzzconfig := config.NewConfig()
		bzzconfig.SyncEnabled = false
		bzzconfig.Path = rawStack.InstanceDir()
		bzzconfig.Init(bzzkey, bzzkey)

		return swarm.NewSwarm(bzzconfig, nil)
	}
	if err := rawStack.Register(bzzSvc); err != nil {
		return nil, fmt.Errorf("pss init: %v", err)
	}
	resultNode = rawStack
	return resultNode, nil
}

func makeDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func getAccount(ks *keystore.KeyStore, passphrase string) (accounts.Account, error) {
	if len(ks.Accounts()) > 0 {
		return ks.Accounts()[0], nil
	}
	return ks.NewAccount(passphrase)
}

//StartNode - start the Swarm node
func StartNode(path, listenAddr, cBootnodeURL, loglevel string) string {
	if swarmNode != nil {
		return "error 0: already started"
	}
	// set logging to stdout
	overrideRootLog(true, loglevel, "", false)

	log.Info("----------- starting node ---------------")
	dir := filepath.Join(path, keyStoreDirName)
	if err := makeDir(dir); err != nil {
		return "error during makeDir: " + err.Error()
	}

	ks := keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)

	account, err := getAccount(ks, defaultPassphrase)
	if err != nil {
		return "error 1.7: " + err.Error()
	}

	log.Info("after getAccount", "account", account)

	swarmNode, err = newNodeWithKeystore(dir, ks, account)
	if err != nil {
		return "error 2: " + err.Error()
	}
	err = swarmNode.Start()
	if err != nil {
		return "error 3: " + err.Error()
	}

	log.Info("----------- node started ---------------")
	return fmt.Sprintf("%v", account.Address.Hex())
}

//StopNode -
func StopNode() string {
	if swarmNode == nil {
		return "node already stopped"
	}
	err := swarmNode.Close()
	if err != nil {
		return "error stopping node: " + err.Error()
	}
	swarmNode = nil
	return "swarmNode closed"

}

// overrideRootLog overrides root logger with file handler, if defined,
// and log level (defaults to INFO).
func overrideRootLog(enabled bool, levelStr string, logFile string, terminal bool) error {
	if !enabled {
		disableRootLog()
		return nil
	}

	return enableRootLog(levelStr, logFile, terminal)
}

func disableRootLog() {
	log.Root().SetHandler(log.DiscardHandler())
}

func enableRootLog(levelStr string, logFile string, terminal bool) error {
	var (
		handler log.Handler
		err     error
	)

	if logFile != "" {
		handler, err = log.FileHandler(logFile, log.LogfmtFormat())
		if err != nil {
			return err
		}
	} else {
		handler = log.StreamHandler(os.Stdout, log.TerminalFormat(terminal))
	}

	if levelStr == "" {
		levelStr = "INFO"
	}

	level, err := log.LvlFromString(strings.ToLower(levelStr))
	if err != nil {
		return err
	}

	glogger := log.NewGlogHandler(handler)
	glogger.Verbosity(log.Lvl(level))

	log.Root().SetHandler(glogger)

	return nil
}
