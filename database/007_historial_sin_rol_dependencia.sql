-- Ejecutar una sola vez después de 006, en la base local de desarrollo.
-- La autorización del actor se resuelve con las fuentes de alcance del MID.
-- No exige que el revisor tenga una vinculación en terceros.vinculacion.
BEGIN;
SET LOCAL lock_timeout = '10s';
LOCK TABLE paz_y_salvos.historial_solicitud_grado,
    paz_y_salvos.historial_soporte_solicitud_grado,
    paz_y_salvos.historial_paz_salvo IN ACCESS EXCLUSIVE MODE;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM paz_y_salvos.historial_solicitud_grado)
        OR EXISTS (SELECT 1 FROM paz_y_salvos.historial_soporte_solicitud_grado)
        OR EXISTS (SELECT 1 FROM paz_y_salvos.historial_paz_salvo) THEN
        RAISE EXCEPTION 'Hay historiales registrados: revisar su conservación antes de retirar rol y dependencia revisora';
    END IF;
END;
$$;

ALTER TABLE paz_y_salvos.historial_solicitud_grado
    DROP COLUMN rol_usuario,
    DROP COLUMN dependencia_oikos_id;
ALTER TABLE paz_y_salvos.historial_soporte_solicitud_grado
    DROP COLUMN rol_usuario,
    DROP COLUMN dependencia_oikos_id;
ALTER TABLE paz_y_salvos.historial_paz_salvo
    DROP COLUMN rol_usuario;

COMMENT ON TABLE paz_y_salvos.historial_solicitud_grado IS
    'Actuaciones sobre la solicitud, incluida radicación, revisión y check final. Conserva tercero, estado, justificación y fecha; el alcance se valida con las fuentes externas del MID sin persistir rol ni dependencia revisora.';
COMMENT ON TABLE paz_y_salvos.historial_soporte_solicitud_grado IS
    'Actuaciones de cargue, revisión y subsanación sobre cada soporte. Conserva el tercero responsable y la observación, sin requerir una vinculación de dependencia del revisor en Terceros.';
COMMENT ON TABLE paz_y_salvos.historial_paz_salvo IS
    'Decisiones sobre una validación de Paz y Salvo. Una corrección agrega otra actuación con justificación y mantiene la decisión anterior; no se persiste rol ni dependencia del actor.';
COMMENT ON TABLE paz_y_salvos.soporte_solicitud_grado IS
    'Documento de una versión del formulario. El cargue y las revisiones se registran en el historial, única fuente de actores, observaciones y estados.';

COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.tercero_id IS
    'ID externo del tercero responsable de esta actuación. No es su número de documento ni un ID de vinculación. La identidad y el alcance de la operación los valida el backend.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.tercero_id IS
    'ID externo del tercero que carga, revisa o subsana el soporte. No exige una fila en terceros.vinculacion para revisores; no almacenar aquí el número de identificación personal.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.tercero_id IS
    'ID externo del tercero que emite o corrige la decisión. El tipo de validación se obtiene de paz_salvo; sus permisos se verifican por el backend usando la fuente de alcance correspondiente.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.dependencia_oikos_id IS
    'ID Oikos del programa del estudiante, conservado junto con programa_academico_id. Permite contrastar la vinculación estudiantil y filtrar solicitudes por alcance; no representa la dependencia de un funcionario revisor.';

COMMIT;
