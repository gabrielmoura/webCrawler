package db

import (
	"context"
	"errors"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

var sess db.Session

// InitDB inicializa a sessão do banco de dados
func InitDB() error {
	settings, err := postgresql.ParseURL(config.Conf.PostgresURI)
	if err != nil {
		return err
	}
	session, err := postgresql.Open(settings)
	if err != nil {
		return err
	}
	sess = session
	return nil
}

// WritePage insere uma nova página no banco de dados
func WritePage(page *data.Page) error {
	_, err := sess.Collection("pages").Insert(page)
	return err
}

// ReadPage recupera uma página do banco de dados por URL
func ReadPage(url string) (*data.Page, error) {
	var page data.Page
	err := sess.Collection("pages").Find(db.Cond{"url": url}).One(&page)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}

// IsVisited verifica se uma URL foi visitada
func IsVisited(url string) bool {
	count, err := sess.Collection("pages").Find(db.And(
		db.Cond{"url": url},
		db.Cond{"visited": true},
	)).Count()
	return err == nil && count > 0
}

// AllVisited recupera todos os URLs visitados
func AllVisited() ([]string, error) {
	var pages []data.Page
	err := sess.Collection("pages").Find(db.Cond{"visited": true}).All(&pages)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(pages))
	for i, page := range pages {
		urls[i] = page.Url
	}

	return urls, nil
}

// SearchByTitleOrDescription pesquisa páginas por título ou descrição
func SearchByTitleOrDescription(ctx context.Context, searchTerm string) ([]data.PageSearch, error) {
	var pages []data.PageSearch
	err := sess.WithContext(ctx).Collection("pages").Find(db.Or(
		db.Cond{"title": db.Cond{"$ilike": "%" + searchTerm + "%"}},
		db.Cond{"description": db.Cond{"$ilike": "%" + searchTerm + "%"}},
	)).All(&pages)
	return pages, err
}

// SearchByContent pesquisa páginas por conteúdo e ordena por frequência
func SearchByContent(ctx context.Context, searchTerm string) ([]data.PageSearchWithFrequency, error) {
	var pages []data.PageSearchWithFrequency
	query := `
		SELECT url, title, (words->>$2)::int AS frequency
		FROM $1
		WHERE words ? $2
		ORDER BY frequency DESC NULLS LAST;
	`
	rows, err := sess.SQL().QueryContext(ctx, query, "pages", searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var page data.PageSearchWithFrequency
		if err := rows.Scan(&page.Url, &page.Title, &page.Frequency); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, rows.Err()
}

// Search pesquisa páginas por título, descrição ou conteúdo
func Search(ctx context.Context, searchTerm string) ([]data.PageSearch, error) {
	var pages []data.PageSearch
	query := `
		SELECT DISTINCT url, title
		FROM (
			SELECT url, title
			FROM $1
			WHERE (words->>$2)::int IS NOT NULL
			UNION
			SELECT url, title
			FROM $1
			WHERE title ILIKE '%' || $2 || '%'
			OR description ILIKE '%' || $2 || '%'
		) AS combined_results
		ORDER BY url;
	`
	rows, err := sess.SQL().QueryContext(ctx, query, "pages", searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var page data.PageSearch
		if err := rows.Scan(&page.Url, &page.Title); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, rows.Err()
}
