CREATE TABLE IF NOT EXISTS public.users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account VARCHAR(200) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.stats (
    short_url VARCHAR(255) PRIMARY KEY,
    today_count INT NOT NULL DEFAULT 0,
    yesterday_count INT NOT NULL DEFAULT 0,
    last_7_days_count INT NOT NULL DEFAULT 0,
    monthly_count INT NOT NULL DEFAULT 0,
    total_count INT NOT NULL DEFAULT 0,
    d_today_count INT NOT NULL DEFAULT 0,
    d_yesterday_count INT NOT NULL DEFAULT 0,
    d_last_7_days_count INT NOT NULL DEFAULT 0,
    d_monthly_count INT NOT NULL DEFAULT 0,
    d_total_count INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public.access_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip VARCHAR(45) NOT NULL
);


CREATE TABLE IF NOT EXISTS public.stats_sum (
    stats_key VARCHAR(255) PRIMARY KEY,
    stats_value INT NOT NULL DEFAULT 0
);


CREATE TABLE IF NOT EXISTS public.daily_stats (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    pv INT NOT NULL DEFAULT 0,
    uv INT NOT NULL DEFAULT 0,
    CONSTRAINT uk_short_url_date UNIQUE (short_url, date)
);


CREATE TABLE IF NOT EXISTS public.short_urls (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    comment TEXT DEFAULT '',
    CONSTRAINT uk_short_urls_short_url UNIQUE (short_url)
);

INSERT INTO public.users (account,password) VALUES ('root','$2a$10$vFZLC5wVjoe1.frG4D3EGe2TsPGqcOrqYcnVCd7zQps6RgKOnG9/q') ON CONFLICT (account) DO NOTHING;