package models

import (
	"strings"
)

func getOrientation(orientation string) int {
	if strings.Compare(orientation, "heterosexual") == 0 {
		return 1
	} else if strings.Compare(orientation, "homosexual") == 0 {
		return 2
	} else if strings.Compare(orientation, "bisexual") == 0 {
		return 3
	}
	return 0
}

func evalMatchSex(sex, oriU, oriM int) bool {
	if sex == 1 {
		if (oriU == 1 || oriU == 3) && (oriM == 1 || oriM == 3) {
			return true
		}
		return false
	} else {
		if (oriU == 2 || oriU == 3) && (oriM == 2 || oriM == 3) {
			return true
		}
		return false
	}
}

func filterByOrientation(user *User, match User) bool {
	var sex int
	if strings.Compare(user.Sex, match.Sex) == 0 {
		sex = 0
	} else {
		sex = 1
	}
	uOrientation := getOrientation(user.Profile.Orientation)
	mOrientation := getOrientation(match.Profile.Orientation)
	return evalMatchSex(sex, uOrientation, mOrientation)
}

func filterUsers(user *User, users []User) []User {
	itr := 0
	var newUsers []User
	for itr < len(users) {
		if !filterByOrientation(user, users[itr]) {
			itr++
			continue
		}
		newUsers = append(newUsers, users[itr])
		itr++
	}
	return newUsers
}
