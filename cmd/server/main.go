package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"local-ca-downloader/internal/certificate"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the password for accessing the application
	password := os.Getenv("AUTH_PASSWORD")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPasswordString := string(hashedPassword)

	// Middleware to validate the password or check the authentication cookie
	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the "authenticated" cookie from the request
			cookie, err := c.Cookie("authenticated")

			// Check if the cookie is not present or the cookie value doesn't match the encoded password
			if err != nil || cookie.Value != hashedPasswordString {
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
					Value:    hashedPasswordString,
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

	// Routes
	e.Any("/", func(c echo.Context) error {
		publicCADetails := certificate.BuildCertificateDetails("certs/public-ca.pem")
		publicCertDetails := certificate.BuildCertificateDetails("certs/cert.pem")

		data := map[string]interface{}{
			"CA":   publicCADetails,
			"Cert": publicCertDetails,
		}

		return c.Render(http.StatusOK, "nav.html", data)
	}, authMiddleware)

	e.GET("/login", func(c echo.Context) error {
		// Render the login page
		return c.Render(http.StatusOK, "login.html", nil)
	})

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

	e.POST("/generate", generateNewCerts, authMiddleware)

	e.POST("/special", func (c echo.Context) error {
		return c.Render(http.StatusOK, "surprise.html", nil)
	}, )

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

func generateNewCerts(c echo.Context) error {
	caPassphrase := c.FormValue("ca_passphrase")
	caTtl := c.FormValue("ca_ttl")
	caSubject := c.FormValue("ca_subject")
	tlsCn := c.FormValue("tls_cn")
	tlsAltNames := c.FormValue("tls_alt_names")
	tlsTtl := c.FormValue("tls_ttl")

	cmd := exec.Command("./scripts/generate_new_certs.sh", caPassphrase, caTtl, caSubject, tlsCn, tlsAltNames, tlsTtl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	publicCADetails := certificate.BuildCertificateDetails("certs/public-ca.pem")
	publicCertDetails := certificate.BuildCertificateDetails("certs/cert.pem")

	if err != nil {
		data := map[string]interface{}{
			"Error": "Error generating new certs, original certs have been restored. Please check your input and try again.",
			"CA":    publicCADetails,
			"Cert":  publicCertDetails,
		}
		return c.Render(http.StatusOK, "nav.html", data)
	}

	data := map[string]interface{}{
		"Success": "Certificate generated successfully.",
		"CA":      publicCADetails,
		"Cert":    publicCertDetails,
	}
	return c.Render(http.StatusOK, "nav.html", data)
}
