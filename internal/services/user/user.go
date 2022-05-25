package user

import (
	"math"
	"wiseman/internal/db"
	"wiseman/internal/entities"
)


var users entities.UsersType

func GetNextLevelMinExperience(u *entities.UserType) uint {
	fLevel := float64(u.CurrentLevel + 1)

	return uint(50 * (math.Pow(fLevel, 3) - 6*math.Pow(fLevel, 2) + 17*fLevel - 12) / 3)
}

func IncreaseExperience(u *entities.UserType, v uint, guildID string) uint {
	// Get original object using ComplexID to avoid injecting other mutated data
	serverMultiplier := db.GetServerMultiplierByGuildId(guildID)

	for {
		if u.CurrentLevelExperience+v < GetNextLevelMinExperience(u) {
			u.CurrentLevelExperience += v * uint(serverMultiplier)
			break
		}

		v -= GetNextLevelMinExperience(u) - u.CurrentLevelExperience
		u.CurrentLevelExperience = 0
		u.CurrentLevel += 1
	}

	// TODO: UpsertUserByID should load from the cache inside the database all data
	// every n minutes (maybe with a cron or a goroutine?)
	db.UpsertUserByID(u.ComplexID, u)

	return u.CurrentLevelExperience
}
