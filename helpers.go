package main

import "github.com/charmbracelet/bubbles/list"

func aeSliceToItem(ae []ActionableElement) []list.Item {
	itemToReturn := make([]list.Item, len(ae))
	for i, elem := range ae {
		itemToReturn[i] = elem
	}
	return itemToReturn
}
