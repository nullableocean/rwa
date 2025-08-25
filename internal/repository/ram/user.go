package ram

import (
	"errors"
	"rwa/internal/models"
	"sync"
	"sync/atomic"
)

type UserRepository struct {
	idCounter    atomic.Int64
	store        map[int64]models.User
	usernamesMap map[string]int64
	emailsMap    map[string]int64

	mu *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		idCounter:    atomic.Int64{},
		store:        make(map[int64]models.User),
		usernamesMap: make(map[string]int64),
		emailsMap:    make(map[string]int64),
		mu:           &sync.RWMutex{},
	}
}

func (r *UserRepository) GetById(id int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exist := r.store[id]
	if !exist {
		return nil, models.ErrNotFound
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exist := r.usernamesMap[username]
	if !exist {
		return nil, models.ErrNotFound
	}

	user := r.store[id]

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exist := r.emailsMap[email]
	if !exist {
		return nil, models.ErrNotFound
	}

	user := r.store[id]

	return &user, nil
}

func (r *UserRepository) Save(user models.User) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exist := r.usernamesMap[user.Username]
	if exist {
		return 0, errors.New("username already taken")
	}
	_, exist = r.emailsMap[user.Email]
	if exist {
		return 0, errors.New("username already taken")
	}

	user.ID = r.newUserID()

	r.store[user.ID] = user
	r.usernamesMap[user.Username] = user.ID
	r.emailsMap[user.Email] = user.ID

	return user.ID, nil
}

func (r *UserRepository) Update(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	currentUserState, exist := r.store[user.ID]
	if !exist {
		return models.ErrNotFound
	}

	if currentUserState.Username != user.Username {
		delete(r.usernamesMap, user.Username)
		r.usernamesMap[user.Username] = user.ID
	}

	if currentUserState.Email != user.Email {
		delete(r.emailsMap, user.Email)
		r.emailsMap[user.Email] = user.ID
	}

	r.store[user.ID] = user

	return nil
}

func (r *UserRepository) Delete(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exist := r.store[user.ID]
	if !exist {
		return models.ErrNotFound
	}

	delete(r.store, user.ID)
	delete(r.usernamesMap, user.Username)
	delete(r.emailsMap, user.Email)

	return nil
}

func (r *UserRepository) newUserID() int64 {
	return r.idCounter.Add(1)
}
