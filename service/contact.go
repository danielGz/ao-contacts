package service

import (
	"accelone-contacts/model"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

// ContactService Interface to abstract contact operations, assuming there could be many stores for contacts or persistence mechanisms
type ContactService interface {
	Get(page, limit int) ([]model.Contact, error)
	Create(c model.Contact) (model.Contact, error)
	GetById(id string) (model.Contact, error)
	Delete(id string) (bool, error)
	Update(c model.Contact) (model.Contact, error)
}

func NewInMemoryContactService() *InMemoryContactService {
	return &InMemoryContactService{
		contacts: make(map[string]model.Contact),
	}
}

// InMemoryContactService implementation of in-memory store for contact, concurrent read/write errors are prevented with a mutex
type InMemoryContactService struct {
	contacts map[string]model.Contact
	lock     sync.RWMutex // Use RWMutex for read/write operations
}

func (i *InMemoryContactService) Get(page, limit int) ([]model.Contact, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	var contactsSlice []model.Contact
	for _, contact := range i.contacts {
		contactsSlice = append(contactsSlice, contact)
	}

	start := (page - 1) * limit
	if start > len(contactsSlice) {
		return []model.Contact{}, nil // Return empty if start is beyond the array length
	}

	end := start + limit
	if end > len(contactsSlice) {
		end = len(contactsSlice)
	}

	if len(contactsSlice) == 0 {
		return []model.Contact{}, nil
	}

	return contactsSlice[start:end], nil
}

func (i *InMemoryContactService) Create(c model.Contact) (model.Contact, error) {
	i.lock.Lock()         // Lock for writing
	defer i.lock.Unlock() // Unlock when the function returns

	if _, ok := i.contacts[c.Id]; ok && c.Id != "" {
		// If the ID is set manually and already exists, return an error
		return model.Contact{}, fmt.Errorf("contact with id '%s' already exists", c.Id)
	}

	// Generate a UUID for the contact if not provided
	if c.Id == "" {
		c.Id = uuid.New().String() // generate uuid of 32 bit
	}

	// Store the contact
	i.contacts[c.Id] = c

	return i.contacts[c.Id], nil
}

func (i *InMemoryContactService) GetById(id string) (model.Contact, error) {
	i.lock.RLock()         // Lock for reading
	defer i.lock.RUnlock() // Unlock when the function returns

	if _, ok := i.contacts[id]; ok {
		return i.contacts[id], nil
	}
	return model.Contact{}, fmt.Errorf("contact not found for Id %s", id)
}

func (i *InMemoryContactService) Delete(id string) (bool, error) {
	i.lock.Lock()         // Lock for writing
	defer i.lock.Unlock() // Unlock when the function returns

	if _, ok := i.contacts[id]; ok {
		delete(i.contacts, id)
		return true, nil
	}
	return false, fmt.Errorf("contact not found for Id %s", id)
}

func (i *InMemoryContactService) Update(c model.Contact) (model.Contact, error) {
	i.lock.Lock()         // Lock for writing
	defer i.lock.Unlock() // Unlock when the function returns

	if _, ok := i.contacts[c.Id]; ok {
		i.contacts[c.Id] = c
		return i.contacts[c.Id], nil
	}
	return model.Contact{}, fmt.Errorf("contact not found for Id %s", c.Id)
}
