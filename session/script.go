package session

import (
	"fmt"

	"github.com/dop251/goja"
)

func executeScript(script string, commands Commands) (result string, err error) {
	vm := goja.New()
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
