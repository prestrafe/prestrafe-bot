package twitchbot

import "fmt"

type chatCommandParameter struct {
	name     string
	required bool
	pattern  string
}

func (p *chatCommandParameter) String() string {
	requiredIndicator := ""
	if !p.required {
		requiredIndicator = "?"
	}

	return fmt.Sprintf("<%s%s>", p.name, requiredIndicator)
}

func (p *chatCommandParameter) getPattern() string {
	requiredFlag := ""
	if !p.required {
		requiredFlag = "?"
	}
	return fmt.Sprintf("(?P<%s>%s)%s", p.name, p.pattern, requiredFlag)
}
