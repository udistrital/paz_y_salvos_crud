-- Documentación del modelo. Ejecutar después de 001 y 002; puede repetirse.
BEGIN;

COMMENT ON TABLE paz_y_salvos.solicitud_grado IS
    'Inscripción del estudiante a grado para un programa y periodo; vincula el proceso con las actividades calendarizadas de inscripción y aprobación.';
COMMENT ON TABLE paz_y_salvos.formulario_solicitud_grado IS
    'Versiones del formulario de inscripción. Cada envío conserva su contenido; las subsanaciones se preparan en una nueva versión.';
COMMENT ON TABLE paz_y_salvos.soporte_solicitud_grado IS
    'Referencias a documentos de una versión del formulario, con estado actual de revisión y vínculo al soporte anterior cuando corresponda.';
COMMENT ON TABLE paz_y_salvos.historial_soporte_solicitud_grado IS
    'Trazabilidad de cargue, observaciones, subsanaciones y decisiones sobre un soporte, identificando actor y unidad revisora.';
COMMENT ON TABLE paz_y_salvos.historial_solicitud_grado IS
    'Trazabilidad de estados de la solicitud, incluida radicación, revisión y check final de Secretaría Académica.';
COMMENT ON TABLE paz_y_salvos.paz_salvo IS
    'Validación individual de una dependencia para una solicitud. El conjunto de validaciones forma la vista lógica del semáforo; ORC no participa.';
COMMENT ON TABLE paz_y_salvos.historial_paz_salvo IS
    'Trazabilidad de decisiones y estados de una validación de Paz y Salvo.';

-- Atributos institucionales comunes a las siete entidades.
DO $$
DECLARE
    tabla text;
BEGIN
    FOREACH tabla IN ARRAY ARRAY[
        'solicitud_grado', 'formulario_solicitud_grado', 'soporte_solicitud_grado',
        'historial_soporte_solicitud_grado', 'historial_solicitud_grado',
        'paz_salvo', 'historial_paz_salvo'
    ] LOOP
        EXECUTE format('COMMENT ON COLUMN paz_y_salvos.%I.id IS %L', tabla,
            'Identificador único interno del registro, generado por la secuencia de la tabla.');
        EXECUTE format('COMMENT ON COLUMN paz_y_salvos.%I.activo IS %L', tabla,
            'Indicador de vigencia lógica del registro. No representa aprobación ni sustituye el estado funcional del proceso.');
        EXECUTE format('COMMENT ON COLUMN paz_y_salvos.%I.fecha_creacion IS %L', tabla,
            'Fecha y hora de creación del registro en hora de Bogotá, almacenada como timestamp sin zona horaria según la convención OAS.');
        EXECUTE format('COMMENT ON COLUMN paz_y_salvos.%I.fecha_modificacion IS %L', tabla,
            'Fecha y hora de Bogotá de la última modificación; inicialmente igual a fecha_creacion. El servicio debe actualizarla en cada modificación.');
    END LOOP;
END;
$$;

COMMENT ON COLUMN paz_y_salvos.solicitud_grado.tercero_id IS
    'Referencia externa al tercero del estudiante solicitante en terceros_crud; debe obtenerse de la identidad autenticada y validarse en el MID.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.codigo_estudiante IS
    'Identificador académico textual sin estructura fija. Conserva ceros iniciales y códigos legados; no derivar programa ni periodo de sus posiciones. Usar string en las APIs.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.periodo_id IS
    'Referencia externa al periodo académico PA seleccionado para grado en parametros_crud; no necesariamente coincide con el periodo de vinculación.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.programa_academico_id IS
    'Referencia externa a proyecto_academico_institucion.id, identificador del programa utilizado por Calendario; no es el ID Oikos.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.estado_solicitud_id IS
    'Referencia externa al parámetro que representa el estado actual de la solicitud. Resolver por códigos en parametros_crud, sin IDs fijos.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.dependencia_oikos_id IS
    'ID de la dependencia Oikos del programa al crear la solicitud, validado contra la vinculación EST/TV y proyecto_academico_institucion.DependenciaId.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.calendario_evento_inscripcion_id IS
    'Referencia externa a la actividad calendario_evento resuelta por INSC_GRADO dentro de PROC_GRAD, para el programa y periodo. No es el ID del catálogo. Consultar sus fechas vigentes en Calendario.';
COMMENT ON COLUMN paz_y_salvos.solicitud_grado.calendario_evento_aprobacion_id IS
    'Referencia externa a la actividad calendario_evento resuelta por APROB_PAZ_SALVO para el mismo proceso, programa y periodo. Regula la ventana de gestión administrativa y check final.';

COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.solicitud_grado_id IS
    'Llave foránea a la solicitud propietaria de esta versión del formulario.';
COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.version IS
    'Número positivo de versión, único dentro de la solicitud. El backend debe asignar la siguiente versión de manera transaccional.';
COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.contenido IS
    'Objeto JSON con los datos variables del formulario, incluida la identificación documental del director consultado en otra fuente. El MID valida estructura y campos obligatorios según el contrato funcional.';
COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.fecha_radicacion IS
    'Fecha y hora de Bogotá al enviar esta versión; NULL indica borrador. Solo se admite un borrador activo por solicitud. El backend debe impedir sobrescribir el contenido radicado y crear una nueva versión para subsanar.';

COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.solicitud_grado_id IS
    'Llave foránea a la solicitud propietaria del soporte; debe coincidir con la solicitud del formulario y del soporte anterior.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.documento_id IS
    'Referencia externa al registro de documento utilizado por la integración documental. El archivo y sus metadatos se administran fuera de este esquema; no almacena el binario ni una URL.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.tipo_documento_id IS
    'Referencia externa al parámetro que clasifica el soporte requerido para grado; no corresponde al tipo de identificación personal del estudiante.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.estado_soporte_id IS
    'Referencia externa al parámetro del estado actual de revisión del soporte. Su cambio debe registrarse junto con el historial en una transacción.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.tercero_id IS
    'Referencia externa al tercero que registra el soporte. Debe validarse contra la identidad de quien realiza el cargue.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.rol_usuario IS
    'Rol con el que el actor registra el soporte, validado por el backend. El texto almacenado no otorga permisos por sí mismo.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.formulario_solicitud_grado_id IS
    'Versión del formulario a la que pertenece el soporte. La FK compuesta exige que ambos pertenezcan a la misma solicitud.';
COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.soporte_anterior_id IS
    'Referencia al soporte sustituido o reutilizado de la misma solicitud y tipo; NULL en el primer cargue. El backend valida que sea de una versión anterior y evita ciclos.';

COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.soporte_solicitud_grado_id IS
    'Llave foránea al soporte objeto de la actuación; permite identificar su documento, versión del formulario y solicitud.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.tercero_id IS
    'Referencia externa al tercero responsable de la actuación sobre el soporte, identificado por el backend.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.estado_soporte_id IS
    'Referencia externa al parámetro del estado resultante del soporte tras esta actuación.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.observacion IS
    'Detalle de la observación, subsanación o decisión sobre el soporte; puede ser NULL cuando la acción no requiere observación.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.rol_usuario IS
    'Rol validado con el que se realizó la actuación, conservado para trazabilidad; no sustituye la autorización del backend.';
COMMENT ON COLUMN paz_y_salvos.historial_soporte_solicitud_grado.dependencia_oikos_id IS
    'Referencia externa a la unidad Oikos del Programa o Secretaría que revisa el soporte; puede ser NULL para actuaciones del estudiante.';

COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.solicitud_grado_id IS
    'Llave foránea a la solicitud cuyo cambio de estado se registra.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.tercero_id IS
    'Referencia externa al tercero responsable del cambio de estado, identificado por el backend.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.estado_solicitud_id IS
    'Referencia externa al parámetro del estado resultante de la solicitud. Registrar junto con el estado actual en una transacción.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.justificacion IS
    'Motivo del cambio de estado u observación general, incluida la decisión de Secretaría Académica cuando corresponda.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.rol_usuario IS
    'Rol validado del actor al cambiar el estado. El check final requiere autorización de Secretaría Académica en el backend.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.formulario_solicitud_grado_id IS
    'Versión objeto de la radicación, revisión o check final. Puede ser NULL para actuaciones previas al formulario; la FK compuesta evita mezclar solicitudes.';
COMMENT ON COLUMN paz_y_salvos.historial_solicitud_grado.dependencia_oikos_id IS
    'Referencia externa a la unidad Oikos responsable de la decisión administrativa; puede ser NULL para actuaciones del estudiante. El MID valida su alcance.';

COMMENT ON COLUMN paz_y_salvos.paz_salvo.solicitud_grado_id IS
    'Llave foránea a la solicitud evaluada por esta validación de Paz y Salvo.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.tipo_paz_salvo_id IS
    'Referencia externa al parámetro que identifica la validación de dependencia. Es único por solicitud; ORC no se configura.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.estado_paz_salvo_id IS
    'Referencia externa al parámetro del estado actual de esta validación; el MID resuelve su significado por código.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.tercero_id IS
    'Referencia externa al tercero responsable de la última evaluación; puede ser NULL mientras no exista evaluación.';
COMMENT ON COLUMN paz_y_salvos.paz_salvo.observacion IS
    'Observación actual de la dependencia sobre esta validación. Las decisiones anteriores se conservan en historial_paz_salvo.';

COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.paz_salvo_id IS
    'Llave foránea a la validación de Paz y Salvo cuyo cambio se registra.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.tercero_id IS
    'Referencia externa al tercero responsable de la actuación, identificado por el backend.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.estado_paz_salvo_id IS
    'Referencia externa al parámetro del estado resultante de la validación. Registrar junto con el estado actual en una transacción.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.justificacion IS
    'Motivo, observación o sustento de la decisión de la dependencia.';
COMMENT ON COLUMN paz_y_salvos.historial_paz_salvo.rol_usuario IS
    'Rol validado con el que se tomó la decisión; conserva contexto de auditoría y no reemplaza los controles de acceso.';

-- Impide considerar completa la documentación si queda una columna sin comentario.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        JOIN pg_attribute a ON a.attrelid = c.oid
        WHERE n.nspname = 'paz_y_salvos' AND c.relkind IN ('r', 'p')
          AND a.attnum > 0 AND NOT a.attisdropped
          AND nullif(btrim(col_description(c.oid, a.attnum)), '') IS NULL
    ) THEN
        RAISE EXCEPTION 'Existen columnas sin comentario en paz_y_salvos';
    END IF;
END;
$$;
COMMIT;
