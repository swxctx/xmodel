package create

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"strings"
	"unsafe"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/henrylee2cn/goutil"
	"github.com/xiaoenai/xmodel/cmd/info"
)

type (
	// Project project Information
	Project struct {
		*tplInfo
		codeFiles    map[string]string
		Name         string
		ImprotPrefix string
	}
	Model struct {
		*structType
		ModelStyle       string
		PrimaryFields    []*field
		UniqueFields     []*field
		Fields           []*field
		IsDefaultPrimary bool
		Doc              string
		Name             string
		SnakeName        string
		LowerFirstName   string
		LowerFirstLetter string
		NameSql          string
		QuerySql         [2]string
		UpdateSql        string
		UpsertSqlSuffix  string
	}
)

// NewProject new project.
func NewProject(src []byte) *Project {
	p := new(Project)
	p.tplInfo = newTplInfo(src).Parse()
	p.Name = info.ProjName()
	p.ImprotPrefix = info.ProjPath()
	p.codeFiles = make(map[string]string)
	for k, v := range tplFiles {
		p.codeFiles[k] = v
	}
	for k := range p.codeFiles {
		p.fillFile(k)
	}
	return p
}

func (p *Project) fillFile(k string) {
	v, ok := p.codeFiles[k]
	if !ok {
		return
	}
	v = strings.Replace(v, "${import_prefix}", p.ImprotPrefix, -1)
	switch k {
	case "model/init.go":
		p.codeFiles[k] = v
	default:
		p.codeFiles[k] = "// Code generated by 'xmodel gen' command.\n// DO NOT EDIT!\n\n" + v
	}
}

func mustMkdirAll(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		erpc.Fatalf("[XModel] %v", err)
	}
}

func hasGenSuffix(name string) bool {
	switch name {
	case "model/init.go":
		return false
	default:
		return true
	}
}

func (p *Project) Generator(force bool) {
	p.gen()
	// make all directorys
	mustMkdirAll("args")
	mustMkdirAll("model")
	// write files
	for k, v := range p.codeFiles {
		if !force && !hasGenSuffix(k) {
			continue
		}
		realName := info.ProjPath() + "/" + k
		f, err := os.OpenFile(k, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			erpc.Fatalf("[XModel] create %s error: %v", realName, err)
		}
		b := formatSource(goutil.StringToBytes(v))
		f.Write(b)
		f.Close()
		fmt.Printf("generate %s\n", realName)
	}
}

// generate all codes
func (p *Project) gen() {
	p.genConstFile()
	p.genTypeFile()
	p.genModelFile()
	p.genAndWriteGoModFile()
}

func (p *Project) genAndWriteGoModFile() {
	f, err := os.OpenFile("./go.mod", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		erpc.Fatalf("[XModel] create go.mod error: %v", err)
	}
	f.WriteString(p.genGoMod())
	f.Close()
	fmt.Printf("generate %s\n", info.ProjPath()+"/go.mod")
}

func (p *Project) fieldsJson(fs []*field) string {
	if len(fs) == 0 {
		return ""
	}
	var text string
	text += "{"
	for _, f := range fs {
		if f.isQuery {
			continue
		}
		fieldName := f.ModelName
		if len(fieldName) == 0 {
			fieldName = goutil.SnakeString(f.Name)
		}
		t := strings.Replace(f.Typ, "*", "", -1)
		var isSlice bool
		if strings.HasPrefix(t, "[]") {
			if t == "[]byte" {
				t = "string"
			} else {
				t = strings.TrimPrefix(t, "[]")
				isSlice = true
			}
		}
		v, ok := baseTypeToJsonValue(t)
		if ok {
			if isSlice {
				text += fmt.Sprintf(`"%s$%d":[%s],`, fieldName, uintptr(unsafe.Pointer(f)), v)
			} else {
				text += fmt.Sprintf(`"%s$%d":%s,`, fieldName, uintptr(unsafe.Pointer(f)), v)
			}
			continue
		}
		if ffs, ok := p.tplInfo.lookupTypeFields(t); ok {
			if isSlice {
				text += fmt.Sprintf(`"%s":[%s],`, fieldName, p.fieldsJson(ffs))
			} else {
				text += fmt.Sprintf(`"%s":%s,`, fieldName, p.fieldsJson(ffs))
			}
			continue
		}
	}
	text = strings.TrimRight(text, ",") + "}"
	return text
}

func baseTypeToJsonValue(t string) (string, bool) {
	if t == "bool" {
		return "false", true
	} else if t == "string" || t == "[]byte" || t == "time.Time" {
		return `""`, true
	} else if strings.HasPrefix(t, "int") || t == "rune" {
		return "-0", true
	} else if strings.HasPrefix(t, "uint") || t == "byte" {
		return "0", true
	} else if strings.HasPrefix(t, "float") {
		return "-0.000000", true
	}
	return "", false
}

func (p *Project) genConstFile() {
	var text string
	for _, s := range p.tplInfo.models.mysql {
		name := s.name + "Sql"
		text += fmt.Sprintf(
			"// %s the statement to create '%s' mysql table\n"+
				"const %s string = ``\n",
			name, goutil.SnakeString(s.name),
			name,
		)
	}
	p.replaceWithLine("args/const.gen.go", "${const_list}", text)
}

func (p *Project) genTypeFile() {
	p.replaceWithLine("args/type.gen.go", "${import_list}", p.tplInfo.TypeImportString())
	p.replaceWithLine("args/type.gen.go", "${type_define_list}", p.tplInfo.TypesString())
}

func (p *Project) genModelFile() {
	for _, m := range p.tplInfo.models.mysql {
		fileName := "model/mysql_" + goutil.SnakeString(m.name) + ".gen.go"
		p.codeFiles[fileName] = newModelString(m)
		p.fillFile(fileName)
	}
	for _, m := range p.tplInfo.models.mongo {
		fileName := "model/mongo_" + goutil.SnakeString(m.name) + ".gen.go"
		p.codeFiles[fileName] = newModelString(m)
		p.fillFile(fileName)
	}
}

func newModelString(s *structType) string {
	model := &Model{
		structType:       s,
		PrimaryFields:    s.primaryFields,
		UniqueFields:     s.uniqueFields,
		IsDefaultPrimary: s.isDefaultPrimary,
		Fields:           s.fields,
		Doc:              s.doc,
		Name:             s.name,
		ModelStyle:       s.modelStyle,
		SnakeName:        goutil.SnakeString(s.name),
	}
	model.LowerFirstLetter = strings.ToLower(model.Name[:1])
	model.LowerFirstName = model.LowerFirstLetter + model.Name[1:]
	switch s.modelStyle {
	case "mysql":
		return model.mysqlString()
	case "mongo":
		return model.mongoString()
	}
	return ""
}

func (mod *Model) mongoString() string {
	mod.NameSql = fmt.Sprintf("`%s`", mod.SnakeName)
	mod.QuerySql = [2]string{}
	mod.UpdateSql = ""
	mod.UpsertSqlSuffix = ""

	var (
		fields               []string
		querySql1, querySql2 string
	)
	for _, field := range mod.fields {
		fields = append(fields, field.ModelName)
	}
	var primaryFields []string
	var primaryFieldMap = make(map[string]bool)
	for _, field := range mod.PrimaryFields {
		primaryFields = append(primaryFields, field.ModelName)
		primaryFieldMap[field.ModelName] = true
	}
	for _, field := range fields {
		if field == "deleted_ts" || primaryFieldMap[field] {
			continue
		}
		querySql1 += fmt.Sprintf("`%s`,", field)
		querySql2 += fmt.Sprintf(":%s,", field)
		if field == "created_at" {
			continue
		}
		mod.UpdateSql += fmt.Sprintf("`%s`=:%s,", field, field)
		mod.UpsertSqlSuffix += fmt.Sprintf("`%s`=VALUES(`%s`),", field, field)
	}
	mod.QuerySql = [2]string{querySql1[:len(querySql1)-1], querySql2[:len(querySql2)-1]}
	mod.UpdateSql = mod.UpdateSql[:len(mod.UpdateSql)-1]
	mod.UpsertSqlSuffix = mod.UpsertSqlSuffix[:len(mod.UpsertSqlSuffix)-1] + ";"

	m, err := template.New("").Parse(mongoModelTpl)
	if err != nil {
		erpc.Fatalf("[XModel] model string: %v", err)
	}
	buf := bytes.NewBuffer(nil)
	err = m.Execute(buf, mod)
	if err != nil {
		erpc.Fatalf("[XModel] model string: %v", err)
	}
	s := strings.Replace(buf.String(), "&lt;", "<", -1)
	return strings.Replace(s, "&gt;", ">", -1)
}

func (mod *Model) mysqlString() string {
	mod.NameSql = fmt.Sprintf("`%s`", mod.SnakeName)
	mod.QuerySql = [2]string{}
	mod.UpdateSql = ""
	mod.UpsertSqlSuffix = ""

	var (
		fields               []string
		querySql1, querySql2 string
	)
	for _, field := range mod.fields {
		fields = append(fields, field.ModelName)
	}
	var primaryFields []string
	var primaryFieldMap = make(map[string]bool)
	for _, field := range mod.PrimaryFields {
		primaryFields = append(primaryFields, field.ModelName)
		primaryFieldMap[field.ModelName] = true
	}
	for _, field := range fields {
		if field == "deleted_ts" || primaryFieldMap[field] {
			continue
		}
		querySql1 += fmt.Sprintf("`%s`,", field)
		querySql2 += fmt.Sprintf(":%s,", field)
		if field == "created_at" {
			continue
		}
		mod.UpdateSql += fmt.Sprintf("`%s`=:%s,", field, field)
		mod.UpsertSqlSuffix += fmt.Sprintf("`%s`=VALUES(`%s`),", field, field)
	}
	mod.QuerySql = [2]string{querySql1[:len(querySql1)-1], querySql2[:len(querySql2)-1]}
	mod.UpdateSql = mod.UpdateSql[:len(mod.UpdateSql)-1]
	mod.UpsertSqlSuffix = mod.UpsertSqlSuffix[:len(mod.UpsertSqlSuffix)-1] + ";"

	m, err := template.New("").Parse(mysqlModelTpl)
	if err != nil {
		erpc.Fatalf("[XModel] model string: %v", err)
	}
	buf := bytes.NewBuffer(nil)
	err = m.Execute(buf, mod)
	if err != nil {
		erpc.Fatalf("[XModel] model string: %v", err)
	}
	s := strings.Replace(buf.String(), "&lt;", "<", -1)
	return strings.Replace(s, "&gt;", ">", -1)
}

func (p *Project) replace(key, placeholder, value string) string {
	a := strings.Replace(p.codeFiles[key], placeholder, value, -1)
	p.codeFiles[key] = a
	return a
}

func (p *Project) replaceWithLine(key, placeholder, value string) string {
	return p.replace(key, placeholder, "\n"+value)
}

func formatSource(src []byte) []byte {
	b, err := format.Source(src)
	if err != nil {
		erpc.Fatalf("[XModel] format error: %v\ncode:\n%s", err, src)
	}
	return b
}

func (p *Project) genGoMod() string {
	r := strings.Replace(__gomod__, "${import_prefix}", p.ImprotPrefix, -1)
	return r
}

