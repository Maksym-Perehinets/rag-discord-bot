CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE public.messages (
    id SERIAL PRIMARY KEY,
    channel_id TEXT,
    message_id TEXT,
    author_id TEXT,
    vectorized_message vector(1024)
);

CREATE INDEX IF NOT EXISTS idx_messages_vectorized_version_hnsw_cosine
    ON public.messages
    USING hnsw (vectorized_message vector_cosine_ops);