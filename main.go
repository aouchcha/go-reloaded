package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) != 3 || os.Args[1] != "sample.txt" || os.Args[2] != "result.txt" {
		fmt.Println("You didn't enter the two important argumments in the terminal")
		return
	}

	InputFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error !!")
		return
	}

	file := string(InputFile)

	//Devide the strings into lines and work eith each line
	slice := strings.Split(file, "\n")

	var result string

	for i := range slice {
		line := slice[i]
		if line != "" {
			temp := strings.Fields(line)

			temp = HandlVowels(temp)

			line = ConvertSliceintoString(temp)

			temp = strings.Fields(line)

			if temp[0] == "(hex)" || temp[0] == "(bin)" || temp[0] == "(up)" || temp[0] == "(cap)" || temp[0] == "(low)" {
				fmt.Println("You give the program a commade to make a change in the first element check it")
				return
			} else if temp[0] == "(up," || temp[0] == "(low," || temp[0] == "(cap," {
				fmt.Println("You give the program a commande to make a numerical change in the first element check it")
				return
			} else if temp[len(temp)-1] == "(up," || temp[len(temp)-1] == "(low," || temp[len(temp)-1] == "(cap," {
				fmt.Println("The flage is incorrect !!!")
				return
			}

			temp = DealWithMarkers(temp)

			temp = DeletEmptyCellules(temp)

			temp = AddWhiteSpaces(temp)

			temp = HandlPonctuation(temp)

			line = HandlSingleQuotes(temp)

			line = strings.Trim(line, " ")

			result += line + "\n"
		} else {
			result += line + "\n"
		}
	}
	result = strings.Trim(result, "\n")

	OutputFile := os.Args[2]

	FinalResult := []byte(result)

	err = os.WriteFile(OutputFile, FinalResult, 0644)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func DealWithMarkers(slice []string) []string {
	for i := 1; i < len(slice); i++ {
		if slice[i] == "(bin)" || slice[i] == "(hex)" {
			if slice[i] == "(bin)" {
				num, err := (strconv.ParseInt(slice[i-1], 2, 64))
				if err == nil {
					slice[i-1] = strconv.Itoa(int(num))
					slice[i] = ""
				} else {
					fmt.Println("Isn't binaire !!!")
					slice[i] = ""
					break
				}
			} else if slice[i] == "(hex)" {
				num, err := (strconv.ParseInt(slice[i-1], 16, 64))
				if err == nil {
					slice[i-1] = strconv.Itoa(int(num))
					slice[i] = ""
				} else {
					fmt.Println("Isn't hexadicimal !!!")
					slice[i] = ""
					break
				}
			}
		} else if slice[i] == "(up)" || slice[i] == "(low)" || slice[i] == "(cap)" {
			//Check if the previous string is valid to be modifyed
			if IsItValideWord(slice[i], slice[i-1]) {
				if slice[i] == "(up)" {
					slice[i-1] = strings.ToUpper(slice[i-1])
					slice[i] = ""

				} else if slice[i] == "(low)" {
					slice[i-1] = strings.ToLower(slice[i-1])
					slice[i] = ""

				} else if slice[i] == "(cap)" {
					slice[i-1] = Capitalize(slice[i-1])
					slice[i] = ""

				}
			} else {
				slice[i] = ""
				break
			}
		} else if slice[i] == "(up," || slice[i] == "(cap," || slice[i] == "(low," {
			// Check if it possible to make an numerical change
			if IsItPossible(slice, slice[i], i) {
				num, _ := strconv.Atoi(slice[i+1][:len(slice[i+1])-1])
				if slice[i] == "(up," {
					for j := 0; j < num; j++ {
						if IsItValideWord(slice[i], slice[i-j-1]) {
							slice[i-j-1] = strings.ToUpper(slice[i-j-1])
						}
					}
				} else if slice[i] == "(low," {
					for j := 0; j < num; j++ {
						if IsItValideWord(slice[i], slice[i-j-1]) {
							slice[i-j-1] = strings.ToLower(slice[i-j-1])
						}
					}
				} else if slice[i] == "(cap," {
					for j := 0; j < num; j++ {
						slice[i-j-1] = strings.ToLower(slice[i-j-1])
						if IsItValideWord(slice[i], slice[i-j-1]) {
							slice[i-j-1] = Capitalize(slice[i-j-1])
						}
					}
				}

				slice[i] = ""
				slice[i+1] = ""

			} else {
				slice[i] = ""
				slice[i+1] = ""
				if i+2 == len(slice)-1 {
					slice[i+2] = ""
				}

			}
		}
	}
	return slice
}

func HandlSingleQuotes(slice []string) string {
	var temp []string
	var slirune []rune
	for _, str := range slice {
		for _, char := range str {
			if char == '\'' {
				slirune = append(slirune, ' ', char, ' ')
			} else {
				slirune = append(slirune, char)
			}
		}
	}
	word := string(slirune)
	slice = strings.Fields(word)
	for i := 0; i < len(slice)-1; i++ {
		if slice[i][len(slice[i])-1] == ':' || slice[i] == " " {
			temp = append(temp, slice[i], " ")
		} else if slice[i+1] == "'" || slice[i] == "'" {
			temp = append(temp, slice[i])
		} else {
			temp = append(temp, slice[i], " ")
		}
	}
	temp = append(temp, slice[len(slice)-1])
	word = ConvertSliceintoString(temp)

	// Handling teh case when we have a lot of single quotes .
	count := 0
	for _, char := range word {
		if char == '\'' {
			count++
		}
	}

	sli := []rune(word)
	first := -1
	last := -1
	var test []rune
	if count%2 == 0 {
		for i := 0; i < len(sli)-1; i++ {
			if sli[i] == '\'' {
				if first == -1 {
					first = i
					test = append(test, ' ', sli[i])
					last = -1
				} else if first != -1 && last == -1 {
					last = i
					first = -1
					if sli[i+1] == ',' || sli[i+1] == '.' || sli[i+1] == ';' || sli[i+1] == ':' || sli[i+1] == '!' || sli[i+1] == '?' {
						test = append(test, sli[i])
					} else {
						test = append(test, sli[i], ' ')
					}
				}
			} else {
				test = append(test, sli[i])
			}
		}
		test = append(test, sli[len(sli)-1])
	} else {
		test = sli
	}

	word = ConvertSliceintoString(AddWhiteSpaces(strings.Fields(string(test))))

	return word
}

func ConvertSliceintoString(slice []string) string {
	word := ""
	for _, str := range slice {
		for _, char := range str {
			word = word + string(char)
		}
	}
	return word
}

func AddWhiteSpaces(slice []string) []string {
	var temp []string
	for i := range slice {
		temp = append(temp, slice[i], " ")
	}
	return temp
}

func HandlPonctuation(slice []string) []string {

	for i := 0; i < len(slice)-1; i++ {
		if slice[i] == " " {
			//check is the next word contain only ponctuation
			if IsTheWordPonctuation(slice[i+1]) {
				for j := i; j < len(slice)-1; j++ {
					slice[j] = slice[j+1]
				}
			}
		}
	}
	word := ConvertSliceintoString((slice))

	//Handle the case when the ponctuation is with the word
	temp := []rune(word)
	var ToReturn []rune
	for i := 0; i < len(temp)-1; i++ {
		if (temp[i+1] == ',' || temp[i+1] == '.' || temp[i+1] == ';' || temp[i+1] == ':' || temp[i+1] == '!' || temp[i+1] == '?') && temp[i] == ' ' {
			for j := i; j < len(temp)-1; j++ {
				temp[j] = temp[j+1]
			}
			ToReturn = append(ToReturn, temp[i], ' ')
		} else {
			ToReturn = append(ToReturn, temp[i])
		}
	}
	if temp[len(temp)-1] != ToReturn[len(ToReturn)-1] {
		ToReturn = append(ToReturn, temp[len(temp)-1])
	}
	slice = AddWhiteSpaces(strings.Fields(string(ToReturn)))
	return slice
}

func IsTheWordPonctuation(str string) bool {
	for _, char := range str {
		if char != ',' && char != '.' && char != ';' && char != ':' && char != '!' && char != '?' {
			return false
		}
	}
	return true
}

func HandlVowels(slice []string) []string {

	for i := 0; i < len(slice)-1; i++ {
		if (len(slice[i]) == 1) && (slice[i] == "a" || slice[i] == "A") && (slice[i+1][0] == 'h' || slice[i+1][0] == 'o' || slice[i+1][0] == 'i' || slice[i+1][0] == 'u' || slice[i+1][0] == 'e' || slice[i+1][0] == 'a' || slice[i+1][0] == 'H' || slice[i+1][0] == 'O' || slice[i+1][0] == 'I' || slice[i+1][0] == 'U' || slice[i+1][0] == 'E' || slice[i+1][0] == 'A') {
			fmt.Println(string(slice[i][len(slice[i])-1]), " ", string(slice[i+1][0]), " ", i)
			slice[i] += "n"
		} else if (len(slice[i]) > 1) && (slice[i][len(slice[i])-1] == 'a' || slice[i][len(slice[i])-1] == 'A') && (slice[i+1][0] == 'h' || slice[i+1][0] == 'o' || slice[i+1][0] == 'i' || slice[i+1][0] == 'u' || slice[i+1][0] == 'e' || slice[i+1][0] == 'a' || slice[i+1][0] == 'H' || slice[i+1][0] == 'O' || slice[i+1][0] == 'I' || slice[i+1][0] == 'U' || slice[i+1][0] == 'E' || slice[i+1][0] == 'A') {
			fmt.Println(string(slice[i][len(slice[i])-1]), " ", string(slice[i+1][0]), " ", i)
			for j := range slice[i][:len(slice[i])-1] {
				if slice[i][j] >= 'a' && slice[i][j] <= 'z' || slice[i][j] >= 'A' && slice[i][j] <= 'Z' {
					break
				} else {
					slice[i] += "n"
				}
			}
		}
	}
	return AddWhiteSpaces(slice)
}

func DeletEmptyCellules(slice []string) []string {
	var temp []string
	for i := range slice {
		if slice[i] != "" {
			temp = append(temp, slice[i])
		}
	}
	return temp
}

func IsItValideWord(Marker, str string) bool {
	if str != "" {
		str1 := strings.ToLower(str)
		if Marker == "(up)" || Marker == "(up," || Marker == "(low)" || Marker == "(low," {
			for _, char := range str {
				if char >= '0' && char <= '9' {
					fmt.Println("The word contain an number")
					return false
				}
			}
		} else if Marker == "(cap)" || Marker == "(cap," {
			if str1[0] < 'a' || str1[0] > 'z' {
				fmt.Println("The the first element of the word isn't a letter")
				return false
			}
		}
	} else {
		fmt.Println("there is two flags one after the other")
		return false
	}
	return true
}

func Capitalize(str string) string {
	str1 := strings.ToLower(str)
	temp := []rune(str1)
	if temp[0] >= 'a' && temp[0] <= 'z' {
		temp[0] = temp[0] - 32
	}
	return string(temp)
}

func IsItPossible(slice []string, Marker string, index int) bool {

	if index == len(slice)-1 {
		return false
	} else if slice[index+1][len(slice[index+1])-1] != ')' {
		fmt.Println("The flag is incorrect !!!")
		return false
	}
	ToConvert := slice[index+1][0 : len(slice[index+1])-1]
	num, _ := strconv.Atoi(ToConvert)
	if num <= 0 || num > len(slice)-(len(slice)-index) {
		fmt.Println("The number That u give is greater than what we have or you didn't enter a number !!!")
		return false
	}
	return true
}
