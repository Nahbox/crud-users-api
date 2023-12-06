package user

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Nahbox/crud-users-api/internal/db"
	"github.com/Nahbox/crud-users-api/internal/handler"
	"github.com/Nahbox/crud-users-api/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

type UserSerialized struct {
	Id   int    `json:"id"`
	Json string `json:"json"`
	Xml  string `json:"xml"`
	Toml string `json:"toml"`
}

const userIDKey = "userID"

func serializeUsers(users []models.User) ([]UserSerialized, error) {
	usersSerialized := make([]UserSerialized, 0, len(users))
	for _, u := range users {
		userSer, err := serializeUser(u)
		if err != nil {
			return nil, err
		}
		usersSerialized = append(usersSerialized, *userSer)
	}
	return usersSerialized, nil
}

// ТЗ: Сериализация в разные форматы должна происходить одновременно.
func serializeUser(user models.User) (*UserSerialized, error) {
	var userSer UserSerialized
	errsCh := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		jsonData, err := json.Marshal(user)
		if err != nil {
			errsCh <- err
			return
		}
		userSer.Json = string(jsonData)
	}()
	go func() {
		defer wg.Done()
		xmlData, err := xml.Marshal(user)
		if err != nil {
			errsCh <- err
			return
		}
		userSer.Xml = string(xmlData)
	}()
	go func() {
		defer wg.Done()
		tomlData := new(bytes.Buffer)
		if err := toml.NewEncoder(tomlData).Encode(user); err != nil {
			errsCh <- err
			return
		}
		userSer.Toml = tomlData.String()
	}()

	wg.Wait()
	close(errsCh)

	for err := range errsCh {
		return nil, err
	}

	userSer.Id = user.Id

	return &userSer, nil
}

func (h *Handler) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")
		if userId == "" {
			render.Render(w, r, handler.ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, r, handler.ErrorRenderer(fmt.Errorf("invalid user ID")))
		}
		ctx := context.WithValue(r.Context(), userIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Tags Users
// @Description Create a new user based on the provided JSON data
// @ID create-user
// @Accept json
// @Produce json
// @Param user body models.NewUser true "User object to be created"
// @Success 200 {object} models.User "Created user object"
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, handler.ErrBadRequest)
		return
	}
	if err := h.db.AddUser(user); err != nil {
		render.Render(w, r, handler.ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
}

// GetAllUsers godoc
// @Summary Get a list of all users
// @Tags Users
// @Description Get a list of all users in the system
// @ID get-all-users
// @Accept json
// @Produce json
// @Success 200 {array} models.User "List of users"
// @Router /users [get]
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println()
	users, err := h.db.GetAllUsers()
	if err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
	out, err := serializeUsers(users)
	if err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
	render.JSON(w, r, out)
}

// GetUser godoc
// @Summary Get user details by ID
// @Tags Users
// @Description Get user details based on the provided user ID
// @ID get-user-by-id
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.User "User details"
// @Router /users/{user_id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	user, err := h.db.GetUserById(userID)
	if err != nil {
		if errors.Is(err, db.ErrNoMatch) {
			render.Render(w, r, handler.ErrNotFound)
		} else {
			render.Render(w, r, handler.ErrorRenderer(err))
		}
		return
	}
	out, err := serializeUser(user)
	if err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
	render.JSON(w, r, out)
}

// UpdateUser godoc
// @Summary Update user details by ID
// @Tags Users
// @Description Update user details based on the provided user ID
// @ID update-user-by-id
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param user body models.NewUser true "Updated user details"
// @Success 200 {object} models.User "Updated user details"
// @Router /users/{user_id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	userData := &models.User{}
	if err := render.Bind(r, userData); err != nil {
		render.Render(w, r, handler.ErrBadRequest)
		return
	}
	item, err := h.db.UpdateUser(userId, userData)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, handler.ErrNotFound)
		} else {
			render.Render(w, r, handler.ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, handler.ServerErrorRenderer(err))
		return
	}
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags Users
// @Description Delete user based on the provided user ID
// @ID delete-user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Router /users/{user_id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	err := h.db.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, db.ErrNoMatch) {
			render.Render(w, r, handler.ErrNotFound)
		} else {
			render.Render(w, r, handler.ServerErrorRenderer(err))
		}
		return
	}
}
