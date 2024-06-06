## Buscando páginas que contêm a palavra "vida" no título ou na descrição.
```sql
SELECT url, title 
FROM pages
WHERE title ILIKE '%' || 'vida' || '%'  -- Case-insensitive search in title
   OR description ILIKE '%' || 'vida' || '%'; -- Case-insensitive search in description
```

## Buscando páginas que contêm a palavra "vida" no conteúdo.
```sql
SELECT url, title, (words->>'vida')::int AS frequency -- Use ->> for text extraction
FROM pages
WHERE words ? 'vida' 
ORDER BY frequency DESC NULLS LAST;
```

## Buscando páginas que contêm a palavra "vida" no título, descrição ou conteúdo.
```sql
SELECT DISTINCT url, title
FROM (
    SELECT url, title
    FROM pages
    WHERE (words->>'vida')::int IS NOT NULL
    UNION
    SELECT url, title
    FROM pages
    WHERE title ILIKE '%' || 'vida' || '%'
       OR description ILIKE '%' || 'vida' || '%'
) AS combined_results
ORDER BY url;
```
## Criando a tabela de páginas.
```sql
CREATE TABLE pages
(
    url         TEXT PRIMARY KEY,
    links       TEXT[],
    title       TEXT,
    description TEXT,
    meta        TEXT[],
    visited     BOOLEAN,
    timestamp   TIMESTAMP WITH TIME ZONE,
    words       JSONB
);

CREATE INDEX idx_words_gin ON pages USING GIN (words);
```