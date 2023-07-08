package constants

type Role uint8

const (
	NONE      Role = 0
	IITK_USER Role = 1
	MOD       Role = 100
	ADMIN     Role = 101
)
