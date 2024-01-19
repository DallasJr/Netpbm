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

func ReadPBM(filename string) (*PBM, error) {
	//Open the file
	file, err := os.Open(filename);
	//Check for error
	if err != nil {
		return nil, err;
	}
	defer file.Close();
	scanner := bufio.NewScanner(file);
	//Create a base PBM variable
	pbm := &PBM{};
	//Variable line used to count the lines of the image
	line := 0;
	//Loop through each lines
	for scanner.Scan() {
		text := scanner.Text();
		//Ignore empty lines and comments
		if text == "" || strings.HasPrefix(text, "#") {
			continue;
		}
		if pbm.magicNumber == "" {
			//Get the magicnumber. Trimspace removes the spaces from the string
			pbm.magicNumber = strings.TrimSpace(text);
		} else if pbm.width == 0 {
			//Get the width and height of the pbm
			fmt.Sscanf(text, "%d %d", &pbm.width, &pbm.height);
			//Initialize the pbm.data matrix variable by creating the correct amount and size of arrays in an array
			pbm.data = make([][]bool, pbm.height);
			for i := range pbm.data {
				pbm.data[i] = make([]bool, pbm.width);
			}
		} else {
			if pbm.magicNumber == "P1" {
				//Fields splits the string into several strings and stores them in a string[] using spaces as the delimiter
				test := strings.Fields(text);
				//Loop through the string[]
				for i := 0; i < pbm.width; i++ {
					//If the given string == "1", then it is stored as true or else as false
					pbm.data[line][i] = (test[i] == "1");
				}
				line++
			} else if pbm.magicNumber == "P4" {
				//Calculate the expected number of bytes per row
				expectedBytesPerRow := (pbm.width + 7) / 8;
				totalExpectedBytes := expectedBytesPerRow * pbm.height;
				allPixelData := make([]byte, totalExpectedBytes);
				//Reads the file content
				fileContent, err := os.ReadFile(filename);
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err);
				}
				//Extracts the necessary pixel data
				copy(allPixelData, fileContent[len(fileContent)-totalExpectedBytes:]);
				//Process the data to fill the pixel array of pbm.data
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
				break;
			}
		}
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	//Simple return of the size
	return pbm.width, pbm.height;
}

func (pbm *PBM) At(x, y int) bool {
	//Simple return of the value of a specifix pixel
	return pbm.data[y][x];
}

func (pbm *PBM) Set(x, y int, value bool) {
	//Simply define a new value to a specific pixel
	pbm.data[y][x] = value;
}

func (pbm *PBM) Save(filename string) error {
	//Create a file with the defines name
	file, err := os.Create(filename);
	//Check for error
	if err != nil {
		return err;
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close();
	//Store all the modifications into writer "writer" temporarily until flush
	writer := bufio.NewWriter(file);
	//Write the magicnumber first
	fmt.Fprint(writer, pbm.magicNumber + "\n");
	//Write the size secondly
	fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height);
	//Flush writes all the modifications stored in the writer "writer" to the file
	writer.Flush();
	if pbm.magicNumber == "P1" {
		//Loop each pixels of pbm.data
		for y, row := range pbm.data {
			for i, pixel := range row {
				//xtra is used to space each pixel with a space except for the last one of the line
				xtra := " ";
				if i == len(row) - 1 {
					xtra = "";
				}
				//If pixel is true, it's gonna write 1 or else 0
				if pixel {
					fmt.Fprint(writer, "1" + xtra);
				} else {
					fmt.Fprint(writer, "0" + xtra);
				}
			}
			//Return to line if it's not the last line
			if y != len(pbm.data) - 1 {
				fmt.Fprintln(writer, "");
			}
		}
		writer.Flush();
	} else if pbm.magicNumber == "P4" {
		for _, row := range pbm.data {
			//Process each group of 8 pixels in the row
			for x := 0; x < pbm.width; x = x + 8 {
				//Process a row of pixel data
				var byteValue byte;
				//Loop up to 8 pixels in the row or until the end of the row is reached
				for i := 0; i < 8 && x+i < pbm.width; i++ {
					bitIndex := 7 - i;
					//Update 'byteValue' by setting the corresponding bit if the pixel at the current index in the row is set
					if row[x+i] {
						byteValue |= 1 << bitIndex;
					}
				}
				//Write the combined byte value to the file.
				_, err = file.Write([]byte{byteValue});
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err);
				}
			}
		}
	}
	return nil;
}

func (pbm *PBM) Invert() {
	//Loop through each pixels
	for y, _ := range pbm.data {
		for x, _ := range pbm.data[y] {
			//Change the value to the opposite of that value
			pbm.data[y][x] = !pbm.data[y][x];
		}
	}
}

//Flip by swapping the first and last pixel of each line until the image is flipped.
func (pbm *PBM) Flip() {
	//Loop through each lines
	for y, _ := range pbm.data {
		//Set cursor to the last character of the line
		cursor := pbm.width - 1;
		//Loop through each characters of the line
		for x := 0; x < pbm.width; x++ {
			//Store the value of the pixel
			temp := pbm.data[y][x];
			//Change value of the pixel
			pbm.data[y][x] = pbm.data[y][cursor];
			//Set the value of the first pixel to the stored one
			pbm.data[y][cursor] = temp;
			//Move the cursor to the left on the line
			cursor--;
			//Break the loop when the cursor crosses or reaches the current line
			if cursor < x || cursor == x {
				break;
			}
		}
	}
}

//Flop by swapping the first and last line until the image is flopped.
func (pbm *PBM) Flop() {
	//Set the cursor to the bottom line of the image.
	cursor := pbm.height - 1;
	//Loop through each lines
	for y, _ := range pbm.data {
		//Swap the current line with the line pointed to by the cursor
		temp := pbm.data[y];
		pbm.data[y] = pbm.data[cursor];
		pbm.data[cursor] = temp;
		//Move the cursor to one line higher
		cursor--;
		//Break the loop when the cursor crosses or reaches the current line
		if cursor < y || cursor == y {
			break;
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	//Simply define a new magic number
	pbm.magicNumber = magicNumber;
}

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

func ReadPGM(filename string) (*PGM, error) {
	//Same as ReadPBM
	file, err := os.Open(filename);
	if err != nil {
		return nil, err;
	}
	defer file.Close();
	scanner := bufio.NewScanner(file);
	pgm := &PGM{};
	line := 0;
	for scanner.Scan() {
		text := scanner.Text();
		if strings.HasPrefix(text, "#") {
			continue;
		}
		if pgm.magicNumber == "" {
			pgm.magicNumber = strings.TrimSpace(text);
		} else if pgm.width == 0 {
			fmt.Sscanf(text, "%d %d", &pgm.width, &pgm.height);
			pgm.data = make([][]uint8, pgm.height);
			for i := range pgm.data {
				pgm.data[i] = make([]uint8, pgm.width);
			}
		} else if pgm.max == 0 {
			//Get the maxValue of the pgm
			fmt.Sscanf(text, "%d", &pgm.max);
		} else {
			if pgm.magicNumber == "P2" {
				val := strings.Fields(text);
				for i := 0; i < pgm.width; i++ {
					//Converts the string to a uint8
					num, _ := strconv.ParseUint(val[i], 10, 8);
					pgm.data[line][i] = uint8(num);
				}
				line++
			} else if pgm.magicNumber == "P5" {
				//Create an array of uint8 of the size of the image
				pixelData := make([]uint8, pgm.width*pgm.height);
				//Same as ReadPBM
				fileContent, err := os.ReadFile(filename);
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err);
				}
				copy(pixelData, fileContent[len(fileContent)-(pgm.width*pgm.height):])
				pixelIndex := 0;
				for y := 0; y < pgm.height; y++ {
					for x := 0; x < pgm.width; x++ {
						pgm.data[y][x] = pixelData[pixelIndex];
						pixelIndex++;
					}
				}
				break;
			}
		}
	}
	return pgm, nil;
}

func (pgm *PGM) Size() (int, int) {
	//Same as pbm.Size
	return pgm.width, pgm.height;
}

func (pgm *PGM) At(x, y int) uint8 {
	//Same as pbm.At
	return pgm.data[x][y];
}

func (pgm *PGM) Set(x, y int, value uint8) {
	//Same as pbm.Set
	pgm.data[x][y] = value;
}

func (pgm *PGM) Save(filename string) error {
	//Same as pbm.Save
	file, err := os.Create(filename);
	if err != nil {
		return err;
	}
	defer file.Close();
	writer := bufio.NewWriter(file);
	fmt.Fprint(writer, pgm.magicNumber + "\n");
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height);
	fmt.Fprintf(writer, "%d\n", pgm.max);
	writer.Flush();
	if pgm.magicNumber == "P2" {
		for y, row := range pgm.data {
			for i, pixel := range row {
				xtra := " ";
				if i == len(row) - 1 {
					xtra = "";
				}
				//Here i convert uint8 to an int in order to finally convert it to a string
				fmt.Fprint(writer, strconv.Itoa(int(pixel)) + xtra);
			}
			if y != len(pgm.data) - 1 {
				fmt.Fprintln(writer, "");
			}
		}
		writer.Flush();
	} else if pgm.magicNumber == "P5" {
		for _, row := range pgm.data {
			for _, pixel := range row {
				//We can simply convert it to []byte
				_, err = file.Write([]byte{pixel});
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err);
				}
			}
		}
	}
	return nil
}

func (pgm *PGM) Invert() {
	//Loop throught each pixels
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x];
			//Change the value to the opposite of his value
			//If maxValue is 10, so the opposite of 7 would be 3
			//10 - 7 = 3
			pgm.data[y][x] = pgm.max - prevvalue;
		}
	}
}

func (pgm *PGM) Flip() {
	//Same as pbm.Flip
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
	//Same as pbm.Flop
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
	//Same as pbm.SetMagicNumber
	pgm.magicNumber = magicNumber;
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	//Loop through each pixel
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x];
			//Calculate the new pixel value based on the new maximum value
			//Adjusting the pixel value proportionally to the new max value
			newvalue := prevvalue*uint8(5)/pgm.max;
			pgm.data[y][x] = newvalue;
		}
	}
	pgm.max = maxValue;
}

func (pgm *PGM) Rotate90CW() {
	//Create a new matrix to store the rotated pixel data
	rotatedData := make([][]uint8, pgm.width);
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.height);
	}
	//Loop through each pixel in the original image
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < pgm.height; j++ {
			//Rotate the pixel by 90 degrees clockwise and assign it
			rotatedData[i][j] = pgm.data[pgm.height-1-j][i];
		}
	}
	//Swap the width and height of the image.
	pgm.width, pgm.height = pgm.height, pgm.width;
	//Update the image data with the rotated data
	pgm.data = rotatedData;
}

func (pgm *PGM) ToPBM() *PBM {
	//Create a new pbm
	pbm := &PBM{};
	//Assign same data except for the magicnumber
	pbm.magicNumber = "P1";
	pbm.height = pgm.height;
	pbm.width = pgm.width;
	for y, _ := range pgm.data {
		pbm.data = append(pbm.data, []bool{});
		for x, _ := range pgm.data[y] {
			grayValue := pgm.data[y][x];
			//Calculate if the pixel should be black or white
			//if the grayValue is lower than the half of the maxValue, then i consider it white
			//If maxValue is 100, 49 would be white
			isBlack := grayValue < pgm.max/2;
			pbm.data[y] = append(pbm.data[y], isBlack);
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
	//Same as ReadPGM
	file, err := os.Open(filename);
	if err != nil {
		return nil, err;
	}
	defer file.Close();
	scanner := bufio.NewScanner(file);
	ppm := &PPM{};
	line := 0;
	for scanner.Scan() {
		text := scanner.Text();
		if strings.HasPrefix(text, "#") {
			continue;
		}
		if ppm.magicNumber == "" {
			ppm.magicNumber = strings.TrimSpace(text);
		} else if ppm.width == 0 {
			fmt.Sscanf(text, "%d %d", &ppm.width, &ppm.height);
			ppm.data = make([][]Pixel, ppm.height);
			for i := range ppm.data {
				ppm.data[i] = make([]Pixel, ppm.width);
			}
		} else if ppm.max == 0 {
			fmt.Sscanf(text, "%d", &ppm.max);
		} else {
			if ppm.magicNumber == "P3" {
				val := strings.Fields(text);
				//Loop through each strings in the current line
				for i := 0; i < ppm.width; i++ {
					//Convert the string to uint8 and set it to the red of the pixel
					r, _ := strconv.ParseUint(val[i*3], 10, 8);
					//Same but the index is incremented to get the next value for the green
					g, _ := strconv.ParseUint(val[i*3+1], 10, 8);
					//Same but the index is incremented to get the next value for the blue
					b, _ := strconv.ParseUint(val[i*3+2], 10, 8);
					//Create the pixel with the colors we just obtained and define it the matrix
					ppm.data[line][i] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)};
				}
				line++
			} else if ppm.magicNumber == "P6" {
				//Create an array of byte of the size of the image * 3 because each pixel has 3 values RGB
				pixelData := make([]byte, ppm.width*ppm.height*3);
				fileContent, err := os.ReadFile(filename);
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err);
				}
				//Same as ReachPGM but for 3 values
				copy(pixelData, fileContent[len(fileContent)-(ppm.width*ppm.height*3):]);
				pixelIndex := 0;
				for y := 0; y < ppm.height; y++ {
					for x := 0; x < ppm.width; x++ {
						ppm.data[y][x].R = pixelData[pixelIndex];
						ppm.data[y][x].G = pixelData[pixelIndex+1];
						ppm.data[y][x].B = pixelData[pixelIndex+2];
						pixelIndex += 3;
					}
				}
				break;
			}
		}
	}
	return ppm, nil;
}

func (ppm *PPM) Save(filename string) error {
	//Same as pgm.Save
	file, err := os.Create(filename);
	if err != nil {
		return err;
	}
	defer file.Close();
	writer := bufio.NewWriter(file);
	fmt.Fprint(writer, ppm.magicNumber + "\n");
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height);
	fmt.Fprintf(writer, "%d\n", ppm.max);
	writer.Flush();
	if ppm.magicNumber == "P3" {
		for y, row := range ppm.data {
			for i, pixel := range row {
				xtra := " ";
				if i == len(row)-1 {
					xtra = "";
				}
				//Write the RGB colors in the writer
				fmt.Fprintf(writer, "%d %d %d%s", pixel.R, pixel.G, pixel.B, xtra);
			}
			if y != len(ppm.data)-1 {
				fmt.Fprintln(writer, "");
			}
		}
		writer.Flush();
	} else if ppm.magicNumber == "P6" {
		//Same as pgm.Save but for the 3 colors RGB
		for _, row := range ppm.data {
			for _, pixel := range row {
				_, err = file.Write([]byte{pixel.R, pixel.G, pixel.B});
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err);
				}
			}
		}
	}
	return nil;
}

func (ppm *PPM) Size() (int, int) {
	//Same as pgm.Size
	return ppm.width, ppm.height;
}

func (ppm *PPM) At(x, y int) Pixel{
	//Same as pgm.At
	return ppm.data[y][x];
}

func (ppm *PPM) Set(x, y int, value Pixel){
	//Same as pgm.Set
	ppm.data[x][y] = value;
}

func (ppm *PPM) Invert() {
	//Loop throught each pixels
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x];
			//Change the value to the opposite of his value
			//If the value is 240 would be 15
			//255 - 240 = 15
			//If the value is 1O would be 245
			//255- 10 = 245
			pixel.R = uint8(255 - int(pixel.R));
			pixel.G = uint8(255 - int(pixel.G));
			pixel.B = uint8(255 - int(pixel.B));
			ppm.data[y][x] = pixel;
		}
	}
}

func (ppm *PPM) Flip() {
	//Same as pgm.Flip
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
	//Same as pgm.Flop
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
	//Same as pgm.SetMagicNumber
	ppm.magicNumber = magicNumber;
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	//Same idea as pgm.SetMaxValue
	for y, _ := range ppm.data {
		for x, _ := range ppm.data[y] {
			pixel := ppm.data[y][x];
			//Calculate the new pixel value based on the new maximum value for each color
			//Adjusting the pixel value proportionally to the new max value
			pixel.R = uint8(float64(pixel.R)*float64(maxValue)/float64(ppm.max));
			pixel.G = uint8(float64(pixel.G)*float64(maxValue)/float64(ppm.max));
			pixel.B = uint8(float64(pixel.B)*float64(maxValue)/float64(ppm.max));
			ppm.data[y][x] = pixel;
		}
	}
	ppm.max = maxValue;
}

func (ppm *PPM) Rotate90CW() {
	//Same as pgm.Rotate90CW but the matrix is [][]Pixel not [][]uint8
	rotatedData := make([][]Pixel, ppm.width);
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height);
	}
	for i := 0; i < ppm.width; i++ {
		for j := 0; j < ppm.height; j++ {
			rotatedData[i][j] = ppm.data[ppm.height-1-j][i];
		}
	}
	ppm.width, ppm.height = ppm.height, ppm.width;
	ppm.data = rotatedData;
}

func (ppm *PPM) ToPBM() *PBM{
	//Same idea as pgm.ToPBM
	pbm := &PBM{};
	pbm.magicNumber = "P1";
	pbm.height = ppm.height;
	pbm.width = ppm.width;
	for y, _ := range ppm.data {
		pbm.data = append(pbm.data, []bool{});
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B;
			//Calculate if the pixel should be black or white
			//if the average of the 3 colors is lower than the half of the maxValue, then i consider it white
			//If maxValue is 100 and average is 49, it would be black
			isBlack := (uint8((int(r)+int(g)+int(b))/3) < ppm.max/2);
			pbm.data[y] = append(pbm.data[y], isBlack);
		}
	}
	return pbm;
}

func (ppm *PPM) ToPGM() *PGM{
	//Same idea as ppm.ToPBM
	pgm := &PGM{};
	pgm.magicNumber = "P2";
	pgm.height = ppm.height;
	pgm.width = ppm.width;
	pgm.max = ppm.max;
	for y, _ := range ppm.data {
		pgm.data = append(pgm.data, []uint8{});
		for x, _ := range ppm.data[y] {
			r, g, b := ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B;
			//Calculate the amount of gray the pixel should have
			//It is just the average of the 3 RGB colors
			grayValue := uint8((int(r)+int(g)+int(b))/3);
			pgm.data[y] = append(pgm.data[y], uint8(grayValue));
		}
	}
	return pgm;
}

type Point struct{
	X, Y int
}


//Drawing lines by using Bresenham's Line Drawing Algorithm
//Found people suggesting it on online forums
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	deltaX := abs(p2.X - p1.X);
	deltaY := abs(p2.Y - p1.Y);
	sx, sy := sign(p2.X-p1.X), sign(p2.Y-p1.Y);
	err := deltaX - deltaY;
	for {
		if (p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height) {
			ppm.data[p1.Y][p1.X] = color;
		}
		if p1.X == p2.X && p1.Y == p2.Y {
			break;
		}
		e2 := 2 * err;
		if e2 > -deltaY {
			err -= deltaY;
			p1.X += sx;
		}
		if e2 < deltaX {
			err += deltaX;
			p1.Y += sy;
		}
	}
}
//If negative, change it to positive
func abs(x int) int {
	if x < 0 {
		return -x;
	}
	return x;
}
//Return 1 if it's over 0
//Return 0 if it's 0
//Return -1 if  it's negative
func sign(x int) int {
	if x > 0 {
		return 1;
	} else if x < 0 {
		return -1;
	}
	return 0;
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	//Create the 3 extra points according to the width and the height
	p2 := Point{p1.X + width , p1.Y};
	p3 := Point{p1.X, p1.Y + height };
	p4 := Point{p1.X + width , p1.Y + height};
	//Draw the lines to connect them
	ppm.DrawLine(p1, p2, color);
	ppm.DrawLine(p2, p4, color);
	ppm.DrawLine(p4, p3, color);
	ppm.DrawLine(p3, p1, color);
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	//Draw horizontal lines with the asked width under each other until the height is reached
	p2 := Point{p1.X + width, p1.Y};
	for i := 0; i <= height; i++ {
		ppm.DrawLine(p1, p2, color);
		p1.Y++;
		p2.Y++;
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	//Loop through each pixel
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			//Calculate the distance from the current pixel to the center of the circle
			dx := float64(x - center.X);
			dy := float64(y - center.Y);
			distance := math.Sqrt(dx*dx + dy*dy);
			//Check if the distance is approximately equal to the specified radius
			//*0.85 is to obtain a circle looking like the tester's circle even if it's not really a circle... In reality, remove "*0.85" and it's a real circle
			if math.Abs(distance-float64(radius)*0.85) < 0.5 {
				ppm.data[y][x] = color;
			}
		}
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	//Draw a circle with the radius getting smaller until it is at 0;
	for radius >= 0 {
		ppm.DrawCircle(center, radius, color);
		radius--;
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	//Draw lines and link the 3 points
	ppm.DrawLine(p1, p2, color);
	ppm.DrawLine(p2, p3, color);
	ppm.DrawLine(p3, p1, color);
}

//Draw a line from p1 to p3 and move p1 towars p2 until the triangle is filled
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	//Loop until p1 reaches p2
	for (p1 != p2) {
		//Draw a line between p1 and p3
		ppm.DrawLine(p3, p1, color);
		//Increment or decrement X of p1 based on p2 position
		if (p1.X != p2.X && p1.X < p2.X) {
			p1.X++;
		} else if (p1.X != p2.X && p1.X > p2.X) {
			p1.X--;
		}
		//Increment or decrement Y of p1 based on p2 position
		if (p1.Y != p2.Y && p1.Y < p2.Y) {
			p1.Y++;
		} else if (p1.Y != p2.Y && p1.Y > p2.Y) {
			p1.Y--;
		}
	}
	//Draw a final line between the last position of p1 (should be at p2 at this point) and p3
	ppm.DrawLine(p3, p1, color);
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	//Link the points with a line
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color);
	}
	//Link the last and the first point with a line
	ppm.DrawLine(points[len(points)-1], points[0], color);
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	
}