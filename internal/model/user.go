package model

type UserType string

func (t UserType) String() string { return string(t) }

const (
	UserAdmin       UserType = "admin"
	UserGuest       UserType = "guest"
	UserService     UserType = "service"
	DefaultUserType UserType = UserGuest
)

var (
	BusroutesUserTypes = []UserType{UserAdmin, UserGuest, UserService}
)

type UserTypes []UserType

// Exists returns true if types exist in t.
func (t UserTypes) Exists(types ...UserType) bool {
	allowedTypes := make(map[UserType]struct{}, len(t))
	for _, ut := range t {
		allowedTypes[ut] = struct{}{}
	}

	for _, typ := range types {
		if _, ok := allowedTypes[typ]; ok {
			return true
		}
	}

	return false
}
