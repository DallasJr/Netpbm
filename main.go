package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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

	image.DrawTriangle(Point{0, 0}, Point{4, 0}, Point{4, 4}, Pixel{200, 200, 200})
	image.DrawCircle(Point{9, 9}, 4, Pixel{100, 100, 100})
	image.DrawRectangle(Point{6, 4}, 3, 3, Pixel{120, 120, 0})
	image.DrawLine(Point{12, 0}, Point{0, 12}, Pixel{0, 200, 0})

	changed := image.ToPBM();
	changedpgm := image.ToPGM();

	err = changedpgm.Save("output.pgm")
	if err != nil {
		fmt.Println("Error saving the image:", err)
		return
	}
	fmt.Println("Image saved successfully.")

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
			grayValue := pgm.data[y][x]
			isBlack := int(grayValue) < pgm.max/2
			pbm.data[y] = append(pbm.data[y], isBlack)
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
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B
			grayValue := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			isBlack := grayValue < 128
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
			grayValue := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			pgm.data[y] = append(pgm.data[y], uint8(grayValue))
		}
	}
	return pgm
}

type Point struct{
    X, Y int
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
    dx := int(math.Abs(float64(p2.X - p1.X)))
    dy := int(math.Abs(float64(p2.Y - p1.Y)))
    var sx, sy int
    if p1.X < p2.X {
        sx = 1
    } else {
        sx = -1
    }
    if p1.Y < p2.Y {
        sy = 1
    } else {
        sy = -1
    }
    err := dx - dy
    for {
        ppm.data[p1.Y][p1.X] = color
        if p1.X == p2.X && p1.Y == p2.Y {
            break
        }
        e2 := 2 * err
        if e2 > -dy {
            err -= dy
            p1.X += sx
        }
        if e2 < dx {
            err += dx
            p1.Y += sy
        }
    }
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
    p2 := Point{p1.X + width - 1, p1.Y}
    p3 := Point{p1.X, p1.Y + height - 1}
    p4 := Point{p1.X + width - 1, p1.Y + height - 1}
    ppm.DrawLine(p1, p2, color)
    ppm.DrawLine(p2, p4, color)
    ppm.DrawLine(p4, p3, color)
    ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
    for y := p1.Y; y < p1.Y+height; y++ {
        for x := p1.X; x < p1.X+width; x++ {
            ppm.data[y][x] = color
        }
    }
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
    x := radius
    y := 0
    err := 0
    for x >= y {
        ppm.data[center.Y+y][center.X+x] = color
        ppm.data[center.Y+x][center.X+y] = color
        ppm.data[center.Y+x][center.X-y] = color
        ppm.data[center.Y+y][center.X-x] = color
        ppm.data[center.Y-y][center.X-x] = color
        ppm.data[center.Y-x][center.X-y] = color
        ppm.data[center.Y-x][center.X+y] = color
        ppm.data[center.Y-y][center.X+x] = color
        y++
        if err <= 0 {
            err += 2*y + 1
        }
        if err > 0 {
            x--
            err -= 2*x + 1
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
    minY := min(p1.Y, min(p2.Y, p3.Y))
    maxY := max(p1.Y, max(p2.Y, p3.Y))
    for y := minY; y <= maxY; y++ {
        x1, x2 := interpolate(p1, p2, p3, y)
        for x := x1; x <= x2; x++ {
            if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
                ppm.data[y][x] = color
            }
        }
    }
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
    for i := 0; i < len(points)-1; i++ {
        ppm.DrawLine(points[i], points[i+1], color)
    }
    ppm.DrawLine(points[len(points)-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
    n := len(points)
    minY := points[0].Y
    maxY := points[0].Y
    for i := 1; i < n; i++ {
        if points[i].Y < minY {
            minY = points[i].Y
        }
        if points[i].Y > maxY {
            maxY = points[i].Y
        }
    }
    for y := minY; y <= maxY; y++ {
        intersections := []int{}
        for i := 0; i < n; i++ {
            p1 := points[i]
            p2 := points[(i+1)%n]
            if (p1.Y <= y && p2.Y > y) || (p2.Y <= y && p1.Y > y) {
                x := int(float64(p1.X) + float64(y-p1.Y)/float64(p2.Y-p1.Y)*float64(p2.X-p1.X))
                intersections = append(intersections, x)
            }
        }
        sort.Ints(intersections)
        for i := 0; i < len(intersections); i += 2 {
            x1 := intersections[i]
            x2 := intersections[i+1]
            for x := x1; x <= x2; x++ {
                if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
                    ppm.data[y][x] = color
                }
            }
        }
    }
}


func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
    p1 := start
    p2 := Point{start.X + width, start.Y}
    height := int(float64(width) * math.Sqrt(3.0) / 2.0)
    p3 := Point{start.X + width / 2, start.Y + height}
    ppm.drawKochCurve(n, p1, p2, color)
    ppm.drawKochCurve(n, p2, p3, color)
    ppm.drawKochCurve(n, p3, p1, color)
}

func (ppm *PPM) drawKochCurve(n int, p1, p2 Point, color Pixel) {
    if n == 0 {
        ppm.DrawLine(p1, p2, color)
    } else {
        v := Point{(2*p1.X + p2.X) / 3, (2*p1.Y + p2.Y) / 3}
        t := Point{(p1.X + 2*p2.X) / 3, (p1.Y + 2*p2.Y) / 3}
        u := calculateKochPoint(v, t)
        ppm.drawKochCurve(n-1, p1, v, color)
        ppm.drawKochCurve(n-1, v, u, color)
        ppm.drawKochCurve(n-1, u, t, color)
        ppm.drawKochCurve(n-1, t, p2, color)
    }
}

func calculateKochPoint(p1, p2 Point) Point {
    angle := math.Pi / 3.0
    dx := float64(p2.X - p1.X)
    dy := float64(p2.Y - p1.Y)
    newX := p1.X + int((dx*math.Cos(angle) - dy*math.Sin(angle))/3.0)
    newY := p1.Y + int((dx*math.Sin(angle) + dy*math.Cos(angle))/3.0)
    return Point{newX, newY}
}

func interpolate(p1, p2, p3 Point, y int) (int, int) {
    var x1, x2 int
    if p1.Y == p2.Y && p2.Y == p3.Y {
        x1 = min(p1.X, min(p2.X, p3.X))
        x2 = max(p1.X, max(p2.X, p3.X))
        return x1, x2
    }
    if p1.Y == p2.Y {
        x1, x2 = interpolateFlatTopTriangle(p1, p2, p3, y)
    } else if p2.Y == p3.Y {
        x1, x2 = interpolateFlatBottomTriangle(p1, p2, p3, y)
    } else {
        alpha := float64(p2.Y-p3.Y) / float64(p1.Y-p3.Y)
        beta := 1.0 - alpha
        x1a, x2a := interpolateFlatTopTriangle(p1, p2, p3, y)
        x1b, x2b := interpolateFlatBottomTriangle(p1, p2, p3, y)
        x1 = int(alpha*float64(x1a) + beta*float64(x1b))
        x2 = int(alpha*float64(x2a) + beta*float64(x2b))
    }
    return x1, x2
}

func interpolateFlatTopTriangle(p1, p2, p3 Point, y int) (int, int) {
    alpha := float64(y-p1.Y) / float64(p2.Y-p1.Y)
    beta := float64(y-p1.Y) / float64(p3.Y-p1.Y)
    x1 := int(float64(p1.X) + alpha*float64(p2.X-p1.X))
    x2 := int(float64(p1.X) + beta*float64(p3.X-p1.X))
    return x1, x2
}

func interpolateFlatBottomTriangle(p1, p2, p3 Point, y int) (int, int) {
    alpha := float64(y-p1.Y) / float64(p3.Y-p1.Y)
    beta := float64(y-p2.Y) / float64(p3.Y-p2.Y)
    x1 := int(float64(p1.X) + alpha*float64(p3.X-p1.X))
    x2 := int(float64(p2.X) + beta*float64(p3.X-p2.X))
    return x1, x2
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}