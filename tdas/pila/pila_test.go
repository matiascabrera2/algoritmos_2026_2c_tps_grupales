package pila_test

import (
	TDAPila "tdas/pila"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	PRIMER_BLOQUE    = 100000 // Cantidad de elementos a apilar en el primer bloque
	SEGUNDO_BLOQUE   = 200000 // Idem. ant. pero en el segundo bloque
	CANT_A_DESAPILAR = 50000  // Cantidad de elementos a desapilar en el primer bloque
)

/* Funciones auxiliarles genéricas */
func apilarYDesapilarEnOrden[T any](t *testing.T, valores []T) {
	pila := TDAPila.CrearPilaDinamica[T]()

	for _, v := range valores {
		pila.Apilar(v)
	}

	for i := len(valores) - 1; i >= 0; i-- {
		require.Equal(t, valores[i], pila.VerTope())
		require.Equal(t, valores[i], pila.Desapilar())
	}

	require.True(t, pila.EstaVacia())
}

func pilaVaciadaSeComportaComoRecienCreada[T any](t *testing.T, valores []T) {
	pila := TDAPila.CrearPilaDinamica[T]()

	for _, v := range valores {
		pila.Apilar(v)
	}

	for i := len(valores) - 1; i >= 0; i-- {
		require.Equal(t, valores[i], pila.Desapilar())
	}

	require.True(t, pila.EstaVacia())

	// Ahora la pila debería comportarse como recién creada
	for _, v := range valores {
		pila.Apilar(v)
		require.Equal(t, v, pila.VerTope())
		require.Equal(t, v, pila.Desapilar())
	}

	require.True(t, pila.EstaVacia())
}

func apilarBloque(t *testing.T, p TDAPila.Pila[int], desde, hasta int) {
	for i := desde; i < hasta; i++ {
		p.Apilar(i)
		require.Equal(t, i, p.VerTope())
	}
}

func desapilarBloque(t *testing.T, p TDAPila.Pila[int], desde, hasta int) {
	for i := hasta - 1; i >= desde; i-- {
		require.Equal(t, i, p.VerTope())
		require.Equal(t, i, p.Desapilar())
	}
}

/*Tests*/
func TestPilaVacia(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	require.True(t, pila.EstaVacia())
	require.Panics(t, func() { pila.VerTope() })
	require.Panics(t, func() { pila.Desapilar() })
}

func TestApilarYDesapilarConEnteros(t *testing.T) {
	apilarYDesapilarEnOrden(t, []int{1, 2, 3})
}

func TestApilarYDesapilarConStrings(t *testing.T) {
	apilarYDesapilarEnOrden(t, []string{"a", "c", "o", "b"})
}

func TestApilarYDesapilarConFloats(t *testing.T) {
	apilarYDesapilarEnOrden(t, []float64{3.14, 2.71})
}

func TestApilarYDesapilarConBools(t *testing.T) {
	apilarYDesapilarEnOrden(t, []bool{true, false, true})
}
func TestApilarYDesapilarConStructs(t *testing.T) {
	type Auto struct {
		Marca string
		Año   int
	}
	apilarYDesapilarEnOrden(t, []Auto{
		{Marca: "Honda", Año: 2020},
		{Marca: "Tesla", Año: 2021},
		{Marca: "BYD", Año: 2025},
	})
}
func TestPilaVaciadaSeComportaComoRecienCreadaConEnteros(t *testing.T) {
	pilaVaciadaSeComportaComoRecienCreada(t, []int{6, 7, 6, 7})
}

func TestPilaVaciadaSeComportaComoRecienCreadaConStrings(t *testing.T) {
	pilaVaciadaSeComportaComoRecienCreada(t, []string{"Hola", "Mundo", "Aguante", "Boca"})
}

func TestPilaVaciadaSeComportaComoRecienCreadaConFloats(t *testing.T) {
	pilaVaciadaSeComportaComoRecienCreada(t, []float64{3.14, 2.71, 3.14, 2.71, 1.61})
}

func TestPilaVaciadaSeComportaComoRecienCreadaConBools(t *testing.T) {
	pilaVaciadaSeComportaComoRecienCreada(t, []bool{true, false, false, true, false})
}
func TestPilaVaciadaSeComportaComoRecienCreadaConStructs(t *testing.T) {
	type Persona struct {
		Nombre string
		Edad   int
	}
	pilaVaciadaSeComportaComoRecienCreada(t, []Persona{
		{Nombre: "Juan", Edad: 30},
		{Nombre: "Román", Edad: 25},
		{Nombre: "Riquelme", Edad: 40},
	})
}
func TestDeVolumen(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	apilarBloque(t, pila, 0, PRIMER_BLOQUE)
	desapilarBloque(t, pila, PRIMER_BLOQUE-CANT_A_DESAPILAR, PRIMER_BLOQUE)
	apilarBloque(t, pila, PRIMER_BLOQUE, PRIMER_BLOQUE+SEGUNDO_BLOQUE)
	desapilarBloque(t, pila, PRIMER_BLOQUE, PRIMER_BLOQUE+SEGUNDO_BLOQUE)
	desapilarBloque(t, pila, 0, PRIMER_BLOQUE-CANT_A_DESAPILAR)

	require.True(t, pila.EstaVacia())
}

func TestPilaDesapiladaDevuelvePanicAlVerTope(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[string]()
	pila.Apilar("Boquita Pasion")
	pila.Desapilar()

	require.Panics(t, func() { pila.VerTope() })
}

func TestPilaDesapiladaDevuelvePanicAlDesapilar(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[string]()
	pila.Apilar("Boquita Pasion")
	pila.Desapilar()

	require.Panics(t, func() { pila.Desapilar() })
}
