-- Actualiza la semántica del modelo después de 007. Solo comentarios.
-- Los códigos de Parámetros están propuestos; este script no los crea.
BEGIN;
SET LOCAL lock_timeout = '10s';

COMMENT ON TABLE paz_y_salvos.solicitud_grado IS
    'Inscripción del estudiante para un programa y periodo. Su revisión documental es distinta de los checks de Paz y Salvos, relacionados mediante paz_salvo.';
COMMENT ON TABLE paz_y_salvos.historial_solicitud_grado IS
    'Actuaciones de inscripción, radicación, observación y aprobación documental de la solicitud. El check final de Paz y Salvos se registra en historial_paz_salvo para el tipo Secretaría Académica, no como sustitución de esta revisión.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.estado_solicitud_id IS
    'ID externo del estado de inscripción/revisión documental alcanzado por la actuación. La última actuación activa determina ese estado, sin representar por sí misma aprobación de Coordinación ni cierre de Paz y Salvos.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.formulario_solicitud_grado_id IS
    'Versión objeto de radicación o revisión documental. Puede ser NULL antes de crear el formulario. La FK compuesta evita mezclar solicitudes; no almacena el check final de Paz y Salvos.';

COMMENT ON TABLE paz_y_salvos.historial_soporte_solicitud_grado IS
    'Cargue, observaciones, subsanaciones y aprobación de cada soporte de inscripción. Aprobar un archivo no equivale a aprobar una validación de Paz y Salvo; conserva actor, estado, observación y fechas.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.estado_soporte_id IS
    'ID externo del estado documental resultante para este soporte. Su última actuación activa determina la revisión del archivo, independientemente del check de Coordinación y del check final de Secretaría.';

COMMENT ON TABLE paz_y_salvos.paz_salvo IS
    'Validaciones de Coordinación, Financiero, Biblioteca, Laboratorios, Bienestar y Urelinter, más la validación final de Secretaría Académica, para una solicitud. Una fila por solicitud y tipo; no contiene ORC ni duplica los estados de los documentos.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.tipo_paz_salvo_id IS
    'ID externo del tipo de validación en Parámetros. Incluye Coordinación (reemplaza el concepto Académico) y Secretaría Académica como check final independiente de la aprobación documental. La condición de cierre y autorización se valida en backend.';
COMMENT ON TABLE paz_y_salvos.historial_paz_salvo IS
    'Decisiones sobre cada check, incluido el final de Secretaría Académica. Estados, actores, motivos y fechas son independientes de la revisión documental. Las correcciones agregan actuaciones conservando las decisiones anteriores.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.estado_paz_salvo_id IS
    'ID externo del estado resultante de un check. La última actuación activa determina su estado; para el tipo Secretaría Académica su aprobación constituye el check final de Paz y Salvos. No se deduce de la aprobación de soportes.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.tercero_id IS
    'ID externo del tercero que decide o corrige un check. Para el check final debe ser el actor autorizado de Secretaría Académica. Su identidad y alcance los valida el backend, sin persistir rol ni dependencia revisora.';

COMMIT;
