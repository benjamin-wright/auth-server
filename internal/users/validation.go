package users

import (
	"github.com/rs/zerolog/log"
)

type Match string

var LOWERCASE = Match("abcdefghijklmnopqrstuvwxyz")
var UPPERCASE = Match("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var NUMBER = Match(`0123456789`)
var SPECIAL = Match(`$&+=?@#|<>^*()%!-`)

func (m *Match) Intersects(value string) bool {
	original := string(*m)

	for _, x := range original {
		for _, y := range value {
			if x == y {
				return true
			}
		}
	}

	return false
}

func CheckPasswordComplexity(password string) bool {
	if len(password) < 8 {
		log.Debug().Msg("password too short")
		return false
	}

	if len(password) > 50 {
		log.Debug().Msg("password too long")
		return false
	}

	if !LOWERCASE.Intersects(password) {
		log.Debug().Msg("password missing lowercase characters")
		return false
	}

	if !UPPERCASE.Intersects(password) {
		log.Debug().Msg("password missing uppercase characters")
		return false
	}

	if !NUMBER.Intersects(password) {
		log.Debug().Msg("password missing numeric characters")
		return false
	}

	if !SPECIAL.Intersects(password) {
		log.Debug().Msg("password missing special characters")
		return false
	}

	return true
}
