package cola_prioridad_test

import (
	"math/rand"
	"strings"
	TDAColaPrioridad "tdas/cola_prioridad"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_PANIC_COLA_VACIA = "La cola esta vacia"

	_VOLUMEN           = 100000
	_PRIMER_BLOQUE     = 100000
	_SEGUNDO_BLOQUE    = 200000
	_CANT_A_DESENCOLAR = 50000
	_VALORES_DISTINTOS = 100 // cuántos valores distintos hay en el volumen con repetidos

	// Heapify de Floyd hace como mucho 2 comparaciones por elemento; encolar de a uno un arreglo
	// ascendente haría n*log(n), así que esta cota distingue una creación lineal de una que no lo es.
	_MAX_COMPARACIONES_POR_ELEMENTO = 2
)

var _TAMANIOS_LINEALIDAD = []int{1000, 10000, _VOLUMEN}

func comparador(primero, segundo int) int {
	return (primero - segundo)
}

func comparadorInvertido(primero, segundo int) int {
	return (segundo - primero)
}

/* Funciones auxiliares */

// secuencia devuelve los números desde..hasta-1 en orden ascendente
func secuencia(desde, hasta int) []int {
	valores := make([]int, hasta-desde)
	for i := range valores {
		valores[i] = desde + i
	}
	return valores
}

// mezclados devuelve los números desde..hasta-1 en un orden mezclado pero reproducible
func mezclados(desde, hasta int) []int {
	valores := rand.New(rand.NewSource(1)).Perm(hasta - desde)
	for i := range valores {
		valores[i] += desde
	}
	return valores
}

// invertido devuelve los elementos en orden inverso al que recibieron
func invertido[T any](elementos []T) []T {
	resultado := make([]T, len(elementos))
	for i, elemento := range elementos {
		resultado[len(elementos)-1-i] = elemento
	}
	return resultado
}

func verificarColaVacia[T any](t *testing.T, cola TDAColaPrioridad.ColaPrioridad[T]) {
	require.True(t, cola.EstaVacia())
	require.Equal(t, 0, cola.Cantidad())
	require.PanicsWithValue(t, _PANIC_COLA_VACIA, func() { cola.VerMax() })
	require.PanicsWithValue(t, _PANIC_COLA_VACIA, func() { cola.Desencolar() })
}

func desencolarTodo[T any](t *testing.T, cola TDAColaPrioridad.ColaPrioridad[T], esperados []T) {
	require.Equal(t, len(esperados), cola.Cantidad())

	for i, esperado := range esperados {
		require.False(t, cola.EstaVacia())
		require.Equal(t, len(esperados)-i, cola.Cantidad())
		require.Equal(t, esperado, cola.VerMax())
		require.Equal(t, esperado, cola.Desencolar())
	}

	verificarColaVacia(t, cola)
}

func desencolarBloque(t *testing.T, cola TDAColaPrioridad.ColaPrioridad[int], desde, hasta int) {
	for i := hasta - 1; i >= desde; i-- {
		cantidad := cola.Cantidad()
		require.False(t, cola.EstaVacia())
		require.Equal(t, i, cola.VerMax())
		require.Equal(t, i, cola.Desencolar())
		require.Equal(t, cantidad-1, cola.Cantidad())
	}
}

func encolarVerificando[T any](t *testing.T, cola TDAColaPrioridad.ColaPrioridad[T], elementos []T, maximoPrevio T, cmp func(T, T) int) {
	maximo := maximoPrevio

	for _, elemento := range elementos {
		cantidad := cola.Cantidad()
		if cantidad == 0 || cmp(elemento, maximo) > 0 {
			maximo = elemento
		}

		cola.Encolar(elemento)
		require.False(t, cola.EstaVacia())
		require.Equal(t, cantidad+1, cola.Cantidad())
		require.Equal(t, maximo, cola.VerMax())
	}
}

func encolarYDesencolarPorPrioridad[T any](t *testing.T, elementos, enOrden []T, cmp func(T, T) int) {
	cola := TDAColaPrioridad.CrearHeap(cmp)
	verificarColaVacia(t, cola)

	var sinMaximoPrevio T
	encolarVerificando(t, cola, elementos, sinMaximoPrevio, cmp)
	desencolarTodo(t, cola, enOrden)
}

func crearHeapArrYDesencolarEnOrden[T any](t *testing.T, arreglo, enOrden []T, cmp func(T, T) int) {
	desencolarTodo(t, TDAColaPrioridad.CrearHeapArr(arreglo, cmp), enOrden)
}

func heapSortYVerificar[T any](t *testing.T, elementos, esperado []T, cmp func(T, T) int) {
	TDAColaPrioridad.HeapSort(elementos, cmp)
	require.Equal(t, esperado, elementos)
}

func verificarCreacionYOrdenamiento[T any](t *testing.T, arreglo, enOrden []T, cmp func(T, T) int) {
	crearHeapArrYDesencolarEnOrden(t, arreglo, invertido(enOrden), cmp)
	heapSortYVerificar(t, arreglo, enOrden, cmp)
}

/* Tests */

func TestColaVacia(t *testing.T) {
	verificarColaVacia(t, TDAColaPrioridad.CrearHeap(comparador))
}

func TestUnElemento(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, []int{7}, []int{7}, comparador)
}

func TestEncolarYDesencolarPorPrioridadConEnteros(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, []int{7, 3, 9, 1, 5, 9}, []int{9, 9, 7, 5, 3, 1}, comparador)
}

func TestEncolarYDesencolarPorPrioridadConStrings(t *testing.T) {
	encolarYDesencolarPorPrioridad(t,
		[]string{"mesa", "auto", "zorro", "barco", "casa"},
		[]string{"zorro", "mesa", "casa", "barco", "auto"},
		strings.Compare)
}

func TestEncolarEnOrdenCreciente(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, secuencia(1, 11), invertido(secuencia(1, 11)), comparador)
}

func TestEncolarEnOrdenDecreciente(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, invertido(secuencia(1, 11)), invertido(secuencia(1, 11)), comparador)
}

func TestComparadorInvertidoSacaPrimeroLosMenores(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, []int{8, 3, 10, 1, 6}, []int{1, 3, 6, 8, 10}, comparadorInvertido)
}

func TestOperacionesMezcladas(t *testing.T) {
	cola := TDAColaPrioridad.CrearHeap(comparador)

	cola.Encolar(5)
	require.False(t, cola.EstaVacia())
	require.Equal(t, 1, cola.Cantidad())
	require.Equal(t, 5, cola.VerMax())

	cola.Encolar(3)
	cola.Encolar(9)
	require.Equal(t, 3, cola.Cantidad())
	require.Equal(t, 9, cola.VerMax())

	require.Equal(t, 9, cola.Desencolar())
	require.Equal(t, 2, cola.Cantidad())
	require.Equal(t, 5, cola.VerMax())

	cola.Encolar(7)
	cola.Encolar(1)
	cola.Encolar(7)
	require.Equal(t, 5, cola.Cantidad())
	require.Equal(t, 7, cola.VerMax())

	require.Equal(t, 7, cola.Desencolar())
	require.Equal(t, 7, cola.Desencolar())
	require.Equal(t, 3, cola.Cantidad())
	require.Equal(t, 5, cola.VerMax())

	cola.Encolar(10)
	require.Equal(t, 4, cola.Cantidad())
	require.Equal(t, 10, cola.VerMax())

	desencolarTodo(t, cola, []int{10, 5, 3, 1})
}

func TestElementosConLaMismaPrioridad(t *testing.T) {
	type tarea struct {
		nombre    string
		prioridad int
	}
	tareas := []tarea{{"a", 2}, {"b", 5}, {"c", 2}, {"d", 5}, {"e", 1}, {"f", 5}}
	cola := TDAColaPrioridad.CrearHeap(func(a, b tarea) int { return a.prioridad - b.prioridad })

	for _, nueva := range tareas {
		cola.Encolar(nueva)
	}
	require.Equal(t, len(tareas), cola.Cantidad())

	// Los de igual prioridad pueden salir en cualquier orden, pero ninguno se pierde ni se repite
	vistos := make(map[string]bool)
	for i, prioridad := range []int{5, 5, 5, 2, 2, 1} {
		require.False(t, cola.EstaVacia())
		require.Equal(t, len(tareas)-i, cola.Cantidad())
		require.Equal(t, prioridad, cola.VerMax().prioridad)

		saliente := cola.Desencolar()
		require.Equal(t, prioridad, saliente.prioridad)
		require.False(t, vistos[saliente.nombre])
		vistos[saliente.nombre] = true
	}

	require.Equal(t, len(tareas), len(vistos))
	verificarColaVacia(t, cola)
}

func TestColaVaciadaSeComportaComoRecienCreada(t *testing.T) {
	cola := TDAColaPrioridad.CrearHeap(comparador)

	for _, elemento := range []int{4, 8, 2} {
		cola.Encolar(elemento)
	}
	desencolarTodo(t, cola, []int{8, 4, 2})

	cola.Encolar(6)
	require.False(t, cola.EstaVacia())
	require.Equal(t, 1, cola.Cantidad())
	require.Equal(t, 6, cola.VerMax())

	cola.Encolar(9)
	desencolarTodo(t, cola, []int{9, 6})
}

func TestDeVolumen(t *testing.T) {
	cola := TDAColaPrioridad.CrearHeap(comparador)

	encolarVerificando(t, cola, mezclados(0, _PRIMER_BLOQUE), 0, comparador)
	desencolarBloque(t, cola, _PRIMER_BLOQUE-_CANT_A_DESENCOLAR, _PRIMER_BLOQUE)
	encolarVerificando(t, cola, mezclados(_PRIMER_BLOQUE, _PRIMER_BLOQUE+_SEGUNDO_BLOQUE), _PRIMER_BLOQUE-_CANT_A_DESENCOLAR-1, comparador)
	desencolarBloque(t, cola, _PRIMER_BLOQUE, _PRIMER_BLOQUE+_SEGUNDO_BLOQUE)
	desencolarBloque(t, cola, 0, _PRIMER_BLOQUE-_CANT_A_DESENCOLAR)

	verificarColaVacia(t, cola)
}

func TestDeVolumenConRepetidos(t *testing.T) {
	repeticiones := _VOLUMEN / _VALORES_DISTINTOS
	cola := TDAColaPrioridad.CrearHeap(comparador)

	for _, i := range mezclados(0, _VOLUMEN) {
		cola.Encolar(i % _VALORES_DISTINTOS)
	}
	require.Equal(t, _VOLUMEN, cola.Cantidad())

	for i := 0; i < _VOLUMEN; i++ {
		esperado := _VALORES_DISTINTOS - 1 - i/repeticiones
		require.Equal(t, _VOLUMEN-i, cola.Cantidad())
		require.Equal(t, esperado, cola.VerMax())
		require.Equal(t, esperado, cola.Desencolar())
	}

	verificarColaVacia(t, cola)
}

/* Tests con CrearHeapArr y con HeapSort */

func TestArregloVacio(t *testing.T) {
	verificarCreacionYOrdenamiento(t, []int{}, []int{}, comparador)
	verificarCreacionYOrdenamiento(t, nil, nil, comparador)
}

func TestArregloConUnElemento(t *testing.T) {
	verificarCreacionYOrdenamiento(t, []int{7}, []int{7}, comparador)
}

func TestArregloDesordenado(t *testing.T) {
	verificarCreacionYOrdenamiento(t, []int{4, 10, 3, 5, 1, 8, 2, 9}, []int{1, 2, 3, 4, 5, 8, 9, 10}, comparador)
}

func TestArregloConRepetidosYNegativos(t *testing.T) {
	verificarCreacionYOrdenamiento(t, []int{5, -2, 9, 0, 5, 3, -7, 9}, []int{-7, -2, 0, 3, 5, 5, 9, 9}, comparador)
}

func TestArregloYaOrdenado(t *testing.T) {
	verificarCreacionYOrdenamiento(t, secuencia(1, 9), secuencia(1, 9), comparador)
}

func TestArregloEnOrdenInverso(t *testing.T) {
	verificarCreacionYOrdenamiento(t, invertido(secuencia(1, 9)), secuencia(1, 9), comparador)
}

func TestTodosLosElementosIguales(t *testing.T) {
	encolarYDesencolarPorPrioridad(t, []int{5, 5, 5, 5}, []int{5, 5, 5, 5}, comparador)
	verificarCreacionYOrdenamiento(t, []int{5, 5, 5, 5}, []int{5, 5, 5, 5}, comparador)
}

func TestArregloConStrings(t *testing.T) {
	verificarCreacionYOrdenamiento(t,
		[]string{"pera", "manzana", "uva", "banana", "kiwi", "manzana"},
		[]string{"banana", "kiwi", "manzana", "manzana", "pera", "uva"},
		strings.Compare)
}

func TestArregloConComparadorInvertido(t *testing.T) {
	verificarCreacionYOrdenamiento(t, []int{3, 4, 8, 9, 2, 10, 5, 6, 7, 1}, invertido(secuencia(1, 11)), comparadorInvertido)
}

func TestArregloDeVolumen(t *testing.T) {
	verificarCreacionYOrdenamiento(t, mezclados(0, _VOLUMEN), secuencia(0, _VOLUMEN), comparador)
}

/* Sólo CrearHeapArr */

func TestCrearHeapArrNoModificaElArregloOriginal(t *testing.T) {
	arreglo := []int{4, 10, 3, 5, 1}
	cola := TDAColaPrioridad.CrearHeapArr(arreglo, comparador)

	cola.Encolar(20)
	require.Equal(t, 20, cola.Desencolar())
	require.Equal(t, 10, cola.Desencolar())

	require.Equal(t, []int{4, 10, 3, 5, 1}, arreglo)
}

func TestCrearHeapArrYLuegoSeguirOperando(t *testing.T) {
	cola := TDAColaPrioridad.CrearHeapArr([]int{4, 10, 3}, comparador)

	cola.Encolar(7)
	require.Equal(t, 4, cola.Cantidad())
	require.Equal(t, 10, cola.VerMax())

	cola.Encolar(12)
	require.Equal(t, 5, cola.Cantidad())
	require.Equal(t, 12, cola.VerMax())

	require.Equal(t, 12, cola.Desencolar())
	desencolarTodo(t, cola, []int{10, 7, 4, 3})
}

func TestCrearHeapArrVacioYLuegoEncolar(t *testing.T) {
	for _, arreglo := range [][]int{{}, nil} {
		cola := TDAColaPrioridad.CrearHeapArr(arreglo, comparador)

		encolarVerificando(t, cola, mezclados(0, 100), 0, comparador)
		desencolarBloque(t, cola, 0, 100)

		verificarColaVacia(t, cola)
	}
}

func TestCrearHeapArrYLuegoEncolarDeVolumen(t *testing.T) {
	cola := TDAColaPrioridad.CrearHeapArr(mezclados(0, _VOLUMEN), comparador)
	require.Equal(t, _VOLUMEN, cola.Cantidad())

	encolarVerificando(t, cola, mezclados(_VOLUMEN, 2*_VOLUMEN), _VOLUMEN-1, comparador)
	desencolarBloque(t, cola, 0, 2*_VOLUMEN)

	verificarColaVacia(t, cola)
}

func TestCrearHeapArrEsLineal(t *testing.T) {
	for _, n := range _TAMANIOS_LINEALIDAD {
		for _, arreglo := range [][]int{secuencia(0, n), mezclados(0, n)} {
			comparaciones := 0
			TDAColaPrioridad.CrearHeapArr(arreglo, func(a, b int) int {
				comparaciones++
				return a - b
			})
			require.LessOrEqual(t, comparaciones, _MAX_COMPARACIONES_POR_ELEMENTO*n)
		}
	}
}

/* Sólo HeapSort */

func TestHeapSortSoloOrdenaLosElementosPasados(t *testing.T) {
	arreglo := []int{9, 8, 7, 6, 5, 4}
	TDAColaPrioridad.HeapSort(arreglo[:4], comparador)

	require.Equal(t, []int{6, 7, 8, 9, 5, 4}, arreglo)
}
