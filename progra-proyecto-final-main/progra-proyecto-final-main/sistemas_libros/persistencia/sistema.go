package persistencia

import "sistemas_libros/sistema"

func GuardarSistema(s *sistema.Sistema) {
	GuardarLibros(s.Libros)
	GuardarUsuarios(s.Usuarios)
	GuardarPrestamos(s.Prestamos)
}

func CargarSistema() *sistema.Sistema {
	return &sistema.Sistema{
		Libros:     CargarLibros(),
		Usuarios:   CargarUsuarios(),
		Prestamos:  CargarPrestamos(),
		NextUserID: calcularSiguienteID(CargarUsuarios()),
	}
}

func calcularSiguienteID(usuarios []sistema.Usuario) int {
	maxID := 0
	for _, u := range usuarios {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	return maxID + 1
}
