package reportrepo

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/internal/core/domain"
	database "mtg-report/internal/sources/databases/mysql"
)

type repository struct {
	db database.Client
}

func New(db database.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertTotalPrice(ctx context.Context) error {
	insertQuery := `
	INSERT INTO prices (old_price, new_price, price_change, last_update)
	SELECT 
		COALESCE((SELECT new_price
		FROM prices
		ORDER BY last_update DESC
		LIMIT 1), 0),
		COALESCE(SUM(subquery.last_price), 0) AS new_price,
		COALESCE(SUM(subquery.last_price), 0) - COALESCE((SELECT new_price
									FROM prices
									ORDER BY last_update DESC
									LIMIT 1), 0),
		NOW() AS last_update
	FROM (
		SELECT 
			cd.last_price AS last_price
		FROM cards c
		LEFT JOIN (
			SELECT *,
				ROW_NUMBER() OVER(PARTITION BY card_id ORDER BY last_update DESC) AS rn
			FROM cards_details
		) cd
		ON c.id = cd.card_id AND cd.rn = 1
	) AS subquery;`

	res, err := r.db.ExecContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("repository failed to exec insert query in insert total price: %w", err)
	}

	err = getRowsAffected(res)
	if err != nil {
		return fmt.Errorf("repository insert total price failed: %w", err)
	}

	return nil
}

func (r *repository) GetCardsReport(ctx context.Context) ([]domain.Cards, error) {
	getCardsQuery := `
	SELECT * FROM 
	(
		SELECT 
			c.id,
			name,
			set_name,
			collector_number,
			foil,
			COALESCE(cd.last_price, 0) as last_price,
			COALESCE(cd.old_price, 0) as old_price,
			COALESCE(cd.price_change, 0) as price_change,
			last_update
		FROM 
			cards c
		LEFT JOIN 
		(
			SELECT *,
				ROW_NUMBER() OVER(PARTITION BY card_id ORDER BY last_update DESC) AS rn
			FROM 
				cards_details
		) cd
		ON 
			c.id = cd.card_id AND cd.rn = 1
	) main 
	ORDER BY price_change DESC LIMIT 20;`

	rows, err := r.db.QueryContext(ctx, getCardsQuery)
	if err != nil {
		return nil, fmt.Errorf("repository failed to exec query in get cards report: %w", err)
	}

	var cardsDomain []domain.Cards

	for rows.Next() {
		var cardDomain domain.Cards
		err := rows.Scan(&cardDomain.ID, &cardDomain.Name, &cardDomain.SetName, &cardDomain.CollectorNumber,
			&cardDomain.Foil, &cardDomain.LastPrice, &cardDomain.OldPrice, &cardDomain.PriceChange, &cardDomain.LastUpdate)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan row in get cards report: %w", err)
		}
		cardsDomain = append(cardsDomain, cardDomain)
	}

	if len(cardsDomain) == 0 {
		return nil, domain.ErrCardNotFound{}
	}

	return cardsDomain, nil
}

func (r *repository) GetTotalPrice(ctx context.Context) (domain.CardsPrice, error) {
	getPriceQuery := `
	SELECT 
		old_price,
		new_price,
		price_change,
		last_update
	FROM 
		prices 
	ORDER by last_update DESC LIMIT 1;`

	row := r.db.QueryRowContext(ctx, getPriceQuery)

	var cardsPriceDomain domain.CardsPrice
	err := row.Scan(&cardsPriceDomain.OldPrice, &cardsPriceDomain.NewPrice, &cardsPriceDomain.PriceChange, &cardsPriceDomain.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.CardsPrice{}, domain.ErrCardNotFound{}
		} else {
			return domain.CardsPrice{}, fmt.Errorf("repository failed to scan row in get total price: %w", err)
		}
	}

	return cardsPriceDomain, nil
}

func getRowsAffected(row sql.Result) error {
	rows, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}
