# Velox Framework Examples

This document provides code samples from a working Velox application to help you understand how the framework is structured and how to use its features effectively.

## Table of Contents

- [Application Structure](#application-structure)
- [Application Initialization](#application-initialization)
- [Routing](#routing)
- [Middleware](#middleware)
- [Handlers](#handlers)
- [Models](#models)
- [Templates and Views](#templates-and-views)
- [Session Management](#session-management)
- [Email Sending](#email-sending)
- [Error Handling](#error-handling)

## Application Structure

A typical Velox application follows this structure:

```
myapp/
├── config/         # Configuration files (database.yml, etc.)
├── data/           # Data models and database code
├── handlers/       # HTTP request handlers
├── migrations/     # Database migrations
├── middleware/     # HTTP middleware
├── public/         # Static assets (JS, CSS, images)
├── views/          # Templates (Jet or Go templates)
├── mail/           # Email templates
├── .env            # Environment configuration
├── go.mod          # Go module definition
├── go.sum          # Go dependencies lock file
├── Makefile        # Build and run commands
└── main.go         # Application entry point
```

## Application Initialization

Here's how a Velox application is initialized in `main.go`:

```go
package main

import (
    "myapp/data"
    "myapp/handlers"
    "myapp/middleware"
    "os"
    "os/signal"
    "sync"
    "syscall"

    "github.com/FernandoJVideira/velox"
)

type application struct {
    App        *velox.Velox
    Handlers   *handlers.Handlers
    Models     data.Models
    Middleware *middleware.Middleware
    wg         sync.WaitGroup
}

func main() {
    v := initApplication()
    go v.listenForShutdwn()
    err := v.App.ListenAndServe()
    if err != nil {
        v.App.ErrorLog.Println(err)
    }
}

func (a *application) shutdown() {
    // Put any cleanup tasks here
    a.App.InfoLog.Println("Shutting down")

    // Block until the waitgroup is empty
    a.wg.Wait()
}

func (a *application) listenForShutdwn() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    s := <-quit

    a.App.InfoLog.Println("Received signal:", s.String())
    a.shutdown()

    os.Exit(0)
}

func initApplication() *application {
    // Create a new application with default settings
    app := &application{}

    // Initialize Velox
    vx := &velox.Velox{}
    err := vx.New(".")
    if err != nil {
        panic(err)
    }
    app.App = vx

    // Set up models
    app.Models = data.New(app.App.DB.Pool)

    // Set up middleware
    app.Middleware = &middleware.Middleware{
        App: app.App,
    }

    // Set up handlers
    app.Handlers = &handlers.Handlers{
        App:    app.App,
        Models: app.Models,
    }

    // Set up application routes
    app.routes()

    return app
}
```

## Routing

Define your routes in a separate function:

```go
func (a *application) routes() {
    // Use middleware for all routes
    a.App.Routes.Use(a.Middleware.CheckRemember)

    // Public routes
    a.App.Routes.Get("/", a.Handlers.Home)
    a.App.Routes.Get("/about", a.Handlers.About)
    a.App.Routes.Get("/contact", a.Handlers.Contact)

    // Auth routes
    a.App.Routes.Get("/login", a.Handlers.LoginPage)
    a.App.Routes.Post("/login", a.Handlers.Login)
    a.App.Routes.Get("/logout", a.Handlers.Logout)
    a.App.Routes.Get("/register", a.Handlers.RegisterPage)
    a.App.Routes.Post("/register", a.Handlers.Register)

    // Protected routes - group with authentication middleware
    mux := a.App.Routes.Group("/admin")
    mux.Use(a.Middleware.Auth)

    mux.Get("/dashboard", a.Handlers.Dashboard)
    mux.Get("/profile", a.Handlers.Profile)

    // Serve static files
    fileServer := http.FileServer(http.Dir("./public"))
    a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))
}
```

## Middleware

Middleware examples from `middleware/auth.go`:

```go
package middleware

import (
    "net/http"

    "github.com/FernandoJVideira/velox"
)

type Middleware struct {
    App *velox.Velox
}

// Auth middleware checks if user is authenticated
func (m *Middleware) Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        if !m.App.Session.Exists(r.Context(), "userID") {
            m.App.Session.Put(r.Context(), "error", "You must log in to access this page")
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// CheckRemember middleware handles "remember me" cookie functionality
func (m *Middleware) CheckRemember(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // If user is not logged in, check for remember token
        if !m.App.Session.Exists(r.Context(), "userID") {
            // Get cookie
            cookie, err := r.Cookie(fmt.Sprintf("_%s_remember", m.App.AppName))
            if err == nil {
                // Cookie exists, validate it
                key := cookie.Value
                parts := strings.Split(key, "|")
                if len(parts) == 2 {
                    // Check if remember token exists in the database
                    userID, _ := strconv.Atoi(parts[0])
                    user := data.User{}

                    if user.CheckForRememberToken(userID, parts[1]) {
                        // Valid remember token, log user in
                        m.App.Session.Put(r.Context(), "userID", userID)
                    }
                }
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

## Handlers

Example handler implementation:

```go
package handlers

import (
    "myapp/data"
    "net/http"

    "github.com/FernandoJVideira/velox"
)

type Handlers struct {
    App    *velox.Velox
    Models data.Models
}

// Home displays the home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
    err := h.App.Render.Page(w, r, "home", nil, nil)
    if err != nil {
        h.App.ErrorLog.Println("Error rendering:", err)
        h.App.Error500(w, r)
    }
}

// LoginPage displays the login page
func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
    err := h.App.Render.Page(w, r, "login", nil, nil)
    if err != nil {
        h.App.ErrorLog.Println("Error rendering:", err)
        h.App.Error500(w, r)
    }
}

// Login handles user login
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        h.App.ErrorLog.Println("Error parsing form:", err)
        h.App.Error500(w, r)
        return
    }

    // Get form data
    email := r.Form.Get("email")
    password := r.Form.Get("password")

    // Get the user from the database
    user, err := h.Models.Users.GetByEmail(email)
    if err != nil {
        h.App.Session.Put(r.Context(), "error", "Invalid credentials")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Check password
    matches, err := user.PasswordMatches(password)
    if err != nil {
        h.App.Session.Put(r.Context(), "error", "Server error. Please try again later.")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if !matches {
        h.App.Session.Put(r.Context(), "error", "Invalid credentials")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Handle "Remember me" functionality
    if r.Form.Get("remember") == "on" {
        // Generate remember token and save it
        token := h.App.RandomString(32)
        user.RememberToken = token
        err = h.Models.Users.Update(*user)
        if err != nil {
            h.App.ErrorLog.Println("Error updating user:", err)
        } else {
            // Set cookie
            cookie := http.Cookie{
                Name:     fmt.Sprintf("_%s_remember", h.App.AppName),
                Value:    fmt.Sprintf("%d|%s", user.ID, token),
                Path:     "/",
                Expires:  time.Now().Add(365 * 24 * time.Hour),
                HttpOnly: true,
                Domain:   h.App.Session.Cookie.Domain,
                Secure:   h.App.Session.Cookie.Secure,
                SameSite: http.SameSiteStrictMode,
            }
            http.SetCookie(w, &cookie)
        }
    }

    // Log user in
    h.App.Session.Put(r.Context(), "userID", user.ID)
    h.App.Session.Put(r.Context(), "success", "You have been logged in successfully")

    // Redirect to intended destination or dashboard
    if h.App.Session.Exists(r.Context(), "intended") {
        http.Redirect(w, r, h.App.Session.Get(r.Context(), "intended").(string), http.StatusSeeOther)
        return
    }

    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
```

## Models

Example User model:

```go
package data

import (
    "errors"
    "time"

    up "github.com/upper/db/v4"
    "golang.org/x/crypto/bcrypt"
)

// User is the type for a user
type User struct {
    ID          int       `db:"id,omitempty"`
    FirstName   string    `db:"first_name"`
    LastName    string    `db:"last_name"`
    Email       string    `db:"email"`
    Active      int       `db:"user_active"`
    Password    string    `db:"password"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
    Token       Token     `db:"-"`
}

// Table returns the table name associated with this model in the database
func (u *User) Table() string {
    return "users"
}

// GetAll returns a slice of all users
func (u *User) GetAll() ([]*User, error) {
    collection := upper.Collection(u.Table())
    var all []*User
    res := collection.Find().OrderBy("last_name")
    err := res.All(&all)
    if err != nil {
        return nil, err
    }
    return all, nil
}

// GetByEmail gets one user, by email
func (u *User) GetByEmail(email string) (*User, error) {
    var theUser User
    collection := upper.Collection(u.Table())
    res := collection.Find(up.Cond{"email =": email})
    err := res.One(&theUser)
    if err != nil {
        return nil, err
    }
    return &theUser, nil
}

// Get gets one user by id
func (u *User) Get(id int) (*User, error) {
    var theUser User
    collection := upper.Collection(u.Table())
    res := collection.Find(up.Cond{"id =": id})
    err := res.One(&theUser)
    if err != nil {
        return nil, err
    }
    return &theUser, nil
}

// Update updates a user record in the database
func (u *User) Update(user User) error {
    user.UpdatedAt = time.Now()
    collection := upper.Collection(u.Table())
    res := collection.Find(user.ID)
    err := res.Update(&user)
    return err
}

// Delete deletes a user by id
func (u *User) Delete(id int) error {
    collection := upper.Collection(u.Table())
    res := collection.Find(id)
    return res.Delete()
}

// PasswordMatches verifies a supplied password against the hash stored in the database
func (u *User) PasswordMatches(plainText string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
    if err != nil {
        switch {
        case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
            // invalid password
            return false, nil
        default:
            // some kind of error occurred
            return false, err
        }
    }
    return true, nil
}

// CheckForRememberToken checks if the remember token exists in the database
func (u *User) CheckForRememberToken(id int, token string) bool {
    var rememberToken RememberToken
    rt := RememberToken{}
    collection := upper.Collection(rt.Table())
    res := collection.Find(up.Cond{"user_id": id, "remember_token": token})
    err := res.One(&rememberToken)
    return err == nil
}
```

## Templates and Views

Velox supports the Jet template engine. Here's an example of a layout and view:

`views/layouts/base.jet`:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{isset(title) ? title : "My App"}}</title>
    <link
      href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"
      rel="stylesheet"
    />
    <link rel="stylesheet" href="/public/css/styles.css" />
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
      <div class="container">
        <a class="navbar-brand" href="/">MyApp</a>
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarNav"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto">
            <li class="nav-item">
              <a class="nav-link" href="/">Home</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="/about">About</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="/contact">Contact</a>
            </li>
          </ul>
          <ul class="navbar-nav">
            {{if authenticated()}}
            <li class="nav-item">
              <a class="nav-link" href="/dashboard">Dashboard</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="/logout">Logout</a>
            </li>
            {{else}}
            <li class="nav-item">
              <a class="nav-link" href="/login">Login</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="/register">Register</a>
            </li>
            {{end}}
          </ul>
        </div>
      </div>
    </nav>

    <div class="container mt-4">
      {{if isset(flashMessage)}}
      <div class="alert alert-{{flashMessage.type}}">
        {{flashMessage.message}}
      </div>
      {{end}} {{yield content}}
    </div>

    <footer class="footer mt-auto py-3 bg-light">
      <div class="container text-center">
        <span class="text-muted"
          >&copy; {{year()}} MyApp. All rights reserved.</span
        >
      </div>
    </footer>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
  </body>
</html>
```

`views/home.jet`:

```html
{{extends "layouts/base.jet"}} {{block content()}}
<div class="jumbotron">
  <h1 class="display-4">Welcome to MyApp!</h1>
  <p class="lead">A sample application built with Velox.</p>
  <hr class="my-4" />
  <p>
    Velox is a powerful Go web framework designed to provide a Laravel-like
    experience with Go's performance.
  </p>
  <a class="btn btn-primary btn-lg" href="/about" role="button">Learn more</a>
</div>
{{end}}
```

## Session Management

Working with sessions in controllers:

```go
// Set session values
h.App.Session.Put(r.Context(), "user_id", user.ID)
h.App.Session.Put(r.Context(), "success", "You've been logged in successfully!")

// Get session values
if h.App.Session.Exists(r.Context(), "user_id") {
    userID := h.App.Session.Get(r.Context(), "user_id").(int)
    // Do something with userID
}

// Remove session values
h.App.Session.Remove(r.Context(), "user_id")

// Flash messages (available for the next request only)
h.App.Session.Put(r.Context(), "flash", "Action completed successfully")
```

## Email Sending

Sending an email:

```go
// Create message data
data := map[string]any{
    "name":   "John Doe",
    "message": "Thank you for registering!",
    "url":    "https://example.com/verify?token=abc123",
}

// Send email using template
msg := velox.MailMessage{
    To:       "recipient@example.com",
    From:     "sender@example.com",
    Subject:  "Welcome to MyApp",
    Template: "welcome",
    Data:     data,
}

// Send the email asynchronously
h.App.Mail.Jobs <- msg
```

## Error Handling

Handling errors in your application:

```go
// Handling errors in handlers
if err != nil {
    h.App.ErrorLog.Println("Database error:", err)
    h.App.Session.Put(r.Context(), "error", "Something went wrong, please try again.")
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

// Using built-in error pages
if user == nil {
    h.App.ErrorStatus(w, http.StatusNotFound)
    return
}

// Creating custom error responses
if !authorized {
    w.WriteHeader(http.StatusForbidden)
    err := h.App.Render.Page(w, r, "errors/403", nil, nil)
    if err != nil {
        h.App.ErrorLog.Println("Error rendering:", err)
        h.App.Error500(w, r)
    }
    return
}

// Performance measurement
defer h.App.LoadTime(time.Now())
```

These examples should give you a good understanding of how to use the Velox framework for building Go web applications.
