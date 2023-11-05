package main

import "strings"

var icons = map[string]string{
	"desktop":      "\uf108",
	"downloads":    "\uf019",
	"pictures":     "\uf03e",
	"node_modules": "\ue5fa",
	"elm-stuff":    "\ue62c",
	".git":         "\ue5fb",
	".github":      "\ue5fd",
}

func GetIcon(dir Directory, isCurrent, displayIcons bool) string {
	if !displayIcons {
		return ""
	}

	if isCurrent {
		return "\uf07c "
	}

	if dir.IsHidden() {
		return "\uf114 "
	}

	if val, ok := icons[strings.ToLower(dir.Name())]; ok {
		return val + " "
	} else {
		return "\ue5fe" + " "
	}
}

func GetOrderIcon(order Order, displayIcons bool) string {
	if !displayIcons {
		return ""
	}

	switch order {
	case ORDER_NAME:
		return " \uf413 "
	case ORDER_TIME:
		return " \ue384 "
	}

	return ""
}
