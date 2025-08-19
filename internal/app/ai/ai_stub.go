//go:build minimal

package ai

func IsAvailable() bool {
	return false
}

func Load(_ string) {
	// blank
}
