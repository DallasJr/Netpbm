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
	image, err := ReadPBM("p4ç.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	DisplayPBM(image);
	fmt.Println("Done loading image");
	width, height := image.Size()
	fmt.Println(image.magicNumber, width, "x", height)

	/*value := image.At(4, 0)
	fmt.Println("Pixel value at (4, 0) is", value)

	image.Set(4, 0, false);
	fmt.Println("Pixel value at (4, 0) changed to false")
	image.Set(4, 4, false);
	fmt.Println("Pixel value at (4, 4) changed to false")
	DisplayPBM(image);

	image.Invert();
	fmt.Println("Image inverted:")
	DisplayPBM(image);

	image.Flip();
	fmt.Println("Image flipped:")
	DisplayPBM(image);

	image.Flop();
	fmt.Println("Image flopped:")
	DisplayPBM(image);*/

	//image.SetMagicNumber("P4");

	err = image.Save("output.pbm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")
}

func DisplayPBM(pbm *PBM) {
	for y, _ := range pbm.data {
		output := "";
		for x, _ := range pbm.data[y] {
			val := pbm.data[y][x];
			if val {
				output = output + "1 ";
			} else {
				output = output + "0 ";
			}
		}
		fmt.Println(output);
	}
	fmt.Println("");
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
				if pbm.magicNumber == "P1" {
					if char == '1' {
						row = append(row, true)
						valid = true;
					} else if char == '0' {
						row = append(row, false)
						valid = true;
					}
				} else if pbm.magicNumber == "P4" {
					b := byte(char - '0')
					for i := 7; i >= 0; i-- {
						bit := (b >> uint(i)) & 1
						row = append(row, bit == 1)
					}
					valid = true
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
	return pbm.data[y][x];
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value;
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, pbm.magicNumber + "\n")
	fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pbm.magicNumber == "P1" {
				if pixel {
					fmt.Fprint(writer, "1")
				} else {
					fmt.Fprint(writer, "0")
				}
			} else if pbm.magicNumber == "P4" {
				for i := 0; i < len(row); i += 8 {
					var byteValue byte
					for j := 0; j < 8 && i+j < len(row); j++ {
						if row[i+j] {
							byteValue |= 1 << uint(7-j)
						}
					}
					fmt.Fprintf(writer, "%c", byteValue)
				}
			}
		}
		fmt.Fprintln(writer)
	}

	writer.Flush()

	return nil
}

func (pbm *PBM) Invert() {
	for y, _ := range pbm.data {
		for x, _ := range pbm.data[y] {
			pbm.data[y][x] = !pbm.data[y][x];
		}
	}
}

func (pbm *PBM) Flip() {
	for y, _ := range pbm.data {
		cursor := pbm.width - 1;
		for x := 0; x < pbm.width; x++ {
			temp := pbm.data[y][x];
			pbm.data[y][x] = pbm.data[y][cursor];
			pbm.data[y][cursor] = temp;
			cursor--;
			if cursor < x || cursor == x {
				break;
			}
		}
	}
}

func (pbm *PBM) Flop() {
	cursor := pbm.height - 1;
	for y, _ := range pbm.data {
		temp := pbm.data[y];
		pbm.data[y] = pbm.data[cursor];
		pbm.data[cursor] = temp;
		cursor--;
		if cursor < y || cursor == y {
			break;
		}
	}
}
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber;
}


/*
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