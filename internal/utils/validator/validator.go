package validator

import (
	"fmt"
	"strconv"
)

func ValidateCpf(cpf string) error {
	verifierDigit, err  := strconv.Atoi(string(cpf[9]))
	if err != nil {
		fmt.Errorf("digit is not a number %c", cpf[9])
	}
	secondVerifierDigit, err  := strconv.Atoi(string(cpf[10]))
	if err != nil {
		fmt.Errorf("digit is not a number %c", cpf[10])
	}
	err = validateDigit(cpf[:len(cpf)-2], 10, verifierDigit)
	if err != nil {
		return err
	}
	err = validateDigit(cpf[:len(cpf)-1], 11, secondVerifierDigit)
	if err != nil {
		return err
	}
	return nil
}
func validateDigit(cpf string, multiplier int, verifierDigit int) error {
	var verifierDigitCompare int = 0
	sumVerifier := 0
	for _, value := range cpf {
		digit,_ := strconv.Atoi(string(value))
		sumVerifier += digit * multiplier
		multiplier--
	}

	restVerifier := sumVerifier % 11
	if restVerifier  < 2 {
		verifierDigitCompare = 0
	} else {
		verifierDigitCompare = 11 - restVerifier
	}
	digitVerified := verifierDigit == verifierDigitCompare
	if !digitVerified {
		return fmt.Errorf("wrong verifier digit")
	}
	return nil
}
