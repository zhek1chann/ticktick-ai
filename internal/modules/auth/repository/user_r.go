package repository

import (
	"context"
	"ticktick-ai/internal/domain"
	"ticktick-ai/internal/modules/auth/repository/row"
	"ticktick-ai/pkg/db"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

var dialect = goqu.Dialect("postgres")

// users — table and columns
var usersTable = goqu.T("users")

const (
	uIDC           = "id"
	uNameC         = "name"
	uPhoneNumberC  = "phone_number"
	uPasswordHashC = "password_hash"
	uRoleC         = "role"
	uCreatedAtC    = "created_at"
	uUpdatedAtC    = "updated_at"
)

func (r *Repo) CreateUser(ctx context.Context, user domain.User) (int, error) {
	query, args, err := dialect.
		Insert(usersTable).
		Cols(uNameC, uPhoneNumberC, uPasswordHashC, uRoleC).
		Vals(goqu.Vals{user.Name, user.PhoneNumber, user.PasswordHash, user.Role}).
		Returning(uIDC).
		Prepared(true).
		ToSQL()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "CreateUser",
		QueryRaw: query,
	}

	var userID int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		return 0, domain.MapError(err)
	}

	return userID, nil
}

func (r *Repo) UserByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error) {
	query, args, err := dialect.
		Select(uIDC, uNameC, uPhoneNumberC, uPasswordHashC, uRoleC, uCreatedAtC, uUpdatedAtC).
		From(usersTable).
		Where(goqu.C(uPhoneNumberC).Eq(phoneNumber)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return domain.User{}, err
	}

	q := db.Query{
		Name:     "UserByPhoneNumber",
		QueryRaw: query,
	}

	var userRow row.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&userRow.ID,
		&userRow.Name,
		&userRow.PhoneNumber,
		&userRow.PasswordHash,
		&userRow.Role,
		&userRow.CreatedAt,
		&userRow.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, domain.MapError(err)
	}

	return row.ToUser(userRow), nil
}

func (r *Repo) UserByID(ctx context.Context, id int) (domain.User, error) {
	query, args, err := dialect.
		Select(uIDC, uNameC, uPhoneNumberC, uPasswordHashC, uRoleC, uCreatedAtC, uUpdatedAtC).
		From(usersTable).
		Where(goqu.C(uIDC).Eq(id)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return domain.User{}, err
	}

	q := db.Query{
		Name:     "UserByID",
		QueryRaw: query,
	}

	var userRow row.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&userRow.ID,
		&userRow.Name,
		&userRow.PhoneNumber,
		&userRow.PasswordHash,
		&userRow.Role,
		&userRow.CreatedAt,
		&userRow.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, domain.MapError(err)
	}

	return row.ToUser(userRow), nil
}

func (r *Repo) UpdateUserRole(ctx context.Context, userID int, role string) error {
	query, args, err := dialect.
		Update(usersTable).
		Set(goqu.Record{
			uRoleC:      role,
			uUpdatedAtC: goqu.L("NOW()"),
		}).
		Where(goqu.C(uIDC).Eq(userID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "UpdateUserRole",
		QueryRaw: query,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *Repo) ListUsers(ctx context.Context) ([]domain.User, error) {
	query, args, err := dialect.
		Select(uIDC, uNameC, uPhoneNumberC, uPasswordHashC, uRoleC, uCreatedAtC, uUpdatedAtC).
		From(usersTable).
		Order(goqu.C(uIDC).Desc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "ListUsers",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var userRow row.User
		err = rows.Scan(
			&userRow.ID,
			&userRow.Name,
			&userRow.PhoneNumber,
			&userRow.PasswordHash,
			&userRow.Role,
			&userRow.CreatedAt,
			&userRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, row.ToUser(userRow))
	}

	return users, nil
}
