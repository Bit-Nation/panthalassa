package module

// Writing a VM extension:
// Every VM extension must follow the "Module" interface

// Every function exposed to the VM should take an callback
// as it's last parameter. This callback will be called with
// at least an error (can be null of course). Those callbacks
// follow must receive the error as the first parameter.
// Every error that can't be passed to the callback
// (e.g. calling a callback can return an error - which can't be passed
// to the VM via the callback (since that's endless recursion))
// should be logged to a logger dedicated to the DApp
// @todo in the future, errors should shut down the DApp
