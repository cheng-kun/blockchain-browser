## Go Nebula AI

Fork from Official golang implementation of the Ethereum protocol.

Automated builds are available for stable releases and the unstable master branch.
Binary archives are published at https://github.com/nebulaai/nbai-node/releases.

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/nebulaai/nbai-node/wiki/Building-Nebula-AI)
on the wiki.

Building gnbai requires both a Go (version 1.7 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gnbai

or, to build the full suite of utilities:

    make all

## Executables

The go-gnbai project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`gnbai`** | Our main Nebula AI CLI client. It is the entry point into the Nebula AI network (main-, test- or private net), capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Ethereum network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `gnbai --help` and the [CLI Wiki page](https://github.com/nebulaai/nbai-node/wiki/Command-Line-Options) for command line options. |
| `abigen` | Source code generator to convert Nebula AI contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain [Ethereum contract ABIs](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI) with expanded functionality if the contract bytecode is also available. However it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](https://github.com/nebulaai/nbai-node/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts) wiki page for details. |
| `bootnode` | Stripped down version of our Nebula AI client implementation that only takes part in the network node discovery protocol, but does not run any of the higher level application protocols. It can be used as a lightweight bootstrap node to aid in finding peers in private networks. |
| `evm` | Developer utility version of the EVM (Ethereum Virtual Machine) that is capable of running bytecode snippets within a configurable environment and execution mode. Its purpose is to allow isolated, fine-grained debugging of EVM opcodes (e.g. `evm --code 60ff60ff --debug`). |
| `swarm`    | Swarm daemon and tools. This is the entrypoint for the Swarm network. `swarm --help` for command line options and subcommands. See [Swarm README](https://github.com/nebulaai/nbai-node/tree/master/swarm) for more information. |
| `puppeth`    | a CLI wizard that aids in creating a new Nebula AI network. |

## Running gnbai

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](https://github.com/nebulaai/nbai-node/wiki/Command-Line-Options)), but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own gnbai instance.

### Full node on the main Nebula AI network

By far the most common scenario is people wanting to simply interact with the Nebula AI network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can fast-sync quickly to the current
state of the network. To do so:

```
$ gnbai console
```

This command will:

 * Start gnbai in fast sync mode (default, can be changed with the `--syncmode` flag), causing it to
   download more data in exchange for avoiding processing the entire history of the Nebula AI network,
   which is very CPU intensive.
 * Start up gnbai's built-in interactive [JavaScript console](https://github.com/nebulaai/nbai-node/wiki/JavaScript-Console),
   (via the trailing `console` subcommand) through which you can invoke all official [`web3` methods](https://github.com/nbaiai/wiki/wiki/JavaScript-API)
   as well as gnbai's own [management APIs](https://github.com/nebulaai/nbai-node/wiki/Management-APIs).
   This tool is optional and if you leave it out you can always attach to an already running gnbai instance
   with `gnbai attach`.

```

### Configuration

As an alternative to passing the numerous flags to the `gnbai` binary, you can also pass a configuration file via:

```
$ gnbai --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ gnbai --your-favourite-flags dumpconfig
```

*Note: This works only with gnbai v0.0.1 and above.*

### Programatically interfacing Gnbai nodes

As a developer, sooner rather than later you'll want to start interacting with gnbai and the Nebula AI
network via your own programs and not manually through the console. To aid this, gnbai has built-in
support for a JSON-RPC based APIs ([standard APIs](https://github.com/nbaiai/wiki/wiki/JSON-RPC) and
[gnbai specific APIs](https://github.com/nebulaai/nbai-node/wiki/Management-APIs)). These can be
exposed via HTTP, WebSockets and IPC (unix sockets on unix based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by gnbai, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 8546)
  * `--wsapi` API's offered over the WS-RPC interface (default: "eth,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,debug,eth,miner,net,personal,shh,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a gnbai node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Nebula AI nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**

### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

#### Defining the private genesis state

First, you'll need to create the genesis state of your networks, which all nodes need to be aware of
and agree upon. This consists of a small JSON file (e.g. call it `genesis.json`):

```json
{
  "config": {
        "chainId": 0,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
    },
  "alloc"      : {},
  "coinbase"   : "0x0000000000000000000000000000000000000000",
  "difficulty" : "0x20000",
  "extraData"  : "",
  "gasLimit"   : "0x2fefd8",
  "nonce"      : "0x0000000000000042",
  "mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp"  : "0x00"
}
```

The above fields should be fine for most purposes, although we'd recommend changing the `nonce` to
some random value so you prevent unknown remote nodes from being able to connect to you. If you'd
like to pre-fund some accounts for easier testing, you can populate the `alloc` field with account
configs:

```json
"alloc": {
  "0x0000000000000000000000000000000000000001": {"balance": "111111111"},
  "0x0000000000000000000000000000000000000002": {"balance": "222222222"}
}
```

With the genesis state defined in the above JSON file, you'll need to initialize **every** gnbai node
with it prior to starting it up to ensure all blockchain parameters are correctly set:

```
$ gnbai init path/to/genesis.json
```

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent gnbai node pointed to the bootnode for peer discovery
via the `--bootnodes` flag. It will probably also be desirable to keep the data directory of your
private network separated, so do also specify a custom `--datadir` flag.

```
$ gnbai --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>
```

*Note: Since your network will be completely cut off from the main and test networks, you'll also
need to configure a miner to process transactions and create new blocks for you.*


## License

The go-gnbai library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-gnbai binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
