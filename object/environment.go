package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

//це зроблено для того аби кожна функція мала свій скоп, і не перезатирала скоп зовнішнішньої
//let a = 10;
// let b = 10;
// fn l(a,b) {
// 		let c = fn (a, b) {
// 			return a + b;
// 		}
// 		let a = 1;
// 		let b = 10;
// 		c(a, b);
//  }
//l(a,b) -> без внутрішнього скоупу викликався б як 1 та 10 бо всередині с ми б затерли зовнішній скоуп

func NewEnclosingEnvironment() *Environment {
	env := NewEnvironment()
	env.outer = env
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
