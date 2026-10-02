-- Ajuste semántico posterior a 009: Extensión tiene check propio antes del cierre.
-- Solo comentarios; el tipo TPS_EXTENSION aún debe provisionarse en Parámetros.
BEGIN;
SET LOCAL lock_timeout = '10s';

COMMENT ON TABLE paz_y_salvos.paz_salvo IS
    'Checks independientes posteriores a aprobación documental: Coordinación, Financiero, Biblioteca, Laboratorios, Bienestar, Urelinter y Extensión; Secretaría Académica emite el octavo check final. Una fila por solicitud y tipo; sin ORC.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.tipo_paz_salvo_id IS
    'ID externo del tipo de check en Parámetros. TPS_EXTENSION corresponde al rol EXTENSION, que solo decide su check y observaciones; Secretaría tiene el check final separado de su revisión documental. Alcance, autorización y orden se validan en backend.';

COMMIT;
