-- Exige referencias verificables a Parámetros y Ubicaciones en formularios radicados.
-- Los borradores permanecen parciales hasta que el MID valida y normaliza el contenido.
BEGIN;
SET LOCAL lock_timeout = '10s';

ALTER TABLE paz_y_salvos.formulario_solicitud_grado
    DROP CONSTRAINT IF EXISTS ck_formulario_solicitud_grado_catalogos;

ALTER TABLE paz_y_salvos.formulario_solicitud_grado
    ADD CONSTRAINT ck_formulario_solicitud_grado_catalogos CHECK (
        fecha_radicacion IS NULL OR
        CASE
            WHEN jsonb_typeof(contenido -> 'modalidad') = 'string'
                AND jsonb_typeof(contenido -> 'lugarExpedicionDocumentoId') = 'number'
                AND jsonb_typeof(contenido -> 'lugarExpedicionDocumento') = 'string'
            THEN btrim(contenido ->> 'modalidad') <> ''
                AND (contenido ->> 'lugarExpedicionDocumentoId') ~ '^[1-9][0-9]*$'
                AND btrim(contenido ->> 'lugarExpedicionDocumento') <> ''
            ELSE false
        END
    );

COMMENT ON CONSTRAINT ck_formulario_solicitud_grado_catalogos
    ON paz_y_salvos.formulario_solicitud_grado IS
    'Al radicar exige el código abreviado de modalidad de parametros_crud e ID y nombre canónico de la ciudad validada por el MID.';

COMMIT;
