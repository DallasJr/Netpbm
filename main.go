package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func main() {
	image, err := ReadPPM("p3.ppm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Done loading image");

	image.Invert();
	fmt.Println("Image inverted:")

	image.Flip();
	fmt.Println("Image flipped:")

	image.Flop();
	fmt.Println("Image flopped:")

	image.Rotate90CW();
	fmt.Println("Image rotated 90° clockwise:")

	changed := image.ToPBM();

	err = changed.Save("output.pbm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")

	err = image.Save("output.ppm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")
	/*image, err := ReadPGM("p2.pgm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Done loading image");
	DisplayPGM(image);
	fmt.Println("Magicnumber: " + image.magicNumber)
	width, height := image.Size()
	fmt.Println(width, "x", height)
	fmt.Println("Max:", image.Max())

	image.Invert();
	fmt.Println("Image inverted:")

	image.Flip();
	fmt.Println("Image flipped:")

	image.Flop();
	fmt.Println("Image flopped:")

	image.Rotate90CW();
	fmt.Println("Image rotated 90° clockwise:")

	changed := image.ToPBM();

	err = changed.Save("output.pbm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")
	err = image.Save("output.pgm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")*/

	/*image, err := ReadPBM("p4.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	DisplayPBM(image);
	fmt.Println("Done loading image");
	width, height := image.Size()
	fmt.Println(image.magicNumber, width, "x", height)

	value := image.At(4, 0)
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
	DisplayPBM(image);

	image.SetMagicNumber("P1");

	err = image.Save("output.pbm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")*/
}

func DisplayPBM(pbm *PBM) {
	for y, _ := range pbm.data {
		output := "";
		for x, _ := range pbm.data[y] {
			val := pbm.data[y][x];
			if val {
				output = output + "■ ";
			} else {
				output = output + "□ ";
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
	return pbm.width, pbm.height;
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

			}
		}
		fmt.Fprintln(writer, "")
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

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pgm := &PGM{}
	var row []uint8
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if pgm.magicNumber == "" {
			pgm.magicNumber = strings.TrimSpace(line)
		} else if pgm.width == 0 {
			fmt.Sscanf(line, "%d %d", &pgm.width, &pgm.height)
		} else if pgm.max == 0 {
			fmt.Sscanf(line, "%d", &pgm.max)
		} else {
			for _, val := range strings.Split(line, " ") {
				if val == "" {
					continue
				}
				if pgm.magicNumber == "P2" {
					num, _ := strconv.ParseUint(val, 10, 8)
					row = append(row, uint8(num));
				} else if pgm.magicNumber == "P5" {
					
				}
				if len(row) == pgm.width {
					pgm.data = append(pgm.data, row)
					row = []uint8{};
				}
			}
			if len(pgm.data) == pgm.height {
				break;
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return pgm, nil
}

func DisplayPGM(pgm *PGM) {
	for y, _ := range pgm.data {
		output := "";
		for x, _ := range pgm.data[y] {
			val := pgm.data[y][x];
			if val > 0 {
				output = output + "■ ";
			} else {
				output = output + "□ ";
			}
		}
		fmt.Println(output);
	}
	fmt.Println("");
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) Max() (int) {
	return pgm.max
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, pgm.magicNumber + "\n")
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)
	for _, row := range pgm.data {
		for _, pixel := range row {
			if pgm.magicNumber == "P2" {
				fmt.Fprint(writer, strconv.Itoa(int(pixel)) + " ")
			} else if pgm.magicNumber == "P5" {
				
			}
		}
		fmt.Fprintln(writer, "")
	}
	writer.Flush()
	return nil
}

func (pgm *PGM) Invert() {
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := int(pgm.data[y][x]);
			mid := pgm.max / 2;
			newvalue := 0; 
			if prevvalue > mid {
				newvalue = mid - (prevvalue - mid)
			} else {
				newvalue = mid + (mid - prevvalue)
			}
			pgm.data[y][x] = uint8(newvalue);
		}
	}
}

func (pgm *PGM) Flip() {
	for y, _ := range pgm.data {
		cursor := pgm.width - 1;
		for x := 0; x < pgm.width; x++ {
			temp := pgm.data[y][x];
			pgm.data[y][x] = pgm.data[y][cursor];
			pgm.data[y][cursor] = temp;
			cursor--;
			if cursor < x || cursor == x {
				break;
			}
		}
	}
}

func (pgm *PGM) Flop() {
	cursor := pgm.height - 1;
	for y, _ := range pgm.data {
		temp := pgm.data[y];
		pgm.data[y] = pgm.data[cursor];
		pgm.data[cursor] = temp;
		cursor--;
		if cursor < y || cursor == y {
			break;
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber;
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue);
}

func (pgm *PGM) Rotate90CW() {
    rotatedData := make([][]uint8, pgm.width)
    for i := range rotatedData {
        rotatedData[i] = make([]uint8, pgm.height)
    }
    for i := 0; i < pgm.width; i++ {
        for j := 0; j < pgm.height; j++ {
            rotatedData[i][j] = pgm.data[pgm.height-1-j][i]
        }
    }
    pgm.width, pgm.height = pgm.height, pgm.width
    pgm.data = rotatedData
}

func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{}
	pbm.magicNumber = "P1";
	pbm.height = pgm.height;
	pbm.width = pgm.width;
	for y, _ := range pgm.data {
		pbm.data = append(pbm.data, []bool{})
		for x, _ := range pgm.data[y] {
			val := int(pgm.data[y][x]);
			if val == 10 {
				pbm.data[y] = append(pbm.data[y], false);
			} else {
				pbm.data[y] = append(pbm.data[y], true);
			}
		}
	}
	return pbm
}

type PPM struct{
    data [][]Pixel
    width, height int
    magicNumber string
    max int
}

type Pixel struct{
    R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error){
    file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ppm := &PPM{}
	for scanner.Scan() {
		line := scanner.Text()
		var row []Pixel
		if strings.HasPrefix(line, "#") {
			continue
		}
		if ppm.magicNumber == "" {
			ppm.magicNumber = strings.TrimSpace(line)
		} else if ppm.width == 0 {
			fmt.Sscanf(line, "%d %d", &ppm.width, &ppm.height)
		} else if ppm.max == 0 {
			fmt.Sscanf(line, "%d", &ppm.max)
		} else {
			var pixel Pixel;
			pcs := 0
			for _, val := range strings.Split(line, " ") {
				if val == "" {
					continue
				}
				if ppm.magicNumber == "P3" {
					num, _ := strconv.ParseUint(val, 10, 8)
					if pcs == 0 {
						pixel.R = uint8(num);
						pcs++
					} else if pcs == 1 {
						pixel.G = uint8(num);
						pcs++
					} else {
						pixel.B = uint8(num);
						row = append(row, pixel);
						pixel = Pixel{};
						pcs = 0;
					}
				} else if ppm.magicNumber == "P5" {
					
				}
				if len(row) == ppm.width {
					ppm.data = append(ppm.data, row)
					row = []Pixel{};
				}
			}
			if len(ppm.data) == ppm.height {
				fmt.Println("all finished")
				break;
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ppm, nil
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, ppm.magicNumber + "\n")
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	fmt.Fprintf(writer, "%d\n", ppm.max)
	for _, row := range ppm.data {
		for _, pixel := range row {
			if ppm.magicNumber == "P3" {
				fmt.Fprint(writer, strconv.Itoa(int(pixel.R)) + " " + strconv.Itoa(int(pixel.G)) + " " + strconv.Itoa(int(pixel.B)) + " ")
			} else if ppm.magicNumber == "P5" {
				
			}
		}
		fmt.Fprintln(writer, "")
	}
	writer.Flush()
	return nil
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) Max() (int) {
	return ppm.max
}

func (ppm *PPM) At(x, y int) Pixel{
	return ppm.data[x][y]
}

func (ppm *PPM) Set(x, y int, value Pixel){
	ppm.data[x][y] = value
}

func (ppm *PPM) Invert() {
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x];
			pixel.R = uint8(255 - int(pixel.R));
			pixel.G = uint8(255 - int(pixel.G));
			pixel.B = uint8(255 - int(pixel.B));
			ppm.data[y][x] = pixel
		}
	}
}

func (ppm *PPM) Flip() {
	for y, _ := range ppm.data {
		cursor := ppm.width - 1;
		for x := 0; x < ppm.width; x++ {
			temp := ppm.data[y][x];
			ppm.data[y][x] = ppm.data[y][cursor];
			ppm.data[y][cursor] = temp;
			cursor--;
			if cursor < x || cursor == x {
				break;
			}
		}
	}
}

func (ppm *PPM) Flop() {
	cursor := ppm.height - 1;
	for y, _ := range ppm.data {
		temp := ppm.data[y];
		ppm.data[y] = ppm.data[cursor];
		ppm.data[cursor] = temp;
		cursor--;
		if cursor < y || cursor == y {
			break;
		}
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber;
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue);
}

func (ppm *PPM) Rotate90CW() {
    rotatedData := make([][]Pixel, ppm.width)
    for i := range rotatedData {
        rotatedData[i] = make([]Pixel, ppm.height)
    }
    for i := 0; i < ppm.width; i++ {
        for j := 0; j < ppm.height; j++ {
            rotatedData[i][j] = ppm.data[ppm.height-1-j][i]
        }
    }
    ppm.width, ppm.height = ppm.height, ppm.width
    ppm.data = rotatedData
}

func (ppm *PPM) ToPBM() *PBM{
	pbm := &PBM{}
	pbm.magicNumber = "P1";
	pbm.height = ppm.height;
	pbm.width = ppm.width;
	for y, _ := range ppm.data {
		pbm.data = append(pbm.data, []bool{})
		for x, _ := range ppm.data[y] {
			val := ppm.data[y][x];
			if int(val.R) == 0 && int(val.G) == 0 && int(val.B) == 0 {
				pbm.data[y] = append(pbm.data[y], true);
			} else {
				pbm.data[y] = append(pbm.data[y], false);
			}
		}
	}
	return pbm
}