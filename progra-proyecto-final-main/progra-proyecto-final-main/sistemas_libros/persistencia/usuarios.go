package persistencia

import (
	"encoding/json"
	"os"
	"sistemas_libros/sistema"
)

func GuardarUsuarios(usuarios []sistema.Usuario) {
	file, err := os.Create("data/usuarios.json")
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(usuarios)
}

func CargarUsuarios() []sistema.Usuario {
	file, err := os.Open("data/usuarios.json")
	if err != nil {
		return []sistema.Usuario{}
	}
	defer file.Close()

	var usuarios []sistema.Usuario
	json.NewDecoder(file).Decode(&usuarios)
	return usuarios
}
