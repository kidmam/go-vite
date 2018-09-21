package main

import (
	"encoding/json"
	"fmt"
	"github.com/vitelabs/go-vite/cmd/utils"
	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/node"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

type NodeManager struct {
	ctx    *cli.Context
	node   *node.Node
	logger log15.Logger
}

func NewNodeManager(ctx *cli.Context) NodeManager {
	return NodeManager{
		ctx:    ctx,
		node:   makeFullNode(ctx),
		logger: log15.New("module", "main/nodemanager"),
	}
}

func (nodeManager *NodeManager) Start() error {

	// Start up the node
	utils.StartNode(nodeManager.node)

	nodeManager.node.Wait()
	return nil
}

func makeFullNode(ctx *cli.Context) *node.Node {

	nodeConfig := makeNodeConfig(ctx)

	node, err := node.New(nodeConfig)

	if err != nil {
		log15.Error("Failed to create the node: %v", err)
	}
	return node
}

func makeNodeConfig(ctx *cli.Context) *node.Config {

	cfg := node.DefaultNodeConfig
	cfg.Name = ClientIdentifier
	cfg.IPCPath = "vite.ipc"

	//Config log to file
	if fileName, e := cfg.NewRunLogDirFile(); e == nil {
		log15.Root().SetHandler(
			log15.LvlFilterHandler(log15.LvlInfo, log15.Must.FileHandler(fileName, log15.TerminalFormat())),
		)
	}

	// Load config file.
	if file := ctx.GlobalString(utils.ConfigFileFlag.Name); file != "" {
		//TODO 注意Logger覆盖
		if jsonConf, err := ioutil.ReadFile(file); err == nil {
			err = json.Unmarshal(jsonConf, &cfg)
			if err != nil {
				log15.Info("cannot unmarshal the config file content, will use the default config", "error", err)
			}
		} else {
			log15.Info("cannot read the config file, will use the default config", "error", err)
		}
	}

	// Apply flags
	mappingNodeConfig(ctx, &cfg)
	return &cfg
}

// SetNodeConfig applies node-related command line flags to the config.
func mappingNodeConfig(ctx *cli.Context, cfg *node.Config) {

	//Global Config
	if dataDir := ctx.GlobalString(utils.DataDirFlag.Name); len(dataDir) > 0 {
		cfg.DataDir = dataDir
	}

	if ctx.GlobalIsSet(utils.KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.GlobalString(utils.KeyStoreDirFlag.Name)
	}

	//Network Config
	if identity := ctx.GlobalString(utils.IdentityFlag.Name); len(identity) > 0 {
		cfg.Name = identity
		cfg.P2P.Name = identity
	}

	if ctx.GlobalIsSet(utils.NetworkIdFlag.Name) {
		cfg.P2P.NetID = ctx.GlobalUint(utils.NetworkIdFlag.Name)
	}

	if ctx.GlobalIsSet(utils.MaxPeersFlag.Name) {
		cfg.P2P.MaxPeers = ctx.GlobalUint(utils.MaxPeersFlag.Name)
	}

	// TODO p2p will use uint
	if ctx.GlobalIsSet(utils.MaxPendingPeersFlag.Name) {
		cfg.P2P.MaxPendingPeers = ctx.GlobalInt(utils.MaxPendingPeersFlag.Name)
	}

	if ctx.GlobalIsSet(utils.ListenPortFlag.Name) {
		cfg.P2P.Addr = fmt.Sprintf(":%d", ctx.GlobalInt(utils.ListenPortFlag.Name))
	}

	if ctx.GlobalIsSet(utils.NodeKeyHexFlag.Name) {
		cfg.P2P.PrivateKey = ctx.GlobalString(utils.NodeKeyHexFlag.Name)
	}

	//Ipc Config
	cfg.IPCEnabled = ctx.GlobalBool(utils.IPCEnabledFlag.Name)

	if ctx.GlobalIsSet(utils.IPCPathFlag.Name) {
		cfg.IPCPath = ctx.GlobalString(utils.IPCPathFlag.Name)
	}

	//Http Config
	cfg.RPCEnabled = ctx.GlobalBool(utils.RPCEnabledFlag.Name)

	if ctx.GlobalIsSet(utils.RPCListenAddrFlag.Name) {
		cfg.HttpHost = ctx.GlobalString(utils.RPCListenAddrFlag.Name)
	}

	if ctx.GlobalIsSet(utils.RPCPortFlag.Name) {
		cfg.HttpPort = ctx.GlobalInt(utils.RPCPortFlag.Name)
	}

	//WS Config
	cfg.WSEnabled = ctx.GlobalBool(utils.WSEnabledFlag.Name)

	if ctx.GlobalIsSet(utils.WSListenAddrFlag.Name) {
		cfg.WSHost = ctx.GlobalString(utils.WSListenAddrFlag.Name)
	}

	if ctx.GlobalIsSet(utils.WSPortFlag.Name) {
		cfg.WSPort = ctx.GlobalInt(utils.WSPortFlag.Name)
	}

}
