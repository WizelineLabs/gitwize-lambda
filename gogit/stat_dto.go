package gogit

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
)

// commitDto data object for commit entity
type commitDto struct {
	RepositoryID int
	Hash         string
	AuthorEmail  string
	AuthorName   string
	Message      string
	NumFiles     int
	AdditionLOC  int
	DeletionLOC  int
	NumParents   int
	LOC          int
	Year         int
	Month        int
	Day          int
	Hour         int
	TimeStamp    string
}

type fileStatDTO struct {
	RepositoryID int
	Hash         string
	AuthorEmail  string
	AuthorName   string
	FileName     string
	AdditionLOC  int
	DeletionLOC  int
	Year         int
	Month        int
	Day          int
	Hour         int
	TimeStamp    string
}

type dtoInterface interface {
	getListValues() []interface{}
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
	dto.LOC = getLineOfCode(c)
	// dto.LOC = 0 // temporary disable getting total loc, to impove perf
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
