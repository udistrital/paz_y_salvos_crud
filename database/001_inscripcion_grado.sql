BEGIN;

DROP SCHEMA IF EXISTS paz_y_salvos CASCADE;
CREATE SCHEMA paz_y_salvos;

CREATE TABLE paz_y_salvos.solicitud_grado (
    id serial4 NOT NULL,
    tercero_id int4 NOT NULL,
    codigo_estudiante varchar(20) NOT NULL,
    periodo_id int4 NOT NULL,
    programa_academico_id int4 NOT NULL,
    estado_solicitud_id int4 NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_solicitud_grado PRIMARY KEY (id)
);

CREATE UNIQUE INDEX uq_solicitud_grado_estudiante_periodo_programa
    ON paz_y_salvos.solicitud_grado (
        codigo_estudiante,
        periodo_id,
        programa_academico_id
    )
    WHERE activo;

CREATE INDEX ix_solicitud_grado_estado
    ON paz_y_salvos.solicitud_grado (estado_solicitud_id);

CREATE INDEX ix_solicitud_grado_programa_periodo
    ON paz_y_salvos.solicitud_grado (programa_academico_id, periodo_id);

COMMENT ON COLUMN paz_y_salvos.solicitud_grado.estado_solicitud_id IS
    'Identificador del estado administrado por parametros_crud.';

CREATE TABLE paz_y_salvos.formulario_solicitud_grado (
    id serial4 NOT NULL,
    solicitud_grado_id int4 NOT NULL,
    version int2 NOT NULL DEFAULT 1,
    contenido jsonb NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_formulario_solicitud_grado PRIMARY KEY (id),
    CONSTRAINT fk_formulario_solicitud_grado_solicitud_grado
        FOREIGN KEY (solicitud_grado_id)
        REFERENCES paz_y_salvos.solicitud_grado (id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT uq_formulario_solicitud_grado_version
        UNIQUE (solicitud_grado_id, version),
    CONSTRAINT ck_formulario_solicitud_grado_version
        CHECK (version > 0)
);

COMMENT ON COLUMN paz_y_salvos.formulario_solicitud_grado.contenido IS
    'Datos variables del formulario, incluido el documento del director consultado en una fuente externa.';

CREATE TABLE paz_y_salvos.soporte_solicitud_grado (
    id serial4 NOT NULL,
    solicitud_grado_id int4 NOT NULL,
    documento_id int4 NOT NULL,
    tipo_documento_id int4 NOT NULL,
    estado_soporte_id int4 NOT NULL,
    tercero_id int4 NOT NULL,
    rol_usuario varchar(100) NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_soporte_solicitud_grado PRIMARY KEY (id),
    CONSTRAINT fk_soporte_solicitud_grado_solicitud_grado
        FOREIGN KEY (solicitud_grado_id)
        REFERENCES paz_y_salvos.solicitud_grado (id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX ix_soporte_solicitud_grado_solicitud
    ON paz_y_salvos.soporte_solicitud_grado (solicitud_grado_id);

COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.documento_id IS
    'Identificador del documento almacenado en el gestor documental.';

COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.tipo_documento_id IS
    'Identificador del tipo de documento administrado por parametros_crud.';

COMMENT ON COLUMN paz_y_salvos.soporte_solicitud_grado.estado_soporte_id IS
    'Identificador del estado del soporte administrado por parametros_crud.';

CREATE TABLE paz_y_salvos.historial_soporte_solicitud_grado (
    id serial4 NOT NULL,
    soporte_solicitud_grado_id int4 NOT NULL,
    tercero_id int4 NOT NULL,
    estado_soporte_id int4 NOT NULL,
    observacion varchar(500) NULL,
    rol_usuario varchar(100) NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_historial_soporte_solicitud_grado PRIMARY KEY (id),
    CONSTRAINT fk_historial_soporte_solicitud_grado_soporte_solicitud_grado
        FOREIGN KEY (soporte_solicitud_grado_id)
        REFERENCES paz_y_salvos.soporte_solicitud_grado (id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX ix_historial_soporte_solicitud_grado_soporte
    ON paz_y_salvos.historial_soporte_solicitud_grado (soporte_solicitud_grado_id);

COMMENT ON TABLE paz_y_salvos.historial_soporte_solicitud_grado IS
    'Registra observaciones, subsanaciones y aprobaciones realizadas por Programa Académico o Secretaría Académica.';

CREATE TABLE paz_y_salvos.historial_solicitud_grado (
    id serial4 NOT NULL,
    solicitud_grado_id int4 NOT NULL,
    tercero_id int4 NOT NULL,
    estado_solicitud_id int4 NOT NULL,
    justificacion varchar(500) NULL,
    rol_usuario varchar(100) NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_historial_solicitud_grado PRIMARY KEY (id),
    CONSTRAINT fk_historial_solicitud_grado_solicitud_grado
        FOREIGN KEY (solicitud_grado_id)
        REFERENCES paz_y_salvos.solicitud_grado (id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX ix_historial_solicitud_grado_solicitud
    ON paz_y_salvos.historial_solicitud_grado (solicitud_grado_id);

COMMENT ON TABLE paz_y_salvos.historial_solicitud_grado IS
    'Incluye el check final emitido por Secretaría Académica mediante el estado parametrizado correspondiente.';

CREATE TABLE paz_y_salvos.paz_salvo (
    id serial4 NOT NULL,
    solicitud_grado_id int4 NOT NULL,
    tipo_paz_salvo_id int4 NOT NULL,
    estado_paz_salvo_id int4 NOT NULL,
    tercero_id int4 NULL,
    observacion varchar(500) NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_paz_salvo PRIMARY KEY (id),
    CONSTRAINT fk_paz_salvo_solicitud_grado
        FOREIGN KEY (solicitud_grado_id)
        REFERENCES paz_y_salvos.solicitud_grado (id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT uq_paz_salvo_solicitud_tipo
        UNIQUE (solicitud_grado_id, tipo_paz_salvo_id)
);

CREATE INDEX ix_paz_salvo_estado
    ON paz_y_salvos.paz_salvo (estado_paz_salvo_id);

COMMENT ON COLUMN paz_y_salvos.paz_salvo.tipo_paz_salvo_id IS
    'Identificador del tipo de Paz y Salvo administrado por parametros_crud; ORC no se configura.';

COMMENT ON COLUMN paz_y_salvos.paz_salvo.estado_paz_salvo_id IS
    'Identificador del estado de Paz y Salvo administrado por parametros_crud.';

CREATE TABLE paz_y_salvos.historial_paz_salvo (
    id serial4 NOT NULL,
    paz_salvo_id int4 NOT NULL,
    tercero_id int4 NOT NULL,
    estado_paz_salvo_id int4 NOT NULL,
    justificacion varchar(500) NULL,
    rol_usuario varchar(100) NOT NULL,
    activo bool NOT NULL DEFAULT true,
    fecha_creacion timestamp NOT NULL DEFAULT now(),
    fecha_modificacion timestamp NOT NULL DEFAULT now(),
    CONSTRAINT pk_historial_paz_salvo PRIMARY KEY (id),
    CONSTRAINT fk_historial_paz_salvo_paz_salvo
        FOREIGN KEY (paz_salvo_id)
        REFERENCES paz_y_salvos.paz_salvo (id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX ix_historial_paz_salvo_paz_salvo
    ON paz_y_salvos.historial_paz_salvo (paz_salvo_id);

COMMIT;
