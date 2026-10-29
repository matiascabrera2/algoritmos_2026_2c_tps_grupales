package diccionario

import (
	"fmt"
	"hash/fnv"
	TDALista "tdas/lista"
)

const (
	_CAPACIDAD_INICIAL  = 8
	_FACTOR_REDIMENSION = 2
	_FACTOR_AGRANDAR    = 2.5
	_FACTOR_ACHICAR     = 0.25

	_PANIC_CLAVE_NO_PERTENECE = "La clave no pertenece al diccionario"
	_PANIC_FIN_ITERADOR       = "El iterador termino de iterar"
)

type par[K comparable, V any] struct {
	clave K
	dato  V
}

type hashAbierto[K comparable, V any] struct {
	tabla []TDALista.Lista[par[K, V]]
	tam   int
	cant  int
}

type iteradorExterno[K comparable, V any] struct {
	hash      *hashAbierto[K, V]
	pos       int
	iterLista TDALista.IteradorLista[par[K, V]]
}

func convertirAbytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

// FNV-1a sacada de la librería standar de go
func funcionDeHash[K comparable](clave K) uint64 {
	h := fnv.New64a()
	h.Write(convertirAbytes(clave)) // Para convertir la clave a bytes
	return h.Sum64()
}

func posicion[K comparable](clave K, tam int) int {
	return int(funcionDeHash(clave) % uint64(tam))
}

func crearTabla[K comparable, V any](tam int) []TDALista.Lista[par[K, V]] {
	tabla := make([]TDALista.Lista[par[K, V]], tam)
	for i := range tabla {
		tabla[i] = TDALista.CrearListaEnlazada[par[K, V]]()
	}
	return tabla
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	return &hashAbierto[K, V]{
		tabla: crearTabla[K, V](_CAPACIDAD_INICIAL),
		tam:   _CAPACIDAD_INICIAL,
		cant:  0,
	}
}

// Devuelve el iterador de la lista parado en el
// elemento con la clave buscada, junto con true, o false si no está.
func posicionarEnClave[K comparable, V any](hash *hashAbierto[K, V], clave K) (TDALista.IteradorLista[par[K, V]], bool) {
	lista := hash.tabla[posicion(clave, hash.tam)]
	iter := lista.Iterador()

	for iter.HayAlgoMas() {
		if iter.VerActual().clave == clave {
			return iter, true
		}
		iter.Avanzar()
	}

	return iter, false
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	iter, encontrado := posicionarEnClave(h, clave)

	if encontrado {
		iter.Borrar()
	} else {
		h.cant++
	}
	iter.Insertar(par[K, V]{clave: clave, dato: dato})

	if float64(h.cant)/float64(h.tam) >= _FACTOR_AGRANDAR {
		h.redimensionar(h.tam * _FACTOR_REDIMENSION)
	}
}

func (h *hashAbierto[K, V]) Pertenece(clave K) bool {
	_, encontrado := posicionarEnClave(h, clave)
	return encontrado
}

func (h *hashAbierto[K, V]) Obtener(clave K) V {
	iter, encontrado := posicionarEnClave(h, clave)
	if !encontrado {
		panic(_PANIC_CLAVE_NO_PERTENECE)
	}
	return iter.VerActual().dato
}

func (h *hashAbierto[K, V]) Borrar(clave K) V {
	iter, encontrado := posicionarEnClave(h, clave)
	if !encontrado {
		panic(_PANIC_CLAVE_NO_PERTENECE)
	}

	dato := iter.Borrar().dato
	h.cant--

	nuevoTam := h.tam / _FACTOR_REDIMENSION
	if float64(h.cant)/float64(h.tam) <= _FACTOR_ACHICAR && nuevoTam >= _CAPACIDAD_INICIAL {
		h.redimensionar(nuevoTam)
	}

	return dato
}

func (h *hashAbierto[K, V]) Cantidad() int {
	return h.cant
}

func (h *hashAbierto[K, V]) redimensionar(nuevoTam int) {
	nuevaTabla := crearTabla[K, V](nuevoTam)

	for _, lista := range h.tabla {
		lista.Iterar(func(p par[K, V]) bool {
			nuevaTabla[posicion(p.clave, nuevoTam)].InsertarPrimero(p)
			return true
		})
	}

	h.tabla = nuevaTabla
	h.tam = nuevoTam
}

func (h *hashAbierto[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	seguir := true

	for i := 0; i < len(h.tabla) && seguir; i++ {
		h.tabla[i].Iterar(func(p par[K, V]) bool {
			seguir = visitar(p.clave, p.dato)
			return seguir
		})
	}
}

func (h *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {
	iter := &iteradorExterno[K, V]{hash: h, pos: 0}
	iter.avanzarHastaElementoValido()
	return iter
}

func (iterador *iteradorExterno[K, V]) avanzarHastaElementoValido() {
	for iterador.pos < len(iterador.hash.tabla) {
		if iterador.iterLista == nil {
			iterador.iterLista = iterador.hash.tabla[iterador.pos].Iterador()
		}

		if iterador.iterLista.HayAlgoMas() {
			return
		}

		iterador.iterLista = nil
		iterador.pos++
	}
}

func (iterador *iteradorExterno[K, V]) HayAlgoMas() bool {
	return iterador.pos < len(iterador.hash.tabla)
}

func (iterador *iteradorExterno[K, V]) VerActual() (K, V) {
	if !iterador.HayAlgoMas() {
		panic(_PANIC_FIN_ITERADOR)
	}
	p := iterador.iterLista.VerActual()
	return p.clave, p.dato
}

func (iterador *iteradorExterno[K, V]) Avanzar() {
	if !iterador.HayAlgoMas() {
		panic(_PANIC_FIN_ITERADOR)
	}

	iterador.iterLista.Avanzar()
	iterador.avanzarHastaElementoValido()
}
