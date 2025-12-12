package libros

import (
	"fmt"
	"sistemas_libros/sistema"
	"sistemas_libros/utils"
	"strings"
	"time"
)

func RegistrarLibro(s *sistema.Sistema) {
	isbn := utils.LeerEntrada("ISBN: ")
	titulo := utils.LeerEntrada("Título: ")
	autor := utils.LeerEntrada("Autor: ")
	categoria := utils.LeerEntrada("Categoría: ")

	if isbn == "" || titulo == "" || autor == "" {
		fmt.Println("El ISBN, título y autor son obligatorios.")
		return
	}

	for _, l := range s.Libros {
		if l.ISBN == isbn {
			fmt.Println("Ya existe un libro con ese ISBN.")
			return
		}
	}

	nuevo := sistema.Libro{
		ISBN:       isbn,
		Titulo:     titulo,
		Autor:      autor,
		Categoria:  categoria,
		Disponible: true,
	}
	s.Libros = append(s.Libros, nuevo)
	fmt.Println("Libro agregado correctamente.")
}

func MenuGestionLibros(s *sistema.Sistema) {
	for {
		fmt.Println("\nGestión de libros")
		fmt.Println("1. Listar libros")
		fmt.Println("2. Buscar libro")
		fmt.Println("3. Actualizar datos")
		fmt.Println("4. Eliminar libro")
		fmt.Println("5. Volver")

		op := utils.LeerEntrada("Elija una opción: ")

		switch op {
		case "1":
			ListarLibros(s)
		case "2":
			busq := utils.LeerEntrada("Ingrese palabra clave del título: ")
			BuscarLibro(s, busq)
		case "3":
			isbn := utils.LeerEntrada("Ingrese ISBN del libro a modificar: ")
			ActualizarLibro(s, isbn)
		case "4":
			isbn := utils.LeerEntrada("Ingrese ISBN del libro a eliminar: ")
			EliminarLibro(s, isbn)
		case "5":
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

func ListarLibros(s *sistema.Sistema) {
	if len(s.Libros) == 0 {
		fmt.Println("No hay libros registrados.")
		return
	}

	fmt.Println("\nCatálogo de libros")
	for _, l := range s.Libros {
		estado := "Disponible"
		if !l.Disponible {
			estado = "Prestado"
		}
		fmt.Printf("ISBN: %s | Título: %s | Autor: %s | Categoría: %s | Estado: %s\n",
			l.ISBN, l.Titulo, l.Autor, l.Categoria, estado)
	}
}

func BuscarLibro(s *sistema.Sistema, palabra string) {
	palabra = strings.ToLower(strings.TrimSpace(palabra))
	if palabra == "" {
		fmt.Println("Debe ingresar un término de búsqueda.")
		return
	}

	encontrado := false
	for _, l := range s.Libros {
		if strings.Contains(strings.ToLower(l.Titulo), palabra) {
			fmt.Printf("ISBN: %s | Título: %s | Autor: %s | Categoría: %s\n",
				l.ISBN, l.Titulo, l.Autor, l.Categoria)
			encontrado = true
		}
	}

	if !encontrado {
		fmt.Println("No se encontraron libros con ese título.")
	}
}

func ActualizarLibro(s *sistema.Sistema, isbn string) {
	for i := range s.Libros {
		if s.Libros[i].ISBN == isbn {
			nuevoTitulo := utils.LeerEntrada("Nuevo título (enter para dejar igual): ")
			if nuevoTitulo != "" {
				s.Libros[i].Titulo = nuevoTitulo
			}
			nuevoAutor := utils.LeerEntrada("Nuevo autor (enter para dejar igual): ")
			if nuevoAutor != "" {
				s.Libros[i].Autor = nuevoAutor
			}
			nuevaCat := utils.LeerEntrada("Nueva categoría (enter para dejar igual): ")
			if nuevaCat != "" {
				s.Libros[i].Categoria = nuevaCat
			}
			fmt.Println("Datos actualizados correctamente.")
			return
		}
	}
	fmt.Println("No se encontró el libro con ese ISBN.")
}

func EliminarLibro(s *sistema.Sistema, isbn string) {
	for i, l := range s.Libros {
		if l.ISBN == isbn {
			s.Libros = append(s.Libros[:i], s.Libros[i+1:]...)
			fmt.Println("Libro eliminado correctamente.")
			return
		}
	}
	fmt.Println("No se encontró un libro con ese ISBN.")
}

func MenuPrestamos(s *sistema.Sistema) {
	for {
		fmt.Println("\nGestión de préstamos")
		fmt.Println("1. Prestar libro")
		fmt.Println("2. Listar préstamos")
		fmt.Println("3. Devolver libro")
		fmt.Println("4. Volver")

		op := utils.LeerEntrada("Seleccione una opción: ")

		switch op {
		case "1":
			isbn := utils.LeerEntrada("ISBN del libro: ")
			idStr := utils.LeerEntrada("ID del usuario: ")
			var id int
			fmt.Sscan(idStr, &id)
			PrestarLibro(s, isbn, id)
		case "2":
			ListarPrestamos(s)
		case "3":
			isbn := utils.LeerEntrada("ISBN del libro a devolver: ")
			DevolverLibro(s, isbn)
		case "4":
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

func PrestarLibro(s *sistema.Sistema, isbn string, usuarioID int) {
	for i := range s.Libros {
		if s.Libros[i].ISBN == isbn {
			if !s.Libros[i].Disponible {
				fmt.Println("El libro ya está prestado.")
				return
			}
			var usuarioExiste bool
			for _, u := range s.Usuarios {
				if u.ID == usuarioID {
					usuarioExiste = true
					break
				}
			}
			if !usuarioExiste {
				fmt.Println("Usuario no encontrado.")
				return
			}
			s.Libros[i].Disponible = false
			prestamo := sistema.Prestamo{
				ISBN:          isbn,
				UsuarioID:     usuarioID,
				FechaPrestamo: time.Now(),
				Devuelto:      false,
			}
			s.Prestamos = append(s.Prestamos, prestamo)
			fmt.Println("Libro prestado correctamente.")
			return
		}
	}
	fmt.Println("No se encontró el libro.")
}

func ListarPrestamos(s *sistema.Sistema) {
	if len(s.Prestamos) == 0 {
		fmt.Println("No hay préstamos registrados.")
		return
	}

	fmt.Println("\nPréstamos activos")
	for _, p := range s.Prestamos {
		if !p.Devuelto {
			var usuarioNombre string
			for _, u := range s.Usuarios {
				if u.ID == p.UsuarioID {
					usuarioNombre = u.Nombre
				}
			}
			fmt.Printf("Libro ISBN: %s | Usuario: %s | Fecha: %s\n",
				p.ISBN, usuarioNombre, p.FechaPrestamo.Format("01-07-2022 18:24"))
		}
	}
}

func DevolverLibro(s *sistema.Sistema, isbn string) {
	for i := range s.Prestamos {
		if s.Prestamos[i].ISBN == isbn && !s.Prestamos[i].Devuelto {
			s.Prestamos[i].Devuelto = true
			for j := range s.Libros {
				if s.Libros[j].ISBN == isbn {
					s.Libros[j].Disponible = true
				}
			}
			fmt.Println("Libro devuelto correctamente.")
			return
		}
	}
	fmt.Println("No se encontró un préstamo activo con ese ISBN.")
}
