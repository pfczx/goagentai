package agent

import (
)

type Agent struct {
	profile *Profile
}

func NewAgent(profile *Profile) *Agent {
	return &Agent{
		profile: profile,
	}
}

func (a *Agent) Ask(input string) (string, error) {
return  "",nil
	
}
