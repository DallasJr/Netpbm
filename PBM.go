package Netpbm

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

func ReadPBM(filename string) (*PBM, error) {
	//Open the file
	file, err := os.Open(filename)
	//Check for error
	if err != nil {
		return nil, err
	}
	//Close the file just before the Save function returns/(finishes its execution). Even if it's an error
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//Create a base PBM variable
	pbm := &PBM{}
	//Variable line used to count the lines of the image
	line := 0
	//Loop through each lines
	for scanner.Scan() {
		text := scanner.Text()
		//Ignore empty lines and comments
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		if pbm.magicNumber == "" {
			//Get the magicnumber. Trimspace removes the spaces from the string
			pbm.magicNumber = strings.TrimSpace(text)
		} else if pbm.width == 0 {
			//Get the width and height of the pbm
			fmt.Sscanf(text, "%d %d", &pbm.width, &pbm.height)
			//Initialize the pbm.data matrix variable by creating the correct amount and size of arrays in an array
			pbm.data = make([][]bool, pbm.height)
			for i := range pbm.data {
				pbm.data[i] = make([]bool, pbm.width)
			}
		} else {
			if pbm.magicNumber == "P1" {
				//Fields splits the string into several strings and stores them in a string[] using spaces as the delimiter
				test := strings.Fields(text)
				//Loop through the string[]
				for i := 0; i < pbm.width; i++ {
					//If the given string == "1", then it is stored as true or else as false
					pbm.data[line][i] = (test[i] == "1")
				}
				line++
			} else if pbm.magicNumber == "P4" {
				//Calculate the expected number of bytes per row
				expectedBytesPerRow := (pbm.width + 7) / 8
				totalExpectedBytes := expectedBytesPerRow * pbm.height
				allPixelData := make([]byte, totalExpectedBytes)
				//Reads the file content
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				//Extracts the necessary pixel data
				copy(allPixelData, fileContent[len(fileContent)-totalExpectedBytes:])
				//Process the data to fill the pixel array of pbm.data
				byteIndex := 0
				for y := 0; y < pbm.height; y++ {
					for x := 0; x < pbm.width; x++ {
						if x%8 == 0 && x != 0 {
							byteIndex++
						}
						pbm.data[y][x] = (allPixelData[byteIndex]>>(7-(x%8)))&1 != 0
					}
					byteIndex++
				}
				break
			}
		}
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	//Simple return of the size
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	//Simple return of the value of a specifix pixel
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	//Simply define a new value to a specific pixel
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
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
	fmt.Fprint(writer, pbm.magicNumber+"\n")
	//Write the size secondly
	fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
	//Flush writes all the modifications stored in the writer "writer" to the file
	writer.Flush()
	if pbm.magicNumber == "P1" {
		//Loop each pixels of pbm.data
		for y, row := range pbm.data {
			for i, pixel := range row {
				//xtra is used to space each pixel with a space except for the last one of the line
				xtra := " "
				if i == len(row)-1 {
					xtra = ""
				}
				//If pixel is true, it's gonna write 1 or else 0
				if pixel {
					fmt.Fprint(writer, "1"+xtra)
				} else {
					fmt.Fprint(writer, "0"+xtra)
				}
			}
			//Return to line if it's not the last line
			if y != len(pbm.data)-1 {
				fmt.Fprintln(writer, "")
			}
		}
		writer.Flush()
	} else if pbm.magicNumber == "P4" {
		for _, row := range pbm.data {
			//Process each group of 8 pixels in the row
			for x := 0; x < pbm.width; x = x + 8 {
				//Process a row of pixel data
				var byteValue byte
				//Loop up to 8 pixels in the row or until the end of the row is reached
				for i := 0; i < 8 && x+i < pbm.width; i++ {
					bitIndex := 7 - i
					//Update 'byteValue' by setting the corresponding bit if the pixel at the current index in the row is set
					if row[x+i] {
						byteValue |= 1 << bitIndex
					}
				}
				//Write the combined byte value to the file.
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
	//Loop through each pixels
	for y, _ := range pbm.data {
		for x, _ := range pbm.data[y] {
			//Change the value to the opposite of that value
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip by swapping the first and last pixel of each line until the image is flipped.
func (pbm *PBM) Flip() {
	//Loop through each lines
	for y, _ := range pbm.data {
		//Set cursor to the last character of the line
		cursor := pbm.width - 1
		//Loop through each characters of the line
		for x := 0; x < pbm.width; x++ {
			//Store the value of the pixel
			temp := pbm.data[y][x]
			//Change value of the pixel
			pbm.data[y][x] = pbm.data[y][cursor]
			//Set the value of the first pixel to the stored one
			pbm.data[y][cursor] = temp
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
func (pbm *PBM) Flop() {
	//Set the cursor to the bottom line of the image.
	cursor := pbm.height - 1
	//Loop through each lines
	for y, _ := range pbm.data {
		//Swap the current line with the line pointed to by the cursor
		temp := pbm.data[y]
		pbm.data[y] = pbm.data[cursor]
		pbm.data[cursor] = temp
		//Move the cursor to one line higher
		cursor--
		//Break the loop when the cursor crosses or reaches the current line
		if cursor < y || cursor == y {
			break
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	//Simply define a new magic number
	pbm.magicNumber = magicNumber
}
