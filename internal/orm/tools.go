package orm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cmelgarejo/go-gql-server/internal/gql/models"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/jinzhu/gorm"
)

// ParseFilters parses the filter and adds the where condition to the transaction
func ParseFilters(db *gorm.DB, filters []*models.QueryFilter) (*gorm.DB, error) {
	for _, f := range filters {
		condition := utils.ToSnakeCase(f.Field) + " " + opToSQL(f.Op)
		switch f.Op {
		case models.OperationTypeBetween:
			if len(f.Values) != 2 {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs an array with exactly two items in [values] field")
			}
			if f.LinkOperation != nil && *f.LinkOperation == models.LinkOperationTypeOr {
				db = db.Or(condition, f.Values[0], f.Values[1])
			} else {
				db = db.Where(condition, f.Values[0], f.Values[1])
			}
		case models.OperationTypeIn, models.OperationTypeNotIn:
			if len(f.Values) < 1 {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs an array with at least 1 item on [values] field")
			}
			if f.LinkOperation != nil && *f.LinkOperation == models.LinkOperationTypeOr {
				db = db.Or(condition, f.Values)
			} else {
				db = db.Where(condition, f.Values)
			}
		case models.OperationTypeMatch:
			if f.LinkOperation != nil && *f.LinkOperation == models.LinkOperationTypeOr {
				db = db.Or("MATCH("+utils.ToSnakeCase(f.Field)+
					") AGAINST (? IN BOOLEAN MODE)", f.Value)
			} else {
				db = db.Where("MATCH("+utils.ToSnakeCase(f.Field)+
					") AGAINST (? IN BOOLEAN MODE)", f.Value)
			}
		case models.OperationTypeIsNotNull:
			fallthrough
		case models.OperationTypeIsNull:
			if f.LinkOperation != nil && *f.LinkOperation == models.LinkOperationTypeOr {
				db = db.Or(condition)
			} else {
				db = db.Where(condition)
			}

		default:
			if f.Value == nil {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs the field [value] to compare")
			}
			if f.LinkOperation != nil && *f.LinkOperation == models.LinkOperationTypeOr {
				db = db.Or(condition, f.Value)
			} else {
				db = db.Where(condition, f.Value)
			}
		}
	}
	return db, db.Error
}

func opToSQL(op models.OperationType) string {
	return map[models.OperationType]string{
		models.OperationTypeEquals:           " = ?",
		models.OperationTypeNotEquals:        " != ?",
		models.OperationTypeLessThan:         " < ?",
		models.OperationTypeLessThanEqual:    " <= ?",
		models.OperationTypeGreaterThan:      " > ?",
		models.OperationTypeGreaterThanEqual: " >= ?",
		models.OperationTypeIs:               " IS ?",
		models.OperationTypeIsNull:           " IS NULL",
		models.OperationTypeIsNotNull:        " IS NOT NULL",
		models.OperationTypeIn:               " IN (?)",
		models.OperationTypeNotIn:            " NOT IN (?)",
		models.OperationTypeLike:             " LIKE ?",
		models.OperationTypeILike:            " ILIKE ?",
		models.OperationTypeNotLike:          " NOT LIKE ?",
		models.OperationTypeBetween:          " BETWEEN ? AND ?",
	}[op]
}

func arrayToString(a interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), "' '", delim, -1), "[]")
}
