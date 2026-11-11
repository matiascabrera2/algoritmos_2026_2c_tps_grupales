package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	TDAVISITANTES "tdas/visitantes"
)

const (
	_AGREGAR_ARCHIVO   = "agregar_archivo"
	_VER_VISITANTES    = "ver_visitantes"
	_VER_MAS_VISITADOS = "ver_mas_visitados"

	_OK                = "OK"
	_ERROR_EN_COMANDO  = "Error en comando"
	_ENCABEZADO_DOS    = "DoS:"
	_ENCABEZADO_IPS    = "Visitantes:"
	_ENCABEZADO_SITIOS = "Sitios más visitados:"
	_SANGRIA           = "\t"
	_SEPARADOR_VISITAS = " - "
)

var errComando = errors.New("El comando es invalido")

type comando struct {
	parametros int
	ejecutar   func(visitantes TDAVISITANTES.Visitantes, parametros []string, salida io.Writer) error
}

var comandos = map[string]comando{
	_AGREGAR_ARCHIVO:   {1, agregarArchivo},
	_VER_VISITANTES:    {2, verVisitantes},
	_VER_MAS_VISITADOS: {1, verMasVisitados},
}

func ejecutarComando(visitantes TDAVISITANTES.Visitantes, campos []string, salida io.Writer) error {
	if len(campos) == 0 {
		return errComando
	}

	comando, existe := comandos[campos[0]]
	if !existe || len(campos)-1 != comando.parametros {
		return errComando
	}

	if err := comando.ejecutar(visitantes, campos[1:], salida); err != nil {
		return err
	}

	fmt.Fprintln(salida, _OK)
	return nil
}

func nombreDe(campos []string) string {
	if len(campos) == 0 {
		return ""
	}
	return campos[0]
}

func agregarArchivo(visitantes TDAVISITANTES.Visitantes, parametros []string, salida io.Writer) error {
	sospechosas, err := visitantes.AgregarArchivo(parametros[0])
	if err != nil {
		return err
	}

	for _, ip := range sospechosas {
		fmt.Fprintln(salida, _ENCABEZADO_DOS, ip)
	}
	return nil
}

func verVisitantes(visitantes TDAVISITANTES.Visitantes, parametros []string, salida io.Writer) error {
	desde, err := TDAVISITANTES.ParsearIP(parametros[0])
	if err != nil {
		return err
	}
	hasta, err := TDAVISITANTES.ParsearIP(parametros[1])
	if err != nil {
		return err
	}

	fmt.Fprintln(salida, _ENCABEZADO_IPS)
	for _, ip := range visitantes.VerVisitantes(desde, hasta) {
		fmt.Fprintln(salida, _SANGRIA+ip.String())
	}
	return nil
}

func verMasVisitados(visitantes TDAVISITANTES.Visitantes, parametros []string, salida io.Writer) error {
	n, err := strconv.Atoi(parametros[0])
	if err != nil || n < 0 {
		return errComando
	}

	fmt.Fprintln(salida, _ENCABEZADO_SITIOS)
	for _, sitio := range visitantes.VerMasVisitados(n) {
		fmt.Fprintln(salida, _SANGRIA+sitio.Ruta+_SEPARADOR_VISITAS+strconv.Itoa(sitio.Visitas))
	}
	return nil
}
