package cpn

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	//    "github.com/deckarep/golang-set"
)

type Token interface {
	Tok2str() string
}

//Tipo Par
type NameTok struct {
	Name string
	Tok  Token
}

func NewNameTok(s string, t Token) *NameTok {
	a := NameTok{Name: s,
		Tok: t}
	return &a
}

//Plaza

type Plaza struct {
	vals   []Token      `default:"[]"`
	tipo   string       `default:""`
	pltype reflect.Type `default:nil`
	myArcs []*Arc       `default:"[]"`
}

func NewPlaza(t string) *Plaza {
	p := Plaza{tipo: t}
	return &p
}

func (p *Plaza) AddArc(a *Arc) {
	p.myArcs = append(p.myArcs, a)
}

func (p *Plaza) Add(t Token) {
	if p.pltype == nil {
		p.pltype = reflect.TypeOf(t)
	}
	var tt = reflect.TypeOf(t).String()
	if p.tipo == tt[strings.LastIndex(tt, ".")+1:] {
		//if p.tipo == tt {
		p.vals = append(p.vals, t)
	} else {
		fmt.Println("no se puedo agregar ", t.Tok2str(), "en", p.tipo, "con", tt)
	}
}

func (p *Plaza) AddTokens(ts []Token) {
	for _, t := range ts {
		p.Add(t)
	}
}

func (p *Plaza) Get(i int) (t Token) {
	return p.vals[i]
}

func (p *Plaza) Delete(t Token) (b bool) {
	ret := false
	for i, pp := range p.vals {
		if t == pp {
			p.vals = append(p.vals[:i], p.vals[i+1:]...)
			ret = true
		}
	}
	return ret
}

func (p *Plaza) GetAll() []Token {
	return p.vals
}

func (p *Plaza) IsEmpty() bool { return len(p.vals) == 0 }

//Arcos
type ArcGuard func(Token) bool

type ArcInject func([]Token) Token

type Arc struct {
	pl       *Plaza
	inout    bool
	varname  string   `default:""`
	varvalue Token    `default:nil`
	guar     ArcGuard `default:"func(t Token) {return true}"`
	injguar  string   `default:""`
	genfunc  ArcInject
}

func NewArc(p *Plaza) *Arc {
	a := Arc{pl: p,
		varvalue: nil}
	p.AddArc(&a)
	return &a
}

func NewInArc(p *Plaza) *Arc {
	a := NewArc(p)
	a.inout = true
	return a
}

func NewOutArc(p *Plaza) *Arc {
	a := NewArc(p)
	a.inout = false
	return a
}

func (a *Arc) IsInput() bool            { return a.inout }
func (a *Arc) IsOutput() bool           { return !a.inout }
func (a *Arc) IsEmpty() bool            { return a.pl.IsEmpty() }
func (a *Arc) Pl() *Plaza               { return a.pl }
func (a *Arc) SetVarname(s string)      { a.varname = s }
func (a *Arc) GetVarname() (s string)   { return a.varname }
func (a *Arc) GetVartype() reflect.Type { return a.pl.pltype }
func (a *Arc) DeleteToken(t Token)      { a.pl.Delete(t) }
func (a *Arc) SetInjGuard(s string)     { a.injguar = s }
func (a *Arc) GetInjGuard() (s string)  { return a.injguar }
func (a *Arc) SetGenFunc(f ArcInject)   { a.genfunc = f }
func (a *Arc) GetGenFunc() ArcInject    { return a.genfunc }

func (a *Arc) GetTokens() (ts *[]Token) {
	ret := make([]Token, 10, 15)
	for _, v := range a.pl.GetAll() {
		if a.guar(v) {
			ret = append(ret, v)
		}
	}
	return &ret
}

func (a *Arc) Absorbe(t Token) {
	a.pl.Delete(t)
}

func (a *Arc) InjectToken(ts []*NameTok) {
	// evaluar la guarda del arco,
	//var sen interface{}
	var err error
	var tss []Token
	fg := a.GetGenFunc()
	vars := make(map[string]interface{})
	for _, v := range ts {
		vars[v.Name] = v.Tok
		tss = append(tss, v.Tok)
	}
	_, err = gval.Evaluate(a.injguar, vars, gval.Full())
	if err != nil {
		fmt.Println("error al evaluar la guarda", a.injguar, err)
	} else {
		to := fg(tss)
		// que genere un token del tipo de la plaza y agregarlo en la plaza
		a.pl.Add(to.(Token))
	}

}

//Transicion
type Tran struct {
	name     string     `default:""`
	inArcs   []*Arc     `default:"[]"`
	outArcs  []*Arc     `default:"[]"`
	sensToks []*NameTok `default:"[]"`
	guard    string     `default:""`
}

func NewTran(s string) *Tran {
	a := Tran{name: s,
		guard: "return true"}
	return &a
}

func (t *Tran) AddInArc(a []*Arc)      { t.inArcs = append(t.inArcs, a...) }
func (t *Tran) AddOutArc(a []*Arc)     { t.outArcs = append(t.outArcs, a...) }
func (t *Tran) InArcs() []*Arc         { return t.inArcs }
func (t *Tran) OutArcs() []*Arc        { return t.outArcs }
func (t *Tran) GetName() string        { return t.name }
func (t *Tran) SetGuard(s string)      { t.guard = s }
func (t *Tran) GetGuard() string       { return t.guard }
func (t *Tran) IsSensible() bool       { return (len(t.sensToks) > 0) }
func (t *Tran) GetSenToks() []*NameTok { return t.sensToks }

func (t *Tran) AbsorbeTokens() {
	t.ComputeSenToks()
	if t.IsSensible() {
		for i, v := range t.GetSenToks() {
			t.InArcs()[i].DeleteToken(v.Tok)
		}
	}
}

func (t *Tran) InjectTokens() {
	if t.IsSensible() {
		for i, _ := range t.OutArcs() {
			t.OutArcs()[i].InjectToken(t.GetSenToks())
		}
	}
}

func (t *Tran) ComputeSenToks() {
	i := len(t.inArcs)
	t.sensToks = []*NameTok{}
	var sen interface{}
	var err error
	vars := make(map[string]interface{})
	switch i {
	case 1:
		{
			t.inArcs[0].varvalue = nil
		out1:
			for _, w := range t.inArcs[0].pl.GetAll() {
				vars[t.inArcs[0].varname] = w
				sen, err = gval.Evaluate(t.guard, vars)
				if sen == true {
					fmt.Println("La guarda ", t.guard, "evaluo con", sen, w.Tok2str())
					t.sensToks = append(t.sensToks, NewNameTok(t.inArcs[0].varname, w))
					break out1
				}
			}
		}
	case 2:
		{
			t.inArcs[0].varvalue = nil
			t.inArcs[1].varvalue = nil
		out2:
			for _, v := range t.inArcs[0].pl.GetAll() {
				for _, w := range t.inArcs[1].pl.GetAll() {
					vars[t.inArcs[0].varname] = v
					vars[t.inArcs[1].varname] = w
					sen, err = gval.Evaluate(t.guard, vars)
					if sen == true {
						fmt.Println("La guarda ", t.guard, "evaluo con", sen, v.Tok2str(), w.Tok2str())
						t.sensToks = append(t.sensToks, NewNameTok(t.inArcs[0].varname, v), NewNameTok(t.inArcs[1].varname, w))
						break out2
					}
				}
			}
		}
	case 3:
		{
			t.inArcs[0].varvalue = nil
			t.inArcs[1].varvalue = nil
			t.inArcs[2].varvalue = nil
		out3:
			for _, u := range t.inArcs[0].pl.GetAll() {
				for _, v := range t.inArcs[1].pl.GetAll() {
					for _, w := range t.inArcs[2].pl.GetAll() {
						vars[t.inArcs[0].varname] = u
						vars[t.inArcs[1].varname] = v
						vars[t.inArcs[2].varname] = w
						sen, err = gval.Evaluate(t.guard, vars)
						if sen == true {
							fmt.Println("La guarda ", t.guard, "evaluo con", sen, u.Tok2str(), v.Tok2str(), w.Tok2str())
							t.sensToks = append(t.sensToks, NewNameTok(t.inArcs[0].varname, u), NewNameTok(t.inArcs[1].varname, v), NewNameTok(t.inArcs[2].varname, w))
							break out3
						}
					}
				}
			}
		}
	case 4:
		{
			t.inArcs[0].varvalue = nil
			t.inArcs[1].varvalue = nil
			t.inArcs[2].varvalue = nil
			t.inArcs[3].varvalue = nil
		out4:
			for _, s := range t.inArcs[0].pl.GetAll() {
				for _, u := range t.inArcs[1].pl.GetAll() {
					for _, v := range t.inArcs[2].pl.GetAll() {
						for _, w := range t.inArcs[3].pl.GetAll() {
							vars[t.inArcs[0].varname] = s
							vars[t.inArcs[1].varname] = u
							vars[t.inArcs[2].varname] = v
							vars[t.inArcs[3].varname] = w
							sen, err = gval.Evaluate(t.guard, vars)
							if sen == true {
								fmt.Println("La guarda ", t.guard, "evaluo con", sen, s.Tok2str(), u.Tok2str(), v.Tok2str(), w.Tok2str())
								t.sensToks = append(t.sensToks, NewNameTok(t.inArcs[0].varname, s), NewNameTok(t.inArcs[1].varname, u), NewNameTok(t.inArcs[2].varname, v), NewNameTok(t.inArcs[3].varname, w))
								break out4
							}
						}
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Println("evaluo", t.guard, sen)
	}
}

func (t *Tran) Fire() {
	rand.Seed(time.Now().UnixNano())

	for {
		sec := rand.Intn(200)
		t.ComputeSenToks()
		if t.IsSensible() {
			fmt.Println("Ejecutando en transicion", t.GetName())
			t.AbsorbeTokens()
			time.Sleep(time.Duration(sec) * time.Millisecond)
			//			time.Sleep(1 * time.Second)
			t.InjectTokens()
		} else {
			time.Sleep(time.Duration(sec) * time.Millisecond)
		}
	}
}

//Modelo
type Model struct {
	trans   []*Tran `default:"[]"`
	firefun []func(*Tran)
}

func NewModel() *Model {
	a := Model{
		trans:   []*Tran{},
		firefun: []func(*Tran){}}
	return &a
}

func (m *Model) AddTran(a []*Tran) {
	m.trans = append(m.trans, a...)
}

func (m *Model) LogStatus() {
	for _, t := range m.trans {
		for _, a := range t.InArcs() {
			fmt.Println("en el arco entrada", t.GetName(), "con varname", a.GetVarname())
			for _, p := range a.Pl().GetAll() {
				fmt.Println("hay :", p.Tok2str())
			}
		}
		/*
			for _, a := range t.OutArcs() {
				fmt.Println("en el arco salida", t.GetName(), "con varname", a.GetVarname())
				for _, p := range a.Pl().GetAll() {
					fmt.Println("hay :", p.Tok2str())
				}
			}
		*/
	}

}

func (m *Model) FireFunc() {
	for _, t := range m.trans {
		go t.Fire()
	}
}
