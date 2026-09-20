package cola_test

import (
	TDACola "tdas/cola"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	PRIMER_BLOQUE     = 100000 // Cantidad de elementos a encolar en el primer bloque
	SEGUNDO_BLOQUE    = 200000 // Idem. ant. pero en el segundo bloque
	CANT_A_DESENCOLAR = 50000  // Cantidad de elementos a desencolar en el primer bloque
)

/* Funciones auxiliares genéricas */
func encolarYDesencolarEnOrden[T any](t *testing.T, valores []T) {
	cola := TDACola.CrearColaEnlazada[T]()

	for _, v := range valores {
		cola.Encolar(v)
	}

	for _, v := range valores {
		require.Equal(t, v, cola.VerPrimero())
		require.Equal(t, v, cola.Desencolar())
	}

	require.True(t, cola.EstaVacia())
}

func colaVaciadaSeComportaComoRecienCreada[T any](t *testing.T, valores []T) {
	cola := TDACola.CrearColaEnlazada[T]()

	for _, v := range valores {
		cola.Encolar(v)
	}

	for _, v := range valores {
		require.Equal(t, v, cola.Desencolar())
	}

	require.True(t, cola.EstaVacia())

	// Ahora la cola debería comportarse como recién creada
	for _, v := range valores {
		cola.Encolar(v)
		require.Equal(t, v, cola.VerPrimero())
		require.Equal(t, v, cola.Desencolar())
	}

	require.True(t, cola.EstaVacia())
}

func encolarBloque(t *testing.T, c TDACola.Cola[int], desde, hasta int) {
	for i := desde; i < hasta; i++ {
		c.Encolar(i)
		require.False(t, c.EstaVacia())
	}
}

func desencolarBloque(t *testing.T, c TDACola.Cola[int], desde, hasta int) {
	for i := desde; i < hasta; i++ {
		require.Equal(t, i, c.VerPrimero())
		require.Equal(t, i, c.Desencolar())
	}
}

/*Tests*/
func TestColaVacia(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	require.True(t, cola.EstaVacia())
	require.Panics(t, func() { cola.VerPrimero() })
	require.Panics(t, func() { cola.Desencolar() })
}

func TestEncolarYDesencolarConEnteros(t *testing.T) {
	encolarYDesencolarEnOrden(t, []int{1, 2, 3})
}

func TestEncolarYDesencolarConStrings(t *testing.T) {
	encolarYDesencolarEnOrden(t, []string{"a", "c", "o", "b"})
}

func TestEncolarYDesencolarConFloats(t *testing.T) {
	encolarYDesencolarEnOrden(t, []float64{3.14, 2.71})
}

func TestEncolarYDesencolarConBools(t *testing.T) {
	encolarYDesencolarEnOrden(t, []bool{true, false, true})
}

func TestEncolarYDesencolarConStructs(t *testing.T) {
	type Auto struct {
		Marca string
		Año   int
	}
	encolarYDesencolarEnOrden(t, []Auto{
		{Marca: "Honda", Año: 2020},
		{Marca: "Tesla", Año: 2021},
		{Marca: "BYD", Año: 2025},
	})
}

func TestColaVaciadaSeComportaComoRecienCreadaConEnteros(t *testing.T) {
	colaVaciadaSeComportaComoRecienCreada(t, []int{6, 7, 6, 7})
}

func TestColaVaciadaSeComportaComoRecienCreadaConStrings(t *testing.T) {
	colaVaciadaSeComportaComoRecienCreada(t, []string{"Hola", "Mundo", "Aguante", "Boca"})
}

func TestColaVaciadaSeComportaComoRecienCreadaConFloats(t *testing.T) {
	colaVaciadaSeComportaComoRecienCreada(t, []float64{3.14, 2.71, 3.14, 2.71, 1.61})
}

func TestColaVaciadaSeComportaComoRecienCreadaConBools(t *testing.T) {
	colaVaciadaSeComportaComoRecienCreada(t, []bool{true, false, false, true, false})
}

func TestColaVaciadaSeComportaComoRecienCreadaConStructs(t *testing.T) {
	type Persona struct {
		Nombre string
		Edad   int
	}
	colaVaciadaSeComportaComoRecienCreada(t, []Persona{
		{Nombre: "Juan", Edad: 30},
		{Nombre: "Román", Edad: 25},
		{Nombre: "Riquelme", Edad: 40},
	})
}

func TestDeVolumen(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	encolarBloque(t, cola, 0, PRIMER_BLOQUE)
	require.Equal(t, 0, cola.VerPrimero())
	desencolarBloque(t, cola, 0, CANT_A_DESENCOLAR)
	encolarBloque(t, cola, PRIMER_BLOQUE, PRIMER_BLOQUE+SEGUNDO_BLOQUE)
	require.Equal(t, CANT_A_DESENCOLAR, cola.VerPrimero())
	desencolarBloque(t, cola, CANT_A_DESENCOLAR, PRIMER_BLOQUE+SEGUNDO_BLOQUE)

	require.True(t, cola.EstaVacia())
}

func TestColaDesencoladaDevuelvePanicAlVerPrimero(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[string]()
	cola.Encolar("Boquita Pasion")
	cola.Desencolar()

	require.Panics(t, func() { cola.VerPrimero() })
}

func TestColaDesencoladaDevuelvePanicAlDesencolar(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[string]()
	cola.Encolar("Boquita Pasion")
	cola.Desencolar()

	require.Panics(t, func() { cola.Desencolar() })
}
