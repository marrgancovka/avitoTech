package repo

import (
	"avitoTech/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetUserBanner(ctx context.Context, tagId, featureId int64) ([]byte, error) {
	var content []byte
	selectContent := `SELECT bb.content
			FROM banner_feature_tag bft
			JOIN banner bb ON bft.id_banner = bb.id
			WHERE bft.id_tag = $1 AND bft.id_feature = $2 AND bb.is_active = true;`
	if err := r.db.QueryRowContext(ctx, selectContent, tagId, featureId).
		Scan(&content); err != nil {
		return nil, err
	}
	return content, nil
}
func (r *PostgresRepo) GetBanners(ctx context.Context, data *models.GetAllBannersRequest) ([]*models.GetAllBannersResponse, error) {
	conditionCount := ""
	conditionTag := ""
	conditionFeature := ""
	response := []*models.GetAllBannersResponse{}

	if data.Limit != 0 {
		conditionCount = conditionCount + " LIMIT " + strconv.FormatInt(data.Limit, 10)
	}
	if data.Offset != 0 {
		conditionCount = conditionCount + " OFFSET " + strconv.FormatInt(data.Offset, 10)
	}
	conditionCount += ";"

	if data.TagId != 0 {
		conditionTag = " HAVING ARRAY_AGG(bft.id_tag) @> ARRAY[" + strconv.FormatInt(data.TagId, 10) + "]"
	}
	if data.FeatureId != 0 {
		conditionFeature = " WHERE bft.id_feature = " + strconv.FormatInt(data.FeatureId, 10)
	}
	selectBanners := `
		select b.id, b.id_feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(bft.id_tag)
		from banner b
         	join banner_feature_tag bft on b.id = bft.id_banner` + conditionFeature + ` group by b.id, b.id_feature` + conditionTag + conditionCount

	rows, err := r.db.QueryContext(ctx, selectBanners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bannerID    int64
			featureID   int64
			contentByte []byte
			createdAt   time.Time
			updatedAt   time.Time
			isActive    bool
			tagIDs      []int64
		)
		if err := rows.Scan(&bannerID, &featureID, &contentByte, &createdAt, &updatedAt, &isActive, pq.Array(&tagIDs)); err != nil {
			return nil, err
		}

		banner := &models.GetAllBannersResponse{
			BannerId:  bannerID,
			FeatureId: featureID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			IsActive:  isActive,
			TagIds:    tagIDs,
		}

		if err := json.Unmarshal(contentByte, &banner.Content); err != nil {
			return nil, err
		}
		response = append(response, banner)
	}
	return response, nil
}

func (r *PostgresRepo) CreateBanner(ctx context.Context, data *models.CreateBannerRequest) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var featureExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM feature WHERE id = $1)`, data.FeatureId).Scan(&featureExists)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	if !featureExists {
		_, err = tx.Exec(`INSERT INTO feature (id) VALUES ($1)`, data.FeatureId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	var newBannerID int64
	content, err := json.Marshal(data.Content)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow("INSERT INTO banner (content, id_feature) VALUES ($1, $2) RETURNING id", content, data.FeatureId).Scan(&newBannerID)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, tagID := range data.TagIds {
		var tagExists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM tag WHERE id = $1)", tagID).Scan(&tagExists)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		if !tagExists {
			_, err = tx.Exec("INSERT INTO tag (id) VALUES ($1)", tagID)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
		}

		_, err = tx.Exec("INSERT INTO banner_feature_tag (id_banner, id_tag, id_feature) VALUES ($1, $2, $3)", newBannerID, tagID, data.FeatureId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	return newBannerID, nil

}
func (r *PostgresRepo) UpdateBanner(ctx context.Context, id int64, data *models.UpdateBannerRequest) error {
	var nullModel models.Content
	var content []byte
	content, err := json.Marshal(data.Content)
	if err != nil {
		panic(err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}

	if data.FeatureId != 0 {
		var exists bool
		err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM feature WHERE id = $1)", data.FeatureId).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			_, err = tx.Exec(`INSERT INTO feature (id) VALUES ($1)`, data.FeatureId)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
		}
		_, err = tx.Exec(`
	   UPDATE banner
	   SET
	       id_feature = $1,
	       updated_at = CURRENT_TIMESTAMP
	   WHERE id = $2
	`, data.FeatureId, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if data.IsActive != nil {
		_, err = tx.Exec(`
	   UPDATE banner
	   SET
	       is_active = $1,
	       updated_at = CURRENT_TIMESTAMP
	   WHERE id = $2
	`, *data.IsActive, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if data.Content != nullModel {
		_, err = tx.Exec(`
	   UPDATE banner
	   SET
	       content = $1,
	       updated_at = CURRENT_TIMESTAMP
	   WHERE id = $2
	`, content, id)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if data.TagIds != nil {
		_, err = tx.Exec("DELETE FROM banner_feature_tag WHERE id_banner = $1", id)
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare("INSERT INTO banner_feature_tag (id_banner, id_tag, id_feature) VALUES ($1, $2, (SELECT id_feature FROM banner WHERE id = $1))")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, tagID := range data.TagIds {
			var tagExists bool
			err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM tag WHERE id = $1)", tagID).Scan(&tagExists)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			if !tagExists {
				_, err = tx.Exec("INSERT INTO tag (id) VALUES ($1)", tagID)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
			}
			_, err = stmt.Exec(id, tagID)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	return nil
}
func (r *PostgresRepo) DeleteBanner(ctx context.Context, id int64) error {
	deleteQuery := `DELETE FROM banner WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, deleteQuery, id); err != nil {
		return err
	}
	return nil
}
