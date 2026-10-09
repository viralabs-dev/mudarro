// Package winapi wraps the few kernel32 calls Mudarro needs on native Windows
// (Job Objects, console modes, console input polling and file locks). It uses
// only the standard syscall package, so Windows support adds no module
// dependency. Every other file in this package is Windows-only.
package winapi
