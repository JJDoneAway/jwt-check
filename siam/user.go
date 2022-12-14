package siam

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

/*
In here we'll extract all relevant user information like the name and the SIAM roles
*/

type User struct {
	Uid   string
	Name  string
	Email string
	Roles []string
}

func (jwt Jwt) GetUser() User {
	user := User{}
	user.Uid = jwt.Payload.Uid
	user.Name = jwt.Payload.FullName
	user.Email = jwt.Payload.Email
	user.Roles = jwt.getSIAMRoles()
	return user
}

func (jwt Jwt) getSIAMRoles() []string {
	regex := *regexp.MustCompile(`(cn=)([^,]+)(,.*)`)

	ret := []string{}

	for _, s := range jwt.Payload.GroupMembership {
		fmt.Println(s)
		res := regex.FindAllStringSubmatch(s, -1)
		ret = append(ret, strings.TrimSpace(res[0][2]))
	}

	return ret
}

func (user User) HasRole(role string) bool {
	idx := slices.IndexFunc(user.Roles, func(r string) bool { return strings.EqualFold(r, role) })
	return idx > -1
}
