package infra

type Config struct {
	Port    int   `json:"port,omitempty"`
	Genesis *bool `json:"genesis,omitempty"`
}
