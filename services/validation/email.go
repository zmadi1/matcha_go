package validation

import (
	"regexp"
	//"strings"
)

//ValidEmail will help us quickly evaluate whether the email that is sent to
//the server is a valid email address.
func ValidEmail(email string) bool {
	/*if strings.Count(email, "@") > 1 || strings.Count(email, "@") == 0 {
		return false
	}
	subemail := strings.Split(email, "@")
	if len(subemail) > 2 || len(subemail) == 1 {
		return false
	}
	re := regexp.MustCompile(`(^[a-z]\b[a-z._-]+)$`)
	if re.MatchString(subemail[0]) == false {
		return false
	}
	re = regexp.MustCompile(`(^[a-z.]+[a-z]+)$`)
	if re.MatchString(subemail[1]) == false {
		return false
	}*/
        re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}
