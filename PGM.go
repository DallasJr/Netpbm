package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

func ReadPGM(filename string) (*PGM, error) {
	//Open the file
	file, err := os.Open(filename)
	//Check for error
	if err != nil {
		return nil, err
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//Create a base PGM variable
	pgm := &PGM{}
	//Variable line used to count the lines of the image
	line := 0
	//Loop through each lines
	for scanner.Scan() {
		text := scanner.Text()
		//Ignore empty lines and comments
		if strings.HasPrefix(text, "#") {
			continue
		}
		if pgm.magicNumber == "" {
			//Get the magicnumber. Trimspace removes the spaces from the string
			pgm.magicNumber = strings.TrimSpace(text)
		} else if pgm.width == 0 {
			//Get the width and height of the pgm
			fmt.Sscanf(text, "%d %d", &pgm.width, &pgm.height)
			//Initialize the pgm.data matrix variable by creating the correct amount and size of arrays in an array
			pgm.data = make([][]uint8, pgm.height)
			for i := range pgm.data {
				pgm.data[i] = make([]uint8, pgm.width)
			}
		} else if pgm.max == 0 {
			//Get the maxValue of the pgm
			fmt.Sscanf(text, "%d", &pgm.max)
		} else {
			if pgm.magicNumber == "P2" {
				//Fields splits the string into several strings and stores them in a string[] using spaces as the delimiter
				val := strings.Fields(text)
				//Loop through the string[]
				for i := 0; i < pgm.width; i++ {
					//Converts the string to a uint8
					num, _ := strconv.ParseUint(val[i], 10, 8)
					pgm.data[line][i] = uint8(num)
				}
				line++
			} else if pgm.magicNumber == "P5" {
				//Create an array of uint8 of the size of the image
				pixelData := make([]uint8, pgm.width*pgm.height)
				//Reads the file content
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				//Extracts the necessary pixel data
				copy(pixelData, fileContent[len(fileContent)-(pgm.width*pgm.height):])
				//Process the data to fill the pixel array of pgm.data
				pixelIndex := 0
				for y := 0; y < pgm.height; y++ {
					for x := 0; x < pgm.width; x++ {
						pgm.data[y][x] = pixelData[pixelIndex]
						pixelIndex++
					}
				}
				break
			}
		}
	}
	return pgm, nil
}

func (pgm *PGM) Size() (int, int) {
	//Simple return of the size
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	//Simple return of the value of a specifix pixel
	return pgm.data[x][y]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	//Simply define a new value to a specific pixel
	pgm.data[x][y] = value
}

func (pgm *PGM) Save(filename string) error {
	//Create a file with the defines name
	file, err := os.Create(filename)
	//Check for error
	if err != nil {
		return err
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close()
	//Store all the modifications into writer "writer" temporarily until flush
	writer := bufio.NewWriter(file)
	//Write the magicnumber first
	fmt.Fprint(writer, pgm.magicNumber+"\n")
	//Write the size secondly
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)
	//Flush writes all the modifications stored in the writer "writer" to the file
	writer.Flush()
	if pgm.magicNumber == "P2" {
		//Loop each pixels of pgm.data
		for y, row := range pgm.data {
			for i, pixel := range row {
				//xtra is used to space each pixel with a space except for the last one of the line
				xtra := " "
				if i == len(row)-1 {
					xtra = ""
				}
				//Here i convert uint8 to an int in order to finally convert it to a string
				fmt.Fprint(writer, strconv.Itoa(int(pixel))+xtra)
			}
			//Return to line if it's not the last line
			if y != len(pgm.data)-1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush()
	} else if pgm.magicNumber == "P5" {
		//For each pixels
		for _, row := range pgm.data {
			for _, pixel := range row {
				//We can simply convert it to []byte
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
	//Loop throught each pixels
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x]
			//Change the value to the opposite of his value
			//If maxValue is 10, so the opposite of 7 would be 3
			//10 - 7 = 3
			pgm.data[y][x] = pgm.max - prevvalue
		}
	}
}

// Flip by swapping the first and last pixel of each line until the image is flipped.
func (pgm *PGM) Flip() {
	//Loop through each lines
	for y, _ := range pgm.data {
		//Set cursor to the last character of the line
		cursor := pgm.width - 1
		//Loop through each characters of the line
		for x := 0; x < pgm.width; x++ {
			//Store the value of the pixel
			temp := pgm.data[y][x]
			//Change value of the pixel
			pgm.data[y][x] = pgm.data[y][cursor]
			//Set the value of the first pixel to the stored one
			pgm.data[y][cursor] = temp
			//Move the cursor to the left on the line
			cursor--
			//Break the loop when the cursor crosses or reaches the current line
			if cursor < x || cursor == x {
				break
			}
		}
	}
}

// Flop by swapping the first and last line until the image is flopped.
func (pgm *PGM) Flop() {
	//Set the cursor to the bottom line of the image.
	cursor := pgm.height - 1
	//Loop through each lines
	for y, _ := range pgm.data {
		//Swap the current line with the line pointed to by the cursor
		temp := pgm.data[y]
		pgm.data[y] = pgm.data[cursor]
		pgm.data[cursor] = temp
		//Move the cursor to one line higher
		cursor--
		//Break the loop when the cursor crosses or reaches the current line
		if cursor < y || cursor == y {
			break
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	//Simply define a new magic number
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	//Loop through each pixel
	for y, _ := range pgm.data {
		for x, _ := range pgm.data[y] {
			prevvalue := pgm.data[y][x]
			//Calculate the new pixel value based on the new maximum value
			//Adjusting the pixel value proportionally to the new max value
			newvalue := prevvalue * uint8(5) / pgm.max
			pgm.data[y][x] = newvalue
		}
	}
	pgm.max = maxValue
}

func (pgm *PGM) Rotate90CW() {
	//Create a new matrix to store the rotated pixel data
	rotatedData := make([][]uint8, pgm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.height)
	}
	//Loop through each pixel in the original image
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < pgm.height; j++ {
			//Rotate the pixel by 90 degrees clockwise and assign it
			rotatedData[i][j] = pgm.data[pgm.height-1-j][i]
		}
	}
	//Swap the width and height of the image.
	pgm.width, pgm.height = pgm.height, pgm.width
	//Update the image data with the rotated data
	pgm.data = rotatedData
}

func (pgm *PGM) ToPBM() *PBM {
	//Create a new pbm
	pbm := &PBM{}
	//Assign same data except for the magicnumber
	pbm.magicNumber = "P1"
	pbm.height = pgm.height
	pbm.width = pgm.width
	for y, _ := range pgm.data {
		pbm.data = append(pbm.data, []bool{})
		for x, _ := range pgm.data[y] {
			grayValue := pgm.data[y][x]
			//Calculate if the pixel should be black or white
			//if the grayValue is lower than the half of the maxValue, then i consider it white
			//If maxValue is 100, 49 would be white
			isBlack := grayValue < pgm.max/2
			pbm.data[y] = append(pbm.data[y], isBlack)
		}
	}
	return pbm
}