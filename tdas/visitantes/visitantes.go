package visitantes

import "errors"

var ErrLogInvalido = errors.New("El log tiene una linea invalida")

type Sitio struct {
	Ruta    string
	Visitas int
}

type Visitantes interface {

	// AgregarArchivo procesa de forma completa el archivo de log de la ruta indicada y suma sus peticiones a las
	// ya registradas. Devuelve, en orden creciente, las IPs sospechosas de DoS: las que hicieron cinco o más
	// peticiones en menos de dos segundos dentro de ese archivo (cada archivo se analiza de forma independiente).
	// Si el archivo no se puede abrir ni leer, o tiene una línea inválida, devuelve error (ErrLogInvalido en el
	// último caso) y las líneas anteriores a la inválida quedan registradas.
	AgregarArchivo(ruta string) ([]IP, error)

	// VerVisitantes devuelve, en orden creciente, las IPs que solicitaron algún recurso en los archivos agregados
	// hasta el momento y que están entre desde y hasta, ambas inclusive.
	VerVisitantes(desde, hasta IP) []IP

	// VerMasVisitados devuelve los n recursos más solicitados en los archivos agregados hasta el momento, de mayor
	// a menor cantidad de visitas (en caso de empate, en cualquier orden). Si hay menos de n recursos devuelve
	// todos, y si n es menor o igual a 0 devuelve una lista vacía.
	VerMasVisitados(n int) []Sitio
}
