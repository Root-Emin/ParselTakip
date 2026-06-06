package iam

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/shared/crypto"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

// UserRepo implements repository.UserRepository with PostgreSQL.
// The encryptor transparently protects PII (national_id / TCKN) at rest.
type UserRepo struct {
	db  *pgxpool.Pool
	enc *crypto.Encryptor
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *pgxpool.Pool, enc *crypto.Encryptor) *UserRepo {
	return &UserRepo{db: db, enc: enc}
}

// userColumns is the canonical SELECT projection; string PII columns are
// COALESCE'd so NULLs scan cleanly into Go strings.
const userColumns = `id, email, password_hash, first_name, last_name, status,
	COALESCE(user_type, 'citizen'), COALESCE(phone, ''), COALESCE(national_id, ''),
	COALESCE(address, ''), COALESCE(city, ''), COALESCE(district, ''),
	company_id, kvkk_consent_at, COALESCE(kvkk_consent_version, ''),
	last_login_at, failed_login_count, locked_until, created_at, updated_at`

func (r *UserRepo) scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	var encNationalID string
	if err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Status,
		&u.UserType, &u.Phone, &encNationalID,
		&u.Address, &u.City, &u.District,
		&u.CompanyID, &u.KVKKConsentAt, &u.KVKKConsentVersion,
		&u.LastLoginAt, &u.FailedLoginCount, &u.LockedUntil, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if encNationalID != "" && r.enc != nil {
		if dec, err := r.enc.Decrypt(encNationalID); err == nil {
			u.NationalID = dec
		}
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.UserType == "" {
		user.UserType = model.UserTypeCitizen
	}

	encNationalID, hash := "", ""
	if user.NationalID != "" && r.enc != nil {
		var err error
		if encNationalID, err = r.enc.Encrypt(user.NationalID); err != nil {
			return domainErr.New(domainErr.ErrInternal, "failed to encrypt national id", err)
		}
		hash = r.enc.BlindIndex(user.NationalID)
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, first_name, last_name, status, user_type,
			phone, national_id, national_id_hash, address, city, district, company_id,
			kvkk_consent_at, kvkk_consent_version, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Status, user.UserType,
		nullable(user.Phone), nullable(encNationalID), nullable(hash), nullable(user.Address),
		nullable(user.City), nullable(user.District), user.CompanyID,
		user.KVKKConsentAt, nullable(user.KVKKConsentVersion), user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create user", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u, err := r.scanUser(r.db.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "user not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get user", err)
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u, err := r.scanUser(r.db.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "user not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get user by email", err)
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now().UTC()
	if user.UserType == "" {
		user.UserType = model.UserTypeCitizen
	}

	encNationalID, hash := "", ""
	if user.NationalID != "" && r.enc != nil {
		var err error
		if encNationalID, err = r.enc.Encrypt(user.NationalID); err != nil {
			return domainErr.New(domainErr.ErrInternal, "failed to encrypt national id", err)
		}
		hash = r.enc.BlindIndex(user.NationalID)
	}

	_, err := r.db.Exec(ctx,
		`UPDATE users SET email=$1, first_name=$2, last_name=$3, status=$4, user_type=$5,
			phone=$6, national_id=$7, national_id_hash=$8, address=$9, city=$10, district=$11,
			company_id=$12, kvkk_consent_at=$13, kvkk_consent_version=$14, updated_at=$15
		 WHERE id=$16`,
		user.Email, user.FirstName, user.LastName, user.Status, user.UserType,
		nullable(user.Phone), nullable(encNationalID), nullable(hash), nullable(user.Address),
		nullable(user.City), nullable(user.District), user.CompanyID,
		user.KVKKConsentAt, nullable(user.KVKKConsentVersion), user.UpdatedAt, user.ID,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update user", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete user", err)
	}
	return nil
}

func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*model.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count users", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+userColumns+` FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list users", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan user", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

// RecordLoginSuccess clears the failed-login counter and stamps the last login.
func (r *UserRepo) RecordLoginSuccess(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET failed_login_count = 0, locked_until = NULL, last_login_at = $2, updated_at = $2 WHERE id = $1`,
		id, time.Now().UTC(),
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to record login success", err)
	}
	return nil
}

// RecordLoginFailure increments the failed-login counter and locks the account
// for lockDuration once it reaches maxAttempts.
func (r *UserRepo) RecordLoginFailure(ctx context.Context, id uuid.UUID, maxAttempts int, lockDuration time.Duration) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users
		 SET failed_login_count = failed_login_count + 1,
		     locked_until = CASE WHEN failed_login_count + 1 >= $2 THEN $3 ELSE locked_until END,
		     updated_at = NOW()
		 WHERE id = $1`,
		id, maxAttempts, time.Now().UTC().Add(lockDuration),
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to record login failure", err)
	}
	return nil
}

// nullable converts an empty string to a SQL NULL to keep optional columns clean.
func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
