-- Para bases que recibieron la primera versión de 005 con vistas.
-- En instalaciones nuevas no se crean vistas; esta migración es inocua.
BEGIN;
SET LOCAL lock_timeout = '10s';
DROP VIEW IF EXISTS paz_y_salvos.solicitud_grado_actual;
DROP VIEW IF EXISTS paz_y_salvos.soporte_solicitud_grado_actual;
DROP VIEW IF EXISTS paz_y_salvos.paz_salvo_actual;
COMMIT;
