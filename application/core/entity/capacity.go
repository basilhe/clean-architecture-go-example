package entity

type Capacity struct {
	Adsl  bool
	Fibre bool
}

func (c Capacity) HasAdslCapacity() bool {
	return c.Adsl
}
func (c Capacity) HasFibreCapacity() bool {
	return c.Fibre
}

