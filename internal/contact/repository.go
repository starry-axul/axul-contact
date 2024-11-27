package contact

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ncostamagna/go-logger-hub/loghub"
	"github.com/ncostamagna/axul_domain/domain"
	"gorm.io/gorm"
)

// Repository is a Repository handler interface
type Repository interface {
	Create(ctx context.Context, contact *domain.Contact) error
	Update(ctx context.Context, id string, firstName, lastName, nickName, gender, phone *string, birthday *time.Time) error
	GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error)
	Get(ctx context.Context, id string) (*domain.Contact, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters Filter) (int, error)
}

type repo struct {
	db  *gorm.DB
	log loghub.Logger
}

// NewRepo is a repositories handler
func NewRepo(db *gorm.DB, logger loghub.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (repo *repo) Create(_ context.Context, contact *domain.Contact) error {

	if err := repo.db.Create(&contact).Error; err != nil {
		repo.log.Error(err)
		return err
	}

	return nil
}

func (repo *repo) GetAll(ctx context.Context, f Filter, offset, limit int) ([]domain.Contact, error) {

	var tx *gorm.DB
	var cs []domain.Contact

	tx = repo.db.WithContext(ctx).Model(&cs)
	tx = applyFilters(tx, f)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Find(&cs)

	for i := range cs {
		year := f.firstDate.Year()
		if cs[i].Birthday.Month() < f.firstDate.Month() {
			year++
		} else if cs[i].Birthday.Month() == f.firstDate.Month() {
			if cs[i].Birthday.Day() < f.firstDate.Day() {
				year++
			}
		}

		bd := time.Date(year, cs[i].Birthday.Month(), cs[i].Birthday.Day(), 0, 0, 0, 0, time.UTC)
		cs[i].Days = int64(bd.Sub(f.firstDate).Hours() / 24)
	}

	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].Days < cs[j].Days
	})

	if result.Error != nil {
		repo.log.Error(result.Error)
		return nil, result.Error
	}

	repo.log.Info(fmt.Sprintf("Row: %d", result.RowsAffected))

	return cs, nil
}

func (repo *repo) Get(_ context.Context, id string) (*domain.Contact, error) {
	contact := domain.Contact{}

	if err := repo.db.Where("id = ?", id).First(&contact).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return &contact, nil
}

func (repo *repo) Update(ctx context.Context, id string, firstName, lastName, nickName, gender, phone *string, birthday *time.Time) error {

	values := make(map[string]interface{})

	if firstName != nil {
		values["firstname"] = *firstName
	}

	if lastName != nil {
		values["lastname"] = *lastName
	}

	if nickName != nil {
		values["nickname"] = *nickName
	}

	if gender != nil {
		values["gender"] = *gender
	}

	if phone != nil {
		values["phone"] = *phone
	}

	if birthday != nil {
		values["birthday"] = *birthday
	}

	result := repo.db.WithContext(ctx).Model(&domain.Contact{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Error(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound{fmt.Sprint(id)}
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	course := domain.Contact{ID: id}
	result := r.db.WithContext(ctx).Delete(&course)

	if result.Error != nil {
		r.log.Error(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound{id}
	}
	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filter) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Contact{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, f Filter) *gorm.DB {

	if f.RangeDays != nil {
		second := f.firstDate.AddDate(0, 0, int(*f.RangeDays)).Add(time.Hour * 20)
		tx = tx.Where("CONCAT('"+strconv.Itoa(f.firstDate.Year())+"',DATE_FORMAT(birthday,'%m%d')) between DATE_FORMAT(?,'%Y%m%d') and DATE_FORMAT(?,'%Y%m%d')", f.firstDate, second)
	}

	if f.Firstname != "" {
		tx = tx.Where("UPPER(firstname) like CONCAT('%',UPPER(?),'%')", f.Firstname)
	}

	if f.Lastname != "" {
		tx = tx.Where("UPPER(lastname) like CONCAT('%',UPPER(?),'%')", f.Lastname)
	}

	if f.Month != 0 {
		tx = tx.Where("MONTH(birthday) = ?", f.Month)
	}

	if f.Birthday != nil {
		date := time.Now().AddDate(0, 0, *f.Birthday)
		day, month := date.Day(), int(date.Month())
		tx = tx.Where("month(birthday) = ? and day(birthday) = ?", month, day)
	}

	return tx
}
