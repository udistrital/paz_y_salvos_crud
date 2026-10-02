-- Ejecutar después de 004. Los historiales son la única fuente de estados.
-- Esta iteración requiere tablas funcionales vacías: no inventa auditorías.
BEGIN;
SET LOCAL lock_timeout = '10s';
LOCK TABLE paz_y_salvos.solicitud_grado, paz_y_salvos.formulario_solicitud_grado,
    paz_y_salvos.soporte_solicitud_grado, paz_y_salvos.paz_salvo,
    paz_y_salvos.historial_solicitud_grado, paz_y_salvos.historial_soporte_solicitud_grado,
    paz_y_salvos.historial_paz_salvo IN ACCESS EXCLUSIVE MODE;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM paz_y_salvos.solicitud_grado)
        OR EXISTS (SELECT 1 FROM paz_y_salvos.soporte_solicitud_grado)
        OR EXISTS (SELECT 1 FROM paz_y_salvos.paz_salvo) THEN
        RAISE EXCEPTION 'Hay datos funcionales: preparar su migración antes de normalizar los historiales';
    END IF;
END;
$$;

ALTER TABLE paz_y_salvos.solicitud_grado DROP COLUMN estado_solicitud_id;
ALTER TABLE paz_y_salvos.soporte_solicitud_grado
    DROP COLUMN tercero_id, DROP COLUMN rol_usuario, DROP COLUMN estado_soporte_id;
ALTER TABLE paz_y_salvos.paz_salvo
    DROP COLUMN tercero_id, DROP COLUMN observacion, DROP COLUMN estado_paz_salvo_id;

-- El CRUD consulta la última actuación activa con estos índices, sin vistas.
CREATE INDEX ix_historial_solicitud_grado_actual
    ON paz_y_salvos.historial_solicitud_grado (solicitud_grado_id, fecha_creacion DESC, id DESC)
    WHERE activo;
CREATE INDEX ix_historial_soporte_grado_actual
    ON paz_y_salvos.historial_soporte_solicitud_grado (soporte_solicitud_grado_id, fecha_creacion DESC, id DESC)
    WHERE activo;
CREATE INDEX ix_historial_paz_salvo_actual
    ON paz_y_salvos.historial_paz_salvo (paz_salvo_id, fecha_creacion DESC, id DESC)
    WHERE activo;

COMMENT ON TABLE paz_y_salvos.paz_salvo IS
    'Validación de una dependencia para una solicitud. Estado, responsable y observación se conservan únicamente en historial_paz_salvo. ORC no participa.';
COMMENT ON TABLE paz_y_salvos.soporte_solicitud_grado IS
    'Documento de una versión del formulario. El cargue y las revisiones se registran en el historial, única fuente de actores, roles y estados.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.estado_solicitud_id IS
    'ID externo del parámetro de estado alcanzado por esta actuación. La última actuación activa determina el estado actual; no se duplica en solicitud_grado.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.estado_soporte_id IS
    'ID externo del parámetro de estado alcanzado por esta actuación. La última actuación activa determina el estado actual; incluir una actuación inicial al registrar el soporte.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.estado_paz_salvo_id IS
    'ID externo del parámetro de estado alcanzado por esta actuación. La última actuación activa determina el estado actual; no se duplica en paz_salvo.';
COMMIT;
