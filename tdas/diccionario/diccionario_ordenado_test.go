package diccionario_test

import (
	"strings"
	TDADiccionario "tdas/diccionario"
	"testing"

	"github.com/stretchr/testify/require"
)

const VOLUMEN = 10000

func comparador(primero, segundo int) int {
	return (primero - segundo)
}

func TestABBVacio(t *testing.T) {
	dic := TDADiccionario.CrearABB[string, string](strings.Compare)
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("A") })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("A") })
}

func TestUnElemento(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(1, "Unico")
	require.True(t, dic.Pertenece(1))
	require.True(t, dic.Cantidad() == 1)
	require.True(t, dic.Obtener(1) == "Unico")
}

func TestIteradorRangoABBChico(t *testing.T) {
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(3, "VER")
	dic.Guardar(22, "TSU")
	dic.Guardar(4, "NOR")
	dic.Guardar(10, "GAS")
	dic.Guardar(30, "LAW")
	dic.Guardar(81, "PIA")
	dic.Guardar(43, "COL")

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

	for i := VOLUMEN; i >= 0; i-- {
		dic.Guardar(i, i)
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
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(3, "VER")
	dic.Guardar(22, "TSU")
	dic.Guardar(4, "NOR")
	dic.Guardar(10, "GAS")
	dic.Guardar(30, "LAW")
	dic.Guardar(81, "PIA")
	dic.Guardar(43, "COL")

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
	dic := TDADiccionario.CrearABB[int, string](comparador)
	dic.Guardar(3, "VER")
	dic.Guardar(22, "TSU")
	dic.Guardar(4, "NOR")
	dic.Guardar(10, "GAS")
	dic.Guardar(30, "LAW")
	dic.Guardar(81, "PIA")
	dic.Guardar(43, "COL")
	inicio := 10
	fin := 43
	var visitados []string
	var llaves []int
	iter := dic.IteradorRango(&inicio, &fin)
	for iter.HayAlgoMas() {
		k, v := iter.VerActual()
		visitados = append(visitados, v)
		llaves = append(llaves, k)
		require.True(t, iter.HayAlgoMas())
		iter.Avanzar()
	}
	require.Equal(t, []string{"GAS", "TSU", "LAW", "COL"}, visitados)
	require.Equal(t, []int{10, 22, 30, 43}, llaves)
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
