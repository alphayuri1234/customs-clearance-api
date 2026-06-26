-- DDL Customs Clearance API
-- Target database: PostgreSQL
-- Jalankan di database customs_clearance.

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'User',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT chk_users_role CHECK (role IN ('User', 'Officer'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS officers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    nip TEXT,
    position TEXT,
    CONSTRAINT fk_officers_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_officers_user_id ON officers (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_officers_nip ON officers (nip);

CREATE TABLE IF NOT EXISTS countries (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(3) NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_countries_code ON countries (code);

CREATE TABLE IF NOT EXISTS ports (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name TEXT NOT NULL,
    country_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_ports_country
        FOREIGN KEY (country_id)
        REFERENCES countries (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ports_code ON ports (code);
CREATE INDEX IF NOT EXISTS idx_ports_country_id ON ports (country_id);

CREATE TABLE IF NOT EXISTS commodities (
    id BIGSERIAL PRIMARY KEY,
    hs_code VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    import_duty_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    vat_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_commodities_hs_code ON commodities (hs_code);

CREATE TABLE IF NOT EXISTS clearances (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    commodity_id BIGINT NOT NULL,
    port_id BIGINT NOT NULL,
    valuation DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'SUBMITTED',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_clearances_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_clearances_commodity
        FOREIGN KEY (commodity_id)
        REFERENCES commodities (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_clearances_port
        FOREIGN KEY (port_id)
        REFERENCES ports (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT chk_clearances_status CHECK (
        status IN (
            'SUBMITTED',
            'INSPECTION',
            'INSPECTION_PASSED',
            'APPROVED',
            'RELEASED',
            'HOLD',
            'GATE_OUT'
        )
    )
);

CREATE INDEX IF NOT EXISTS idx_clearances_user_id ON clearances (user_id);
CREATE INDEX IF NOT EXISTS idx_clearances_commodity_id ON clearances (commodity_id);
CREATE INDEX IF NOT EXISTS idx_clearances_port_id ON clearances (port_id);
CREATE INDEX IF NOT EXISTS idx_clearances_status ON clearances (status);

CREATE TABLE IF NOT EXISTS risk_profiles (
    id BIGSERIAL PRIMARY KEY,
    clearance_id BIGINT NOT NULL,
    level TEXT NOT NULL DEFAULT 'LOW',
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    reason TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_risk_profiles_clearance
        FOREIGN KEY (clearance_id)
        REFERENCES clearances (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT chk_risk_profiles_level CHECK (level IN ('HIGH', 'LOW'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_risk_profiles_clearance_id ON risk_profiles (clearance_id);
CREATE INDEX IF NOT EXISTS idx_risk_profiles_level ON risk_profiles (level);

CREATE TABLE IF NOT EXISTS inspection_results (
    id BIGSERIAL PRIMARY KEY,
    clearance_id BIGINT NOT NULL,
    officer_id BIGINT NOT NULL,
    result TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_inspection_results_clearance
        FOREIGN KEY (clearance_id)
        REFERENCES clearances (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_inspection_results_officer
        FOREIGN KEY (officer_id)
        REFERENCES officers (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT chk_inspection_results_result CHECK (result IN ('PASS', 'FAIL'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_inspection_results_clearance_id ON inspection_results (clearance_id);
CREATE INDEX IF NOT EXISTS idx_inspection_results_officer_id ON inspection_results (officer_id);

CREATE TABLE IF NOT EXISTS release_orders (
    id BIGSERIAL PRIMARY KEY,
    clearance_id BIGINT NOT NULL,
    release_no TEXT NOT NULL,
    officer_id BIGINT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_release_orders_clearance
        FOREIGN KEY (clearance_id)
        REFERENCES clearances (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_release_orders_officer
        FOREIGN KEY (officer_id)
        REFERENCES officers (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_release_orders_clearance_id ON release_orders (clearance_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_release_orders_release_no ON release_orders (release_no);
CREATE INDEX IF NOT EXISTS idx_release_orders_officer_id ON release_orders (officer_id);
