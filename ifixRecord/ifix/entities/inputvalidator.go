package entities

//***************************//
// Package entities
// Date Of Creation: 18/12/2020
// Authour Name: Moloy Mondal
// History: N/A
// Synopsis: This file is used to define commonly used validation functions.
// Functions: LengthChecker,
// LengthChecker() Parameter:  (<string>,<int>)
// LengthRangeChecker() Parameter:  (<string>,<int>,<int>)
// Global Variable: N/A
// Version: 1.0.0
//***************************//

// LengthChecker function is used to satisfy given text length
func LengthChecker(textData string, textLength int) bool {
	return len(textData) == textLength
}

// LengthRangeChecker function is used to satisfy given text length range
func LengthRangeChecker(textData string, textStartLength int, textEndLength int) bool {
	return (len(textData) >= textStartLength && len(textData) <= textEndLength)
}
