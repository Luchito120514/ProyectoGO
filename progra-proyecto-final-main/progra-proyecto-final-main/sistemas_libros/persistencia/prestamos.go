package persistencia

import (
	"encoding/json"
	"os"
	"sistemas_libros/sistema"
)

func GuardarPrestamos(prestamos []sistema.Prestamo) {
	file, err := os.Create("data/prestamos.json")
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(prestamos)
}

func CargarPrestamos() []sistema.Prestamo {
	file, err := os.Open("data/prestamos.json")
	if err != nil {
		return []sistema.Prestamo{}
	}
	defer file.Close()

	var prestamos []sistema.Prestamo
	json.NewDecoder(file).Decode(&prestamos)
	return prestamos
}
