package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LeerEntrada(mensaje string) string {
	fmt.Print(mensaje)
	reader := bufio.NewReader(os.Stdin)
	texto, _ := reader.ReadString('\n')
	return strings.TrimSpace(texto)
}
