package fakes

type Users struct {
	store map[string]User
}

func NewUsers() *Users {
	return &Users{
		store: make(map[string]User),
	}
}

func (u Users) Add(user User) {
	u.store[user.ID] = user
}

func (u Users) Update(user User) {
	u.store[user.ID] = user
}

func (u Users) Get(id string) (User, bool) {
	user, ok := u.store[id]
	return user, ok
}

func (u Users) GetByName(name string) (User, bool) {
	for _, user := range u.store {
		if user.UserName == name {
			return user, true
		}
	}

	return User{}, false
}

func (u Users) Delete(id string) bool {
	_, ok := u.store[id]
	delete(u.store, id)
	return ok
}

func (u *Users) Clear() {
	u.store = make(map[string]User)
}
