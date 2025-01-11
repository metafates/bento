//go:build windows
// +build windows

package bento

// listenForResize is not available on windows because windows does not
// implement syscall.SIGWINCH.
func (a *appRunner) listenForResize(done chan struct{}) {
	close(done)
}
