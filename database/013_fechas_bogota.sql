BEGIN;

ALTER TABLE paz_y_salvos.solicitud_grado
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.formulario_solicitud_grado
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.soporte_solicitud_grado
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.historial_soporte_solicitud_grado
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.historial_solicitud_grado
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.paz_salvo
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

ALTER TABLE paz_y_salvos.historial_paz_salvo
    ALTER COLUMN fecha_creacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP),
    ALTER COLUMN fecha_modificacion SET DEFAULT timezone('America/Bogota', CURRENT_TIMESTAMP);

COMMIT;
