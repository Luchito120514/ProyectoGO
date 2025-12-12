package persistencia

import (
	"encoding/json"
	"os"
	"sistemas_libros/sistema"
)

func GuardarLibros(libros []sistema.Libro) {
	file, err := os.Create("data/libros.json")
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(libros)
}

func CargarLibros() []sistema.Libro {
	file, err := os.Open("data/libros.json")
	if err != nil {
		return []sistema.Libro{}
	}
	defer file.Close()

	var libros []sistema.Libro
	json.NewDecoder(file).Decode(&libros)
	return libros
}
