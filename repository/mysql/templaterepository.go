package mysql

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"goacs/models/templates"
	"goacs/repository"
)

type TemplateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(connection *sqlx.DB) *TemplateRepository {
	return &TemplateRepository{
		db: connection,
	}
}

func (r *TemplateRepository) Find(id int64) (*templates.Template, error) {
	templateInstance := new(templates.Template)
	err := r.db.Unsafe().Get(templateInstance, "SELECT * FROM templates WHERE id=? LIMIT 1", id)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound

	}

	return templateInstance, nil
}

func (r *TemplateRepository) FindByName(name string) (*templates.Template, error) {
	templateInstance := new(templates.Template)

	err := r.db.Unsafe().Get(templateInstance, "SELECT id, name FROM templates WHERE name=? LIMIT 1", name)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound

	}

	return templateInstance, nil
}

func (r *TemplateRepository) List(request repository.PaginatorRequest) ([]templates.Template, int) {
	var total int
	var templates = make([]templates.Template, 0)
	_ = r.db.Get(&total, "SELECT count(*) FROM templates")
	err := r.db.Unsafe().Select(&templates, "SELECT * FROM templates LIMIT ?,?", request.CalcOffset(), request.PerPage)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, 0
	}

	return templates, total
}

func (r *TemplateRepository) GetParametersForTemplate(template_id int64) ([]templates.TemplateParameter, error) {
	var parameters = []templates.TemplateParameter{}

	err := r.db.Unsafe().Select(&parameters,
		"SELECT * FROM templates_parameters WHERE template_id=? LIMIT 1",
		template_id,
	)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound
	}

	return parameters, nil
}

func (r *TemplateRepository) HydrateTemplatesParameters(templatesData []templates.Template) []templates.Template {

	var parameters []templates.TemplateParameter
	var ids []int64

	for _, template := range templatesData {
		ids = append(ids, template.Id)
	}

	dialect := goqu.Dialect("mysql")
	selectSql, _, _ := dialect.From("templates_parameters").
		Where(goqu.C("template_id").In(ids)).ToSQL()

	err := r.db.Select(&parameters, selectSql)

	if err != nil {
		return templatesData
	}

	for templateIdx, template := range templatesData {
		for _, parameter := range parameters {
			if parameter.TemplateId == template.Id {
				templatesData[templateIdx].Parameters = parameters
			}
		}
	}

	return templatesData
}
