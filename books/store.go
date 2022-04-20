package books

import (
	"context"
	"errors"

	"github.com/circleci/ex/db"
	"github.com/circleci/ex/o11y"
	"github.com/google/uuid"
)

var (
	ErrNotFound = o11y.NewWarning("no update or results")
)

type Store struct {
	txm *db.TxManager
}

func NewStore(txm *db.TxManager) *Store {
	return &Store{
		txm: txm,
	}
}

func mapError(err error, to error) error {
	if errors.Is(err, db.ErrNop) {
		return to
	}
	return err
}

type Book struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Price string    `db:"price"`
}

func (s *Store) ByID(ctx context.Context, id uuid.UUID) (book *Book, err error) {
	ctx, span := o11y.StartSpan(ctx, "store: by_id")
	defer o11y.End(span, &err)
	span.AddField("id", id)

	err = s.txm.WithTx(ctx, func(ctx context.Context, q db.Querier) (err error) {
		book, err = queryGetBookByID(ctx, q, id)
		return err
	})

	return book, mapError(err, ErrNotFound)
}

func (s *Store) All(ctx context.Context) (books []Book, err error) {
	ctx, span := o11y.StartSpan(ctx, "store: all")
	defer o11y.End(span, &err)

	err = s.txm.WithTx(ctx, func(ctx context.Context, q db.Querier) (err error) {
		books, err = queryGetAllBooks(ctx, q)
		return err
	})

	return books, mapError(err, ErrNotFound)
}

type ToAdd struct {
	Name  string `db:"name"`
	Price string `db:"price"`
}

func (s *Store) Add(ctx context.Context, toAdd ToAdd) (id uuid.UUID, err error) {
	ctx, span := o11y.StartSpan(ctx, "store: add")
	defer o11y.End(span, &err)
	span.AddField("name", toAdd.Name)
	span.RecordMetric(o11y.Incr("book-added"))

	err = s.txm.WithTx(ctx, func(ctx context.Context, q db.Querier) (err error) {
		id, err = queryInsertBook(ctx, q, toAdd)
		return err
	})

	return id, mapError(err, ErrNotFound)
}

func (s *Store) AddMany(ctx context.Context, manyToAdd []ToAdd) (ids []uuid.UUID, err error) {
	ctx, span := o11y.StartSpan(ctx, "store: addMany")
	defer o11y.End(span, &err)

	var id uuid.UUID
	ids = []uuid.UUID{}
	for _, toAdd := range manyToAdd {
		id, err = s.Add(ctx, toAdd)
		if err != nil {
			break
		}

		ids = append(ids, id)
	}

	return ids, mapError(err, ErrNotFound)
}
