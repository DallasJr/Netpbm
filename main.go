package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func main() {
	image, err := ReadPBM("p1.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	for y, _ := range image.data {
		output := "";
		for x, _ := range image.data[y] {
			val := image.data[y][x];
			if val {
				output = output + "1";
			} else {
				output = output + "0";
			}
		}
		fmt.Println(output);
	}
	fmt.Println("Done loading image");

	width, height := image.Size()
	fmt.Println("Image Size:", width, "x", height)

	value := image.At(2, 3)
	fmt.Println("Pixel value at (2, 3):", value)
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pbm := &PBM{}
	var row []bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if pbm.magicNumber == "" {
			pbm.magicNumber = strings.TrimSpace(line)
		} else if pbm.width == 0 {
			fmt.Sscanf(line, "%d %d", &pbm.width, &pbm.height)
		} else {
			for _, char := range line {
				valid := false;
				if char == '1' {
					row = append(row, true)
					valid = true;
				} else if char == '0' {
					row = append(row, false)
					valid = true;
				}
				if valid {
					if len(row) == pbm.width {
						pbm.data = append(pbm.data, row)
						row = []bool{};
					}
				}
			}
			if len(pbm.data) == pbm.height {
				break;
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width;
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y];
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value;
}

/*func (pbm *PBM) Save(filename string) error {
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

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	// ...
}

func (pgm *PGM) Size() (int, int) {
	// ...
}

func (pgm *PGM) At(x, y int) uint8 {
	// ...
}

func (pgm *PGM) Set(x, y int, value uint8) {
	// ...
}

func (pgm *PGM) Save(filename string) error {
	// ...
}

func (pgm *PGM) Invert() {
	// ...
}

func (pgm *PGM) Flip() {
	// ...
}

func (pgm *PGM) Flop() {
	// ...
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	// ...
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	// ...
}

func (pgm *PGM) Rotate90CW() {
	// ...
}

func (pgm *PGM) ToPBM() *PBM {
	// ...
}*/