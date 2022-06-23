package helpers

func GetStringPointer(str string) *string {
	var strp *string
	strp = nil
	if str != "" {
		strp = &str
	}
	return strp
}
