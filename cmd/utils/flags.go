package utils

import (
	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/config"
	"gopkg.in/urfave/cli.v1"
)

var (
	// General settings
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "use for store all files",
		Value: DirectoryString{common.DefaultDataDir()}, // TODO Distinguish different environmental addresses
	}

	KeyStoreDirFlag = DirectoryFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}

	// Config settings
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Json configuration file",
	}

	// Network Settings
	IdentityFlag = cli.StringFlag{
		Name:  "identity", //mapping:p2p.Name
		Usage: "Custom node name",
	}
	NetworkIdFlag = cli.UintFlag{
		Name: "networkid", //mapping:p2p.NetID
		Usage: "Network identifier (integer," +
			" 1=MainNet," +
			" 2=Aquarius," +
			" 3=Pisces," +
			" 4=Aries," +
			" 5=Taurus," +
			" 6=Gemini," +
			" 7=Cancer," +
			" 8=Leo," +
			" 9=Virgo," +
			" 10=Libra," +
			" 11=Scorpio," +
			" 12=Sagittarius," +
			" 13=Capricorn,)",
		Value: config.GlobalConfig.NetID,
	}
	MaxPeersFlag = cli.UintFlag{
		Name:  "maxpeers", //mapping:p2p.MaxPeers
		Usage: "Maximum number of network peers (network disabled if set to 0)",
		Value: config.GlobalConfig.MaxPeers,
	}
	MaxPendingPeersFlag = cli.UintFlag{
		Name:  "maxpendpeers", //mapping:p2p.MaxPendingPeers
		Usage: "Maximum number of pending connection attempts (defaults used if set to 0)",
		Value: config.GlobalConfig.MaxPeers,
	}
	ListenPortFlag = cli.IntFlag{
		Name:  "port", //mapping:p2p.Addr
		Usage: "Network listening port",
		Value: common.DefaultP2PPort,
	}
	NodeKeyHexFlag = cli.StringFlag{
		Name:  "nodekeyhex", //mapping:p2p.PrivateKey
		Usage: "P2P node key as hex",
	}

	//IPC Settings
	IPCEnabledFlag = cli.BoolFlag{
		Name:  "ipc",
		Usage: "Enable the IPC-RPC server",
	}
	IPCPathFlag = DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}

	//HTTP Settings
	RPCEnabledFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable the HTTP-RPC server",
	}
	RPCListenAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC server listening interface",
		Value: common.DefaultHTTPHost,
	}
	RPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port",
		Value: common.DefaultHTTPPort,
	}

	//WS Settings
	WSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface",
		Value: common.DefaultWSHost,
	}
	WSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port",
		Value: common.DefaultWSPort,
	}
)
