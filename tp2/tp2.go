package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	TDAVISITANTES "tdas/visitantes"
)

func main() {
	salida := bufio.NewWriter(os.Stdout)
	defer salida.Flush()

	visitantes := TDAVISITANTES.CrearVisitantes()
	lector := bufio.NewScanner(os.Stdin)

	for lector.Scan() {
		campos := strings.Fields(lector.Text())

		if err := ejecutarComando(visitantes, campos, salida); err != nil {
			salida.Flush()
			fmt.Fprintln(os.Stderr, _ERROR_EN_COMANDO, nombreDe(campos))
			return
		}
	}
}
