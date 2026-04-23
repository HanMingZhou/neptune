package product

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductPrice struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ProductID uint      `json:"productId" gorm:"column:product_id;index:idx_product_price_unique,unique;comment:产品ID"`
	PriceType int       `json:"priceType" gorm:"column:price_type;index:idx_product_price_unique,unique;comment:价格类型(1-小时 2-天 3-周 4-月)"`
	Price     float64   `json:"price" gorm:"column:price;type:decimal(20,6);comment:价格"`
}

func (ProductPrice) TableName() string {
	return "product_prices"
}

type ComputePriceValues struct {
	Hourly  float64
	Daily   float64
	Weekly  float64
	Monthly float64
}

func ComputePriceValuesFromItems(items []ProductPrice) ComputePriceValues {
	values := ComputePriceValues{}
	for _, item := range items {
		switch item.PriceType {
		case ChargeTypeHourly:
			values.Hourly = item.Price
		case ChargeTypeDaily:
			values.Daily = item.Price
		case ChargeTypeWeekly:
			values.Weekly = item.Price
		case ChargeTypeMonthly:
			values.Monthly = item.Price
		}
	}
	return values
}

func GetPriceByTypeFromItems(items []ProductPrice, priceType int) float64 {
	values := ComputePriceValuesFromItems(items)
	switch NormalizeComputePriceType(priceType) {
	case ChargeTypeHourly:
		return values.Hourly
	case ChargeTypeDaily:
		return values.Daily
	case ChargeTypeWeekly:
		return values.Weekly
	case ChargeTypeMonthly:
		return values.Monthly
	default:
		return values.Hourly
	}
}

func HasSameComputePrices(a, b *Product) bool {
	if a == nil || b == nil {
		return false
	}
	left := ComputePriceValuesFromItems(a.PriceItems)
	right := ComputePriceValuesFromItems(b.PriceItems)
	return left.Hourly == right.Hourly &&
		left.Daily == right.Daily &&
		left.Weekly == right.Weekly &&
		left.Monthly == right.Monthly
}

func (p *Product) GetPriceByType(priceType int) float64 {
	return GetPriceByTypeFromItems(p.PriceItems, priceType)
}

func (p *Product) GetPrice(priceType int64) float64 {
	return GetPriceByTypeFromItems(p.PriceItems, int(priceType))
}

func LoadPriceItems(ctx context.Context, db *gorm.DB, product *Product) error {
	if product == nil || product.ID == 0 {
		return nil
	}

	var items []ProductPrice
	if err := db.WithContext(ctx).
		Where("product_id = ?", product.ID).
		Order("price_type ASC, id ASC").
		Find(&items).Error; err != nil {
		return err
	}

	product.PriceItems = append([]ProductPrice(nil), items...)
	return nil
}

func LoadPriceItemsForProducts(ctx context.Context, db *gorm.DB, products []Product) error {
	if len(products) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(products))
	seen := make(map[uint]struct{}, len(products))
	for _, product := range products {
		if product.ID == 0 {
			continue
		}
		if _, ok := seen[product.ID]; ok {
			continue
		}
		seen[product.ID] = struct{}{}
		ids = append(ids, product.ID)
	}
	if len(ids) == 0 {
		return nil
	}

	var items []ProductPrice
	if err := db.WithContext(ctx).
		Where("product_id IN ?", ids).
		Order("product_id ASC, price_type ASC, id ASC").
		Find(&items).Error; err != nil {
		return err
	}

	priceMap := make(map[uint][]ProductPrice, len(ids))
	for _, item := range items {
		priceMap[item.ProductID] = append(priceMap[item.ProductID], item)
	}

	for idx := range products {
		products[idx].PriceItems = append([]ProductPrice(nil), priceMap[products[idx].ID]...)
	}
	return nil
}

func UpsertPriceItems(ctx context.Context, tx *gorm.DB, items []ProductPrice) error {
	if len(items) == 0 {
		return nil
	}

	return tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "product_id"},
				{Name: "price_type"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"price", "updated_at"}),
		}).
		Create(&items).Error
}

func DeletePriceItems(ctx context.Context, tx *gorm.DB, productID uint) error {
	if productID == 0 {
		return nil
	}
	return tx.WithContext(ctx).Where("product_id = ?", productID).Delete(&ProductPrice{}).Error
}

type legacyComputePriceRow struct {
	ID           uint    `gorm:"column:id"`
	ProductType  int     `gorm:"column:product_type"`
	PriceHourly  float64 `gorm:"column:price_hourly"`
	PriceDaily   float64 `gorm:"column:price_daily"`
	PriceWeekly  float64 `gorm:"column:price_weekly"`
	PriceMonthly float64 `gorm:"column:price_monthly"`
}

func (legacyComputePriceRow) TableName() string {
	return "products"
}

func SyncLegacyComputePriceItems(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return nil
	}
	if !db.Migrator().HasColumn(&legacyComputePriceRow{}, "PriceHourly") {
		return nil
	}

	var rows []legacyComputePriceRow
	if err := db.WithContext(ctx).
		Where("product_type = ?", ProductTypeCompute).
		Find(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	items := make([]ProductPrice, 0, len(rows)*len(ComputePriceTypes()))
	for _, row := range rows {
		items = append(items,
			ProductPrice{ProductID: row.ID, PriceType: ChargeTypeHourly, Price: row.PriceHourly},
			ProductPrice{ProductID: row.ID, PriceType: ChargeTypeDaily, Price: row.PriceDaily},
			ProductPrice{ProductID: row.ID, PriceType: ChargeTypeWeekly, Price: row.PriceWeekly},
			ProductPrice{ProductID: row.ID, PriceType: ChargeTypeMonthly, Price: row.PriceMonthly},
		)
	}

	return UpsertPriceItems(ctx, db, items)
}
