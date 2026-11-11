package visitantes

import (
	"bufio"
	"os"
	"strings"
	TDACola "tdas/cola"
	TDAColaPrioridad "tdas/cola_prioridad"
	TDADiccionario "tdas/diccionario"
	"time"
)

const (
	_FORMATO_FECHA    = "2006-01-02T15:04:05-07:00"
	_SEPARADOR_CAMPOS = "\t"
	_CANTIDAD_CAMPOS  = 4
	_CAMPO_IP         = 0
	_CAMPO_FECHA      = 1
	_CAMPO_RECURSO    = 3

	_PETICIONES_DOS = 5
	_VENTANA_DOS    = 2 * time.Second // se multiplica por time.Second para que la constante sea de tipo Duration y no int
)

type visitantesDoS struct {
	visitantes TDADiccionario.DiccionarioOrdenado[IP, struct{}]
	sitios     TDADiccionario.Diccionario[string, int]
}

type entradaLog struct {
	ip      IP
	momento time.Time
	recurso string
}

// peticionesRecientes guarda los momentos de las últimas peticiones de una IP,
// sin superar las que hacen falta para sospechar de un DoS
type peticionesRecientes struct {
	momentos TDACola.Cola[time.Time]
	cantidad int
}

func CrearVisitantes() Visitantes {
	return &visitantesDoS{
		visitantes: TDADiccionario.CrearABB[IP, struct{}](CompararIP),
		sitios:     TDADiccionario.CrearHash[string, int](),
	}
}

func (v *visitantesDoS) AgregarArchivo(ruta string) ([]IP, error) {
	archivo, err := os.Open(ruta)

	if err != nil {
		return nil, err
	}

	defer archivo.Close()

	recientes := TDADiccionario.CrearHash[IP, *peticionesRecientes]()
	sospechosas := TDADiccionario.CrearABB[IP, struct{}](CompararIP)

	lector := bufio.NewScanner(archivo)
	for lector.Scan() {
		entrada, err := parsearEntrada(lector.Text())

		if err != nil {
			return nil, err
		}

		v.sumarVisita(entrada)

		if peticionesDe(recientes, entrada.ip).agregar(entrada.momento) {
			sospechosas.Guardar(entrada.ip, struct{}{})
		}
	}

	if err := lector.Err(); err != nil {
		return nil, err
	}

	return ipsEnRango(sospechosas, nil, nil), nil
}

func (v *visitantesDoS) VerVisitantes(desde, hasta IP) []IP {
	return ipsEnRango(v.visitantes, &desde, &hasta)
}

func (v *visitantesDoS) VerMasVisitados(n int) []Sitio {
	sitios := make([]Sitio, 0, v.sitios.Cantidad())

	v.sitios.Iterar(func(ruta string, visitas int) bool {
		sitios = append(sitios, Sitio{Ruta: ruta, Visitas: visitas})

		return true
	})

	cola := TDAColaPrioridad.CrearHeapArr(sitios, compararVisitas)
	masVisitados := make([]Sitio, 0)
	for len(masVisitados) < n && !cola.EstaVacia() {
		masVisitados = append(masVisitados, cola.Desencolar())
	}

	return masVisitados
}

func (v *visitantesDoS) sumarVisita(entrada entradaLog) {
	v.visitantes.Guardar(entrada.ip, struct{}{})

	visitas := 0
	if v.sitios.Pertenece(entrada.recurso) {
		visitas = v.sitios.Obtener(entrada.recurso)
	}
	v.sitios.Guardar(entrada.recurso, visitas+1)
}

// agregar suma una petición y devuelve true si las últimas cinco ocurrieron en menos de dos segundos
func (recientes *peticionesRecientes) agregar(momento time.Time) bool {
	recientes.momentos.Encolar(momento)
	recientes.cantidad++

	if recientes.cantidad > _PETICIONES_DOS {
		recientes.momentos.Desencolar()
		recientes.cantidad--
	}

	return recientes.cantidad == _PETICIONES_DOS && momento.Sub(recientes.momentos.VerPrimero()) < _VENTANA_DOS
}

func peticionesDe(recientes TDADiccionario.Diccionario[IP, *peticionesRecientes], ip IP) *peticionesRecientes {
	if !recientes.Pertenece(ip) {
		recientes.Guardar(ip, &peticionesRecientes{momentos: TDACola.CrearColaEnlazada[time.Time]()})
	}

	return recientes.Obtener(ip)
}

func parsearEntrada(linea string) (entradaLog, error) {
	campos := strings.Split(linea, _SEPARADOR_CAMPOS)
	if len(campos) != _CANTIDAD_CAMPOS {
		return entradaLog{}, ErrLogInvalido
	}

	ip, errIP := ParsearIP(campos[_CAMPO_IP])
	momento, errFecha := time.Parse(_FORMATO_FECHA, campos[_CAMPO_FECHA])

	if errIP != nil || errFecha != nil {
		return entradaLog{}, ErrLogInvalido
	}

	return entradaLog{ip: ip, momento: momento, recurso: campos[_CAMPO_RECURSO]}, nil
}

func ipsEnRango(dic TDADiccionario.DiccionarioOrdenado[IP, struct{}], desde, hasta *IP) []IP {
	ips := make([]IP, 0)
	dic.IterarRango(desde, hasta, func(ip IP, _ struct{}) bool {
		ips = append(ips, ip)
		return true
	})

	return ips
}

func compararVisitas(a, b Sitio) int {
	return a.Visitas - b.Visitas
}
