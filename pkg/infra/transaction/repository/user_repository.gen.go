// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package repository

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/scylladb/go-set/strset"
	"google.golang.org/grpc/codes"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
	repository "github.com/xhayamix/proto-gen-spanner/pkg/domain/repository/transaction"
	cspanner "github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/infra/transaction/repository/base"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

/*
func (r *userRepository) extractQueryCache(ctx context.Context) (base.UserSearchResultCache, base.UserMutationWaitBuffer) {
	return base.ExtractUserSearchResultCache(ctx), base.ExtractUserMutationWaitBuffer(ctx)
}
*/

func (r *userRepository) LoadByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (*transaction.User, error) {
	row, err := r.SelectByPK(ctx, tx, pk)
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	if row == nil {
		return nil, cerrors.Newf(cerrors.InvalidArgument, "Userが見つかりません。 pk = %s", pk)
	}

	return row, nil
}

func (r *userRepository) LoadByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (transaction.UserSlice, error) {
	rows, err := r.SelectByPKs(ctx, tx, pks)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	set := strset.NewWithSize(len(rows))
	for _, row := range rows {
		set.Add(row.GetPK().Key())
	}

	notFoundPKs := make(transaction.UserPKs, 0, len(pks))
	for _, pk := range pks {
		if !set.Has(pk.Key()) {
			notFoundPKs = append(notFoundPKs, pk)
		}
	}
	if len(notFoundPKs) > 0 {
		return nil, cerrors.Newf(cerrors.InvalidArgument, "Userが見つかりません。 pks = %s", notFoundPKs)
	}

	return rows, nil
}

func (r *userRepository) SelectByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (entity *transaction.User, err error) {
	rtx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	row, err := rtx.ReadRow(ctx, transaction.UserTableName, spanner.Key(pk.Generate()), transaction.UserColumnNameSlice)
	if err != nil {
		// PKでのSelect結果がなかった場合
		if spanner.ErrCode(err) == codes.NotFound {
			return nil, nil
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	entity, err = r.decodeAllColumns(row)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	return entity, nil
}

func (r *userRepository) SelectByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (rows transaction.UserSlice, err error) {
	if len(pks) == 0 {
		return transaction.UserSlice{}, nil
	}

	rtx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	keySets := make([]spanner.KeySet, 0, len(pks))
	for _, pk := range pks {
		keySets = append(keySets, spanner.Key(pk.Generate()))
	}
	ri := rtx.Read(ctx, transaction.UserTableName, spanner.KeySets(keySets...), transaction.UserColumnNameSlice)
	rows = make(transaction.UserSlice, 0)
	keySet := strset.New()
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.UserSlice, 0, ri.RowCount)
			keySet = strset.NewWithSize(int(ri.RowCount))
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		keySet.Add(entity.GetPK().Key())
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}

func (r *userRepository) SelectAll(ctx context.Context, tx database.ROTx, limit, offset int32) (rows transaction.UserSlice, err error) {
	roTx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	sql, params := base.NewUserQueryBuilder().
		SelectAllFromUser().
		OrderBy(base.OrderPairs{{"UserID", base.OrderTypeASC}}).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL:    sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.UserSlice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.UserSlice, 0, ri.RowCount)
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}

func (r *userRepository) Search(ctx context.Context, tx database.ROTx, search string, limit, offset int32) (rows transaction.UserSlice, err error) {
	roTx, err := cspanner.ExtractROTx(tx)
	if err != nil {
		return nil, cerrors.Stack(err)
	}

	searchParam := search

	sql, params := base.NewUserQueryBuilder().
		SelectAllFromUser().
		Where().UserIDEq(searchParam).
		OrderBy(base.OrderPairs{{"UserID", base.OrderTypeASC}}).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL:    sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.UserSlice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.UserSlice, 0, ri.RowCount)
		}
		entity, err := r.decodeAllColumns(row)
		if err != nil {
			return cerrors.Stack(err)
		}
		rows = append(rows, entity)
		return nil
	}); err != nil {
		if err, ok := cerrors.As(err); ok {
			return nil, cerrors.Stack(err)
		}
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	return rows, nil
}

func (r *userRepository) Insert(ctx context.Context, tx database.RWTx, entity *transaction.User) (err error) {
	now := time.Now()
	entity.CreatedTime = now
	entity.UpdatedTime = now

	rwTx, err := cspanner.ExtractRWTx(tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	mutation := spanner.Insert(transaction.UserTableName,
		[]string{
			"UserID",
			"PublicUserID",
			"CreatedTime",
			"UpdatedTime",
		},
		[]interface{}{
			entity.UserID,
			entity.PublicUserID,
			entity.CreatedTime,
			entity.UpdatedTime,
		},
	)
	if err := rwTx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
		return err
	}

	return nil
}

// 気が向いたら作る(必要になったら)
func (r *userRepository) BulkInsert(ctx context.Context, tx database.RWTx, entities transaction.UserSlice) (err error) {
	if len(entities) == 0 {
		return nil
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, tx database.RWTx, entity *transaction.User) (err error) {
	now := time.Now()
	entity.UpdatedTime = now

	rwTx, err := cspanner.ExtractRWTx(tx)
	if err != nil {
		return cerrors.Stack(err)
	}
	mutation := spanner.Update(transaction.UserTableName,
		[]string{
			"UserID",
			"PublicUserID",
			"CreatedTime",
			"UpdatedTime",
		},
		[]interface{}{
			entity.UserID,
			entity.PublicUserID,
			entity.CreatedTime,
			entity.UpdatedTime,
		},
	)
	if err := rwTx.BufferWrite([]*spanner.Mutation{mutation}); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Save(ctx context.Context, tx database.RWTx, entity *transaction.User) error {
	var err error
	if entity.CreatedTime.IsZero() {
		err = r.Insert(ctx, tx, entity)
	} else {
		err = r.Update(ctx, tx, entity)
	}
	if err != nil {
		return cerrors.Stack(err)
	}

	return nil
}

// 気が向いたら作る(必要になったら)
func (r *userRepository) Delete(ctx context.Context, tx database.RWTx, pk *transaction.UserPK) (err error) {

	return nil
}

// 気が向いたら作る(必要になったら)
func (r *userRepository) BulkDelete(ctx context.Context, tx database.RWTx, pks transaction.UserPKs) (err error) {

	return nil
}

func (r *userRepository) decodeAllColumns(row *spanner.Row) (*transaction.User, error) {
	var userID spanner.NullString
	var publicUserID spanner.NullString
	var createdTime spanner.NullTime
	var updatedTime spanner.NullTime

	if err := row.Columns(
		&userID,
		&publicUserID,
		&createdTime,
		&updatedTime,
	); err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	var result transaction.User
	if userID.Valid {
		result.UserID = userID.StringVal
	}
	if publicUserID.Valid {
		result.PublicUserID = publicUserID.StringVal
	}
	if createdTime.Valid {
		result.CreatedTime = createdTime.Time.In(time.Local)
	}
	if updatedTime.Valid {
		result.UpdatedTime = updatedTime.Time.In(time.Local)
	}
	return &result, nil
}

func (r *userRepository) diffEntity(source, target *transaction.User) map[string]any {
	result := make(map[string]any)

	// PKの差分は取らない
	if source.PublicUserID != target.PublicUserID {
		result["PublicUserID"] = target.PublicUserID
	}
	if !source.CreatedTime.Equal(target.CreatedTime) {
		result["CreatedTime"] = target.CreatedTime
	}
	if !source.UpdatedTime.Equal(target.UpdatedTime) {
		result["UpdatedTime"] = target.UpdatedTime
	}

	return result
}
