package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
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
	file, err := os.Open(filename);
	if err != nil {
		return nil, err;
	}
	defer file.Close();
	scanner := bufio.NewScanner(file);
	pbm := &PBM{};
	line := 0;
	for scanner.Scan() {
		text := scanner.Text();
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		if pbm.magicNumber == "" {
			pbm.magicNumber = strings.TrimSpace(text);
		} else if pbm.width == 0 {
			fmt.Sscanf(text, "%d %d", &pbm.width, &pbm.height);
			pbm.data = make([][]bool, pbm.height);
			for i := range pbm.data {
				pbm.data[i] = make([]bool, pbm.width);
			}
		} else {
			if pbm.magicNumber == "P1" {
				test := strings.Fields(text);
				for i := 0; i < pbm.width; i++ {
					pbm.data[line][i] = test[i] == "1";
				}
				line++
			} else if pbm.magicNumber == "P4" {
				expectedBytesPerRow := (pbm.width + 7) / 8;
				totalExpectedBytes := expectedBytesPerRow * pbm.height;
				allPixelData := make([]byte, totalExpectedBytes);
				fileContent, err := os.ReadFile(filename);
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err);
				}
				copy(allPixelData, fileContent[len(fileContent)-totalExpectedBytes:]);
				byteIndex := 0;
				for y := 0; y < pbm.height; y++ {
					for x := 0; x < pbm.width; x++ {
						if x%8 == 0 && x != 0 {
							byteIndex++;
						}
						pbm.data[y][x] = (allPixelData[byteIndex]>>(7-(x%8)))&1 != 0;
					}
					byteIndex++;
				}
			}
		}
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
	writer.Flush();
	if pbm.magicNumber == "P1" {
		for y, row := range pbm.data {
			for i, pixel := range row {
				xtra := " ";
				if i == len(row) - 1 {
					xtra = "";
				}
				if pixel {
					fmt.Fprint(writer, "1" + xtra)
				} else {
					fmt.Fprint(writer, "0" + xtra)
				}
			}
			if y != len(pbm.data) - 1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush();
	} else if pbm.magicNumber == "P4" {
        for _, row := range pbm.data {
            for x := 0; x < pbm.width; x += 8 {
                var byteValue byte;
                for i := 0; i < 8 && x+i < pbm.width; i++ {
                    bitIndex := 7 - i;
                    if row[x+i] {
                        byteValue |= 1 << bitIndex;
                    }
                }
                _, err = file.Write([]byte{byteValue})
                if err != nil {
                    return fmt.Errorf("error writing pixel data: %v", err)
                }
            }
        }
    }
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
	max           uint8
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pgm := &PGM{}
	line := 0;
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		if pgm.magicNumber == "" {
			pgm.magicNumber = strings.TrimSpace(text)
		} else if pgm.width == 0 {
			fmt.Sscanf(text, "%d %d", &pgm.width, &pgm.height)
			pgm.data = make([][]uint8, pgm.height);
			for i := range pgm.data {
				pgm.data[i] = make([]uint8, pgm.width);
			}
		} else if pgm.max == 0 {
			fmt.Sscanf(text, "%d", &pgm.max)
		} else {
			if pgm.magicNumber == "P2" {
				val := strings.Fields(text);
				for i := 0; i < pgm.width; i++ {
					num, _ := strconv.ParseUint(val[i], 10, 8)
					pgm.data[line][i] = uint8(num);
				}
				line++
			} else if pgm.magicNumber == "P5" {
				pixelData := make([]uint8, pgm.width*pgm.height)
                fileContent, err := os.ReadFile(filename)
                if err != nil {
                    return nil, fmt.Errorf("couldn't read file: %v", err)
                }
                copy(pixelData, fileContent[len(fileContent)-(pgm.width*pgm.height):])
                pixelIndex := 0
                for y := 0; y < pgm.height; y++ {
                    for x := 0; x < pgm.width; x++ {
                        pgm.data[y][x] = pixelData[pixelIndex]
                        pixelIndex++
                    }
                }
			}
		}
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

func (pgm *PGM) Max() (uint8) {
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
	writer.Flush();
	if pgm.magicNumber == "P2" {
		for y, row := range pgm.data {
			for i, pixel := range row {
				xtra := " ";
				if i == len(row) - 1 {
					xtra = "";
				}
				fmt.Fprint(writer, strconv.Itoa(int(pixel)) + xtra)
			}
			if y != len(pgm.data) - 1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush();
	} else if pgm.magicNumber == "P5" {
        for _, row := range pgm.data {
            for _, pixel := range row {
                _, err = file.Write([]byte{pixel})
                if err != nil {
                    return fmt.Errorf("error writing pixel data: %v", err)
                }
            }
        }
    }
	return nil
}

func (pgm *PGM) Invert() {
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x];
			pgm.data[y][x] = pgm.max - prevvalue;
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
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x];
			newvalue := prevvalue*uint8(5)/pgm.max
			pgm.data[y][x] = newvalue;
		}
	}
	pgm.max = maxValue;
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
			grayValue := pgm.data[y][x]
			isBlack := grayValue < pgm.max/2
			pbm.data[y] = append(pbm.data[y], isBlack)
		}
	}
	return pbm
}

type PPM struct{
    data [][]Pixel
    width, height int
    magicNumber string
    max uint8
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
	line := 0;
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		if ppm.magicNumber == "" {
			ppm.magicNumber = strings.TrimSpace(text)
		} else if ppm.width == 0 {
			fmt.Sscanf(text, "%d %d", &ppm.width, &ppm.height)
			ppm.data = make([][]Pixel, ppm.height);
			for i := range ppm.data {
				ppm.data[i] = make([]Pixel, ppm.width);
			}
		} else if ppm.max == 0 {
			fmt.Sscanf(text, "%d", &ppm.max)
		} else {
			if ppm.magicNumber == "P3" {
                val := strings.Fields(text)
                for i := 0; i < ppm.width; i++ {
                    r, _ := strconv.ParseUint(val[i*3], 10, 8)
                    g, _ := strconv.ParseUint(val[i*3+1], 10, 8)
                    b, _ := strconv.ParseUint(val[i*3+2], 10, 8)
                    ppm.data[line][i] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
                }
                line++
            } else if ppm.magicNumber == "P6" {
                pixelData := make([]byte, ppm.width*ppm.height*3)
                fileContent, err := os.ReadFile(filename)
                if err != nil {
                    return nil, fmt.Errorf("couldn't read file: %v", err)
                }
                copy(pixelData, fileContent[len(fileContent)-(ppm.width*ppm.height*3):])
                pixelIndex := 0
                for y := 0; y < ppm.height; y++ {
                    for x := 0; x < ppm.width; x++ {
                        ppm.data[y][x].R = pixelData[pixelIndex]
                        ppm.data[y][x].G = pixelData[pixelIndex+1]
                        ppm.data[y][x].B = pixelData[pixelIndex+2]
                        pixelIndex += 3
                    }
                }
            }
		}
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
	writer.Flush();
	if ppm.magicNumber == "P3" {
        for y, row := range ppm.data {
            for i, pixel := range row {
                xtra := " "
                if i == len(row)-1 {
                    xtra = ""
                }
                fmt.Fprintf(writer, "%d %d %d%s", pixel.R, pixel.G, pixel.B, xtra)
            }
            if y != len(ppm.data)-1 {
                fmt.Fprintln(writer, "")
            }
        }
        writer.Flush()
    } else if ppm.magicNumber == "P6" {
        for _, row := range ppm.data {
            for _, pixel := range row {
                _, err = file.Write([]byte{pixel.R, pixel.G, pixel.B})
                if err != nil {
                    return fmt.Errorf("error writing pixel data: %v", err)
                }
            }
        }
    }
	return nil
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) Max() (uint8) {
	return ppm.max
}

func (ppm *PPM) At(x, y int) Pixel{
	return ppm.data[y][x]
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
			ppm.data[y][x] = pixel;
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
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x];
			pixel.R = uint8(float64(pixel.R)*float64(maxValue)/float64(ppm.max));
			pixel.G = uint8(float64(pixel.G)*float64(maxValue)/float64(ppm.max));
			pixel.B = uint8(float64(pixel.B)*float64(maxValue)/float64(ppm.max));
			ppm.data[y][x] = pixel
		}
	}
	ppm.max = maxValue;
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
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			isBlack := (uint8((int(r)+int(g)+int(b))/3) < ppm.max/2)
			pbm.data[y] = append(pbm.data[y], isBlack)
		}
	}
	return pbm
}

func (ppm *PPM) ToPGM() *PGM{
	pgm := &PGM{}
	pgm.magicNumber = "P2";
	pgm.height = ppm.height;
	pgm.width = ppm.width;
	pgm.max = ppm.max;
	for y, _ := range ppm.data {
		pgm.data = append(pgm.data, []uint8{})
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			grayValue := uint8((int(r)+int(g)+int(b))/3)
			pgm.data[y] = append(pgm.data[y], uint8(grayValue))
		}
	}
	return pgm
}

type Point struct{
    X, Y int
}


// Bresenham's Line Drawing Algorithm
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
    deltaX := abs(p2.X - p1.X)
    deltaY := abs(p2.Y - p1.Y)
    sx, sy := sign(p2.X-p1.X), sign(p2.Y-p1.Y)
    err := deltaX - deltaY
    for {
        if (p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height) {
            ppm.data[p1.Y][p1.X] = color
        }
        if p1.X == p2.X && p1.Y == p2.Y {
            break
        }
        e2 := 2 * err
        if e2 > -deltaY {
            err -= deltaY
            p1.X += sx
        }
        if e2 < deltaX {
            err += deltaX
            p1.Y += sy
        }
    }
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func sign(x int) int {
    if x > 0 {
        return 1
    } else if x < 0 {
        return -1
    }
    return 0
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
    p2 := Point{p1.X + width , p1.Y}
    p3 := Point{p1.X, p1.Y + height }
    p4 := Point{p1.X + width , p1.Y + height }
    ppm.DrawLine(p1, p2, color)
    ppm.DrawLine(p2, p4, color)
    ppm.DrawLine(p4, p3, color)
    ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	for i := 0; i <= height; i++ {
		ppm.DrawLine(p1, p2, color);
		p1.Y++;
		p2.Y++;
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			dx := float64(x - center.X)
			dy := float64(y - center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)
			if math.Abs(distance-float64(radius-1)) < 0.5 {
				ppm.data[y][x] = color
			}
		}
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
    for y := -radius; y <= radius; y++ {
        for x := -radius; x <= radius; x++ {
            if x*x+y*y <= radius*radius {
                px := center.X + x
                py := center.Y + y
                if px >= 0 && px < ppm.width && py >= 0 && py < ppm.height {
                    ppm.data[py][px] = color
                }
            }
        }
    }
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
    ppm.DrawLine(p1, p2, color)
    ppm.DrawLine(p2, p3, color)
    ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {

}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
    for i := 0; i < len(points)-1; i++ {
        ppm.DrawLine(points[i], points[i+1], color)
    }
    ppm.DrawLine(points[len(points)-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	
}