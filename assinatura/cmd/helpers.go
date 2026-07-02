package cmd

import "fmt"

func flagError(err error, name string) error {
	if err != nil {
		return fmt.Errorf("falha ao ler flag %s: %w", name, err)
	}
	return nil
}
