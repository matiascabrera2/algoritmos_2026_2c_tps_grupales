package lista

type Lista[T any] interface {

	// EstaVacia devuelve true si la lista no tiene elementos, false en caso contrario.
	EstaVacia() bool

	// InsertarPrimero agrega un elemento al principio de la lista.
	// Si la lista está vacía, el elemento agregado será también el último.
	InsertarPrimero(T)

	// InsertarUltimo agrega un elemento al final de la lista.
	// Si la lista está vacía, el elemento agregado será también el primero.
	InsertarUltimo(T)

	// BorrarPrimero saca el primer elemento de la lista y lo devuelve. Si la lista
	// está vacía, entra en pánico con un mensaje "La lista esta vacia".
	BorrarPrimero() T

	// VerPrimero devuelve el valor del primer elemento de la lista. Si la lista
	// está vacía, entra en pánico con un mensaje "La lista esta vacia".
	VerPrimero() T

	// VerUltimo devuelve el valor del último elemento de la lista. Si la lista
	// está vacía, entra en pánico con un mensaje "La lista esta vacia".
	VerUltimo() T

	// Largo devuelve la cantidad de elementos de la lista.
	Largo() int

	// Iterar aplica la función visitar a cada elemento de la lista, de primero a
	// último, hasta recorrerla entera o hasta que visitar devuelva false.
	Iterar(visitar func(T) bool)

	// Iterador devuelve un iterador externo de tipo IteradorLista para recorrer la lista de primero a último.
	Iterador() IteradorLista[T]
}

type IteradorLista[T any] interface {

	// VerActual devuelve el elemento en la posición actual del iterador. Si ya
	// terminó de iterar, entra en pánico con un mensaje "El iterador termino de iterar".
	VerActual() T

	// HayAlgoMas devuelve true si el iterador todavía no terminó de iterar.
	HayAlgoMas() bool

	// Avanzar mueve el iterador a la siguiente posición. Si ya terminó de iterar,
	// entra en pánico con un mensaje "El iterador termino de iterar".
	Avanzar()

	// Insertar agrega un elemento en la posición actual del iterador, y el
	// iterador queda parado en el elemento recién insertado.
	Insertar(T)

	// Borrar saca el elemento en la posición actual y lo devuelve. El iterador
	// queda parado en el elemento siguiente. Si ya terminó de iterar, entra en
	// pánico con un mensaje "El iterador termino de iterar".
	Borrar() T
}
