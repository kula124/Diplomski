package utils

type OperatingMode int

const (
	Encryption OperatingMode = 0
	Decryption OperatingMode = 1
	Unset      OperatingMode = 2
)

type RequiredType int

const (
	Required   RequiredType = 0
	RequiredOr RequiredType = 1
	Optional   RequiredType = 2
)

func (mode OperatingMode) String() string {
	values := [...]string{
		"Encryption",
		"Decryption",
		"Unset",
	}
	if mode < Encryption || mode > Unset {
		return "Unknown" // should throw I TODO
	}
	return values[mode]
}

func (required RequiredType) String() string {
	values := [...]string{
		"Required",
		"RequiredOr",
		"Optional",
	}
	if required < Required || required > Optional {
		return "Unknown" // should throw I TODO
	}
	return values[required]
}
