package services

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/udistrital/paz_y_salvos_crud/models"
)

var (
	ErrRevisionDocumentalInvalida = errors.New("revisión documental inválida")
	ErrSolicitudFueraAlcance      = errors.New("solicitud fuera del alcance autorizado")
)

// RevisionDocumentalRepository desacopla las reglas de revisión de Beego ORM
// y de la forma concreta en que la capa de datos ejecuta la transacción.
type RevisionDocumentalRepository interface {
	Consultar(id, dependenciaID int, estadosPermitidos []int) (*models.BorradorInscripcionGrado, error)
	Listar(dependenciaID int, estadosPermitidos []int, limit, offset int, filtros ...int) ([]models.BorradorInscripcionGrado, error)
	Guardar(id, dependenciaID, estadoSolicitudID int, entrada models.RevisarDocumentacionGrado) (*models.BorradorInscripcionGrado, error)
}

func idsRevisionValidos(ids []int, cantidad int) bool {
	if len(ids) != cantidad {
		return false
	}
	vistos := make(map[int]bool, len(ids))
	for _, id := range ids {
		if id <= 0 || vistos[id] {
			return false
		}
		vistos[id] = true
	}
	return true
}

func textoRevisionValido(texto string, obligatorio bool) bool {
	texto = strings.TrimSpace(texto)
	return (!obligatorio || texto != "") && utf8.RuneCountInString(texto) <= 500
}

func idsSoportesRevision(soportes []models.DecisionSoporteGrado) []int {
	ids := make([]int, len(soportes))
	for i := range soportes {
		ids[i] = soportes[i].SoporteId
	}
	return ids
}

func validarRevisionDocumental(entrada *models.RevisarDocumentacionGrado) error {
	entrada.Justificacion = strings.TrimSpace(entrada.Justificacion)
	if entrada.FormularioId <= 0 || entrada.TerceroId <= 0 || entrada.EstadoRadicadaId <= 0 ||
		entrada.EstadoObservadaId <= 0 || entrada.EstadoDocumentacionAprobadaId <= 0 ||
		entrada.EstadoSoporteObservadoId <= 0 || entrada.EstadoSoporteAprobadoId <= 0 ||
		entrada.EstadoObservadaId == entrada.EstadoDocumentacionAprobadaId ||
		entrada.EstadoSoporteObservadoId == entrada.EstadoSoporteAprobadoId || len(entrada.Soportes) < 3 || len(entrada.Soportes) > 4 ||
		!idsRevisionValidos(idsSoportesRevision(entrada.Soportes), len(entrada.Soportes)) {
		return ErrRevisionDocumentalInvalida
	}
	if !textoRevisionValido(entrada.Justificacion, !entrada.Aprobada) {
		return fmt.Errorf("%w: la devolución requiere un comentario general válido", ErrRevisionDocumentalInvalida)
	}
	observados := 0
	for i := range entrada.Soportes {
		entrada.Soportes[i].Observacion = strings.TrimSpace(entrada.Soportes[i].Observacion)
		decision := entrada.Soportes[i]
		if !textoRevisionValido(decision.Observacion, false) {
			return ErrRevisionDocumentalInvalida
		}
		switch decision.EstadoSoporteId {
		case entrada.EstadoSoporteAprobadoId:
		case entrada.EstadoSoporteObservadoId:
			observados++
			if entrada.Aprobada || !textoRevisionValido(decision.Observacion, true) {
				return fmt.Errorf("%w: cada soporte observado requiere comentario", ErrRevisionDocumentalInvalida)
			}
		default:
			return ErrRevisionDocumentalInvalida
		}
	}
	if entrada.Aprobada {
		if observados != 0 || entrada.EstadoPazSalvoPendienteId <= 0 || !idsRevisionValidos(entrada.TiposPazSalvoId, 8) {
			return ErrRevisionDocumentalInvalida
		}
	} else if len(entrada.TiposPazSalvoId) != 0 || entrada.EstadoPazSalvoPendienteId != 0 {
		return ErrRevisionDocumentalInvalida
	}
	return nil
}

func ConsultarSolicitudRevisionGrado(repo RevisionDocumentalRepository, id, dependenciaID int, estadosPermitidos []int) (*models.BorradorInscripcionGrado, error) {
	if repo == nil || id <= 0 || dependenciaID <= 0 || len(estadosPermitidos) == 0 || !idsRevisionValidos(estadosPermitidos, len(estadosPermitidos)) {
		return nil, ErrRevisionDocumentalInvalida
	}
	return repo.Consultar(id, dependenciaID, estadosPermitidos)
}

func ListarSolicitudesRevisionGrado(repo RevisionDocumentalRepository, dependenciaID int, estadosPermitidos []int, limit, offset int, filtros ...int) ([]models.BorradorInscripcionGrado, error) {
	if repo == nil || dependenciaID <= 0 || len(estadosPermitidos) == 0 || !idsRevisionValidos(estadosPermitidos, len(estadosPermitidos)) ||
		limit <= 0 || limit > 100 || offset < 0 || len(filtros) > 2 {
		return nil, ErrRevisionDocumentalInvalida
	}
	for _, filtro := range filtros {
		if filtro < 0 {
			return nil, ErrRevisionDocumentalInvalida
		}
	}
	return repo.Listar(dependenciaID, estadosPermitidos, limit, offset, filtros...)
}

func RevisarDocumentacionGrado(repo RevisionDocumentalRepository, id, dependenciaID int, entrada models.RevisarDocumentacionGrado) (*models.BorradorInscripcionGrado, error) {
	if repo == nil || id <= 0 || dependenciaID <= 0 {
		return nil, ErrRevisionDocumentalInvalida
	}
	if err := validarRevisionDocumental(&entrada); err != nil {
		return nil, err
	}
	estadoSolicitud := entrada.EstadoObservadaId
	if entrada.Aprobada {
		estadoSolicitud = entrada.EstadoDocumentacionAprobadaId
	}
	return repo.Guardar(id, dependenciaID, estadoSolicitud, entrada)
}
