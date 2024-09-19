{{ template "autogen_comment" }}
package preference

import (
	"context"

	"github.com/xhayamix/proto-gen-spanner/pkg/cerrors"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/database"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/entity/transaction"
	"github.com/xhayamix/proto-gen-spanner/pkg/domain/enum"
	transactionrepository "github.com/xhayamix/proto-gen-spanner/pkg/domain/repository/transaction"
)

type Interactor interface {
	Update(ctx context.Context, userID string, types enum.PreferenceTypeSlice, pref *Preference) error
}

type interactor struct {
	txManager                database.UserTxManager
	userPreferenceRepository transactionrepository.UserPreferenceRepository
}

func New(
	txManager database.UserTxManager,
	userPreferenceRepository transactionrepository.UserPreferenceRepository,
) Interactor {
	return &interactor{
		txManager:                txManager,
		userPreferenceRepository: userPreferenceRepository,
	}
}

type Preference struct {
	{{- range .Elements }}
	{{ .PascalName }} {{ .Type }}
	{{- end }}
}

func (i *interactor) Update(ctx context.Context, userID string, types enum.PreferenceTypeSlice, pref *Preference) error {
	if pref == nil {
		return nil
	}

	if err := i.txManager.Transaction(ctx, func(ctx context.Context, tx database.RWTx) error {
		userPreference, err := i.userPreferenceRepository.SelectByPK(ctx, tx, &transaction.UserPreferencePK{UserID: userID})
		if err != nil {
			return cerrors.Stack(err)
		}
		if userPreference == nil {
			userPreference = &transaction.UserPreference{
				UserID: userID,
			}
		}

		for _, typ := range types {
			switch typ {
			{{- range .Elements }}
			case enum.PreferenceType_{{ .PascalName }}:
				userPreference.{{ .PascalName }} = pref.{{ .PascalName }}
			{{- end }}
			}
		}

		if err := i.userPreferenceRepository.Save(ctx, tx, userPreference); err != nil {
			return cerrors.Stack(err)
		}

		return nil
	}); err != nil {
		return cerrors.Stack(err)
	}

	return nil
}
