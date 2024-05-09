package zones

type Federation struct {
	clocks Clock
	zones  []DBM
}

func (federation Federation) Foo() {

}
