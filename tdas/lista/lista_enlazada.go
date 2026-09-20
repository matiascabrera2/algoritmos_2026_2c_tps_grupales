package lista

type nodoLista[T any] struct {
	valor     T
	siguiente *nodoLista[T]
}

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

type iteradorExterno[T any] struct {
	lista    *listaEnlazada[T]
	actual   *nodoLista[T]
	anterior *nodoLista[T]
}

func CrearListaEnlazada[T any]() Lista[T] {
	return &listaEnlazada[T]{
		primero: nil,
		ultimo:  nil,
		largo:   0,
	}
}

func crearNodo[T any](valor T) *nodoLista[T] {
	return &nodoLista[T]{valor: valor}
}

func (lista *listaEnlazada[T]) EstaVacia() bool {
	return lista.largo == 0
}

func (lista *listaEnlazada[T]) InsertarPrimero(valor T) {
	nuevoNodo := crearNodo(valor)
	nuevoNodo.siguiente = lista.primero
	lista.primero = nuevoNodo

	if lista.EstaVacia() {
		lista.ultimo = nuevoNodo
	}

	lista.largo++
}

func (lista *listaEnlazada[T]) InsertarUltimo(valor T) {
	nuevoNodo := crearNodo(valor)

	if !lista.EstaVacia() {
		lista.ultimo.siguiente = nuevoNodo
	} else {
		lista.primero = nuevoNodo
	}

	lista.ultimo = nuevoNodo
	lista.largo++
}

func (lista *listaEnlazada[T]) BorrarPrimero() T {
	if lista.EstaVacia() {
		panic("La lista está vacía")
	}

	valor := lista.primero.valor
	lista.primero = lista.primero.siguiente

	if lista.primero == nil {
		lista.ultimo = nil
	}

	lista.largo--
	return valor
}

func (lista *listaEnlazada[T]) VerPrimero() T {
	if lista.EstaVacia() {
		panic("La lista está vacía")
	}

	return lista.primero.valor
}

func (lista *listaEnlazada[T]) VerUltimo() T {
	if lista.EstaVacia() {
		panic("La lista está vacía")
	}

	return lista.ultimo.valor
}

func (lista *listaEnlazada[T]) Largo() int {
	return lista.largo
}

func (lista *listaEnlazada[T]) Iterar(visitar func(T) bool) {
	actual := lista.primero

	for actual != nil {
		if !visitar(actual.valor) {
			break
		}
		actual = actual.siguiente
	}
}

func (lista *listaEnlazada[T]) Iterador() IteradorLista[T] {
	return &iteradorExterno[T]{
		lista:    lista,
		actual:   lista.primero,
		anterior: nil,
	}
}
