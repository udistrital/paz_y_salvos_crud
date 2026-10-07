package services

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/udistrital/paz_y_salvos_crud/models"
)

var (
	ErrPazSalvoGradoInvalido     = errors.New("decisión de Paz y Salvo inválida")
	ErrPazSalvoGradoNoEncontrado = errors.New("Paz y Salvo no encontrado")
	ErrPazSalvoGradoConflicto    = errors.New("el Paz y Salvo cambió o no admite la operación")
)

type PazSalvoGradoRepository interface {
	Consultar(solicitudID, estadoDocumentacionAprobadaID int, tiposEsperados []int) (*models.PazSalvosSolicitudGrado, error)
	Decidir(solicitudID int, entrada models.DecidirPazSalvoGrado) (*models.PazSalvosSolicitudGrado, error)
}

type PazSalvoGradoListadoRepository interface {
	Listar(estadoDocumentacionAprobadaID int, tiposEsperados, dependencias []int, limit, offset, periodoID, programaID, terceroID int, codigo string) (*models.PaginaPazSalvosGrado, error)
}

func tiposPazSalvoValidos(entrada models.DecidirPazSalvoGrado) bool {
	if len(entrada.TiposPreviosId) != 7 || entrada.TipoSecretariaId <= 0 {
		return false
	}
	vistos := map[int]bool{entrada.TipoSecretariaId: true}
	tipoPermitido := entrada.TipoPazSalvoId == entrada.TipoSecretariaId
	for _, tipo := range entrada.TiposPreviosId {
		if tipo <= 0 || vistos[tipo] {
			return false
		}
		vistos[tipo] = true
		if tipo == entrada.TipoPazSalvoId {
			tipoPermitido = true
		}
	}
	return tipoPermitido
}

func validarDecisionPazSalvo(entrada *models.DecidirPazSalvoGrado) error {
	entrada.Justificacion = strings.TrimSpace(entrada.Justificacion)
	if entrada.TerceroId <= 0 || entrada.TipoPazSalvoId <= 0 || entrada.EstadoPazSalvoId <= 0 ||
		entrada.EstadoDocumentacionAprobadaId <= 0 || entrada.EstadoPendienteId <= 0 ||
		entrada.EstadoAprobadoId <= 0 || entrada.EstadoDesaprobadoId <= 0 ||
		entrada.EstadoPendienteId == entrada.EstadoAprobadoId ||
		entrada.EstadoPendienteId == entrada.EstadoDesaprobadoId ||
		entrada.EstadoAprobadoId == entrada.EstadoDesaprobadoId || !tiposPazSalvoValidos(*entrada) {
		return ErrPazSalvoGradoInvalido
	}
	if entrada.EstadoPazSalvoId != entrada.EstadoPendienteId && entrada.EstadoPazSalvoId != entrada.EstadoAprobadoId &&
		entrada.EstadoPazSalvoId != entrada.EstadoDesaprobadoId {
		return fmt.Errorf("%w: estado no permitido", ErrPazSalvoGradoInvalido)
	}
	if utf8.RuneCountInString(entrada.Justificacion) > 500 ||
		((entrada.EstadoPazSalvoId == entrada.EstadoPendienteId || entrada.EstadoPazSalvoId == entrada.EstadoDesaprobadoId) && entrada.Justificacion == "") {
		return fmt.Errorf("%w: reabrir o desaprobar requiere una justificación válida", ErrPazSalvoGradoInvalido)
	}
	return nil
}

func ConsultarPazSalvosGrado(repo PazSalvoGradoRepository, solicitudID, estadoDocumentacionAprobadaID int, tiposEsperados []int) (*models.PazSalvosSolicitudGrado, error) {
	if repo == nil || solicitudID <= 0 || estadoDocumentacionAprobadaID <= 0 || !idsRevisionValidos(tiposEsperados, 8) {
		return nil, ErrPazSalvoGradoInvalido
	}
	return repo.Consultar(solicitudID, estadoDocumentacionAprobadaID, tiposEsperados)
}

func ListarPazSalvosGrado(repo PazSalvoGradoListadoRepository, estadoDocumentacionAprobadaID int, tiposEsperados, dependencias []int, limit, offset, periodoID, programaID, terceroID int, codigo string) (*models.PaginaPazSalvosGrado, error) {
	codigo = strings.TrimSpace(codigo)
	if repo == nil || estadoDocumentacionAprobadaID <= 0 || !idsRevisionValidos(tiposEsperados, 8) ||
		limit <= 0 || limit > 100 || offset < 0 || periodoID < 0 || programaID < 0 || terceroID < 0 ||
		utf8.RuneCountInString(codigo) > 50 ||
		(len(dependencias) > 0 && !idsRevisionValidos(dependencias, len(dependencias))) {
		return nil, ErrPazSalvoGradoInvalido
	}
	return repo.Listar(estadoDocumentacionAprobadaID, tiposEsperados, dependencias, limit, offset, periodoID, programaID, terceroID, codigo)
}

func DecidirPazSalvoGrado(repo PazSalvoGradoRepository, solicitudID int, entrada models.DecidirPazSalvoGrado) (*models.PazSalvosSolicitudGrado, error) {
	if repo == nil || solicitudID <= 0 {
		return nil, ErrPazSalvoGradoInvalido
	}
	if err := validarDecisionPazSalvo(&entrada); err != nil {
		return nil, err
	}
	return repo.Decidir(solicitudID, entrada)
}
