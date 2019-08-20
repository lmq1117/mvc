package repositories

//todo 检查sync 加锁及解锁是否正确
import (
	"errors"
	"mvc/datamodels"
	"sync"
)

type Query func(user datamodels.User) bool

type UserRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (user datamodels.User, found bool)
	SelectMany(query Query, limit int) (results []datamodels.User)
	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	Delete(query Query, limit int) (deleted bool)
}

type userMemoryRepository struct {
	source map[int64]datamodels.User
	mu     sync.RWMutex
}

const (
	ReadOnlyMode = iota
	ReadWriteMode
)

func NewUserRepository(source map[int64]datamodels.User) UserRepository {
	return &userMemoryRepository{source: source}
}

//Exec
func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.Unlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}
	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break
				}
			}
		}
	}

	return
}

//Select
func (r *userMemoryRepository) Select(query Query) (user datamodels.User, found bool) {
	found = r.Exec(query, func(m datamodels.User) bool {
		user = m
		return true
	}, 1, ReadOnlyMode)

	if !found {
		user = datamodels.User{}
	}
	return
}

//SelectMany
func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.User) {
	r.Exec(query, func(m datamodels.User) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)
	return
}

//InsertOrUpdate
func (r *userMemoryRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	id := user.ID
	if id == 0 { //Insert
		var lastID int64
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()
		id = lastID + 1
		user.ID = id
		r.mu.Lock()
		r.source[id] = user
		r.mu.Unlock()
		return user, nil

	}

	//Update
	current, exists := r.Select(func(m datamodels.User) bool {
		return m.ID == id
	})
	if !exists {
		return datamodels.User{}, errors.New("failed to update a nonexistent user")
	}

	if user.Username != "" {
		current.Username = user.Username
	}
	if user.Firstname != "" {
		current.Firstname = user.Firstname
	}
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()
	return user, nil
}

func (r *userMemoryRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m datamodels.User) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}
