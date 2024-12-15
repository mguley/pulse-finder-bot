package dto

import (
	"domain/vacancy/entity"
	"sync"
	"time"
)

var vacancyPool = newVacancyPool()

// newVacancyPool creates and initializes a singleton sync.Pool for managing reusable Vacancy objects.
// It uses a once.Do mechanism to ensure the pool is only created once during runtime.
func newVacancyPool() func() *sync.Pool {
	var once sync.Once
	var pool *sync.Pool

	return func() *sync.Pool {
		once.Do(func() {
			pool = &sync.Pool{
				New: func() interface{} {
					return &Vacancy{}
				},
			}
		})
		return pool
	}
}

// Vacancy represents a job vacancy with details.
type Vacancy struct {
	Title       string    // The title of the job vacancy.
	Company     string    // The company offering the job vacancy.
	Description string    // A brief description of the job vacancy.
	PostedAt    time.Time // The timestamp when the job was posted.
	Location    string    // The location of the job vacancy.
}

// Reset clears all fields of the Vacancy object and resets them to their zero values.
func (v *Vacancy) Reset() *Vacancy {
	v.Title = ""
	v.Company = ""
	v.Description = ""
	v.PostedAt = time.Time{}
	v.Location = ""
	return v
}

// Release returns the Vacancy instance to the pool after resetting its fields.
func (v *Vacancy) Release() {
	vacancyPool().Put(v.Reset())
}

// GetVacancy retrieves a Vacancy object from the pool, ensuring it is reset before use.
// If no object is available, a new one is created.
func GetVacancy() *Vacancy {
	return vacancyPool().Get().(*Vacancy).Reset()
}

// ToEntity maps the fields of the Vacancy object to an entity.Vacancy object.
func (v *Vacancy) ToEntity(e *entity.Vacancy) {
	e.Title = v.Title
	e.Company = v.Company
	e.Description = v.Description
	e.PostedAt = v.PostedAt
	e.Location = v.Location
}
