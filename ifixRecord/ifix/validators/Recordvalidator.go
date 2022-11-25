package validators

import "ifixRecord/ifix/entities"

// ValidateAddRecordData function is used to validate the input data of ValidateAddRecordData insert operation
func ValidateAddRecordData(recordData entities.RecordEntity) []string {
	var responseError []string

	// clientid blank checking
	if recordData.ClientID == 0 {
		responseError = append(responseError, "input field code required")
	}

	return responseError

}
