package mysql

import (
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/thoas/go-funk"
	"goacs/acs/types"
	"goacs/models/cpe"
	"goacs/repository"
	"log"
	"strings"
	"time"
)

type CPERepository struct {
	db *sqlx.DB
}

func NewCPERepository(connection *sqlx.DB) CPERepository {
	return CPERepository{
		db: connection,
	}
}

func (r *CPERepository) All() ([]*cpe.CPE, error) {
	var cpes = []*cpe.CPE{}
	err := r.db.Unsafe().Select(&cpes, "SELECT * FROM cpe")

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound
	}

	return cpes, nil
}

func (r *CPERepository) Count() (cpe_count int64) {
	err := r.db.Unsafe().Get(&cpe_count, "SELECT count(uuid) FROM cpe")

	if err != nil {
		return 0
	}

	return cpe_count
}

func (r *CPERepository) List(request repository.PaginatorRequest) ([]cpe.CPE, int) {
	var total int
	var cpes = make([]cpe.CPE, 0)
	_ = r.db.Get(&total, "SELECT count(*) FROM cpe")
	err := r.db.Unsafe().Select(&cpes, "SELECT * FROM cpe LIMIT ?,?", request.CalcOffset(), request.PerPage)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, 0
	}

	return cpes, total
}

func (r *CPERepository) Find(uuid string) (*cpe.CPE, error) {
	cpeInstance := new(cpe.CPE)
	err := r.db.Unsafe().Get(cpeInstance, "SELECT * FROM cpe WHERE uuid=? LIMIT 1", uuid)

	if err == sql.ErrNoRows {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound
	}

	return cpeInstance, nil
}

func (r *CPERepository) FindBySerial(serial string) (*cpe.CPE, error) {
	cpeInstance := new(cpe.CPE)
	err := r.db.Unsafe().Get(cpeInstance, "SELECT * FROM cpe WHERE serial_number=? LIMIT 1", serial)

	if err != nil {
		fmt.Println("Error while fetching query results")
		fmt.Println(err.Error())
		return nil, repository.ErrNotFound
	}

	return cpeInstance, nil
}

func (r *CPERepository) Create(cpe *cpe.CPE) (bool, error) {
	uuidInstance, _ := uuid.NewRandom()
	uuidString := uuidInstance.String()

	_, err := r.db.Exec(`INSERT INTO cpe SET uuid=?, 
			serial_number=?, 
			hardware_version=?, 
			software_version=?, 
      connection_request_url=?,
      connection_request_user=?,
      connection_request_password=?,              
			created_at=?, 
			updated_at=?
			`,
		uuidString,
		cpe.SerialNumber,
		cpe.HardwareVersion,
		cpe.SoftwareVersion,
		cpe.ConnectionRequestUrl,
		cpe.ConnectionRequestUser,
		cpe.ConnectionRequestPassword,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		fmt.Println(err)
		return false, repository.ErrInserting
	}

	cpe.UUID = uuidInstance.String()
	//cpe.NewInACS = true

	return true, nil
}

func (r *CPERepository) DeleteDevice(cpe *cpe.CPE) {
	dialect := goqu.Dialect("mysql")

	query, args, _ := dialect.Delete("cpe").
		Prepared(true).
		Where(goqu.C("cpe_uuid").Eq(cpe.UUID)).
		ToSQL()

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Cannot delete device ", err)
	}

}

func (r *CPERepository) UpdateOrCreate(cpe *cpe.CPE) (result bool, cpeExist bool, err error) {

	dbCPE, _ := r.FindBySerial(cpe.SerialNumber)

	if dbCPE == nil {
		result, err = r.Create(cpe)
		cpeExist = false
	} else {
		fmt.Println("Updating CPE")
		stmt, _ := r.db.Prepare(`UPDATE cpe SET 
               hardware_version=?, 
               software_version=?, 
               connection_request_url=?,   
               updated_at=? 
			   WHERE uuid=?`)

		_, err := stmt.Exec(
			cpe.HardwareVersion,
			cpe.SoftwareVersion,
			cpe.ConnectionRequestUrl,
			time.Now(),
			dbCPE.UUID,
		)
		cpe.UUID = dbCPE.UUID

		if err != nil {
			log.Println("error while updatng cpe " + err.Error())
			return false, false, repository.ErrUpdating
		}

		result = true
		err = nil
		cpeExist = true

	}

	return
}

func (r *CPERepository) FindParameter(cpe *cpe.CPE, parameterKey string) (*types.ParameterValueStruct, error) {
	parameterValueStruct := new(types.ParameterValueStruct)
	err := r.db.Unsafe().Get(parameterValueStruct, "SELECT *  FROM cpe_parameters WHERE cpe_uuid=? AND name=? LIMIT 1", cpe.UUID, parameterKey)

	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}

	return parameterValueStruct, nil
}

func (r *CPERepository) CreateParameter(cpe *cpe.CPE, parameter types.ParameterValueStruct) (bool, error) {
	var query string = `INSERT INTO cpe_parameters (cpe_uuid, name, value, type, flags, created_at, updated_at) 
						VALUES (?, ?, ?, ?, ?, ?, ?)`

	stmt, _ := r.db.Prepare(query)

	_, err := stmt.Exec(
		cpe.UUID,
		parameter.Name,
		parameter.Value,
		parameter.Type,            //TODO: NORMALIZE
		parameter.Flag.AsString(), //TODO: Flags support (R - Read, W - Write and more...)
		time.Now(),
		time.Now(),
	)

	if err != nil {
		fmt.Println(repository.ErrParameterCreating, err.Error())
		return false, err
	}

	return true, nil
}

func (r *CPERepository) BulkInsertOrUpdateParameters(cpe *cpe.CPE, parameters []types.ParameterValueStruct) bool {
	tx, err := r.db.Begin()

	if err != nil {
		log.Println("Cannot create TX for BulkInsertOrUpdateParameters ", err.Error())
		return false
	}

	chunks := funk.Chunk(parameters, 300)
	for _, chunk := range chunks.([][]types.ParameterValueStruct) {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, parameter := range chunk {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, cpe.UUID)
			valueArgs = append(valueArgs, parameter.Name)
			valueArgs = append(valueArgs, parameter.Value)
			valueArgs = append(valueArgs, parameter.Type)
			valueArgs = append(valueArgs, parameter.Flag.AsString())
		}

		stmt := fmt.Sprintf("INSERT INTO cpe_parameters(cpe_uuid,name,value,type,flags) VALUES %s "+
			"ON DUPLICATE KEY UPDATE name=values(name),value=values(value), type=values(type), flags=values(flags)", strings.Join(valueStrings, ","))
		_, err := tx.Exec(stmt, valueArgs...)

		if err != nil {
			_ = tx.Rollback()
			fmt.Println(err.Error())
			return false
		}
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func (r *CPERepository) UpdateOrCreateParameter(cpe *cpe.CPE, parameter types.ParameterValueStruct) (result bool, err error) {
	//log.Println("UoCP ", parameter.Name)
	//log.Println("UoCP ", parameter.Value)

	existParameter, err := r.FindParameter(cpe, parameter.Name)

	if existParameter == nil {
		//fmt.Println("non exist param", existParameter)
		result, err = r.CreateParameter(cpe, parameter)
	} else {
		//fmt.Println("param exist", existParameter)
		result, err = r.UpdateParameter(cpe, parameter)
	}

	return
}

func (r *CPERepository) UpdateParameter(cpe *cpe.CPE, parameter types.ParameterValueStruct) (result bool, err error) {
	query := "UPDATE cpe_parameters SET value=?, type=?, flags=?, updated_at=? WHERE cpe_uuid=? and name = ?"
	stmt, _ := r.db.Prepare(query)

	//log.Println("Parameter flags ", parameter.Flag.AsString())

	_, err = stmt.Exec(
		parameter.Value,
		parameter.Type,
		parameter.Flag.AsString(),
		time.Now(),
		cpe.UUID,
		parameter.Name,
	)

	if err != nil {
		fmt.Println("ERROR", err.Error())
		result = false
	}

	return
}

func (r *CPERepository) SaveParameters(cpe *cpe.CPE) (bool, error) {

	for _, parameterValue := range cpe.ParameterValues {
		//fmt.Println("param value", parameterValue)
		_, err := r.UpdateOrCreateParameter(cpe, parameterValue)

		if err != nil {
			fmt.Println(repository.ErrParameterCreating, err.Error())
			return false, err
		}
	}

	return true, nil
}

func (r *CPERepository) GetCPEParameters(cpe *cpe.CPE) ([]types.ParameterValueStruct, error) {
	var parameters = []types.ParameterValueStruct{}

	err := r.db.Select(&parameters, "SELECT * FROM cpe_parameters WHERE cpe_uuid=?", cpe.UUID)

	if err != nil {
		log.Println(err.Error())
		log.Println("CPE UUID ", cpe.UUID)
		log.Println(parameters)
		return nil, repository.ErrNotFound
	}

	return parameters, nil
}

func (r *CPERepository) ListCPEParameters(cpe *cpe.CPE, request repository.PaginatorRequest) ([]types.ParameterValueStruct, int) {
	dialect := goqu.Dialect("mysql")

	baseBulder := dialect.From("cpe_parameters").
		Where(goqu.C("cpe_uuid").Eq(cpe.UUID))

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
	parameters := make([]types.ParameterValueStruct, 0)
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

func (r *CPERepository) LoadParameters(cpe *cpe.CPE) (bool, error) {
	var err error
	cpe.ParameterValues, err = r.GetCPEParameters(cpe)

	return err == nil, err
}

func (r *CPERepository) DeleteAllParameters(cpe *cpe.CPE) {
	dialect := goqu.Dialect("mysql")

	query, args, _ := dialect.Delete("cpe_parameters").
		Prepared(true).
		Where(goqu.C("cpe_uuid").Eq(cpe.UUID)).
		ToSQL()

	_, err := r.db.Exec(query, args...)

	if err != nil {
		log.Println("Cannot delete device ", err)
	}
}
