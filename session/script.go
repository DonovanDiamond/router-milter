package session

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dop251/goja"
)

func executeScript(script string, commands Commands) (result string, err error) {
	vm := goja.New()
	err = vm.Set("sha256", func(s string) string {
		hash := sha256.Sum256([]byte(s))
		return hex.EncodeToString(hash[:])
	})
	if err != nil {
		return "", fmt.Errorf("error setting sha256 function: %w", err)
	}
	_, err = vm.RunString(script)
	if err != nil {
		return "", fmt.Errorf("error running script: %w", err)
	}
	execute, ok := goja.AssertFunction(vm.Get("execute"))
	if !ok {
		return "", fmt.Errorf("script does not define an execute function")
	}
	r, err := execute(goja.Undefined(), vm.ToValue(commands))
	if err != nil {
		return "", err
	}
	return r.String(), nil
}
