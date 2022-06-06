package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Actor struct {
	name        string
	secret      string
	permissions uint8
}

// Input:  <NAME>:<SECRET>:<PERMISSIONS>,<NAME>:<SECRET>:<PERMISSIONS>,<NAME>:<SECRET>:<PERMISSIONS>,...
func GetTrustedActors() ([]Actor, error) {
	actorsString := os.Getenv("TRUSTED_ACTORS")
	if actorsString == "" {
		return []Actor{}, errors.New("TRUSTED_ACTORS is required")
	}

	actors := strings.Split(actorsString, ",")
	if len(actors) == 0 {
		return []Actor{}, errors.New("no trusted actors defined")
	}

	var trustedActors []Actor
	for idx, actor := range actors {
		splitted := strings.Split(actor, ":")
		if len(splitted) != 3 {
			return []Actor{}, fmt.Errorf("malformed actor provided at index: %d", idx)
		}

		if len(splitted[0]) == 0 || len(splitted[1]) == 0 {
			return []Actor{}, fmt.Errorf("invalid actor metadata index: %d", idx)
		}

		i, err := strconv.ParseUint(splitted[2], 10, 8)
		if err != nil {
			return []Actor{}, fmt.Errorf("invalid actor permissions provided at index: %d", idx)
		}

		trustedActors = append(trustedActors, Actor{
			name:        splitted[0],
			secret:      splitted[1],
			permissions: uint8(i),
		})
	}

	return trustedActors, nil
}
