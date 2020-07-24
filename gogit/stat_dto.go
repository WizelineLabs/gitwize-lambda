package gogit

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"reflect"
)

// commitDto data object for commit entity
type commitDto struct {
	RepositoryID   int    `json:"repository_id"`
	Hash           string `json:"hash"`
	AuthorEmail    string `json:"author_email"`
	AuthorName     string `json:"author_name"`
	Message        string `json:"message"`
	NumFiles       int    `json:"num_files"`
	AdditionLOC    int    `json:"addition_loc"`
	DeletionLOC    int    `json:"deletion_loc"`
	NumParents     int    `json:"num_parents"`
	InsertionPoint int    `json:"insertion_point"`
	LOC            int    `json:"total_loc"`
	Year           int    `json:"year"`
	Month          int    `json:"month"`
	Day            int    `json:"day"`
	Hour           int    `json:"hour"`
	TimeStamp      string `json:"commit_time_stamp"`
}

type fileStatDTO struct {
	RepositoryID int    `json:"repository_id"`
	Hash         string `json:"hash"`
	AuthorEmail  string `json:"author_email"`
	AuthorName   string `json:"author_name"`
	FileName     string `json:"file_name"`
	AdditionLOC  int    `json:"addition_loc"`
	DeletionLOC  int    `json:"deletion_loc"`
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Day          int    `json:"day"`
	Hour         int    `json:"hour"`
	TimeStamp    string `json:"commit_time_stamp"`
}

type dtoInterface interface {
	getListValues() []interface{}
	getFieldNames() []string
}

func getFieldNames(item interface{}) []string {
	val := reflect.ValueOf(item)
	names := make([]string, val.Type().NumField())
	for i := 0; i < len(names); i++ {
		names[i] = val.Type().Field(i).Tag.Get("json")
	}
	return names
}

func (dto commitDto) getFieldNames() []string {
	return getFieldNames(dto)
}

func (fdto fileStatDTO) getFieldNames() []string {
	return getFieldNames(fdto)
}

func (dto commitDto) getListValues() []interface{} {
	return []interface{}{
		dto.RepositoryID,
		dto.Hash,
		dto.AuthorEmail,
		dto.AuthorName,
		dto.Message,
		dto.NumFiles,
		dto.AdditionLOC,
		dto.DeletionLOC,
		dto.NumParents,
		dto.InsertionPoint,
		dto.LOC,
		dto.Year,
		dto.Month,
		dto.Day,
		dto.Hour,
		dto.TimeStamp,
	}
}

func (dto fileStatDTO) getListValues() []interface{} {
	return []interface{}{
		dto.RepositoryID,
		dto.Hash,
		dto.AuthorEmail,
		dto.AuthorName,
		dto.FileName,
		dto.AdditionLOC,
		dto.DeletionLOC,
		dto.Year,
		dto.Month,
		dto.Day,
		dto.Hour,
		dto.TimeStamp,
	}
}

func getCommitDTO(c *object.Commit) commitDto {
	dto := commitDto{}
	dto.Hash = c.Hash.String()
	dto.AuthorEmail = c.Author.Email
	dto.AuthorName = c.Author.Name
	dto.Message = c.Message
	dto.Year = c.Author.When.UTC().Year()
	dto.Month = int(c.Author.When.UTC().Month())
	dto.Day = c.Author.When.UTC().Day()
	dto.Hour = c.Author.When.UTC().Hour()
	dto.TimeStamp = c.Author.When.UTC().String()
	// dto.LOC = getLineOfCode(c)
	dto.LOC = 0 //  disable getting total loc as this not used, need to remove go routine to avoid issue too many open file
	fileStats, err := c.Stats()
	if err != nil {
		log.Panicln(err)
	}
	dto.NumFiles = len(fileStats)
	for _, file := range fileStats {
		dto.AdditionLOC += file.Addition
		dto.DeletionLOC += file.Deletion
	}
	dto.NumParents = c.NumParents()
	return dto
}

func getLineOfCode(c *object.Commit) (loc int) {
	fileIter, err := c.Files()
	if err != nil {
		panic(err.Error())
	}
	err = fileIter.ForEach(func(f *object.File) error {
		lines, _ := f.Lines()
		loc += len(lines)
		return nil
	})
	return loc
}

func getFileStatDTO(c *object.Commit, rID int) []fileStatDTO {
	fileStats, err := c.Stats()
	if err != nil {
		log.Panicln(err)
	}
	dtos := make([]fileStatDTO, len(fileStats))
	for i, file := range fileStats {
		dto := fileStatDTO{}
		dto.RepositoryID = rID
		dto.Hash = c.Hash.String()
		dto.AuthorEmail = c.Author.Email
		dto.AuthorName = c.Author.Name
		dto.FileName = file.Name
		dto.AdditionLOC = file.Addition
		dto.DeletionLOC = file.Deletion
		dto.Year = c.Author.When.UTC().Year()
		dto.Month = int(c.Author.When.UTC().Month())
		dto.Day = c.Author.When.UTC().Day()
		dto.Hour = c.Author.When.UTC().Hour()
		dto.TimeStamp = c.Author.When.UTC().String()
		dtos[i] = dto
	}
	return dtos
}

func convertFileDtosToDtoInterfaces(fdtos []fileStatDTO) []dtoInterface {
	result := make([]dtoInterface, len(fdtos))
	for i, v := range fdtos {
		result[i] = dtoInterface(v)
	}
	return result
}
