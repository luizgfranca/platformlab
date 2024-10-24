package service

import (
	"fmt"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/util"

	"gorm.io/gorm"
)

type Tool struct {
	Db *gorm.DB
}

func (t *Tool) FindAll() *[]model.Tool {
	var tools []model.Tool

	result := t.Db.Find(&tools)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &tools
}

func (t *Tool) FindAssociatedWithProject(project *model.Project) *[]model.Tool {
	var tools []model.Tool

	result := t.Db.Find(&tools).Where("projectId = ?", project.ID)
	if result.Error != nil {
		panic(fmt.Sprintf("unable to query database: %s", result.Error.Error()))
	}

	return &tools
}

func (t *Tool) Create(tool *model.Tool) (*model.Tool, error) {
	var result *gorm.DB
	var maybeExisting *model.Tool

	util.Probe(tool.Acronym)

	result = t.Db.Where("acronym = ?", tool.Acronym).First(&maybeExisting)
	if result.Error == nil {
		return nil, &model.GenericLogicError{
			Message: fmt.Sprintf("element with acronym %s already exists", tool.Acronym),
		}
	}

	result = t.Db.Create(tool)
	if result.Error != nil {
		return nil, result.Error
	}

	var created *model.Tool
	result = t.Db.Where("acronym = ?", tool.Acronym).First(&created)
	if result.Error != nil {
		return nil, result.Error
	}
	if created == nil {
		panic("created item in database, but it was not found after insertion")
	}

	return created, nil
}
