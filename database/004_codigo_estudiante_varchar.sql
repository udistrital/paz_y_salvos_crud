-- Permite códigos legados y futuros sin asumir longitud o estructura fija.
BEGIN;
SET LOCAL lock_timeout = '10s';
ALTER TABLE paz_y_salvos.solicitud_grado
    ALTER COLUMN codigo_estudiante TYPE varchar;
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.codigo_estudiante IS
    'Identificador académico textual sin límite de longitud declarado ni estructura fija. Conserva ceros iniciales y códigos legados; no derivar programa ni periodo de sus posiciones. Usar string en las APIs.';
COMMIT;
