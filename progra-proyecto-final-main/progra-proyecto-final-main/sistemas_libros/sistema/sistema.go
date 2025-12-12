package sistema

import (
	"encoding/json"
	"strconv"
	"time"
)

type Sistema struct {
	Libros     []Libro
	Usuarios   []Usuario
	Prestamos  []Prestamo
	NextUserID int
}

type Entidad interface {
	GetID() string
	ToJSON() string
}

type Libro struct {
	ISBN       string `json:"isbn"`
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	Categoria  string `json:"categoria"`
	Disponible bool   `json:"disponible"`
}

func (l Libro) GetID() string {
	return l.ISBN
}

func (l Libro) ToJSON() string {
	data, _ := json.Marshal(l)
	return string(data)
}

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

func (u Usuario) GetID() string {
	return strconv.Itoa(u.ID)
}

func (u Usuario) ToJSON() string {
	data, _ := json.Marshal(u)
	return string(data)
}

type Prestamo struct {
	ISBN          string    `json:"isbn"`
	UsuarioID     int       `json:"usuario_id"`
	Activo        bool      `json:"activo"`
	Devuelto      bool      `json:"devuelto"`
	FechaPrestamo time.Time `json:"fecha_prestamo"`
}

func (p Prestamo) GetID() string {
	return p.ISBN
}

func (p Prestamo) ToJSON() string {
	data, _ := json.Marshal(p)
	return string(data)
}
