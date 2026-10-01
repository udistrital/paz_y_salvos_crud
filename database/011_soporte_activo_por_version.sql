-- Posterior a 010. Conserva todos los soportes e historiales existentes.
-- Ante duplicados activos se detiene; no decide qué documento debe conservarse.
BEGIN;
SET LOCAL lock_timeout = '10s';
CREATE UNIQUE INDEX IF NOT EXISTS uq_soporte_grado_version_tipo_activo
    ON paz_y_salvos.soporte_solicitud_grado (formulario_solicitud_grado_id, tipo_documento_id)
    WHERE activo;
COMMIT;
