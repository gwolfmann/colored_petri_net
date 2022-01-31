package main

import (
	"cpn"
	"fmt"
	"fpdomains"
	"time"
	//    "github.com/deckarep/golang-set"
)

func main() {

	var t1 fpdomains.Fork = 0
	var t2 fpdomains.Fork = 1
	var t3 fpdomains.Fork = 2
	var t4 fpdomains.Fork = 3
	var t5 fpdomains.Fork = 4

	var f1 fpdomains.Philo = 0
	var f2 fpdomains.Philo = 1
	var f3 fpdomains.Philo = 2
	var f4 fpdomains.Philo = 3
	var f5 fpdomains.Philo = 4

	piensa := cpn.NewPlaza("Philo")

	come := cpn.NewPlaza("Philo")

	tene := cpn.NewPlaza("Fork")

	piensa.AddTokens([]cpn.Token{f1, f2, f3, f4, f5})
	tene.AddTokens([]cpn.Token{t1, t2, t3, t4, t5})

	ai1 := cpn.NewInArc(piensa)
	ai1.SetVarname("p")
	ai2 := cpn.NewInArc(come)
	ai2.SetVarname("p")
	ai3 := cpn.NewInArc(tene)
	ai3.SetVarname("i")
	ai4 := cpn.NewInArc(tene)
	ai4.SetVarname("d")

	ao1 := cpn.NewOutArc(piensa)
	ao1.SetVarname("p")
	ao1.SetInjGuard("p")
	ao1f := func(ts []cpn.Token) cpn.Token {
		return ts[0].(fpdomains.Philo)
	}
	ao1.SetGenFunc(cpn.ArcInject(ao1f))

	ao2 := cpn.NewOutArc(come)
	ao2.SetVarname("p")
	ao2.SetInjGuard("p")
	ao2f := func(ts []cpn.Token) cpn.Token {
		return ts[0].(fpdomains.Philo)
		//		return fpdomains.Fork(ts[0].(fpdomains.Philo))
	}
	ao2.SetGenFunc(cpn.ArcInject(ao2f))

	ao3 := cpn.NewOutArc(tene)
	ao3.SetVarname("i")
	ao3.SetInjGuard("p")
	ao3f := func(ts []cpn.Token) cpn.Token {
		//return ts[0].(fpdomains.Philo)
		return fpdomains.Fork(ts[0].(fpdomains.Philo))
	}
	ao3.SetGenFunc(cpn.ArcInject(ao3f))

	ao4 := cpn.NewOutArc(tene)
	ao4.SetVarname("d")
	ao4.SetInjGuard("p")
	//	ao4.SetInjGuard("(p+1)%5")
	ao4f := func(ts []cpn.Token) cpn.Token {
		a := fpdomains.Fork((ts[0].(fpdomains.Philo) + 1) % 5)
		return a
	}
	ao4.SetGenFunc(cpn.ArcInject(ao4f))

	starteat := cpn.NewTran("empieza")
	endeat := cpn.NewTran("termina")

	starteat.AddInArc([]*cpn.Arc{ai1, ai3, ai4})
	starteat.AddOutArc([]*cpn.Arc{ao2})
	starteat.SetGuard("(i==p) && (d==(p+1)%5)")

	endeat.AddInArc([]*cpn.Arc{ai2})
	endeat.AddOutArc([]*cpn.Arc{ao1, ao3, ao4})
	endeat.SetGuard("1==1")

	m := cpn.NewModel()
	m.AddTran([]*cpn.Tran{starteat, endeat})

	fmt.Println("transicion empieza con guarda", starteat.GetGuard())
	fmt.Println("transicion termina con guarda", endeat.GetGuard())
	/*
		for i := 1; i <= 5; i++ {
			fmt.Println("Status antes del step", i)
			m.LogStatus()
			m.Step()
		}
		fmt.Println("Status final")
	*/
	m.FireFunc()
	time.Sleep(6 * time.Second)

	defer m.LogStatus()
}
