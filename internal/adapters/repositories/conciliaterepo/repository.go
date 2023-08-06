package conciliaterepo

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/adapters/factories"
	"mtg-report/internal/core/domain"
	database "mtg-report/internal/sources/databases/mysql"
	"strings"
)

type repository struct {
	db database.Client
}

func New(db database.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertCardDetails(ctx context.Context, cardDetails []domain.CardsDetails) error {
	if len(cardDetails) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(cardDetails))
	valueArgs := make([]interface{}, 0, len(cardDetails)*5)

	for _, card := range cardDetails {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, card.CardID, card.LastPrice, card.OldPrice, card.PriceChange, card.LastUpdate)
	}

	insertCardQuery := fmt.Sprintf("INSERT INTO cards_details (card_id, last_price, old_price, price_change, last_update) VALUES %s", strings.Join(valueStrings, ", "))

	res, err := r.db.ExecContext(ctx, insertCardQuery, valueArgs...)
	if err != nil {
		return fmt.Errorf("repository failed to execute insert statement: %w", err)
	}

	err = getRowsAffected(res)
	if err != nil {
		return fmt.Errorf("repository insert card details failed: %w", err)
	}

	return nil
}

func (r *repository) GetCardsForUpdate(ctx context.Context, offset int, limit int) ([]domain.Cards, error) {
	cards := []entities.MysqlCardInfo{}

	getQuery := `
	SELECT 
		c.id,
		c.name,
		c.set_name,
		c.collector_number,
		c.foil,
		cd.last_price
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
	LIMIT ?, ?;
	`
	rows, err := r.db.QueryContext(ctx, getQuery, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("repository failed to query in get cards for update: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var card entities.MysqlCardInfo
		err = rows.Scan(&card.ID, &card.Name, &card.SetName, &card.CollectorNumber, &card.Foil, &card.LastPrice)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan rows in get cards for update: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository failed after iterating rows in get cards for update: %w", err)
	}

	return factories.CardsInfoToCardsDomain(cards), nil
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
