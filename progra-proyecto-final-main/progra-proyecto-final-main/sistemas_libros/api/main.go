package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
)

type Libro struct {
	ISBN       string `json:"isbn"`
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	Categoria  string `json:"categoria"`
	Disponible bool   `json:"disponible"`
}

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Rol    string `json:"rol"`
}

type Prestamo struct {
	ISBN          string    `json:"isbn"`
	UsuarioID     int       `json:"usuario_id"`
	Activo        bool      `json:"activo"`
	Devuelto      bool      `json:"devuelto"`
	FechaPrestamo time.Time `json:"fecha_prestamo"`
}

// Slices
var libros []Libro
var usuarios []Usuario
var prestamos []Prestamo

// Persistencia de datos
func guardarDatos() {
	save("libros.json", libros)
	save("usuarios.json", usuarios)
	save("prestamos.json", prestamos)
}

func cargarDatos() {
	load("libros.json", &libros)
	load("usuarios.json", &usuarios)
	load("prestamos.json", &prestamos)
}

func save(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creando", filename, ":", err)
		return
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(data)
}

func load(filename string, target interface{}) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(target)
}

// Helper
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func nextUsuarioID() int {
	max := 0
	for _, u := range usuarios {
		if u.ID > max {
			max = u.ID
		}
	}
	return max + 1
}

// Conectores
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API Biblioteca conectada a carpeta data/")
}

// Libros
func getLibros(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, libros)
}

func postLibros(w http.ResponseWriter, r *http.Request) {
	var libro Libro
	if err := json.NewDecoder(r.Body).Decode(&libro); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	libro.ISBN = strings.TrimSpace(libro.ISBN)
	libro.Titulo = strings.TrimSpace(libro.Titulo)
	if libro.ISBN == "" || libro.Titulo == "" {
		http.Error(w, "ISBN y Título obligatorios", http.StatusBadRequest)
		return
	}
	for _, l := range libros {
		if l.ISBN == libro.ISBN {
			http.Error(w, "ISBN duplicado", http.StatusBadRequest)
			return
		}
	}
	libro.Disponible = true
	libros = append(libros, libro)
	guardarDatos()
	writeJSON(w, http.StatusCreated, libro)
}

// Eliminar Libros
func resetLibrosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	libros = []Libro{}
	guardarDatos()
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "Todos los libros han sido eliminados"})
}

func getUsuarios(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, usuarios)
}

func postUsuarios(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	usuario.Nombre = strings.TrimSpace(usuario.Nombre)
	usuario.Email = strings.TrimSpace(usuario.Email)
	if usuario.Nombre == "" || usuario.Email == "" {
		http.Error(w, "Nombre y Email obligatorios", http.StatusBadRequest)
		return
	}
	for _, u := range usuarios {
		if u.Email == usuario.Email {
			http.Error(w, "Email duplicado", http.StatusBadRequest)
			return
		}
	}
	usuario.ID = nextUsuarioID()
	if usuario.Rol == "" {
		usuario.Rol = "usuario"
	}
	usuarios = append(usuarios, usuario)
	guardarDatos()
	writeJSON(w, http.StatusCreated, usuario)
}

func deleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	for i, u := range usuarios {
		if u.ID == id {
			usuarios = append(usuarios[:i], usuarios[i+1:]...)
			guardarDatos()
			writeJSON(w, http.StatusOK, map[string]string{"mensaje": "Usuario eliminado"})
			return
		}
	}
	http.Error(w, "Usuario no encontrado", http.StatusNotFound)
}

// Prestamo
func getPrestamos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, prestamos)
}

func postPrestamos(w http.ResponseWriter, r *http.Request) {
	var prestamo Prestamo
	if err := json.NewDecoder(r.Body).Decode(&prestamo); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if prestamo.ISBN == "" || prestamo.UsuarioID == 0 {
		http.Error(w, "ISBN y usuario_id obligatorios", http.StatusBadRequest)
		return
	}

	for i := range libros {
		if libros[i].ISBN == prestamo.ISBN {
			if !libros[i].Disponible {
				http.Error(w, "Libro no disponible", http.StatusBadRequest)
				return
			}
			libros[i].Disponible = false
			break
		}
	}
	prestamo.FechaPrestamo = time.Now()
	prestamo.Activo = true
	prestamo.Devuelto = false
	prestamos = append(prestamos, prestamo)
	guardarDatos()
	writeJSON(w, http.StatusCreated, prestamo)
}

func putPrestamosDevolver(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	for i := range prestamos {
		if prestamos[i].ISBN == isbn && prestamos[i].Activo {
			prestamos[i].Activo = false
			prestamos[i].Devuelto = true

			for j := range libros {
				if libros[j].ISBN == isbn {
					libros[j].Disponible = true
					break
				}
			}
			guardarDatos()
			writeJSON(w, http.StatusOK, map[string]string{"mensaje": "Préstamo devuelto", "isbn": isbn})
			return
		}
	}
	http.Error(w, "Préstamo no encontrado", http.StatusNotFound)
}

// Inicio Sesion
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
	}
	_ = json.NewDecoder(r.Body).Decode(&data)
	for _, u := range usuarios {
		if (data.Nombre != "" && strings.EqualFold(u.Nombre, data.Nombre)) ||
			(data.Email != "" && strings.EqualFold(u.Email, data.Email)) {
			writeJSON(w, http.StatusOK, map[string]string{
				"mensaje": "Sesión iniciada",
				"usuario": u.Nombre,
				"rol":     u.Rol,
			})
			return
		}
	}
	http.Error(w, "Usuario no encontrado", http.StatusNotFound)
}

func main() {
	cargarDatos()

	if len(usuarios) == 0 {
		usuarios = []Usuario{
			{ID: 1, Nombre: "Matías", Email: "matias@gmail.com", Rol: "usuario"},
			{ID: 2, Nombre: "Luis", Email: "luis@gmail.com", Rol: "admin"},
		}
		guardarDatos()
	}

	http.HandleFunc("/", rootHandler)

	// Libros
	http.HandleFunc("/libros", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getLibros(w, r)
		case http.MethodPost:
			postLibros(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/libros/reset", resetLibrosHandler)

	// Usuarios
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsuarios(w, r)
		case http.MethodPost:
			postUsuarios(w, r)
		case http.MethodDelete:
			deleteUsuarioHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Prestamo
	http.HandleFunc("/prestamos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPrestamos(w, r)
		case http.MethodPost:
			postPrestamos(w, r)
		case http.MethodPut:
			putPrestamosDevolver(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Inicio Sesion
	http.HandleFunc("/login", loginHandler)

	// Configuración CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5500", "http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowCredentials: true,
	})
	handler := c.Handler(http.DefaultServeMux)

	fmt.Println("Servidor escuchando en http://localhost:8080")
	_ = http.ListenAndServe(":8080", handler)
}
