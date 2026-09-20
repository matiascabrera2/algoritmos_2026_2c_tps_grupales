package pila

/* Definición del struct pila proporcionado por la cátedra. */

type pilaDinamica[T any] struct {
	datos    []T
	cantidad int
}

const (
	TAM_INICIAL      = 10
	FACTOR_AUMENTO   = 2
	FACTOR_REDUCCION = 4
)

/* Funciones auxiliares */

func redimensionar[T any](p *pilaDinamica[T], nuevaCapacidad int) {

	if nuevaCapacidad < TAM_INICIAL {
		nuevaCapacidad = TAM_INICIAL
	}

	nuevosDatos := make([]T, nuevaCapacidad)
	copy(nuevosDatos, p.datos[:p.cantidad])
	p.datos = nuevosDatos
}

/* Implementación de la pila */

func CrearPilaDinamica[T any]() Pila[T] {
	return &pilaDinamica[T]{
		datos:    make([]T, TAM_INICIAL),
		cantidad: 0,
	}
}

func (p *pilaDinamica[T]) EstaVacia() bool {
	return p.cantidad == 0
}

func (p *pilaDinamica[T]) VerTope() T {
	if p.EstaVacia() {
		panic("La pila esta vacia")
	}
	return p.datos[p.cantidad-1]
}

func (p *pilaDinamica[T]) Apilar(v T) {
	if p.cantidad == cap(p.datos) {
		redimensionar(p, cap(p.datos)*FACTOR_AUMENTO)
	}
	p.datos[p.cantidad] = v
	p.cantidad++
}

func (p *pilaDinamica[T]) Desapilar() T {
	if p.EstaVacia() {
		panic("La pila esta vacia")
	}
	tope := p.VerTope()
	p.cantidad--

	if p.cantidad*FACTOR_REDUCCION <= cap(p.datos) {
		redimensionar(p, cap(p.datos)/FACTOR_AUMENTO)
	}

	return tope
}
