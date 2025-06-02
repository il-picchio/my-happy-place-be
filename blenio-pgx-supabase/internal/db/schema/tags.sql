CREATE TABLE tags (
  id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  name        TEXT        NOT NULL UNIQUE,
  type        TEXT        NOT NULL
                   CHECK (type IN ('Category','Attribute','Feature')),
  display_name   JSONB    NOT NULL, -- e.g. { "en": "Restaurant", "it": "Ristorante", "de": "Restaurant", "fr": "Restaurant" }
  description JSONB       NOT NULL,
  deprecated  BOOLEAN     NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
