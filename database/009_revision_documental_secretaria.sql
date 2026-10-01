-- Precisión del alcance posterior a 008: solo Secretaría revisa documentos.
-- Actualiza comentarios; no crea parámetros ni implementa autorización SQL.
BEGIN;
SET LOCAL lock_timeout = '10s';

COMMENT ON TABLE paz_y_salvos.solicitud_grado IS
    'Inscripción del estudiante para un programa y periodo. Secretaría Académica aprueba u observa su documentación; solo la aprobación documental inicia la gestión de Paz y Salvos, con checks independientes.';
COMMENT ON TABLE paz_y_salvos.historial_solicitud_grado IS
    'Actuaciones de borrador, radicación, subsanación y decisión documental de la solicitud. Solo Secretaría Académica emite la aprobación u observación documental; el check final de Paz y Salvos queda en historial_paz_salvo.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.estado_solicitud_id IS
    'ID externo del estado de inscripción/revisión documental. Una observación de Secretaría requiere subsanación y nueva radicación del estudiante. Su aprobación documental habilita iniciar los checks, pero no los aprueba.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.tercero_id IS
    'ID externo del actor de esta actuación: estudiante al preparar, radicar o subsanar; funcionario autorizado de Secretaría al observar o aprobar documentación. La identidad y los permisos los valida el backend; no se persiste rol.';

COMMENT ON TABLE paz_y_salvos.historial_soporte_solicitud_grado IS
    'Cargue y subsanación del estudiante y revisión de cada soporte exclusivamente por Secretaría Académica. Registra actor, estado, observación y fechas, conservando las versiones y decisiones anteriores.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.estado_soporte_id IS
    'ID externo del estado del soporte. Secretaría lo observa o aprueba; un archivo corregido se presenta para nueva revisión. La aprobación de documentos no equivale al check de Coordinación ni al check final.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.tercero_id IS
    'ID externo del estudiante que carga o subsana, o del funcionario de Secretaría que revisa. No es su documento de identidad ni requiere vinculación de revisor en terceros.vinculacion.';

COMMENT ON TABLE paz_y_salvos.paz_salvo IS
    'Validaciones iniciadas después de la aprobación documental de Secretaría: Coordinación, Financiero, Biblioteca, Laboratorios, Bienestar y Urelinter, más el check final independiente de Secretaría. Una fila por solicitud y tipo; sin ORC.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.tipo_paz_salvo_id IS
    'ID externo del tipo de check en Parámetros. Coordinación participa en Paz y Salvos, no en la aprobación documental de este alcance. Secretaría tiene un check final distinto de su revisión de documentos; autorización y orden se validan en backend.';

COMMIT;
