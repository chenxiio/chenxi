package main

import "github.com/hashicorp/go-plugin"

type MyPlugin struct{}

func (p *MyPlugin) Add(a, b int) int {
	return a + b
}
func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "MY_PLUGIN",
			MagicCookieValue: "12345",
		},
		Plugins: map[string]plugin.Plugin{
			//"myplugin": &MyPlugin{},
		},
	})
}
