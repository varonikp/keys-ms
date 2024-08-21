package domain

type Software struct {
	id   int
	name string
}

type NewSoftwareData struct {
	ID   int
	Name string
}

func NewSoftware(data NewSoftwareData) Software {
	return Software{
		id:   data.ID,
		name: data.Name,
	}
}

func (s Software) ID() int {
	return s.id
}

func (s Software) Name() string {
	return s.name
}
