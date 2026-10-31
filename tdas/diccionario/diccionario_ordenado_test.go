package diccionario_test

import (
	"math/rand"
	"strings"
	TDADiccionario "tdas/diccionario"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	VOLUMEN                  = 10000
	PANIC_CLAVE_NO_PERTENECE = "La clave no pertenece al diccionario"
	PANIC_FIN_ITERADOR       = "El iterador termino de iterar"
)

func comparador(primero, segundo int) int {
	return (primero - segundo)
}

/* Funciones auxiliares */

// crearABBPilotos arma el árbol de pilotos que usamos en varios tests. Con este orden de inserción
// el árbol queda así: la raíz 3 tiene sólo hijo derecho (22), que tiene a 4 (con hijo derecho 10)
// a su izquierda y a 30 (con hijo derecho 81, que tiene a 43 a su izquierda) a su derecha
func crearABBPilotos() TDADiccionario.DiccionarioOrdenado[int, string] {
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(3, "VER")
	dic.Guardar(22, "TSU")
	dic.Guardar(4, "NOR")
	dic.Guardar(10, "GAS")
	dic.Guardar(30, "LAW")
	dic.Guardar(81, "PIA")
	dic.Guardar(43, "COL")
	return dic
}

func verificarClaves[V any](t *testing.T, dic TDADiccionario.DiccionarioOrdenado[int, V], esperadas []int) {
	require.Equal(t, len(esperadas), dic.Cantidad())

	claves := make([]int, 0)
	dic.Iterar(func(k int, v V) bool {
		claves = append(claves, k)
		return true
	})
	require.Equal(t, esperadas, claves)

	for _, clave := range esperadas {
		require.True(t, dic.Pertenece(clave))
	}
}

func verificarIterador[V any](t *testing.T, iter TDADiccionario.IterDiccionario[int, V], claves []int, datos []V) {
	for i := range claves {
		require.True(t, iter.HayAlgoMas())
		clave, dato := iter.VerActual()
		require.Equal(t, claves[i], clave)
		require.Equal(t, datos[i], dato)
		iter.Avanzar()
	}

	require.False(t, iter.HayAlgoMas())
	require.PanicsWithValue(t, PANIC_FIN_ITERADOR, func() { iter.VerActual() })
	require.PanicsWithValue(t, PANIC_FIN_ITERADOR, func() { iter.Avanzar() })
}

func verificarRango(t *testing.T, desde, hasta *int, claves []int) {
	dic := crearABBPilotos()

	datos := make([]string, 0, len(claves))
	for _, clave := range claves {
		datos = append(datos, dic.Obtener(clave))
	}

	verificarIterador(t, dic.IteradorRango(desde, hasta), claves, datos)

	clavesInterno := make([]int, 0)
	datosInterno := make([]string, 0)
	dic.IterarRango(desde, hasta, func(k int, v string) bool {
		clavesInterno = append(clavesInterno, k)
		datosInterno = append(datosInterno, v)
		return true
	})
	require.Equal(t, claves, clavesInterno)
	require.Equal(t, datos, datosInterno)
}

func iterarRangoHastaCortar(dic TDADiccionario.DiccionarioOrdenado[int, string], desde, hasta *int, cortarEn int) []int {
	claves := make([]int, 0)
	dic.IterarRango(desde, hasta, func(k int, v string) bool {
		claves = append(claves, k)
		return k != cortarEn
	})
	return claves
}

/* Tests */

func TestABBVacio(t *testing.T) {
	dic := TDADiccionario.CrearABB[string, string](strings.Compare)
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece("A"))
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener("A") })
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Borrar("A") })
}

func TestUnElemento(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(1, "Unico")
	require.True(t, dic.Pertenece(1))
	require.True(t, dic.Cantidad() == 1)
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(2) })
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(2) })
}
func TestAgregar(t *testing.T) {
	dic := crearABBPilotos()
	require.True(t, dic.Cantidad() == 7)
	dic.Guardar(11, "PER")
	require.True(t, dic.Cantidad() == 8)
	require.True(t, dic.Obtener(11) == "PER")
}

func TestBorrar(t *testing.T) {
	dic := crearABBPilotos()
	require.True(t, dic.Cantidad() == 7)
	dic.Borrar(22)
	require.True(t, dic.Cantidad() == 6)
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(22) })
}

func TestVolumenSinIter(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, int](comparador)
	for i := 1; i <= VOLUMEN; i++ {
		dic.Guardar(i, i*2)
		require.Equal(t, i*2, dic.Obtener(i))
		require.True(t, dic.Cantidad() == i)
	}
	require.True(t, dic.Cantidad() == VOLUMEN)
	for i := VOLUMEN; i > 0; i-- {
		require.True(t, dic.Cantidad() == i)
		dic.Borrar(i)
		require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(i) })
	}

}

func TestIteradorRangoABBChico(t *testing.T) {
	dic := crearABBPilotos()

	require.True(t, dic.Cantidad() == 7)
	inicio := 10
	fin := 30
	var pilotos []string
	dic.IterarRango(&inicio, &fin, func(k int, v string) bool {
		pilotos = append(pilotos, v)
		return true
	})

	require.Equal(t, []string{"GAS", "TSU", "LAW"}, pilotos)
}

func TestVolumenIteradorInterno(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, int](comparador)

	// Se inserta mezclado para que no quede como una lista
	for _, clave := range rand.New(rand.NewSource(1)).Perm(VOLUMEN + 1) {
		dic.Guardar(clave, clave)
	}
	var datos []int
	var check []int
	for i := 0; i <= VOLUMEN; i++ {
		check = append(check, i)
	}
	dic.IterarRango(nil, nil, func(k int, v int) bool {
		datos = append(datos, v)
		return true
	})
	require.Equal(t, check, datos)
	require.Equal(t, VOLUMEN+1, dic.Cantidad())
}

func TestCorteDeUnSoloLado(t *testing.T) {
	dic := crearABBPilotos()

	require.True(t, dic.Cantidad() == 7)
	inicio := 30
	var pilotos []string
	dic.IterarRango(&inicio, nil, func(k int, v string) bool {
		pilotos = append(pilotos, v)
		return true
	})

	require.Equal(t, []string{"LAW", "COL", "PIA"}, pilotos)
}

func TestIteradorExterno(t *testing.T) {
	dic := crearABBPilotos()
	inicio := 10
	fin := 43

	verificarIterador(t, dic.IteradorRango(&inicio, &fin), []int{10, 22, 30, 43}, []string{"GAS", "TSU", "LAW", "COL"})
}

func TestOrdenStrings(t *testing.T) {
	dic := TDADiccionario.CrearABB[string, int](strings.Compare)
	dic.Guardar("HAM", 7)
	dic.Guardar("ALO", 2)
	dic.Guardar("NOR", 1)
	dic.Guardar("VER", 4)
	dic.Guardar("FAN", 5)
	var campeones []string
	var cantidades []int
	dic.IterarRango(nil, nil, func(k string, v int) bool {
		campeones = append(campeones, k)
		cantidades = append(cantidades, v)
		return true
	})
	require.Equal(t, []string{"ALO", "FAN", "HAM", "NOR", "VER"}, campeones)
	for i := 0; i < len(campeones); i++ {
		require.Equal(t, cantidades[i], dic.Obtener(campeones[i]))
	}

}

/* Test de Iteración*/

func TestIteradoresSinRango(t *testing.T) {
	dic := crearABBPilotos()
	claves := []int{3, 4, 10, 22, 30, 43, 81}
	datos := []string{"VER", "NOR", "GAS", "TSU", "LAW", "COL", "PIA"}

	verificarIterador(t, dic.Iterador(), claves, datos)
	verificarIterador(t, dic.IteradorRango(nil, nil), claves, datos)

	porIterar := make([]int, 0)
	dic.Iterar(func(k int, v string) bool {
		porIterar = append(porIterar, k)
		return true
	})
	require.Equal(t, claves, porIterar)

	porRangoNil := make([]int, 0)
	dic.IterarRango(nil, nil, func(k int, v string) bool {
		porRangoNil = append(porRangoNil, k)
		return true
	})
	require.Equal(t, claves, porRangoNil)
}

func TestIterarSinRangoConCorte(t *testing.T) {
	dic := crearABBPilotos()

	porIterar := make([]int, 0)
	dic.Iterar(func(k int, v string) bool {
		porIterar = append(porIterar, k)
		return k != 22
	})
	require.Equal(t, []int{3, 4, 10, 22}, porIterar)

	require.Equal(t, []int{3, 4, 10, 22}, iterarRangoHastaCortar(dic, nil, nil, 22))
	require.Equal(t, []int{3}, iterarRangoHastaCortar(dic, nil, nil, 3))
	require.Equal(t, []int{3, 4, 10, 22, 30, 43, 81}, iterarRangoHastaCortar(dic, nil, nil, 81))
}

func TestIterarRangoConCorte(t *testing.T) {
	dic := crearABBPilotos()
	inicio := 4
	fin := 43

	require.Equal(t, []int{4}, iterarRangoHastaCortar(dic, &inicio, &fin, 4))                  // en el primero
	require.Equal(t, []int{4, 10, 22}, iterarRangoHastaCortar(dic, &inicio, &fin, 22))         // en el medio
	require.Equal(t, []int{4, 10, 22, 30, 43}, iterarRangoHastaCortar(dic, &inicio, &fin, 43)) // en el último
	require.Equal(t, []int{4, 10, 22, 30, 43}, iterarRangoHastaCortar(dic, &inicio, &fin, 81)) // clave que no está en el rango

	inicio = 22
	require.Equal(t, []int{22, 30, 43, 81}, iterarRangoHastaCortar(dic, &inicio, nil, 81))
}

func TestIteradorRangoCasosBorde(t *testing.T) {
	var inicio, fin int

	// Sólo desde (una clave que está y una que no)
	inicio = 30
	verificarRango(t, &inicio, nil, []int{30, 43, 81})
	inicio = 11
	verificarRango(t, &inicio, nil, []int{22, 30, 43, 81})

	// Sólo hasta (idem. anterior)
	fin = 22
	verificarRango(t, nil, &fin, []int{3, 4, 10, 22})
	fin = 21
	verificarRango(t, nil, &fin, []int{3, 4, 10})

	// Rango vacío
	inicio, fin = 30, 10
	verificarRango(t, &inicio, &fin, []int{})

	// Rango sin elementos
	inicio, fin = 11, 21
	verificarRango(t, &inicio, &fin, []int{})
	inicio, fin = 0, 2
	verificarRango(t, &inicio, &fin, []int{})
	inicio, fin = 82, 100
	verificarRango(t, &inicio, &fin, []int{})

	// Límites que no son claves del árbol
	inicio, fin = 5, 40
	verificarRango(t, &inicio, &fin, []int{10, 22, 30})
	inicio, fin = 0, 100
	verificarRango(t, &inicio, &fin, []int{3, 4, 10, 22, 30, 43, 81})

	// desde == hasta
	inicio, fin = 22, 22
	verificarRango(t, &inicio, &fin, []int{22})
	inicio, fin = 23, 23
	verificarRango(t, &inicio, &fin, []int{})
}

func TestGuardarClaveExistenteActualizaElDato(t *testing.T) {
	dic := crearABBPilotos()

	for _, dato := range []string{"ANT", "ALO", "HAM"} {
		dic.Guardar(22, dato)
		require.Equal(t, dato, dic.Obtener(22))
		require.Equal(t, 7, dic.Cantidad())
	}

	dic.Guardar(3, "RUS")
	dic.Guardar(43, "BOT")
	require.Equal(t, 7, dic.Cantidad())

	verificarIterador(t, dic.Iterador(), []int{3, 4, 10, 22, 30, 43, 81}, []string{"RUS", "NOR", "GAS", "HAM", "LAW", "BOT", "PIA"})
}

func TestBorrarHojas(t *testing.T) {
	dic := crearABBPilotos()

	require.Equal(t, "GAS", dic.Borrar(10))
	require.False(t, dic.Pertenece(10))
	verificarClaves(t, dic, []int{3, 4, 22, 30, 43, 81})

	require.Equal(t, "COL", dic.Borrar(43))
	require.False(t, dic.Pertenece(43))
	verificarClaves(t, dic, []int{3, 4, 22, 30, 81})

	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(10) })
}

func TestBorrarNodoConSoloHijoIzquierdo(t *testing.T) {
	dic := crearABBPilotos()

	require.Equal(t, "PIA", dic.Borrar(81))
	require.False(t, dic.Pertenece(81))
	verificarClaves(t, dic, []int{3, 4, 10, 22, 30, 43})
	require.Equal(t, "COL", dic.Obtener(43))
}

func TestBorrarNodoConSoloHijoDerecho(t *testing.T) {
	dic := crearABBPilotos()

	require.Equal(t, "LAW", dic.Borrar(30))
	require.False(t, dic.Pertenece(30))
	verificarClaves(t, dic, []int{3, 4, 10, 22, 43, 81})

	require.Equal(t, "NOR", dic.Borrar(4))
	require.False(t, dic.Pertenece(4))
	verificarClaves(t, dic, []int{3, 10, 22, 43, 81})
}

func TestBorrarNodoConDosHijos(t *testing.T) {
	dic := crearABBPilotos()

	require.Equal(t, "TSU", dic.Borrar(22))
	require.False(t, dic.Pertenece(22))
	verificarClaves(t, dic, []int{3, 4, 10, 30, 43, 81})
	require.Equal(t, "LAW", dic.Obtener(30))
}

func TestBorrarNodoConDosHijosYSucesorProfundo(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, int](comparador)
	for _, clave := range []int{50, 30, 70, 60, 80, 55, 57} {
		dic.Guardar(clave, clave*10)
	}

	require.Equal(t, 500, dic.Borrar(50))
	require.False(t, dic.Pertenece(50))
	verificarClaves(t, dic, []int{30, 55, 57, 60, 70, 80})

	require.Equal(t, 550, dic.Borrar(55))
	require.False(t, dic.Pertenece(55))
	verificarClaves(t, dic, []int{30, 57, 60, 70, 80})
	require.Equal(t, 570, dic.Obtener(57))
}

func TestBorrarRaizConUnHijo(t *testing.T) {
	dic := crearABBPilotos()

	require.Equal(t, "VER", dic.Borrar(3))
	require.False(t, dic.Pertenece(3))
	verificarClaves(t, dic, []int{4, 10, 22, 30, 43, 81})
	require.Equal(t, "TSU", dic.Obtener(22))
}

func TestBorrarTodasLasClavesVaciaElArbol(t *testing.T) {
	dic := crearABBPilotos()
	claves := []int{22, 3, 81, 10, 30, 4, 43}
	datos := []string{"TSU", "VER", "PIA", "GAS", "LAW", "NOR", "COL"}

	for i, clave := range claves {
		require.Equal(t, datos[i], dic.Borrar(clave))
		require.False(t, dic.Pertenece(clave))
		require.Equal(t, len(claves)-i-1, dic.Cantidad())
		for _, restante := range claves[i+1:] {
			require.True(t, dic.Pertenece(restante))
		}
	}

	require.Equal(t, 0, dic.Cantidad())
	verificarIterador(t, dic.Iterador(), []int{}, []string{})
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(22) })
	require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Borrar(22) })

	dic.Guardar(5, "NUEVO")
	verificarClaves(t, dic, []int{5})
	require.Equal(t, "NUEVO", dic.Obtener(5))
}

func TestClaveInexistenteEnABBNoVacio(t *testing.T) {
	dic := crearABBPilotos()

	for _, clave := range []int{0, 5, 25, 50, 100} {
		require.False(t, dic.Pertenece(clave))
		require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Obtener(clave) })
		require.PanicsWithValue(t, PANIC_CLAVE_NO_PERTENECE, func() { dic.Borrar(clave) })
		verificarClaves(t, dic, []int{3, 4, 10, 22, 30, 43, 81})
	}
}
