package main

import (
	"fmt"

	"net/url"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func LoginCommand(store KVStore, address string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to OTF",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := url.URL{
				Scheme: "https",
				Host:   address,
				Path:   "/app/settings/tokens",
			}

			if err := browser.OpenURL(u.String()); err != nil {
				return err
			}

			var token string

			fmt.Printf("Enter token: ")
			if _, err := fmt.Scanln(&token); err != nil {
				return err
			}

			if err := store.Save(address, token); err != nil {
				return err
			}

			fmt.Printf("Successfully added credentials for %s to %s\n", address, store)

			return nil
		},
	}

	return cmd
}
