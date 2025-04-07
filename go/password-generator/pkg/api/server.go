package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ss49919201/ai-kata/go/password-generator/internal/config"
	"github.com/ss49919201/ai-kata/go/password-generator/pkg/generator"
)

type Server struct {
	Router    *mux.Router
	Generator *generator.Generator
	Config    *config.Config
}

type PasswordRequest struct {
	Length    *int  `json:"length,omitempty"`
	Uppercase *bool `json:"uppercase,omitempty"`
	Lowercase *bool `json:"lowercase,omitempty"`
	Digits    *bool `json:"digits,omitempty"`
	Special   *bool `json:"special,omitempty"`
}

type PasswordResponse struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(cfg *config.Config) *Server {
	gen, err := generator.NewGenerator(
		cfg.Generator.MinLength,
		cfg.Generator.MaxLength,
		generator.CharacterSet{
			Uppercase: cfg.Generator.UseUppercase,
			Lowercase: cfg.Generator.UseLowercase,
			Digits:    cfg.Generator.UseDigits,
			Special:   cfg.Generator.UseSpecial,
		},
	)
	if err != nil {
		panic(err) // In a real app, handle this more gracefully
	}

	s := &Server{
		Router:    mux.NewRouter(),
		Generator: gen,
		Config:    cfg,
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	s.Router.HandleFunc("/api/v1/password", s.handleGeneratePassword).Methods("GET")
	s.Router.HandleFunc("/api/v1/password", s.handleGeneratePasswordJSON).Methods("POST")
	
	s.Router.HandleFunc("/health", s.handleHealthCheck).Methods("GET")
	
	s.Router.HandleFunc("/", s.handleDocumentation).Methods("GET")
}

func (s *Server) handleGeneratePassword(w http.ResponseWriter, r *http.Request) {
	length := s.Config.Generator.DefaultLength
	if lenParam := r.URL.Query().Get("length"); lenParam != "" {
		parsedLen, err := strconv.Atoi(lenParam)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid length parameter")
			return
		}
		length = parsedLen
	}

	charSet := generator.CharacterSet{
		Uppercase: getBoolParam(r, "uppercase", s.Config.Generator.UseUppercase),
		Lowercase: getBoolParam(r, "lowercase", s.Config.Generator.UseLowercase),
		Digits:    getBoolParam(r, "digits", s.Config.Generator.UseDigits),
		Special:   getBoolParam(r, "special", s.Config.Generator.UseSpecial),
	}

	password, err := s.generatePassword(length, charSet)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, PasswordResponse{
		Password: password,
		Length:   length,
	})
}

func (s *Server) handleGeneratePasswordJSON(w http.ResponseWriter, r *http.Request) {
	var req PasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	length := s.Config.Generator.DefaultLength
	if req.Length != nil {
		length = *req.Length
	}

	charSet := generator.CharacterSet{
		Uppercase: getPointerBoolOrDefault(req.Uppercase, s.Config.Generator.UseUppercase),
		Lowercase: getPointerBoolOrDefault(req.Lowercase, s.Config.Generator.UseLowercase),
		Digits:    getPointerBoolOrDefault(req.Digits, s.Config.Generator.UseDigits),
		Special:   getPointerBoolOrDefault(req.Special, s.Config.Generator.UseSpecial),
	}

	password, err := s.generatePassword(length, charSet)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, PasswordResponse{
		Password: password,
		Length:   length,
	})
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleDocumentation(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Password Generator API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; line-height: 1.6; }
        h1, h2, h3 { color: #333; }
        pre { background-color: #f5f5f5; padding: 10px; border-radius: 5px; overflow-x: auto; }
        code { background-color: #f5f5f5; padding: 2px 5px; border-radius: 3px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { text-align: left; padding: 8px; border-bottom: 1px solid #ddd; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Password Generator API</h1>
    <p>This API allows you to generate secure random passwords with customizable options.</p>
    
    <h2>Endpoints</h2>
    
    <h3>GET /api/v1/password</h3>
    <p>Generate a password using query parameters.</p>
    <h4>Query Parameters:</h4>
    <ul>
        <li><code>length</code>: Length of the password (default: configured default length)</li>
        <li><code>uppercase</code>: Include uppercase letters (true/false)</li>
        <li><code>lowercase</code>: Include lowercase letters (true/false)</li>
        <li><code>digits</code>: Include digits (true/false)</li>
        <li><code>special</code>: Include special characters (true/false)</li>
    </ul>
    <h4>Example:</h4>
    <pre>GET /api/v1/password?length=12&uppercase=true&lowercase=true&digits=true&special=false</pre>
    
    <h3>POST /api/v1/password</h3>
    <p>Generate a password using JSON request body.</p>
    <h4>Request Body:</h4>
    <pre>
{
  "length": 12,
  "uppercase": true,
  "lowercase": true,
  "digits": true,
  "special": false
}
    </pre>
    <h4>Response:</h4>
    <pre>
{
  "password": "generated-password",
  "length": 12
}
    </pre>
    
    <h3>GET /health</h3>
    <p>Health check endpoint.</p>
    <h4>Response:</h4>
    <pre>
{
  "status": "ok"
}
    </pre>
    
    <h2>Configuration</h2>
    <p>The server can be configured using a YAML configuration file:</p>
    <pre>
server:
  port: 8080
generator:
  defaultLength: 12
  minLength: 8
  maxLength: 64
  useUppercase: true
  useLowercase: true
  useDigits: true
  useSpecial: true
    </pre>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *Server) generatePassword(length int, charSet generator.CharacterSet) (string, error) {
	gen, err := generator.NewGenerator(s.Generator.MinLength, s.Generator.MaxLength, charSet)
	if err != nil {
		return "", err
	}
	
	return gen.Generate(length)
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func getBoolParam(r *http.Request, name string, defaultValue bool) bool {
	param := r.URL.Query().Get(name)
	if param == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(param)
	if err != nil {
		return defaultValue
	}
	return value
}

func getPointerBoolOrDefault(value *bool, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}
