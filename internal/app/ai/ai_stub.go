//go:build minimal

package ai

var Model = ""

func IsAvailable() bool {
	return false
}

func Load(_ string) {
	// blank
}
