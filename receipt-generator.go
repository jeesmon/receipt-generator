package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jeesmon/receipt-generator/num2words"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"gopkg.in/yaml.v2"
)

const (
	RECEIPT_YEAR         = "receiptYear"
	RECEIPT_DATE         = "receiptDate"
	RECEIPT_START_NUMBER = "receiptStartNumber"
	PAYMENTS_FILE        = "paymentsFile"
	PROJECTS_FILE        = "projectsFile"
	OUTPUT_FOLDER        = "outputFolder"

	ORG_NAME        = "orgName"
	ORG_ADDRESS     = "orgAddress"
	ORG_EIN         = "orgEIN"
	ORG_EMAIL       = "orgEmail"
	ORG_WEBSITE     = "orgWebsite"
	TREASURER_NAME  = "treasurerName"
	TREASURER_PHONE = "treasurerPhone"
	TREASURER_EMAIL = "treasurerEmail"
	ORG_LOGO        = "orgLogo"

	RECEIPT_TITLE     = "receiptTitle"
	ITEMS_TABLE_TITLE = "itemsTableTitle"
	TOTAL_TEXT        = "totalText"
	TABLE_COLUMNS     = "tableColumns"
	FOOTER1_TEXT      = "footer1Text"
	FOOTER2_TEXT      = "footer2Text"
	FOOTER3_TEXT      = "footer3Text"

	DEFAULT_RECEIPT_START_NUMBER = 100001
	DEFAULT_PAYMENTS_FILE        = "payments.csv"
	DEFAULT_PROJECTS_FILE        = "projects.csv"
	DEFAULT_OUTPUT_FOLDER        = "."
)

var (
	configFile string
	config     map[string]interface{}

	DEFAULT_RECEIPT_YEAR = time.Now().Format("2021")
	DEFAULT_RECEIPT_DATE = time.Now().Format("01/02/2006")
)

func main() {
	flag.StringVar(&configFile, "config", "config.yaml", "config file")
	flag.Parse()

	config = readConfig(configFile)

	pm := getProjects()
	py := getPayments()

	keys := make([]string, 0, len(py))
	for k := range py {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rn := config[RECEIPT_START_NUMBER].(int)
	for _, v := range keys {
		c, taxamt, totamt, name := getContents(pm, py[v])
		render(rn, c, taxamt, totamt, name)
		rn++
	}
}

func getPayments() map[string][][]string {
	f, err := os.Open(config[PAYMENTS_FILE].(string))
	defer f.Close()
	if err != nil {
		panic(err)
	}
	lines, err := csv.NewReader(f).ReadAll()

	m := make(map[string][][]string)

	for _, l := range lines {
		l[0] = strings.TrimSpace(l[0])
		if _, ok := m[l[0]]; ok {
			m[l[0]] = append(m[l[0]], l)
		} else {
			m[l[0]] = [][]string{l}
		}
	}

	return m
}

func getProjects() map[string]string {
	f, err := os.Open(config[PROJECTS_FILE].(string))
	defer f.Close()
	if err != nil {
		panic(err)
	}
	lines, err := csv.NewReader(f).ReadAll()

	m := make(map[string]string)

	for _, l := range lines {
		m[strings.TrimSpace(l[0])] = strings.TrimSpace(l[1])
	}

	return m
}

func render(rn int, contents [][]string, taxamt, totaamt float64, name string) {
	begin := time.Now()

	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	blackColor := getBlackColor()
	header := getHeader()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				val, _ := config[ORG_LOGO]
				if val != nil {
					_ = m.Base64Image(val.(string), consts.Png, props.Rect{
						Center:  true,
						Percent: 80,
					})
				}
			})

			m.Col(9, func() {
				val, _ := config[ORG_NAME]
				if val != nil {
					m.Text(val.(string), props.Text{
						Size:        16,
						Style:       consts.Bold,
						Align:       consts.Center,
						Extrapolate: false,
						Color:       blackColor,
					})
				}

				val, _ = config[ORG_ADDRESS]
				if val != nil {
					line := val.(string)
					ein, _ := config[ORG_EIN]
					if ein != nil {
						line = fmt.Sprintf("%s. EIN: %v", line, ein)
					}
					m.Text(line, props.Text{
						Top:   8,
						Style: consts.Bold,
						Size:  10,
						Align: consts.Center,
						Color: blackColor,
					})
				}
				val, _ = config[ORG_EMAIL]
				if val != nil {
					m.Text(fmt.Sprintf("Email: %s", val), props.Text{
						Top:   12,
						Style: consts.Bold,
						Size:  10,
						Align: consts.Center,
						Color: blackColor,
					})
				}
				val, _ = config[ORG_WEBSITE]
				if val != nil {
					m.Text(fmt.Sprintf("Website: %s", val), props.Text{
						Top:   16,
						Style: consts.Bold,
						Size:  10,
						Align: consts.Center,
						Color: blackColor,
					})
				}
			})
		})
	})

	m.RegisterFooter(func() {
		m.Row(30, func() {
			m.Col(12, func() {
				val, _ := config[TREASURER_NAME]
				if val != nil {
					m.Text(val.(string), props.Text{
						Top:   12,
						Style: consts.Bold,
						Size:  10,
						Align: consts.Left,
						Color: blackColor,
					})
					m.Text(fmt.Sprintf("Treasurer, %v", config[ORG_NAME]), props.Text{
						Top:   18,
						Style: consts.Bold,
						Size:  10,
						Align: consts.Left,
						Color: blackColor,
					})
					val, _ = config[TREASURER_PHONE]
					if val != nil {
						m.Text(fmt.Sprintf("Phone: %v", val), props.Text{
							Top:   22,
							Style: consts.Bold,
							Size:  10,
							Align: consts.Left,
							Color: blackColor,
						})
					}
					val, _ = config[TREASURER_EMAIL]
					if val != nil {
						m.Text(fmt.Sprintf("Email: %v", val), props.Text{
							Top:   26,
							Style: consts.Bold,
							Size:  10,
							Align: consts.Left,
							Color: blackColor,
						})
					}
				}
			})
		})

		m.Row(20, func() {
			m.Col(12, func() {
				val, _ := config[FOOTER1_TEXT]
				if val != nil {
					m.Text(val.(string), props.Text{
						Top:         10,
						Size:        8,
						Style:       consts.Italic,
						Align:       consts.Center,
						Extrapolate: false,
						Color:       blackColor,
					})
				}
				val, _ = config[FOOTER2_TEXT]
				if val != nil {
					m.Text(val.(string), props.Text{
						Top:         14,
						Size:        8,
						Style:       consts.Italic,
						Align:       consts.Center,
						Extrapolate: false,
						Color:       blackColor,
					})
				}
				val, _ = config[FOOTER3_TEXT]
				if val != nil {
					m.Text(val.(string), props.Text{
						Top:         18,
						Size:        8,
						Style:       consts.Italic,
						Align:       consts.Center,
						Extrapolate: false,
						Color:       blackColor,
					})
				}
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("%v", config[RECEIPT_TITLE]), props.Text{
				Top:         6,
				Size:        24,
				Style:       consts.Bold,
				Align:       consts.Center,
				Extrapolate: false,
				Color:       blackColor,
			})
		})
	})

	m.Row(15, func() {
		m.Col(3, func() {
			m.Text(fmt.Sprintf("Receipt #: %d", rn), props.Text{
				Top:         10,
				Size:        12,
				Style:       consts.Bold,
				Align:       consts.Left,
				Extrapolate: false,
				Color:       blackColor,
			})
		})

		m.ColSpace(6)

		m.Col(3, func() {
			m.Text(fmt.Sprintf("Date: %s", config[RECEIPT_DATE].(string)), props.Text{
				Top:         10,
				Size:        12,
				Style:       consts.Bold,
				Align:       consts.Right,
				Extrapolate: false,
				Color:       blackColor,
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			totaamtstr := fmt.Sprintf("%.2f", totaamt)
			fields := strings.Split(totaamtstr, ".")
			dollar, _ := strconv.Atoi(fields[0])
			cents, _ := strconv.Atoi(fields[1])

			words := num2words.ConvertNum2Words(dollar) + " Dollars"
			if cents > 0 {
				words += " and " + num2words.ConvertNum2Words(cents) + " Cents"
			}

			m.Text("Payment received from "+name+" in the amount of "+words+" for Sponsored Projects listed below.", props.Text{
				Top:         15,
				Size:        10,
				Style:       consts.Normal,
				Align:       consts.Left,
				Extrapolate: false,
				Color:       blackColor,
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("%v in %d", config[ITEMS_TABLE_TITLE], config[RECEIPT_YEAR]), props.Text{
				Top:         10,
				Size:        10,
				Style:       consts.Bold,
				Align:       consts.Center,
				Extrapolate: false,
				Color:       blackColor,
			})
		})
	})

	m.SetBackgroundColor(grayColor)

	columns, _ := config[TABLE_COLUMNS].([]interface{})

	m.Row(7, func() {
		m.Col(1, func() {
			m.Text(columns[0].(string), props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(columns[1].(string), props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(columns[2].(string), props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(5, func() {
			m.Text(columns[3].(string), props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(columns[4].(string), props.Text{
				Top:   1.5,
				Size:  9,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	m.SetBackgroundColor(whiteColor)

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      0.1,
			GridSizes: []uint{1, 2, 2, 5, 2},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{1, 2, 2, 5, 2},
		},
		Align:              consts.Left,
		HeaderContentSpace: 0,
		Line:               true,
	})

	m.TableList([]string{"", "", "", "", ""}, [][]string{
		{"", "", "", fmt.Sprintf("%v for %d", config[TOTAL_TEXT], config[RECEIPT_YEAR]), fmt.Sprintf("%10.02f", taxamt)}}, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      1,
			GridSizes: []uint{1, 2, 2, 5, 2},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{1, 2, 2, 5, 2},
			Style:     consts.Bold,
		},
		Align:              consts.Left,
		HeaderContentSpace: 1,
		Line:               true,
	})

	err := m.OutputFileAndClose(fmt.Sprintf("%s/%d-%s.pdf", config[OUTPUT_FOLDER].(string), rn, sanitize(name)))
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
}

func getHeader() []string {
	return []string{"", "", "", "", ""}
}

func getContents(pm map[string]string, v [][]string) ([][]string, float64, float64, string) {
	c := [][]string{}

	taxamt := 0.00
	totamt := 0.00
	n := ""
	snum := 0
	for _, r := range v {
		n = strings.TrimSpace(r[1])
		d := strings.TrimSpace(r[2])
		p := strings.TrimSpace(r[3])
		a := strings.TrimSpace(strings.ReplaceAll(r[4], ",", ""))
		taxded := strings.ToUpper(strings.TrimSpace(r[5])) == "Y"
		ai, _ := strconv.ParseFloat(a, 32)
		if taxded {
			taxamt += ai
		} else {
			continue
		}
		totamt += ai

		snum++
		c = append(c, []string{fmt.Sprintf("%d", snum), d, p, pm[p], fmt.Sprintf("%10.2f", ai)})
	}

	return c, taxamt, totamt, n
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlackColor() color.Color {
	return color.Color{
		Red:   0,
		Green: 0,
		Blue:  0,
	}
}

func sanitize(s string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(s, "")
}

func readConfig(config string) map[string]interface{} {
	data, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		panic(err)
	}

	val, ok := m[RECEIPT_YEAR]
	if val == nil || !ok {
		m[RECEIPT_YEAR] = DEFAULT_RECEIPT_YEAR
	}

	val, ok = m[RECEIPT_DATE]
	if val == nil || !ok {
		m[RECEIPT_DATE] = DEFAULT_RECEIPT_DATE
	}

	val, ok = m[RECEIPT_START_NUMBER]
	if val == nil || !ok {
		m[RECEIPT_START_NUMBER] = DEFAULT_RECEIPT_START_NUMBER
	}

	val, ok = m[PAYMENTS_FILE]
	if val == nil || !ok {
		m[PAYMENTS_FILE] = DEFAULT_PAYMENTS_FILE
	}

	val, ok = m[PROJECTS_FILE]
	if val == nil || !ok {
		m[PROJECTS_FILE] = DEFAULT_PROJECTS_FILE
	}

	val, ok = m[OUTPUT_FOLDER]
	if val == nil || !ok {
		m[OUTPUT_FOLDER] = DEFAULT_OUTPUT_FOLDER
	}

	return m
}
