# node-monitoring

Node monitoring tool - a lightweight monitoring tool for mx-chain-go node.

It is designed to be able to define multiple clients(plugins) to fetch the node info, and to be able to push notification events to one or multiple notifiers.

For now it has a single client which checks a node (or nodes) temp rating via multiversx api (select node by public key).

# How to use

* Compile the binary:
```bash
cd cmd/node && go build -o node-monitoring
```
OR just using a single make command:
```bash
make build
```

* Update config file at `cmd/node/config/config.toml`

* Start the app:
```bash
cd cmd/node && ./node-monitoring
```
OR
```bash
make run
```

## TODO/Improvements

- handle logging (file if needed) in a better way
- add more simple push notifiers (slack, telegram)
- evaluate adding separate config files for separate users
- ssh integration, in case the tool is to be run close to the node/nodes (for more specific monitoring)

- fetch nodes also by identifier, not only by bls keys
