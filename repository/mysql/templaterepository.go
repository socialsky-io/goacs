package mysql

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"goacs/models/cpe"
	"goacs/models/templates"
	"goacs/repository"
	"log"
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

func (r *TemplateRepository) ListTemplateParameters(template *templates.Template, request repository.PaginatorRequest) ([]templates.TemplateParameter, int) {
	dialect := goqu.Dialect("mysql")

	baseBulder := dialect.From("templates_parameters").
		Where(goqu.C("template_id").Eq(template.Id))

	if len(request.Filter) > 0 {
		for key, value := range request.Filter {
			baseBulder = baseBulder.Where(goqu.Ex{
				key: goqu.Op{"ilike": "%" + value + "%"},
			})
		}
	}

	totalSql, _, _ := baseBulder.
		Select(goqu.COUNT("*")).
		ToSQL()

	var total int
	_ = r.db.Get(&total, totalSql)
	var parameters []templates.TemplateParameter
	parametersBuilder := baseBulder.
		Offset(uint(request.CalcOffset())).
		Limit(uint(request.PerPage))

	log.Println(request.Filter)

	parametersSql, _, _ := parametersBuilder.ToSQL()

	log.Println(parametersSql)

	err := r.db.Unsafe().Select(&parameters, parametersSql)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, 0
	}

	return parameters, total
}

func (r *TemplateRepository) GetTemplatesForCPE(cpe *cpe.CPE) []templates.CPETemplate {
	dialect := goqu.Dialect("mysql")

	query, _, _ := dialect.From("cpe_to_templates").Join(
		goqu.T("templates").As("t"),
		goqu.On(goqu.Ex{"cpe_to_templates.template_id": goqu.I("t.id")}),
	).Where(goqu.Ex{"cpe_to_templates.cpe_uuid": cpe.UUID}).
		Order(goqu.I("priority").Desc()).ToSQL()

	var cpeTemplates []templates.CPETemplate
	err := r.db.Select(&cpeTemplates, query)
	if err != nil {
		log.Println(err)
		return []templates.CPETemplate{}
	}

	return cpeTemplates
}
