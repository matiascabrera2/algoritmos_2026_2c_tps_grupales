package diccionario

import (
	"fmt"
	"hash/fnv"
	TDALista "tdas/lista"
)

const (
	CAPACIDAD_INICIAL  = 8
	FACTOR_REDIMENSION = 2
	FACTOR_AGRANDAR    = 2.5
	FACTOR_ACHICAR     = 0.25
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
		tabla: crearTabla[K, V](CAPACIDAD_INICIAL),
		tam:   CAPACIDAD_INICIAL,
		cant:  0,
	}
}

// Devuelve el iterador de la lista parado en el
// elemento con la clave buscada, junto con true, o false si no está.
func posicionarEnClave[K comparable, V any](bucket TDALista.Lista[par[K, V]], clave K) (TDALista.IteradorLista[par[K, V]], bool) {
	iter := bucket.Iterador()

	for iter.HayAlgoMas() {
		if iter.VerActual().clave == clave {
			return iter, true
		}
		iter.Avanzar()
	}

	return iter, false
}

func (h *hashAbierto[K, V]) Guardar(clave K, dato V) {
	bucket := h.tabla[posicion(clave, h.tam)]

	if iter, encontrado := posicionarEnClave(bucket, clave); encontrado {
		iter.Borrar()
	} else {
		h.cant++
	}

	bucket.InsertarPrimero(par[K, V]{clave: clave, dato: dato})

	if float64(h.cant)/float64(h.tam) >= FACTOR_AGRANDAR {
		h.redimensionar(h.tam * FACTOR_REDIMENSION)
	}
}

func (h *hashAbierto[K, V]) Pertenece(clave K) bool {
	_, encontrado := posicionarEnClave(h.tabla[posicion(clave, h.tam)], clave)
	return encontrado
}

func (h *hashAbierto[K, V]) Obtener(clave K) V {
	iter, encontrado := posicionarEnClave(h.tabla[posicion(clave, h.tam)], clave)
	if !encontrado {
		panic("La clave no pertenece al diccionario")
	}
	return iter.VerActual().dato
}

func (h *hashAbierto[K, V]) Borrar(clave K) V {
	bucket := h.tabla[posicion(clave, h.tam)]

	iter, encontrado := posicionarEnClave(bucket, clave)
	if !encontrado {
		panic("La clave no pertenece al diccionario")
	}

	dato := iter.VerActual().dato
	iter.Borrar()
	h.cant--

	nuevoTam := h.tam / FACTOR_REDIMENSION
	if float64(h.cant)/float64(h.tam) <= FACTOR_ACHICAR && nuevoTam >= CAPACIDAD_INICIAL {
		h.redimensionar(nuevoTam)
	}

	return dato
}

func (h *hashAbierto[K, V]) Cantidad() int {
	return h.cant
}

func (h *hashAbierto[K, V]) redimensionar(nuevoTam int) {
	nuevaTabla := crearTabla[K, V](nuevoTam)

	for _, bucket := range h.tabla {
		bucket.Iterar(func(p par[K, V]) bool {
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
		panic("El iterador termino de iterar")
	}
	p := iterador.iterLista.VerActual()
	return p.clave, p.dato
}

func (iterador *iteradorExterno[K, V]) Avanzar() {
	if !iterador.HayAlgoMas() {
		panic("El iterador termino de iterar")
	}

	iterador.iterLista.Avanzar()

	if !iterador.iterLista.HayAlgoMas() {
		iterador.iterLista = nil
		iterador.pos++
	}

	iterador.avanzarHastaElementoValido()
}
