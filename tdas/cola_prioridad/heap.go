package cola_prioridad

const (
	_PANIC_COLA_VACIA = "La cola esta vacia"

	TAM_INICIAL      = 10
	FACTOR_AUMENTO   = 2
	FACTOR_REDUCCION = 4
)

type heap[T any] struct {
	datos    []T
	cantidad int
	cmp      func(T, T) int
}

func CrearHeap[T any](funcion_cmp func(T, T) int) ColaPrioridad[T] {
	return &heap[T]{datos: make([]T, TAM_INICIAL), cmp: funcion_cmp}
}

func CrearHeapArr[T any](arreglo []T, funcion_cmp func(T, T) int) ColaPrioridad[T] {
	capacidad := len(arreglo)
	if capacidad < TAM_INICIAL {
		capacidad = TAM_INICIAL
	}

	datos := make([]T, capacidad)
	copy(datos, arreglo)
	heapify(datos[:len(arreglo)], funcion_cmp)

	return &heap[T]{datos: datos, cantidad: len(arreglo), cmp: funcion_cmp}
}

func (h *heap[T]) EstaVacia() bool {
	return h.cantidad == 0
}

func (h *heap[T]) Cantidad() int {
	return h.cantidad
}

func (h *heap[T]) VerMax() T {
	if h.EstaVacia() {
		panic(_PANIC_COLA_VACIA)
	}
	return h.datos[0]
}

func (h *heap[T]) Encolar(elemento T) {
	if h.cantidad == len(h.datos) {
		h.redimensionar(len(h.datos) * FACTOR_AUMENTO)
	}

	h.datos[h.cantidad] = elemento
	h.cantidad++
	upheap(h.datos[:h.cantidad], h.cantidad-1, h.cmp)
}

func (h *heap[T]) Desencolar() T {
	maximo := h.VerMax()

	h.cantidad--
	h.datos[0] = h.datos[h.cantidad]

	var vacio T
	h.datos[h.cantidad] = vacio // para no retener referencias al elemento que salió
	downheap(h.datos[:h.cantidad], 0, h.cmp)

	nuevaCapacidad := len(h.datos) / FACTOR_AUMENTO
	if h.cantidad*FACTOR_REDUCCION <= len(h.datos) && nuevaCapacidad >= TAM_INICIAL {
		h.redimensionar(nuevaCapacidad)
	}

	return maximo
}

func (h *heap[T]) redimensionar(nuevaCapacidad int) {
	nuevosDatos := make([]T, nuevaCapacidad)
	copy(nuevosDatos, h.datos[:h.cantidad])
	h.datos = nuevosDatos
}

func HeapSort[T any](elementos []T, funcion_cmp func(T, T) int) {
	heapify(elementos, funcion_cmp)

	for fin := len(elementos) - 1; fin > 0; fin-- {
		elementos[0], elementos[fin] = elementos[fin], elementos[0]
		downheap(elementos[:fin], 0, funcion_cmp)
	}
}

func heapify[T any](datos []T, cmp func(T, T) int) {
	for pos := len(datos)/2 - 1; pos >= 0; pos-- {
		downheap(datos, pos, cmp)
	}
}

func upheap[T any](datos []T, pos int, cmp func(T, T) int) {
	padre := (pos - 1) / 2

	for pos > 0 && cmp(datos[pos], datos[padre]) > 0 {
		datos[pos], datos[padre] = datos[padre], datos[pos]
		pos = padre
		padre = (pos - 1) / 2
	}
}

func downheap[T any](datos []T, pos int, cmp func(T, T) int) {
	hijo := hijoMayor(datos, pos, cmp)

	for hijo != -1 && cmp(datos[hijo], datos[pos]) > 0 {
		datos[pos], datos[hijo] = datos[hijo], datos[pos]
		pos = hijo
		hijo = hijoMayor(datos, pos, cmp)
	}
}

// hijoMayor devuelve la posición del hijo de pos con mayor prioridad o -1 si no tiene hijos
func hijoMayor[T any](datos []T, pos int, cmp func(T, T) int) int {
	izquierdo := 2*pos + 1
	derecho := 2*pos + 2

	if izquierdo >= len(datos) {
		return -1
	}
	if derecho < len(datos) && cmp(datos[derecho], datos[izquierdo]) > 0 {
		return derecho
	}
	return izquierdo
}
