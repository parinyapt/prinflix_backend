package utilsUUID

import "github.com/google/uuid"

func GenerateUUIDv4() uuid.UUID {
	return uuid.New()
}

func GenerateUUIDv4String() string {
	return GenerateUUIDv4().String()
}

func ParseUUIDfromString(uuidStr string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedUUID, nil 
}