package util

type Delegate struct {
	Data interface{}

	fns []func()
}

func (d *Delegate) Add(f func()) {
	d.fns = append(d.fns, f)
}

func (d Delegate) Invoke() {
	for _, f := range d.fns {
		f()
	}
}
