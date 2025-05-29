package tasks

import "sync"

type UniqueUsers struct {
	mu    sync.Mutex
	users map[int]struct{}
}

func NewUniqueUsers() *UniqueUsers {
	return &UniqueUsers{users: make(map[int]struct{})}
}

func (u *UniqueUsers) AddUser(id int) {
	u.mu.Lock()
	u.users[id] = struct{}{}
	u.mu.Unlock()
}

func (u *UniqueUsers) Count() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.users)
}