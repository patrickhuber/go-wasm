package types

type ComponentInstance struct {
	MayEnter bool
	MayLeave bool
	Handles  *HandleTables
}
