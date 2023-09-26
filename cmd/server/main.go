package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"local-ca-downloader/internal/certificate"
	"time"

	// "local-ca-downloader/internal/certificate"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the password for accessing the application
	password := os.Getenv("AUTH_PASSWORD")

	// Middleware to validate the password or check the authentication cookie
	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the "authenticated" cookie from the request
			cookie, err := c.Cookie("authenticated")

			// Check if the cookie is not present or the cookie value doesn't match the encoded password
			if err != nil || cookie.Value != encodeBase64(password) {
				// Get the password from the request form
				reqPassword := c.FormValue("password")

				// Check if the request password is empty or doesn't match the defined password
				if reqPassword == "" || reqPassword != password {
					// Determine the error message based on the condition
					errorMsg := "Invalid password"
					if reqPassword == "" {
						errorMsg = ""
					}

					// Prepare the data for rendering the login page with the error message
					data := map[string]interface{}{
						"Error": errorMsg,
					}

					// Render the login page with the error message
					return c.Render(http.StatusOK, "login.html", data)
				}

				// Set the authentication cookie with the base64 encoded password and expiration time
				cookie := &http.Cookie{
					Name:     "authenticated",
					Value:    encodeBase64(password),
					Path:     "/",
					HttpOnly: true,
					MaxAge:   900, // 15 minutes (in seconds)
				}
				c.SetCookie(cookie)
			}

			// Continue to the next handler
			return next(c)
		}
	}

	// Register the HTML template renderer
	e.Renderer = createRenderer()

	// Build Cert-Details for all Certs
	publicCADetails := certificate.BuildCertificateDetails("certs/public-ca.pem")
	publicCertDetails := certificate.BuildCertificateDetails("certs/cert.pem")
	// privateCertDetails := certificate.BuildCertificateDetails("certs/cert-key.pem")

	// Routes
	e.Any("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "nav.html", publicCADetails)
	}, authMiddleware)

	e.GET("/login", func(c echo.Context) error {
		// Render the login page
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.GET("/ca", func(c echo.Context) error {
		return c.Render(http.StatusOK, "nav.html", publicCADetails)
	}, authMiddleware)

	e.GET("/cert", func(c echo.Context) error {
		return c.Render(http.StatusOK, "nav.html", publicCertDetails)
	}, authMiddleware)

	// e.GET("/certKey", func(c echo.Context) error {
	// 	return c.Render(http.StatusOK, "nav.html", privateCertDetails)
	// }, authMiddleware)

	e.GET("/download/ca", func(c echo.Context) error {
		return c.Attachment("certs/public-ca.pem", "public-ca.pem")
	}, authMiddleware)

	e.GET("/download/cert", func(c echo.Context) error {
		return c.Attachment("certs/cert.pem", "cert.pem")
	}, authMiddleware)

	e.GET("/download/certKey", func(c echo.Context) error {
		return c.Attachment("certs/cert-key.pem", "cert-key.pem")
	}, authMiddleware)

	e.POST("/logout", deleteCookieHandler)

	// Start the server
	address := ":8443"
	fmt.Printf("Server listening on %s\n", address)
	e.Logger.Fatal(e.StartTLS(address, "certs/cert.pem", "certs/cert-key.pem"))
}

// Template struct for rendering HTML templates
type Template struct {
	templates *template.Template
}

// Render method to render templates
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Execute the template and write the output to a buffer
	buf := new(bytes.Buffer)
	if err := t.templates.ExecuteTemplate(buf, name, data); err != nil {
		return err
	}

	// Write the rendered HTML to the response
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

// Function to create and return the HTML template renderer
func createRenderer() echo.Renderer {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	return t
}

// Function to base64 encode a string
func encodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Function to delete a cookie (Basically it just clears the cookie out)
func deleteCookieHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "authenticated"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.Path = "/"

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}