package core

func NewFact(name, value, source string, kind FactKind) *Fact {
	return &Fact{
		Name:   name,
		Value:  value,
		Kind:   kind,
		Source: source,
	}

}
