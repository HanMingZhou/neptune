package system

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/request"
	"gin-vue-admin/model/system"
	systemReq "gin-vue-admin/model/system/request"

	"github.com/pkg/errors"
)

//@author: [granty1](https://github.com/granty1)
//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

type OperationRecordService struct{}

var OperationRecordServiceApp = new(OperationRecordService)

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids request.IdsReq
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecordByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	if err != nil {
		return errors.Errorf("OperationRecordService.DeleteSysOperationRecordByIds global.GVA_DB.Delete failed: %w", err)
	}
	return nil
}

//@author: [granty1](https://github.com/granty1)
//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	err = global.GVA_DB.Delete(&sysOperationRecord).Error
	if err != nil {
		return errors.Errorf("OperationRecordService.DeleteSysOperationRecord global.GVA_DB.Delete failed: %w", err)
	}
	return nil
}

//@author: [granty1](https://github.com/granty1)
//@function: GetSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: sysOperationRecord system.SysOperationRecord, err error

func (operationRecordService *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord system.SysOperationRecord, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&sysOperationRecord).Error
	if err != nil {
		return sysOperationRecord, errors.Errorf("OperationRecordService.GetSysOperationRecord global.GVA_DB.First failed: %w", err)
	}
	return
}

//@author: [granty1](https://github.com/granty1)
//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info systemReq.SysOperationRecordSearch
//@return: list interface{}, total int64, err error

func (operationRecordService *OperationRecordService) GetSysOperationRecordInfoList(info systemReq.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&system.SysOperationRecord{})
	var sysOperationRecords []system.SysOperationRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, errors.Errorf("OperationRecordService.GetSysOperationRecordInfoList db.Count failed: %w", err)
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	if err != nil {
		return nil, 0, errors.Errorf("OperationRecordService.GetSysOperationRecordInfoList db.Find failed: %w", err)
	}
	return sysOperationRecords, total, err
}
