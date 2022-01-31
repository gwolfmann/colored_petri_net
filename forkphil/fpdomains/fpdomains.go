package fpdomains

import (
	"strconv"
)

type Token interface {
	Tok2str() string
}

//tipo enumerado
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Dir2 int

const (
	Norte Dir2 = iota
	Sur
	Este
	Oeste
)

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func (d Dir2) String() string {
	return [...]string{"Norte", "Sur", "Este", "Oeste"}[d]
}

//no exiten los parametric type
//usar interface nula

func (d Direction) Tok2str() string {
	return d.String()
}

func (d Dir2) Tok2str() string {
	return d.String()
}

//dominio Filosofo y Tenedor
type Philo int

type Fork int

func (p Philo) String() string {
	return strconv.Itoa(int(p))
}

func (f Fork) String() string {
	return strconv.Itoa(int(f))
}

func (p Philo) Tok2str() string {
	return p.String()
}

func (f Fork) Tok2str() string {
	return f.String()
}
