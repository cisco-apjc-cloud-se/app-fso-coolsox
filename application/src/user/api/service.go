package api

// service.go contains the definition and implementation (business logic) of the
// user service. Everything here is agnostic to the transport (HTTP).

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/microservices-demo/user/db"
	"github.com/microservices-demo/user/users"

	// AppDynamics Go SDK Agent
	appd "appdynamics"
	// Note: Delete() not instrumented
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

// Service is the user service, providing operations for users to login, register, and retrieve customer information.
type Service interface {
	Login(username, password string) (users.User, error) // GET /login
	Register(username, password, email, first, last string) (string, error)
	GetUsers(id string) ([]users.User, error)
	PostUser(u users.User) (string, error)
	GetAddresses(id string) ([]users.Address, error)
	PostAddress(u users.Address, userid string) (string, error)
	GetCards(id string) ([]users.Card, error)
	PostCard(u users.Card, userid string) (string, error)
	Delete(entity, id string) error
	Health() []Health // GET /health
}

// NewFixedService returns a simple implementation of the Service interface,
func NewFixedService() Service {
	return &fixedService{}
}

type fixedService struct{}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func (s *fixedService) Login(username, password string) (users.User, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB Login", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	u, err := db.GetUserByName(username)
	if err != nil {
		return users.New(), err
		// AppD
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}
	if u.Password != calculatePassHash(password, u.Salt) {
		return users.New(), ErrUnauthorized
		// AppD
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "User Unauthorized Error", true)
	}
	db.GetUserAttributes(&u)
	u.MaskCCs()

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return u, nil

}

func (s *fixedService) Register(username, password, email, first, last string) (string, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB Register", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	u := users.New()
	u.Username = username
	u.Password = calculatePassHash(password, u.Salt)
	u.Email = email
	u.FirstName = first
	u.LastName = last
	err := db.CreateUser(&u)

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return u.UserID, err
}

func (s *fixedService) GetUsers(id string) ([]users.User, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB GetUsers", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	if id == "" {
		us, err := db.GetUsers()
		// AppD Error
		if err != nil {
			appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
		}
		for k, u := range us {
			u.AddLinks()
			us[k] = u
		}
		// AppD End BT & Exit Call
		appd.EndExitcall(exitHandle)
		appd.EndBT(btHandle)

		return us, err
	}
	u, err := db.GetUser(id)
	u.AddLinks()

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return []users.User{u}, err
}

func (s *fixedService) PostUser(u users.User) (string, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB PostUser", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	u.NewSalt()
	u.Password = calculatePassHash(u.Password, u.Salt)
	err := db.CreateUser(&u)

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return u.UserID, err
}

func (s *fixedService) GetAddresses(id string) ([]users.Address, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB GetAddresses", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	if id == "" {
		as, err := db.GetAddresses()

		// AppD Error
		if err != nil {
			appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
		}

		for k, a := range as {
			a.AddLinks()
			as[k] = a
		}
		// AppD End BT & Exit Call
		appd.EndExitcall(exitHandle)
		appd.EndBT(btHandle)

		return as, err
	}
	a, err := db.GetAddress(id)
	a.AddLinks()

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return []users.Address{a}, err
}

func (s *fixedService) PostAddress(add users.Address, userid string) (string, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB PostAddress", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	err := db.CreateAddress(&add, userid)

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return add.ID, err
}

func (s *fixedService) GetCards(id string) ([]users.Card, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB GetCards", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	if id == "" {
		cs, err := db.GetCards()

		// AppD Error
		if err != nil {
			appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
		}

		for k, c := range cs {
			c.AddLinks()
			cs[k] = c
		}

		// AppD End BT & Exit Call
		appd.EndExitcall(exitHandle)
		appd.EndBT(btHandle)

		return cs, err
	}
	c, err := db.GetCard(id)
	c.AddLinks()

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return []users.Card{c}, err
}

func (s *fixedService) PostCard(card users.Card, userid string) (string, error) {
	// AppD Start BT & Exit Call
	btHandle := appd.StartBT("MongoDB PostCard", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")

	err := db.CreateCard(&card, userid)

	// AppD Error
	if err != nil {
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Database Error", true)
	}

	// AppD End BT & Exit Call
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	return card.ID, err
}

func (s *fixedService) Delete(entity, id string) error {
	return db.Delete(entity, id)
}

func (s *fixedService) Health() []Health {
	var health []Health
	dbstatus := "OK"

	err := db.Ping()
	if err != nil {
		dbstatus = "err"
	}

	app := Health{"user", "OK", time.Now().String()}
	db := Health{"user-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
