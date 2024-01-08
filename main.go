package main

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func main() {

}

func ReadPBM(filename string) (*PBM, error) {
	// ...
}

func (pbm *PBM) Size() (int, int) {
	// ...
}

func (pbm *PBM) At(x, y int) bool {
	// ...
}

func (pbm *PBM) Set(x, y int, value bool) {
	// ...
}

func (pbm *PBM) Save(filename string) error {
	// ...
}

func (pbm *PBM) Invert() {
	// ...
}

func (pbm *PBM) Flip() {
	// ...
}

func (pbm *PBM) Flop() {
	// ...
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	// ...
}