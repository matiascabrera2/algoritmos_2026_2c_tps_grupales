package lista_test

import (
	TDALista "tdas/lista"
	"testing"

	"github.com/stretchr/testify/require"
)

const VOLUMEN = 10000

func TestCrearLista(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	require.True(t, lista.EstaVacia())
	require.PanicsWithValue(t, "La lista está vacía", func() { lista.VerPrimero() })
	require.PanicsWithValue(t, "La lista está vacía", func() { lista.VerUltimo() })

}

func TestAgregar(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	lista.InsertarPrimero(43)
	lista.InsertarPrimero(7)

	require.True(t, lista.VerPrimero() == 7)
	require.True(t, lista.VerUltimo() == 43)
	require.False(t, lista.EstaVacia())
}

func TestVaciar(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	lista.InsertarPrimero(43)
	lista.BorrarPrimero()
	require.True(t, lista.EstaVacia())
	require.PanicsWithValue(t, "La lista está vacía", func() { lista.VerPrimero() })
	require.PanicsWithValue(t, "La lista está vacía", func() { lista.VerUltimo() })

}

func TestStrings(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[string]()

	lista.InsertarPrimero("Pruebas")
	lista.InsertarPrimero("de la")
	lista.InsertarPrimero("Lista")

	require.True(t, lista.VerPrimero() == "Lista")
	require.True(t, lista.VerUltimo() == "Pruebas")
	require.False(t, lista.EstaVacia())
}
func TestVolumen(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	for i := 0; i <= VOLUMEN; i++ {
		lista.InsertarPrimero(i)
		require.True(t, lista.VerPrimero() == i)
		require.True(t, lista.Largo() == i+1)
	}
}
func TestCrearIterador(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	iterador := lista.Iterador()
	iterador.Insertar(10)
	//revisa que insertar con el iter al principio equivalga a InsertarPrimero
	require.True(t, lista.VerPrimero() == 10)
	require.False(t, iterador.HayAlgoMas())
	require.True(t, iterador.VerActual() == 10)
	require.PanicsWithValue(t, "Llegaste al final de la lista", func() { iterador.Avanzar() })
	iterador.Insertar(43)
	require.True(t, iterador.VerActual() == 43)
	require.False(t, iterador.HayAlgoMas())
}
func TestIteradorFinal(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	lista.InsertarPrimero(43)
	lista.InsertarPrimero(7)
	iterador := lista.Iterador()
	iterador.Avanzar()
	iterador.Insertar(10)
	require.True(t, iterador.VerActual() == 10)
	require.True(t, lista.VerUltimo() == 10)

}
func TestVolumenIterador(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	for i := 0; i <= VOLUMEN; i++ {
		lista.InsertarUltimo(i)
	}

	iterador := lista.Iterador()

	for i := 0; i < VOLUMEN; i++ {
		require.True(t, iterador.VerActual() == i)
		iterador.Avanzar()
	}
	require.True(t, iterador.VerActual() == VOLUMEN)
}

func TestBorrarIterador(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	iterador := lista.Iterador()

	iterador.Insertar(43)
	iterador.Insertar(7)
	iterador.Insertar(22)

	require.True(t, iterador.VerActual() == 22)
	iterador.Borrar()
	require.True(t, iterador.VerActual() == 7)

}
func TestBorrarMedio(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarPrimero(10)
	lista.InsertarPrimero(20)
	lista.InsertarPrimero(30)
	iterador := lista.Iterador()
	iterador.Avanzar()
	require.True(t, iterador.VerActual() == 20)
	iterador.Borrar()
	require.True(t, iterador.VerActual() == 10)

}
func TestIteradorInterno(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()
	suma := 0
	lista.InsertarPrimero(10)
	lista.InsertarPrimero(20)
	lista.InsertarPrimero(30)
	lista.InsertarPrimero(40)
	lista.Iterar(func(item int) bool {
		suma += item
		return true
	})
	require.True(t, suma == 100)
	require.False(t, lista.EstaVacia())

}

func TestIteradorInternoCorte(t *testing.T) {
	suma := 0
	lista := TDALista.CrearListaEnlazada[int]()
	lista.InsertarPrimero(10)
	lista.InsertarPrimero(20)
	lista.InsertarPrimero(30)
	lista.InsertarPrimero(40)
	lista.Iterar(func(item int) bool {
		suma += item
		return item > 25
	})
	require.True(t, suma == 90)
	require.False(t, lista.EstaVacia())

}
