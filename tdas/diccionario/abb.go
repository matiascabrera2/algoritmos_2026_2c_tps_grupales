package diccionario

import (
	TDAPila "tdas/pila"
)

const (
	_PANIC_CLAVE_NO_PERTENECE = "La clave no pertenece al diccionario"
	_PANIC_FIN_ITERADOR       = "El iterador termino de iterar"
)

type nodoAbb[K comparable, V any] struct {
	izquierdo *nodoAbb[K, V]
	derecho   *nodoAbb[K, V]
	clave     K
	dato      V
}

type abb[K comparable, V any] struct {
	raiz     *nodoAbb[K, V]
	cantidad int
	cmp      func(K, K) int
}

type iteradorAbb[K comparable, V any] struct {
	pila  TDAPila.Pila[*nodoAbb[K, V]]
	hasta *K
	cmp   func(K, K) int
}

func CrearABB[K comparable, V any](funcion_cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: funcion_cmp}
}

// buscarPunteroAlNodo devuelve el puntero al nodo donde está la clave,
// o si no está, el enlace vacío donde habría que insertarla.
func (a *abb[K, V]) buscarPunteroAlNodo(enlace **nodoAbb[K, V], clave K) **nodoAbb[K, V] {
	if *enlace == nil {
		return enlace
	}

	comparacion := a.cmp(clave, (*enlace).clave)

	if comparacion < 0 {
		return a.buscarPunteroAlNodo(&(*enlace).izquierdo, clave)
	}
	if comparacion > 0 {
		return a.buscarPunteroAlNodo(&(*enlace).derecho, clave)
	}
	return enlace
}

func (a *abb[K, V]) Guardar(clave K, dato V) {
	enlace := a.buscarPunteroAlNodo(&a.raiz, clave)

	if *enlace == nil {
		*enlace = &nodoAbb[K, V]{clave: clave, dato: dato}
		a.cantidad++
	} else {
		(*enlace).dato = dato
	}
}

func (a *abb[K, V]) Pertenece(clave K) bool {
	return *a.buscarPunteroAlNodo(&a.raiz, clave) != nil
}

func (a *abb[K, V]) Obtener(clave K) V {
	nodo := *a.buscarPunteroAlNodo(&a.raiz, clave)
	if !a.Pertenece(clave) {
		panic(_PANIC_CLAVE_NO_PERTENECE)
	}
	return nodo.dato
}

func (a *abb[K, V]) Borrar(clave K) V {
	enlace := a.buscarPunteroAlNodo(&a.raiz, clave)
	nodo := *enlace
	if !a.Pertenece(clave) {
		panic(_PANIC_CLAVE_NO_PERTENECE)
	}

	dato := nodo.dato

	if nodo.izquierdo == nil {
		*enlace = nodo.derecho
	} else if nodo.derecho == nil {
		*enlace = nodo.izquierdo
	} else {
		// Dos hijos: el nodo pasa a tener la clave y el dato de su sucesor inmediato
		// (el mínimo del subárbol derecho), y se desengancha ese nodo del árbol.
		enlaceSucesor := &nodo.derecho
		for (*enlaceSucesor).izquierdo != nil {
			enlaceSucesor = &(*enlaceSucesor).izquierdo
		}
		sucesor := *enlaceSucesor

		nodo.clave, nodo.dato = sucesor.clave, sucesor.dato
		*enlaceSucesor = sucesor.derecho
	}

	a.cantidad--
	return dato
}

func (a *abb[K, V]) Cantidad() int {
	return a.cantidad
}

func (a *abb[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	a.IterarRango(nil, nil, visitar)
}

func (a *abb[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	a.iterarRango(a.raiz, desde, hasta, visitar)
}

// iterarRango recorre in-order el subárbol, sin entrar a las ramas que quedan fuera del rango.
// Devuelve false si la iteración se cortó (visitar devolvió false).
func (a *abb[K, V]) iterarRango(nodo *nodoAbb[K, V], desde, hasta *K, visitar func(K, V) bool) bool {
	if nodo == nil {
		return true
	}

	seguir := true

	if desde == nil || a.cmp(*desde, nodo.clave) < 0 {
		seguir = a.iterarRango(nodo.izquierdo, desde, hasta, visitar)
	}

	if seguir && a.enRango(nodo.clave, desde, hasta) {
		seguir = visitar(nodo.clave, nodo.dato)
	}

	if seguir && (hasta == nil || a.cmp(nodo.clave, *hasta) < 0) {
		seguir = a.iterarRango(nodo.derecho, desde, hasta, visitar)
	}

	return seguir
}

func (a *abb[K, V]) enRango(clave K, desde, hasta *K) bool {
	return (desde == nil || a.cmp(clave, *desde) >= 0) && (hasta == nil || a.cmp(clave, *hasta) <= 0)
}

func (a *abb[K, V]) Iterador() IterDiccionario[K, V] {
	return a.IteradorRango(nil, nil)
}

// IteradorRango deja en la pila el camino hasta la primera clave >= desde: el tope es siempre
// el próximo elemento a visitar.
func (a *abb[K, V]) IteradorRango(desde *K, hasta *K) IterDiccionario[K, V] {
	iter := &iteradorAbb[K, V]{
		pila:  TDAPila.CrearPilaDinamica[*nodoAbb[K, V]](),
		hasta: hasta,
		cmp:   a.cmp,
	}

	nodo := a.raiz

	for nodo != nil {
		if desde == nil || a.cmp(nodo.clave, *desde) >= 0 {
			iter.pila.Apilar(nodo)
			nodo = nodo.izquierdo
		} else {
			nodo = nodo.derecho
		}
	}

	return iter
}

func (it *iteradorAbb[K, V]) HayAlgoMas() bool {
	return !it.pila.EstaVacia() && (it.hasta == nil || it.cmp(it.pila.VerTope().clave, *it.hasta) <= 0)
}

func (it *iteradorAbb[K, V]) VerActual() (K, V) {
	if !it.HayAlgoMas() {
		panic(_PANIC_FIN_ITERADOR)
	}

	nodo := it.pila.VerTope()
	return nodo.clave, nodo.dato
}

func (it *iteradorAbb[K, V]) Avanzar() {
	if !it.HayAlgoMas() {
		panic(_PANIC_FIN_ITERADOR)
	}

	nodo := it.pila.Desapilar().derecho

	for nodo != nil {
		it.pila.Apilar(nodo)
		nodo = nodo.izquierdo
	}
}
