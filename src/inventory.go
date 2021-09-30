package main

import (
	"github.com/newrelic/infra-integrations-sdk/data/inventory"
)

// populateInventory This integration does not have inventory, so I simply return nil.
func populateInventory(inventory *inventory.Inventory, args argumentList) error {
	_, _ = inventory, args
	return nil
}
