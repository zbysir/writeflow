CREATE TABLE IF NOT EXISTS "public"."document"
(
    "id"               BIGSERIAL PRIMARY KEY NOT NULL,
    "book_id"          BIGINT                NULL,
    "category_id"      BIGINT                NULL,
    "title"            VARCHAR(255)          NULL,
    "source"           jsonb                 NULL,
    "embedding_status" VARCHAR(255)          NULL,
    "created_at"       TIMESTAMP             NULL,
    "updated_at"       TIMESTAMP             NULL
);

CREATE TABLE IF NOT EXISTS "public"."fragment"
(
    "id"          BIGSERIAL PRIMARY KEY NOT NULL,
    "document_id" BIGINT                NULL,
    "book_id"     BIGINT                NULL,
    "body"        VARCHAR               NULL,
    "start_index" INTEGER               NULL,
    "end_index"   INTEGER               NULL,
    "vector"      vector(1536),
    "md5"         VARCHAR(255)          NULL,
    "created_at"  TIMESTAMP             NULL,
    "updated_at"  TIMESTAMP             NULL
);

-- drop table if exists "public"."fragment";
-- drop table if exists "public"."document";