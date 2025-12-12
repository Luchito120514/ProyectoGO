package usuarios

import (
	"fmt"
	"sistemas_libros/sistema"
	"sistemas_libros/utils"
)

func RegistrarUsuario(s *sistema.Sistema) {
	nombre := utils.LeerEntrada("Nombre del usuario: ")
	email := utils.LeerEntrada("Email del usuario: ")

	if nombre == "" || email == "" {
		fmt.Println("Debe ingresar nombre y correo.")
		return
	}

	u := sistema.Usuario{
		ID:     s.NextUserID,
		Nombre: nombre,
		Email:  email,
	}
	s.Usuarios = append(s.Usuarios, u)
	s.NextUserID++
	fmt.Printf("Usuario '%s' registrado con ID %d.\n", nombre, u.ID)
}
