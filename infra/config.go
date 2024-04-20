package infra

const (
	GenesisNodeAddr = "/ip4/127.0.0.1/tcp/50000/p2p/QmPHZR5AdqSpKtcSZUCA5AqEni9sKgwPh7R6LxjZCLbXav"
	GenesisPort     = 50000
)

type Config struct {
	Port    int   `json:"port,omitempty"`
	Genesis *bool `json:"genesis,omitempty"`
}
