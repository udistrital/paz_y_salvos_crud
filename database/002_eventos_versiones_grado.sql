-- Ejecutar una sola vez después de 001_inscripcion_grado.sql.
-- No recrea el esquema ni asigna IDs externos ficticios a datos existentes.
BEGIN;
SET LOCAL lock_timeout = '10s';

ALTER TABLE paz_y_salvos.solicitud_grado
    ADD COLUMN dependencia_oikos_id integer NOT NULL,
    ADD COLUMN calendario_evento_inscripcion_id integer NOT NULL,
    ADD COLUMN calendario_evento_aprobacion_id integer NOT NULL,
    ADD CONSTRAINT ck_solicitud_grado_referencia_externa CHECK (
        dependencia_oikos_id > 0 AND calendario_evento_inscripcion_id > 0
        AND calendario_evento_aprobacion_id > 0
        AND calendario_evento_inscripcion_id <> calendario_evento_aprobacion_id
    );

COMMENT ON COLUMN paz_y_salvos.solicitud_grado.dependencia_oikos_id IS
    'ID Oikos validado contra la vinculación EST/TV y DependenciaId del programa al crear la solicitud.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.calendario_evento_inscripcion_id IS
    'Referencia externa a calendario_evento; resolver INSC_GRADO dentro de PROC_GRAD por programa y periodo. Sin ID fijo ni FK entre servicios.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.calendario_evento_aprobacion_id IS
    'Referencia externa a calendario_evento; resolver APROB_PAZ_SALVO para el mismo proceso, programa y periodo. Las fechas vigentes se consultan en Calendario.';

ALTER TABLE paz_y_salvos.formulario_solicitud_grado
    ADD COLUMN fecha_radicacion timestamp NULL,
    ADD CONSTRAINT ck_formulario_solicitud_grado_contenido CHECK (jsonb_typeof(contenido) = 'object'),
    ADD CONSTRAINT uq_formulario_solicitud_grado_id_solicitud UNIQUE (id, solicitud_grado_id);

CREATE UNIQUE INDEX uq_formulario_solicitud_grado_borrador
    ON paz_y_salvos.formulario_solicitud_grado (solicitud_grado_id)
    WHERE activo AND fecha_radicacion IS NULL;

COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.fecha_radicacion IS
    'NULL indica borrador; al radicar se registra la fecha del servidor. Para subsanar se crea otra versión, sin sobrescribir el envío previo.';

ALTER TABLE paz_y_salvos.soporte_solicitud_grado
    ADD COLUMN formulario_solicitud_grado_id integer NOT NULL,
    ADD COLUMN soporte_anterior_id integer NULL,
    ADD CONSTRAINT uq_soporte_solicitud_grado_id_solicitud_tipo
        UNIQUE (id, solicitud_grado_id, tipo_documento_id),
    ADD CONSTRAINT fk_soporte_solicitud_grado_formulario_solicitud_grado
        FOREIGN KEY (formulario_solicitud_grado_id, solicitud_grado_id)
        REFERENCES paz_y_salvos.formulario_solicitud_grado (id, solicitud_grado_id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    ADD CONSTRAINT fk_soporte_solicitud_grado_soporte_anterior
        FOREIGN KEY (soporte_anterior_id, solicitud_grado_id, tipo_documento_id)
        REFERENCES paz_y_salvos.soporte_solicitud_grado (id, solicitud_grado_id, tipo_documento_id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    ADD CONSTRAINT ck_soporte_solicitud_grado_anterior CHECK (soporte_anterior_id <> id);

CREATE INDEX ix_soporte_solicitud_grado_formulario
    ON paz_y_salvos.soporte_solicitud_grado (formulario_solicitud_grado_id, solicitud_grado_id);
CREATE INDEX ix_soporte_solicitud_grado_anterior
    ON paz_y_salvos.soporte_solicitud_grado (soporte_anterior_id, solicitud_grado_id, tipo_documento_id)
    WHERE soporte_anterior_id IS NOT NULL;

COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.formulario_solicitud_grado_id IS
    'Versión del formulario a la que pertenece este soporte. La FK compuesta impide mezclar solicitudes.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.soporte_anterior_id IS
    'Soporte reemplazado o reutilizado en un nuevo envío, de la misma solicitud y tipo documental. NULL en el primer soporte.';

ALTER TABLE paz_y_salvos.historial_solicitud_grado
    ADD COLUMN formulario_solicitud_grado_id integer NULL,
    ADD COLUMN dependencia_oikos_id integer NULL,
    ADD CONSTRAINT ck_historial_solicitud_grado_dependencia CHECK (dependencia_oikos_id > 0),
    ADD CONSTRAINT fk_historial_solicitud_grado_formulario_solicitud_grado
        FOREIGN KEY (formulario_solicitud_grado_id, solicitud_grado_id)
        REFERENCES paz_y_salvos.formulario_solicitud_grado (id, solicitud_grado_id)
        ON DELETE RESTRICT ON UPDATE RESTRICT;

CREATE INDEX ix_historial_solicitud_grado_formulario
    ON paz_y_salvos.historial_solicitud_grado (formulario_solicitud_grado_id, solicitud_grado_id)
    WHERE formulario_solicitud_grado_id IS NOT NULL;

ALTER TABLE paz_y_salvos.historial_soporte_solicitud_grado
    ADD COLUMN dependencia_oikos_id integer NULL,
    ADD CONSTRAINT ck_historial_soporte_grado_dependencia CHECK (dependencia_oikos_id > 0);

COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.formulario_solicitud_grado_id IS
    'Versión objeto de radicación, revisión o check final; puede ser NULL en eventos previos al formulario.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.dependencia_oikos_id IS
    'Unidad revisora del Programa o Secretaría Académica. La identidad, el rol y el alcance los valida el MID.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.dependencia_oikos_id IS
    'Unidad Oikos responsable de la revisión documental; NULL para acciones del estudiante.';

COMMIT;
