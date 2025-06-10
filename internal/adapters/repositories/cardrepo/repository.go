package cardrepo

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/internal/adapters/entities"
	"mtg-report/internal/adapters/factories"
	"mtg-report/internal/core/domain"
	database "mtg-report/internal/sources/databases/mysql"
	"mtg-report/internal/sources/logger/logrus"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type repository struct {
	db  database.Client
	log logrus.Logger
}

func New(db database.Client, log logrus.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (r *repository) InsertCard(ctx context.Context, card domain.Cards) (domain.Cards, error) {
	insertCardQuery := `
	INSERT INTO cards 
		(name, set_name, collector_number, foil) 
	VALUES 
		(?, ?, ?, ?);`

	res, err := r.db.ExecContext(ctx, insertCardQuery, card.Name, card.SetName, card.CollectorNumber, card.Foil)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				return domain.Cards{}, domain.ErrCardAlreadyExists{}
			}
		}
		return domain.Cards{}, fmt.Errorf("repository failed to exec insert query in insert card: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to get last inserted id in insert card: %w", err)
	}

	card.ID = id

	return card, nil
}

func (r *repository) DeleteCard(ctx context.Context, id string) error {
	DeleteCardQuery := `
	DELETE FROM 
		cards 
	WHERE
		id = ?`

	_, err := r.db.ExecContext(ctx, DeleteCardQuery, id)
	if err != nil {
		return fmt.Errorf("repository failed to exec delete query in delete card: %w", err)
	}

	return nil
}

func (r *repository) GetCardbyID(ctx context.Context, id string) (domain.Cards, error) {
	getCardQuery := `
	SELECT 
		c.id,
		name,
		set_name,
		collector_number,
		foil,
		COALESCE(cd.last_price, 0) as last_price,
		COALESCE(cd.old_price, 0) as old_price,
		COALESCE(cd.price_change, 0) as price_change,
		cd.last_update
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
	WHERE 
		c.id = ?;`

	row := r.db.QueryRowContext(ctx, getCardQuery, id)

	var cardDomain domain.Cards
	err := row.Scan(&cardDomain.ID, &cardDomain.Name, &cardDomain.SetName, &cardDomain.CollectorNumber,
		&cardDomain.Foil, &cardDomain.LastPrice, &cardDomain.OldPrice, &cardDomain.PriceChange, &cardDomain.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Cards{}, domain.ErrCardNotFound{}
		} else {
			return domain.Cards{}, fmt.Errorf("repository failed to scan row in get card by id: %w", err)
		}
	}

	return cardDomain, nil
}

func (r *repository) GetCards(ctx context.Context, filters map[string]string) ([]domain.Cards, error) {
	getCardsQuery := `
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
    `

	var first bool = true
	var values []interface{}

	for key, value := range filters {
		if !first {
			getCardsQuery += " AND "
		} else {
			getCardsQuery += " WHERE "
		}
		getCardsQuery += fmt.Sprintf("%s = ?", key)
		values = append(values, value)
		first = false
	}

	getCardsQuery += " ORDER BY last_price DESC"

	rows, err := r.db.QueryContext(ctx, getCardsQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("repository failed to exec query in get cards: %w", err)
	}

	var cardsDomain []domain.Cards

	for rows.Next() {
		var cardDomain domain.Cards
		err := rows.Scan(&cardDomain.ID, &cardDomain.Name, &cardDomain.SetName, &cardDomain.CollectorNumber,
			&cardDomain.Foil, &cardDomain.LastPrice, &cardDomain.OldPrice, &cardDomain.PriceChange, &cardDomain.LastUpdate)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan row in get cards: %w", err)
		}
		cardsDomain = append(cardsDomain, cardDomain)
	}

	if len(cardsDomain) == 0 {
		return nil, domain.ErrCardNotFound{}
	}

	return cardsDomain, nil
}

func (r *repository) InsertCards(ctx context.Context, cards []domain.Cards) error {
	valueStrings := make([]string, 0, len(cards))
	valueArgs := make([]interface{}, 0, len(cards)*6)
	for _, card := range cards {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, card.Name)
		valueArgs = append(valueArgs, card.SetName)
		valueArgs = append(valueArgs, card.CollectorNumber)
		valueArgs = append(valueArgs, card.Foil)
	}

	stmt := fmt.Sprintf(`
	INSERT INTO cards 
		(name, set_name, collector_number, foil) 
	VALUES 
		%s 
	ON DUPLICATE KEY UPDATE 
		name = name;`,
		strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("repository failed to exec insert query in insert cards: %w", err)
	}

	return nil
}

func (r *repository) GetCardHistory(ctx context.Context, id string) ([]domain.Cards, error) {
	cards := []entities.MysqlCardPriceHistory{}

	getQuery := `
	SELECT 
		c.id,
		c.name,
		c.set_name,
		c.collector_number,
		c.foil,
		COALESCE(cd.last_price, 0),
		COALESCE(cd.old_price, 0),
		COALESCE(cd.price_change, 0),
		cd.last_update
	FROM 
		cards c 
	LEFT JOIN 
		cards_details cd 
	ON 
		c.id = cd.card_id
	WHERE 
		c.id = ?
	ORDER BY 
		last_update DESC;
	`
	rows, err := r.db.QueryContext(ctx, getQuery, id)
	if err != nil {
		return nil, fmt.Errorf("repository failed to query in get cards history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var card entities.MysqlCardPriceHistory
		err = rows.Scan(&card.ID, &card.Name, &card.SetName, &card.CollectorNumber, &card.Foil, &card.LastPrice, &card.OldPrice, &card.PriceChange, &card.LastUpdate)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan rows in get cards history: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository failed after iterating rows in get cards for update: %w", err)
	}

	if len(cards) == 0 {
		return nil, domain.ErrCardNotFound{}
	}

	return factories.CardPriceHistoryToCardsDomain(cards), nil
}

func (r *repository) UpdateCard(ctx context.Context, card domain.UpdateCard) (domain.Cards, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to begin transaction in update card: %w", err)
	}
	defer tx.Rollback()

	checkCardQuery := `
	SELECT 
		COUNT(*) 
	FROM 
		cards 
	WHERE 
		id = ?;`

	var count int
	err = tx.QueryRowContext(ctx, checkCardQuery, card.ID).Scan(&count)
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to query card count in update card: %w", err)
	}

	if count == 0 {
		return domain.Cards{}, domain.ErrCardNotFound{}
	}

	updateCardQuery := `
	UPDATE cards 
	SET 
		name = ? 
	WHERE
		id = ?;`

	result, err := tx.ExecContext(ctx, updateCardQuery, card.Name, card.ID)
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to exec update query in update card: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to get rows affected in update card: %w", err)
	}

	getCardQuery := `
	SELECT 
		c.id,
		name,
		set_name,
		collector_number,
		foil
	FROM 
		cards c
	WHERE 
		c.id = ?;`

	row := tx.QueryRowContext(ctx, getCardQuery, card.ID)
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to exec update query in update card: %w", err)
	}

	var cardDomain domain.Cards
	err = row.Scan(&cardDomain.ID, &cardDomain.Name, &cardDomain.SetName, &cardDomain.CollectorNumber, &cardDomain.Foil)
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to scan row in update card: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Cards{}, fmt.Errorf("repository failed to commit transaction in update card: %w", err)
	}

	return cardDomain, nil
}

func (r *repository) GetCardsPaginated(ctx context.Context, filters map[string]string, offset, limit int) ([]domain.Cards, error) {
	getCardsQuery := `
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
    `

	var first bool = true
	var values []interface{}

	for key, value := range filters {
		if !first {
			getCardsQuery += " AND "
		} else {
			getCardsQuery += " WHERE "
		}
		getCardsQuery += fmt.Sprintf("%s = ?", key)
		values = append(values, value)
		first = false
	}

	getCardsQuery += " ORDER BY last_price DESC LIMIT ? OFFSET ?"
	values = append(values, limit, offset)

	rows, err := r.db.QueryContext(ctx, getCardsQuery, values...)
	if err != nil {
		return nil, fmt.Errorf("repository failed to exec query in get cards paginated: %w", err)
	}
	defer rows.Close()

	var cardsDomain []domain.Cards

	for rows.Next() {
		var cardDomain domain.Cards
		err := rows.Scan(&cardDomain.ID, &cardDomain.Name, &cardDomain.SetName, &cardDomain.CollectorNumber,
			&cardDomain.Foil, &cardDomain.LastPrice, &cardDomain.OldPrice, &cardDomain.PriceChange, &cardDomain.LastUpdate)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan row in get cards paginated: %w", err)
		}
		cardsDomain = append(cardsDomain, cardDomain)
	}

	return cardsDomain, nil
}

func (r *repository) GetCardsCount(ctx context.Context, filters map[string]string) (int64, error) {
	countQuery := `
    SELECT COUNT(*)
    FROM cards c
    `

	var first bool = true
	var values []interface{}

	for key, value := range filters {
		if !first {
			countQuery += " AND "
		} else {
			countQuery += " WHERE "
		}
		countQuery += fmt.Sprintf("%s = ?", key)
		values = append(values, value)
		first = false
	}

	row := r.db.QueryRowContext(ctx, countQuery, values...)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository failed to scan count in get cards count: %w", err)
	}

	return count, nil
}

func (r *repository) GetCardHistoryPaginated(ctx context.Context, id string, offset, limit int) ([]domain.Cards, error) {
	cards := []entities.MysqlCardPriceHistory{}

	getQuery := `
	SELECT 
		c.id,
		c.name,
		c.set_name,
		c.collector_number,
		c.foil,
		COALESCE(cd.last_price, 0),
		COALESCE(cd.old_price, 0),
		COALESCE(cd.price_change, 0),
		cd.last_update
	FROM 
		cards c 
	LEFT JOIN 
		cards_details cd 
	ON 
		c.id = cd.card_id
	WHERE 
		c.id = ?
	ORDER BY 
		last_update DESC
	LIMIT ? OFFSET ?;
	`
	rows, err := r.db.QueryContext(ctx, getQuery, id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository failed to query in get cards history paginated: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var card entities.MysqlCardPriceHistory
		err = rows.Scan(&card.ID, &card.Name, &card.SetName, &card.CollectorNumber, &card.Foil, &card.LastPrice, &card.OldPrice, &card.PriceChange, &card.LastUpdate)
		if err != nil {
			return nil, fmt.Errorf("repository failed to scan rows in get cards history paginated: %w", err)
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository failed after iterating rows in get cards history paginated: %w", err)
	}

	return factories.CardPriceHistoryToCardsDomain(cards), nil
}

func (r *repository) GetCardHistoryCount(ctx context.Context, id string) (int64, error) {
	countQuery := `
	SELECT COUNT(*)
	FROM 
		cards c 
	LEFT JOIN 
		cards_details cd 
	ON 
		c.id = cd.card_id
	WHERE 
		c.id = ?;
	`

	row := r.db.QueryRowContext(ctx, countQuery, id)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository failed to scan count in get card history count: %w", err)
	}

	return count, nil
}

func (r *repository) GetCollectionStats(ctx context.Context) (domain.CollectionStats, error) {
	statsQuery := `
	SELECT 
		COUNT(*) as total_cards,
		SUM(CASE WHEN foil = true THEN 1 ELSE 0 END) as foil_cards,
		COUNT(DISTINCT set_name) as unique_sets,
		COALESCE(SUM(cd.last_price), 0) as total_value
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
		c.id = cd.card_id AND cd.rn = 1;
	`

	row := r.db.QueryRowContext(ctx, statsQuery)

	var stats domain.CollectionStats
	err := row.Scan(&stats.TotalCards, &stats.FoilCards, &stats.UniqueSets, &stats.TotalValue)
	if err != nil {
		return domain.CollectionStats{}, fmt.Errorf("repository failed to scan collection stats: %w", err)
	}

	return stats, nil
}
