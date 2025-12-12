package main

import (
	"fmt"
	"sistemas_libros/libros"
	"sistemas_libros/persistencia"
	"sistemas_libros/usuarios"
	"sistemas_libros/utils"
)

func main() {
	s := persistencia.CargarSistema()

	for {
		fmt.Println("\n===== MENÚ PRINCIPAL =====")
		fmt.Println("1. Registro de libros")
		fmt.Println("2. Gestión de libros")
		fmt.Println("3. Registro de usuarios")
		fmt.Println("4. Préstamos")
		fmt.Println("5. Salir")

		op := utils.LeerEntrada("Seleccione una opción: ")

		switch op {
		case "1":
			libros.RegistrarLibro(s)
		case "2":
			libros.MenuGestionLibros(s)
		case "3":
			usuarios.RegistrarUsuario(s)
		case "4":
			libros.MenuPrestamos(s)
		case "5":
			persistencia.GuardarSistema(s)
			fmt.Println("Saliendo del sistema...")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}
