package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"istomyang.github.com/like-iam/log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// csvPump outputs log to file in local storage.
// like: ~/workspaces/some_path/iam/2022-09-09-23.csv
type csvPump struct {
	fileDir    string
	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter

	ctx    context.Context
	cancel context.CancelFunc
}

func New() pumps.Pump {
	var cp = csvPump{}
	return &cp
}

func (c *csvPump) Init(ctx context.Context, config map[string]any) error {
	a, ex := config["file-dir"]
	fd := a.(string)
	if !ex || fd == "" {
		pwd, _ := os.Getwd()
		fd = pwd + "/iam-pump"
	}
	c.fileDir = fd
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *csvPump) Run() error {
	if err := os.MkdirAll(c.fileDir, 0777); err != nil {
		return err
	}
	return nil
}

func (c *csvPump) Close() error {
	c.cancel()
	return nil
}

func (c *csvPump) GetName() string {
	return "CSV Pump"
}

func (c *csvPump) SetFilter(filter *pumps.Filter) {
	c.filter = filter
}

func (c *csvPump) GetFilter() *pumps.Filter {
	return c.filter
}

func (c *csvPump) SetTimeout(duration time.Duration) {
	c.timeout = duration
}

func (c *csvPump) GetTimeout() time.Duration {
	return c.timeout
}

func (c *csvPump) SetOmitDetail(b bool) {
	c.omitDetail = b
}

func (c *csvPump) GetOmitDetail() bool {
	return c.omitDetail
}

func (c *csvPump) Write(i []interface{}) error {

	file, newFile, err := c.openFile(c.genFilePath())
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	w := csv.NewWriter(file)

	if newFile {
		if err := w.Write(c.header(i)); err != nil {
			log.Errorf("%s: fail to append header to file, %s", c.GetName(), err.Error())
			return err
		}
	}

	for _, d := range c.body(i) {
		if err := w.Write(d); err != nil {
			log.Errorf("%s: fail to write to file, %s", c.GetName(), err.Error())
			return err
		}
	}

	w.Flush()

	return nil
}

func (c *csvPump) genFilePath() string {
	t0 := time.Now()
	fileName := fmt.Sprintf("%d-%d-%d-%d.csv", t0.Year(), t0.Month(), t0.Day(), t0.Hour())
	return path.Join(c.fileDir, fileName)
}

func (c *csvPump) openFile(fp string) (*os.File, bool, error) {
	var file *os.File
	var newOne bool
	var err error

	if _, err = os.Stat(fp); os.IsNotExist(err) {
		err = nil
		file, err = os.Create(fp)
		newOne = true
		if err != nil {
			log.Errorf("%s: can't create csv file, %s", c.GetName(), err)
			return nil, newOne, err
		}
	} else {
		err = nil
		file, err = os.OpenFile(fp, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Errorf("%s: can't open csv file, %s", c.GetName(), err)
			return nil, newOne, err
		}
	}

	return file, newOne, nil
}

func (c *csvPump) header(data []interface{}) []string {
	infoType := reflect.TypeOf(data).Elem()
	var r = make([]string, infoType.NumField())

	for i := 0; i < infoType.NumField(); i++ {
		r[i] = infoType.Field(i).Name
	}
	return r
}

func (c *csvPump) body(dataset []interface{}) [][]string {
	var r = make([][]string, len(dataset))

	for i, data := range dataset {
		v := reflect.ValueOf(data).Elem() // pointer
		var rr = make([]string, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			var valueField = v.Field(i)
			var thisVal string
			switch valueField.Type().String() {
			case "int":
				thisVal = strconv.Itoa(int(valueField.Int()))
			case "int64":
				thisVal = strconv.Itoa(int(valueField.Int()))
			case "[]string":
				tmpVal, _ := valueField.Interface().([]string)
				thisVal = strings.Join(tmpVal, ";")
			case "time.Time":
				tmpVal, _ := valueField.Interface().(time.Time)
				thisVal = tmpVal.String()
			case "time.Month":
				tmpVal, _ := valueField.Interface().(time.Month)
				thisVal = tmpVal.String()
			default:
				thisVal = valueField.String()
			}
			rr[i] = thisVal
		}

		r[i] = rr
	}

	return r
}

var _ pumps.Pump = &csvPump{}
