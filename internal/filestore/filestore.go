package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// ErrNotFound is returned when the ID can't be found.
var ErrNotFound = errors.New("not found")

// Policy refers to a list of notifiers.
type Policy struct {
	NotifierIDs []string
}

// Notifier defines how to notify.
type Notifier struct {
	Email string
}

type db struct {
	Notifiers map[string]Notifier
	Policies  map[string]Policy
}

func open() (*db, error) {
	contents, err := ioutil.ReadFile("db")
	if os.IsNotExist(err) {
		return &db{}, nil
	}

	var db db
	if err := json.Unmarshal(contents, &db); err != nil {
		return nil, err
	}

	return &db, nil
}

func (d *db) save() error {
	out, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("db", out, 0666)
}

func ReadPolicy(id string) (Policy, error) {
	db, err := open()
	if err != nil {
		return Policy{}, err
	}

	svc, ok := db.Policies[id]
	if !ok {
		return Policy{}, ErrNotFound
	}

	return svc, nil
}

func WritePolicy(id string, policy Policy, create bool) error {
	db, err := open()
	if err != nil {
		return err
	}

	if _, exists := db.Policies[id]; create && exists {
		return fmt.Errorf("policy with id %q already exists", id)
	}

	// Foreign key check.
	for _, notifierID := range policy.NotifierIDs {
		if _, ok := db.Notifiers[notifierID]; !ok {
			return fmt.Errorf("referencing a non-existent notifer: %q (have %#v)", notifierID, db.Notifiers)
		}
	}

	if db.Policies == nil {
		db.Policies = make(map[string]Policy)
	}

	db.Policies[id] = policy
	return db.save()
}

func DeletePolicy(id string) error {
	db, err := open()
	if err != nil {
		return err
	}

	if _, ok := db.Policies[id]; !ok {
		return ErrNotFound
	}

	delete(db.Policies, id)
	return db.save()
}

func ReadNotifier(id string) (Notifier, error) {
	db, err := open()
	if err != nil {
		return Notifier{}, err
	}

	svc, ok := db.Notifiers[id]
	if !ok {
		return Notifier{}, ErrNotFound
	}

	return svc, nil
}

func WriteNotifier(id string, svc Notifier, create bool) error {
	db, err := open()
	if err != nil {
		return err
	}

	if _, exists := db.Notifiers[id]; create && exists {
		return fmt.Errorf("notifier with id %q already exists", id)
	}

	if db.Notifiers == nil {
		db.Notifiers = make(map[string]Notifier)
	}

	db.Notifiers[id] = svc
	return db.save()
}

func DeleteNotifier(id string) error {
	db, err := open()
	if err != nil {
		return err
	}

	// Foreign key check.
	var policies []string
	for policyID, policy := range db.Policies {
		for _, notifierID := range policy.NotifierIDs {
			if id == notifierID {
				policies = append(policies, policyID)
			}
		}
	}
	if len(policies) > 0 {
		sort.Strings(policies)
		return fmt.Errorf("cannot delete notifier, as policies %v still refer to it", policies)
	}

	if _, ok := db.Notifiers[id]; !ok {
		return ErrNotFound
	}

	delete(db.Notifiers, id)
	return db.save()
}
