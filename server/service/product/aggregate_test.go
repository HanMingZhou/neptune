package product

import (
	"context"
	"fmt"
	"testing"

	"gin-vue-admin/global"
	productModel "gin-vue-admin/model/product"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestReleaseResourceAllocationsIsIdempotent(t *testing.T) {
	db := newTestProductDB(t)
	svc := &ProductService{}

	product := productModel.Product{
		ID:           1,
		Name:         "p1",
		ProductType:  productModel.ProductTypeCompute,
		Status:       productModel.ProductStatusEnabled,
		MaxInstances: 4,
		UsedCapacity: 2,
	}
	require.NoError(t, db.Create(&product).Error)
	require.NoError(t, db.Create([]productModel.ResourceAllocation{
		{
			InstanceType:      "training",
			InstanceID:        100,
			ClusterID:         1,
			TemplateProductID: 1,
			ProductID:         1,
			NodeName:          "node-a",
			ScheduleStrategy:  "BALANCED",
			ReplicaIndex:      0,
			TaskRole:          "worker",
			ReservedCount:     1,
		},
		{
			InstanceType:      "training",
			InstanceID:        100,
			ClusterID:         1,
			TemplateProductID: 1,
			ProductID:         1,
			NodeName:          "node-a",
			ScheduleStrategy:  "BALANCED",
			ReplicaIndex:      1,
			TaskRole:          "worker",
			ReservedCount:     1,
		},
	}).Error)

	released, err := svc.ReleaseResourceAllocations(context.Background(), "training", 100)
	require.NoError(t, err)
	require.True(t, released)

	var current productModel.Product
	require.NoError(t, db.First(&current, product.ID).Error)
	require.Equal(t, int64(0), current.UsedCapacity)

	var allocationCount int64
	require.NoError(t, db.Model(&productModel.ResourceAllocation{}).
		Where("instance_type = ? AND instance_id = ?", "training", 100).
		Count(&allocationCount).Error)
	require.Equal(t, int64(0), allocationCount)

	released, err = svc.ReleaseResourceAllocations(context.Background(), "training", 100)
	require.NoError(t, err)
	require.False(t, released)

	require.NoError(t, db.First(&current, product.ID).Error)
	require.Equal(t, int64(0), current.UsedCapacity)
}

func TestLoadPriceItemsForProductsHydratesComputePrices(t *testing.T) {
	db := newTestProductDB(t)

	products := []productModel.Product{
		{
			ID:           11,
			Name:         "gpu-a",
			ProductType:  productModel.ProductTypeCompute,
			Status:       productModel.ProductStatusEnabled,
			MaxInstances: 2,
		},
		{
			ID:           12,
			Name:         "gpu-b",
			ProductType:  productModel.ProductTypeCompute,
			Status:       productModel.ProductStatusEnabled,
			MaxInstances: 4,
		},
	}
	require.NoError(t, db.Create(&products).Error)
	require.NoError(t, productModel.UpsertPriceItems(context.Background(), db, []productModel.ProductPrice{
		{ProductID: 11, PriceType: productModel.ChargeTypeHourly, Price: 1.25},
		{ProductID: 11, PriceType: productModel.ChargeTypeDaily, Price: 10.5},
		{ProductID: 11, PriceType: productModel.ChargeTypeWeekly, Price: 60},
		{ProductID: 11, PriceType: productModel.ChargeTypeMonthly, Price: 180},
		{ProductID: 12, PriceType: productModel.ChargeTypeHourly, Price: 2.5},
		{ProductID: 12, PriceType: productModel.ChargeTypeMonthly, Price: 260},
	}))

	var loaded []productModel.Product
	require.NoError(t, db.Order("id ASC").Find(&loaded).Error)
	require.NoError(t, productModel.LoadPriceItemsForProducts(context.Background(), db, loaded))

	require.Len(t, loaded[0].PriceItems, 4)
	prices0 := productModel.ComputePriceValuesFromItems(loaded[0].PriceItems)
	prices1 := productModel.ComputePriceValuesFromItems(loaded[1].PriceItems)
	require.Equal(t, 1.25, prices0.Hourly)
	require.Equal(t, 10.5, prices0.Daily)
	require.Equal(t, 60.0, prices0.Weekly)
	require.Equal(t, 180.0, prices0.Monthly)
	require.Equal(t, 2.5, prices1.Hourly)
	require.Equal(t, 0.0, prices1.Daily)
	require.Equal(t, 0.0, prices1.Weekly)
	require.Equal(t, 260.0, prices1.Monthly)
}

func TestReleaseResourceAllocationsRollsBackWhenCapacityIsInvalid(t *testing.T) {
	db := newTestProductDB(t)
	svc := &ProductService{}

	product := productModel.Product{
		ID:           2,
		Name:         "p2",
		ProductType:  productModel.ProductTypeCompute,
		Status:       productModel.ProductStatusEnabled,
		MaxInstances: 4,
		UsedCapacity: 0,
	}
	require.NoError(t, db.Create(&product).Error)
	require.NoError(t, db.Create(&productModel.ResourceAllocation{
		InstanceType:      "training",
		InstanceID:        200,
		ClusterID:         1,
		TemplateProductID: 2,
		ProductID:         2,
		NodeName:          "node-b",
		ScheduleStrategy:  "BALANCED",
		ReplicaIndex:      0,
		TaskRole:          "worker",
		ReservedCount:     1,
	}).Error)

	released, err := svc.ReleaseResourceAllocations(context.Background(), "training", 200)
	require.Error(t, err)
	require.False(t, released)

	var current productModel.Product
	require.NoError(t, db.First(&current, product.ID).Error)
	require.Equal(t, int64(0), current.UsedCapacity)

	var allocationCount int64
	require.NoError(t, db.Model(&productModel.ResourceAllocation{}).
		Where("instance_type = ? AND instance_id = ?", "training", 200).
		Count(&allocationCount).Error)
	require.Equal(t, int64(1), allocationCount)
}

func newTestProductDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&productModel.Product{}, &productModel.ProductPrice{}, &productModel.ResourceAllocation{}))

	previous := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = previous
	})

	redisServer, err := miniredis.Run()
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	previousRedis := global.GVA_REDIS
	global.GVA_REDIS = redisClient
	t.Cleanup(func() {
		require.NoError(t, redisClient.Close())
		redisServer.Close()
		global.GVA_REDIS = previousRedis
	})

	return db
}
