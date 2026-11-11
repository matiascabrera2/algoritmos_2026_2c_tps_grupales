package visitantes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	_CANTIDAD_BLOQUES  = 4
	_SEPARADOR_BLOQUES = "."
)

var ErrIPInvalida = errors.New("La IP no es valida")

// IP es una dirección IPv4: cuatro bloques de un byte, del más al menos significativo
type IP [_CANTIDAD_BLOQUES]uint8

// ParsearIP convierte un texto con formato a.b.c.d en una IP. Si el texto no tiene
// exactamente cuatro bloques numéricos entre 0 y 255, devuelve ErrIPInvalida
func ParsearIP(texto string) (IP, error) {
	bloques := strings.Split(texto, _SEPARADOR_BLOQUES)
	if len(bloques) != _CANTIDAD_BLOQUES {
		return IP{}, ErrIPInvalida
	}

	var ip IP
	for i, bloque := range bloques {
		numero, err := strconv.ParseUint(bloque, 10, 8)
		if err != nil {
			return IP{}, ErrIPInvalida
		}
		ip[i] = uint8(numero)
	}
	return ip, nil
}

// CompararIP devuelve un entero menor que 0 si a es anterior a b, mayor que 0 si es posterior
// y 0 si son iguales, comparando bloque a bloque por valor numérico
func CompararIP(a, b IP) int {
	for i := range a {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return 0
}

func (ip IP) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}
