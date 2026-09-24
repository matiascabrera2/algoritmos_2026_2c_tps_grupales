package lista_test

import (
	TDALista "tdas/lista"
	"testing"

	"github.com/stretchr/testify/require"
)

const VOLUMEN = 10000

/* Funciones auxiliares */

// verificarContenido recorre la lista con Iterar ycompara el resultado contra lo esperado.
func verificarContenido[T any](t *testing.T, lista TDALista.Lista[T], esperado []T) {
	obtenido := make([]T, 0, len(esperado))
	lista.Iterar(func(v T) bool {
		obtenido = append(obtenido, v)
		return true
	})
	require.Equal(t, esperado, obtenido)
}

func insertarUltimoYBorrarPrimeroEnOrden[T any](t *testing.T, valores []T) {
	lista := TDALista.CrearListaEnlazada[T]()

	for _, v := range valores {
		lista.InsertarUltimo(v)
	}

	require.Equal(t, len(valores), lista.Largo())
	require.Equal(t, valores[0], lista.VerPrimero())
	require.Equal(t, valores[len(valores)-1], lista.VerUltimo())
	verificarContenido(t, lista, valores)

	for _, v := range valores {
		require.Equal(t, v, lista.VerPrimero())
		require.Equal(t, v, lista.BorrarPrimero())
	}

	require.True(t, lista.EstaVacia())
}

/* Casos básicos */

func TestListaVacia(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	require.True(t, lista.EstaVacia())
	require.Equal(t, 0, lista.Largo())
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerPrimero() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerUltimo() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.BorrarPrimero() })
}

func TestInsertarUltimoYBorrarPrimeroConEnteros(t *testing.T) {
	insertarUltimoYBorrarPrimeroEnOrden(t, []int{4, 7, 3})
}

func TestInsertarUltimoYBorrarPrimeroConStrings(t *testing.T) {
	insertarUltimoYBorrarPrimeroEnOrden(t, []string{"a", "c", "o", "b"})
}

func TestInsertarUltimoYBorrarPrimeroConFloats(t *testing.T) {
	insertarUltimoYBorrarPrimeroEnOrden(t, []float64{3.14, 2.71})
}

func TestInsertarUltimoYBorrarPrimeroConBools(t *testing.T) {
	insertarUltimoYBorrarPrimeroEnOrden(t, []bool{true, false, true})
}

func TestInsertarUltimoYBorrarPrimeroConStructs(t *testing.T) {
	type Auto struct {
		Marca string
		Año   int
	}
	insertarUltimoYBorrarPrimeroEnOrden(t, []Auto{
		{Marca: "Honda", Año: 2017},
		{Marca: "Tesla", Año: 2021},
		{Marca: "Ford", Año: 2023},
		{Marca: "Chevrolet", Año: 2024},
		{Marca: "BYD", Año: 2025},
	})
}

func TestInsertarPrimero(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	lista.InsertarPrimero(1)
	lista.InsertarPrimero(2)
	lista.InsertarPrimero(3)

	verificarContenido(t, lista, []int{3, 2, 1})
	require.Equal(t, 3, lista.VerPrimero())
	require.Equal(t, 1, lista.VerUltimo())
}

func TestListaVaciadaSeComportaComoRecienCreada(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	for !lista.EstaVacia() {
		lista.BorrarPrimero()
	}

	require.True(t, lista.EstaVacia())
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerPrimero() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerUltimo() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.BorrarPrimero() })

	lista.InsertarPrimero(9)
	require.Equal(t, 9, lista.VerPrimero())
	require.Equal(t, 9, lista.VerUltimo())
	require.Equal(t, 1, lista.Largo())
}

func TestVolumen(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	for i := 0; i < VOLUMEN; i++ {
		lista.InsertarUltimo(i)
		require.Equal(t, 0, lista.VerPrimero())
		require.Equal(t, i, lista.VerUltimo())
		require.Equal(t, i+1, lista.Largo())
	}

	for i := 0; i < VOLUMEN; i++ {
		require.Equal(t, i, lista.VerPrimero())
		require.Equal(t, i, lista.BorrarPrimero())
	}

	require.True(t, lista.EstaVacia())
}

/* Iterador interno */

func TestIterarInterno(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(10)
	lista.InsertarUltimo(20)
	lista.InsertarUltimo(30)
	lista.InsertarUltimo(40)

	var visitados []int
	lista.Iterar(func(v int) bool {
		visitados = append(visitados, v)
		return true
	})

	require.Equal(t, []int{10, 20, 30, 40}, visitados)
}

func TestIterarInternoConCorte(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(10)
	lista.InsertarUltimo(20)
	lista.InsertarUltimo(30)
	lista.InsertarUltimo(40)

	var visitados []int
	lista.Iterar(func(v int) bool {
		visitados = append(visitados, v)
		return v < 25
	})

	require.Equal(t, []int{10, 20, 30}, visitados)
}

func TestIterarInternoConListaVacia(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	llamadas := 0
	lista.Iterar(func(v int) bool {
		llamadas++
		return true
	})

	require.Equal(t, 0, llamadas)
}

/* Iterador externo */

func TestIteradorInsertarAlPrincipio(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	iter.Insertar(99)

	verificarContenido(t, lista, []int{99, 1, 2, 3})
	require.Equal(t, 99, iter.VerActual())
	require.Equal(t, 99, lista.VerPrimero())
}

func TestIteradorInsertarAlFinal(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	for iter.HayAlgoMas() {
		iter.Avanzar()
	}

	iter.Insertar(99)

	verificarContenido(t, lista, []int{1, 2, 3, 99})
	require.Equal(t, 99, lista.VerUltimo())
	require.Equal(t, 99, iter.VerActual())
}

func TestIteradorInsertarEnElMedio(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	iter.Avanzar()

	iter.Insertar(99)

	verificarContenido(t, lista, []int{1, 99, 2, 3})
	require.Equal(t, 99, iter.VerActual())
}

func TestIteradorBorrarPrimero(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	borrado := iter.Borrar()

	require.Equal(t, 1, borrado)
	verificarContenido(t, lista, []int{2, 3})
	require.Equal(t, 2, lista.VerPrimero())
}

func TestIteradorBorrarUltimo(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	iter.Avanzar()
	iter.Avanzar()

	borrado := iter.Borrar()

	require.Equal(t, 3, borrado)
	verificarContenido(t, lista, []int{1, 2})
	require.Equal(t, 2, lista.VerUltimo())
}

func TestIteradorBorrarEnElMedio(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)
	lista.InsertarUltimo(2)
	lista.InsertarUltimo(3)

	iter := lista.Iterador()
	iter.Avanzar()

	borrado := iter.Borrar()

	require.Equal(t, 2, borrado)
	verificarContenido(t, lista, []int{1, 3})
	require.Equal(t, 3, iter.VerActual())
}

func TestIteradorConListaVacia(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	iter := lista.Iterador()

	require.False(t, iter.HayAlgoMas())
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.VerActual() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Avanzar() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Borrar() })

	iter.Insertar(1)
	verificarContenido(t, lista, []int{1})
}

func TestIteradorArrojaPanicLuegoDeIterar(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarUltimo(1)

	iter := lista.Iterador()
	iter.Avanzar()

	require.False(t, iter.HayAlgoMas())
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.VerActual() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Avanzar() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Borrar() })
}

func TestVolumenIterador(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	for i := 0; i < VOLUMEN; i++ {
		lista.InsertarUltimo(i)
	}

	iter := lista.Iterador()
	for i := 0; i < VOLUMEN; i++ {
		require.Equal(t, i, iter.VerActual())
		iter.Avanzar()
	}

	require.False(t, iter.HayAlgoMas())
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.VerActual() })
}

func TestVolumenIteradorInsertarYBorrarEnElMedio(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	for i := 0; i < VOLUMEN; i++ {
		lista.InsertarUltimo(i)
	}

	iter := lista.Iterador()
	for i := 0; i < VOLUMEN/2; i++ {
		iter.Avanzar()
	}

	iter.Insertar(-1)
	require.Equal(t, -1, iter.VerActual())

	borrado := iter.Borrar()
	require.Equal(t, -1, borrado)
	require.Equal(t, VOLUMEN, lista.Largo())

	esperado := make([]int, VOLUMEN)
	for i := range esperado {
		esperado[i] = i
	}
	verificarContenido(t, lista, esperado)
}
