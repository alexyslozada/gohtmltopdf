package gohtmltopdf

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type DIAN struct {
	isDebug bool
}

func NewDIAN(isDebug bool) DIAN {
	return DIAN{isDebug: isDebug}
}

func (d DIAN) CreateDIANForm220(data DIANForms220Relation) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to generate PDF")
	}

	var m core.Maroto
	switch data[0].Year {
	case 2022:
		m = d.dIAN2022(data)
	default:
		log.Printf("Year %d not supported", data[0].Year)
		return nil, ErrorProcess{Msg: fmt.Sprintf("year %d not supported", data[0].Year)}
	}

	document, err := m.Generate()
	if err != nil {
		log.Println("Error on generate PDF", err)
		return nil, err
	}

	if d.isDebug {
		// Save the metrics report
		err = document.GetReport().Save(fmt.Sprintf("./report-maroto-%s.txt", time.Now().Format("2006-01-02-15-04-05")))
		if err != nil {
			// We don't need to stop the process if the report can't be saved
			log.Println("Error on save metrics report", err)
		}
	}

	return document.GetBytes(), nil
}

// dIAN2022 Structure of the DIAN 220 form for the year 2022
func (d DIAN) dIAN2022(data DIANForms220Relation) core.Maroto {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.Letter).
		WithTopMargin(5).
		WithBottomMargin(5).
		WithMaxGridSize(28).
		WithDebug(d.isDebug).
		Build()

	mrt := maroto.New(cfg)
	if d.isDebug {
		// Add a metrics report to the maroto instance
		mrt = maroto.NewMetricsDecorator(mrt)
	}

	colorWhite := props.Color{Red: 255, Green: 255, Blue: 255}
	colorBlue := props.Color{Red: 65, Green: 95, Blue: 126}
	colorLightBlue := props.Color{Red: 242, Green: 245, Blue: 248}

	cellStyleFullBorder := props.Cell{
		BorderColor: &colorBlue,
		BorderType:  border.Full,
		LineStyle:   linestyle.Solid,
	}

	cellStyleLeftBorder := props.Cell{
		BorderColor: &colorBlue,
		BorderType:  border.Left,
		LineStyle:   linestyle.Solid,
	}

	cellStyleBgBlueLeftBorder := props.Cell{
		BackgroundColor: &colorBlue,
		BorderColor:     &colorBlue,
		BorderType:      border.Left,
		LineStyle:       linestyle.Solid,
	}

	cellStyleBgLigthBlueFullBorder := props.Cell{
		BackgroundColor: &colorLightBlue,
		BorderColor:     &colorBlue,
		BorderType:      border.Full,
		LineStyle:       linestyle.Solid,
	}

	cellStyleBgLightBlueLeftBorder := props.Cell{
		BackgroundColor: &colorLightBlue,
		BorderColor:     &colorBlue,
		BorderType:      border.Left,
		LineStyle:       linestyle.Solid,
	}

	textPropTitle := props.Text{Size: 8, Align: align.Center, Top: 2}
	textPropSubtitle := props.Text{Size: 8, Top: 7, Bottom: 4, Align: align.Center}
	textPropWarning := props.Text{Size: 6, Align: align.Center, Top: 3, Bottom: 2, Left: 8, Right: 8}
	textPropNumberForm := props.Text{Size: 8, Align: align.Center, Top: 4, Bottom: 1.5}
	textPropLabel := props.Text{Size: 5, Top: 1, Left: 1}
	textPropLabelSmall := props.Text{Size: 5, Top: 1, Left: 1}
	textPropLabelCenter := props.Text{Size: 5, Top: 1, Align: align.Center}
	textPropLabelTitleBgBlue := props.Text{Color: &colorWhite, Size: 5, Top: 1, Left: 1, Bottom: 1.5}
	textPropRetenedor := props.Text{Size: 5, Top: 4, Bottom: 1.5, Left: 1}
	textPropRetenedorCenter := props.Text{Size: 5, Top: 4, Bottom: 1.5, Align: align.Center}
	textPropLabelConcept := props.Text{Size: 5, Top: 0.9, Left: 1, Bottom: 0.9}
	textPropLabelConceptCenter := props.Text{Size: 5, Top: 0.9, Bottom: 1, Align: align.Center}
	textPropLabelConceptRight := props.Text{Size: 5, Top: 0.9, Right: 1, Bottom: 1, Align: align.Right}
	textPropLabelConceptTitleBgLightBlue := props.Text{Size: 5, Top: 1, Bottom: 1.5, Align: align.Center, Style: fontstyle.Bold}
	textPropLabelConceptTitle := props.Text{Size: 5, Top: 1, Bottom: 1.5, Align: align.Center, Style: fontstyle.Bold}
	textPropDependent := props.Text{Size: 5, Top: 4, Left: 1, Bottom: 0.9}
	textPropLabelBigBox := props.Text{Size: 5, Top: 1, Left: 1, Bottom: 1.5, Align: align.Left, VerticalPadding: 1.25}
	textPropLabelDisclaimer := props.Text{Size: 5, Top: 0.6, Left: 1, Bottom: 0.6}

	printerSpanish := message.NewPrinter(language.Spanish)
	for _, item := range data {
		pageData := page.New().Add(
			row.New().Add(
				image.NewFromFileCol(6, "./logo_dian.png", props.Rect{Top: 1, Left: 1, Percent: 95}).WithStyle(&cellStyleFullBorder),
				col.New(16).Add(
					text.New("Certificado de Ingresos y Retenciones por Rentas de Trabajo y de Pensiones", textPropTitle),
					text.New("Año gravable 2022", textPropSubtitle),
				).WithStyle(&cellStyleFullBorder),
				image.NewFromFileCol(6, "./form_220.png", props.Rect{Top: 1, Left: 1, Percent: 95}).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				text.NewCol(14, "Antes de diligenciar este formulario lea cuidadosamente las instrucciones", textPropWarning).WithStyle(&cellStyleFullBorder),
				col.New(14).Add(
					text.New("4. Número de formulario", textPropLabel),
					text.New(strconv.Itoa(int(item.Sequence)), textPropNumberForm),
				).WithStyle(&cellStyleFullBorder),
			),
			// *********************
			// Retenedor
			// *********************
			row.New().Add(
				col.New(1).WithStyle(&cellStyleLeftBorder),
				col.New(10).Add(
					text.New("5. Número de identificación tributaria (NIT)", textPropLabel),
					text.New(item.Nit, textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(1).Add(
					text.New("6. DV", textPropLabel),
					text.New(item.Dv, textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("7. Primer apellido", textPropLabel),
					text.New("", textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("8. Segundo apellido", textPropLabel),
					text.New("", textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("9. Primer nombre", textPropLabel),
					text.New("", textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("10. Otros nombres", textPropLabel),
					text.New("", textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				col.New(1).WithStyle(&cellStyleLeftBorder),
				col.New(27).Add(
					text.New("11. Razón social", textPropLabel),
					text.New(strings.ToUpper(item.BusinessName), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				col.New(1).WithStyle(&cellStyleFullBorder),
				col.New(3).Add(
					text.New("24. Tipo documento", textPropLabel),
					text.New(strconv.Itoa(int(item.IdentificationTypeCode)), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("25. Número Identificación", textPropLabel),
					text.New(item.IdentificationNumber, textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(5).Add(
					text.New("26. Primer apellido", textPropLabel),
					text.New(strings.ToUpper(item.LastName), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(5).Add(
					text.New("27. Segundo apellido", textPropLabel),
					text.New(strings.ToUpper(item.Surname), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(5).Add(
					text.New("28. Primer nombre", textPropLabel),
					text.New(strings.ToUpper(item.FirstName), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(5).Add(
					text.New("29. Otros nombres", textPropLabel),
					text.New(strings.ToUpper(item.MiddleName), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				col.New(8).Add(
					text.New("Periodo de la certificación", textPropLabelCenter),
					text.New(fmt.Sprintf("30. DE: %s    31. A: %s", item.BeginsAt.Format(time.DateOnly), item.EndsAt.Format(time.DateOnly)), textPropRetenedorCenter),
				).WithStyle(&cellStyleFullBorder),
				col.New(5).Add(
					text.New("32. Fecha de expedición", textPropLabelCenter),
					text.New("2023-03-31", textPropRetenedorCenter),
				).WithStyle(&cellStyleFullBorder),
				col.New(10).Add(
					text.New("33. Lugar donde se practicó la retención", textPropLabel),
					text.New(strings.ToUpper(item.Place), textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(2).Add(
					text.New("34. Cód. Dpto.", textPropLabelSmall),
					text.New(item.DepartmentCode, textPropRetenedor),
				).WithStyle(&cellStyleFullBorder),
				col.New(3).Add(
					text.New("35. Cód. Ciudad/Municipio", textPropLabelSmall),
					text.New(item.MunicipalityCode, textPropRetenedorCenter),
				).WithStyle(&cellStyleFullBorder),
			),
			// *********************
			// Sección de ingresos
			// *********************
			row.New().Add(
				text.NewCol(20, "Concepto de los ingresos", textPropLabelConceptTitleBgLightBlue).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(8, "Valor", textPropLabelConceptTitleBgLightBlue).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(20, "Pagos por salarios o emolumentos eclesiásticos", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "36", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["36"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(20, "Pagos realizados con bonos electrónicos o de papel de servicio, cheques, tarjetas, vales, etc.", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "36", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["37"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(20, "Pagos por honorarios", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "38", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["38"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por servicios - 39
			row.New().Add(
				text.NewCol(20, "Pagos por servicios", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "39", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["39"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por comisiones - 40
			row.New().Add(
				text.NewCol(20, "Pagos por comisiones", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "40", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["40"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por prestaciones sociales - 41
			row.New().Add(
				text.NewCol(20, "Pagos por prestaciones sociales", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "41", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["41"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por viáticos - 42
			row.New().Add(
				text.NewCol(20, "Pagos por viáticos", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "42", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["42"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por gastos de representación - 43
			row.New().Add(
				text.NewCol(20, "Pagos por gastos de representación", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "43", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["43"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pagos por compensaciones por el trabajo asociado cooperativo - 44
			row.New().Add(
				text.NewCol(20, "Pagos por compensaciones por el trabajo asociado cooperativo", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "44", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["44"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Otros pagos - 45
			row.New().Add(
				text.NewCol(20, "Otros pagos", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "45", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["45"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Cesantías e intereses de cesantías efectivamente pagadas al empleado - 46
			row.New().Add(
				text.NewCol(20, "Cesantías e intereses de cesantías efectivamente pagadas al empleado", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "46", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["46"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Cesantías consignadas al fondo de cesantias - 47
			row.New().Add(
				text.NewCol(20, "Cesantías consignadas al fondo de cesantias", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "47", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["47"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Pensiones de jubilación, vejez o invalidez - 48
			row.New().Add(
				text.NewCol(20, "Pensiones de jubilación, vejez o invalidez", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "48", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["48"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Total de ingresos brutos (Sume 36 a 48) - 49
			row.New().Add(
				text.NewCol(20, "Total de ingresos brutos (Sume 36 a 48)", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "49", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["49"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// *********************
			// Sección Concepto de los aportes
			// *********************
			row.New().Add(
				text.NewCol(20, "Concepto de los aportes", textPropLabelConceptTitle).WithStyle(&cellStyleFullBorder),
				text.NewCol(8, "Valor", textPropLabelConceptTitle).WithStyle(&cellStyleFullBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Aportes obligatorios por salud a cargo del trabajador - 50
			row.New().Add(
				text.NewCol(20, "Aportes obligatorios por salud a cargo del trabajador", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "50", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["50"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Aportes obligatorios a fondos de pensiones y solidaridad pensional a cargo del trabajador - 51
			row.New().Add(
				text.NewCol(20, "Aportes obligatorios a fondos de pensiones y solidaridad pensional a cargo del trabajador", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "51", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["51"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Cotizaciones voluntarias al régimen de ahorro individual con solidaridad - RAIS - 52
			row.New().Add(
				text.NewCol(20, "Cotizaciones voluntarias al régimen de ahorro individual con solidaridad - RAIS", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "52", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["52"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Aportes voluntarios a fondos de pensiones - 53
			row.New().Add(
				text.NewCol(20, "Aportes voluntarios a fondos de pensiones", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "53", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["53"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Aportes a cuentas AFC o AVC - 54
			row.New().Add(
				text.NewCol(20, "Aportes a cuentas AFC o AVC", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "54", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["54"]), textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Valor de la retención en la fuente por ingresos laborales y de pensiones - 55
			row.New().Add(
				text.NewCol(20, "Valor de la retención en la fuente por ingresos laborales y de pensiones", textPropLabelTitleBgBlue).WithStyle(&cellStyleBgBlueLeftBorder),
				text.NewCol(1, "55", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, printerSpanish.Sprintf("%.0f", item.RowsMap["55"]), textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Datos del pagador
			row.New().Add(
				col.New(28).Add(
					text.New("Nombre del pagador o agente retenedor: Este documento no requiere para su validez firma autógrafa de acuerdo con el artículo 10 del Decreto 836 de 1991, recopilado en el artículo 1.6.1.12.12 del DUT 1625 de octubre 11 de 2016, que regula el contenido del certificado de retenciones a título de renta.", textPropLabelBigBox),
					text.New(strings.ToUpper(item.BusinessName), props.Text{Size: 6, Top: 8, Align: align.Center}),
					text.New(fmt.Sprintf("NIT: %s - %s", item.IdentificationNumber, item.Dv), props.Text{Size: 6, Top: 11, Bottom: 1.5, Align: align.Center}),
				).WithStyle(&cellStyleBgLigthBlueFullBorder),
			),
			// Datos a cargo del trabajador o pensionado
			row.New().Add(
				text.NewCol(28, "Datos a cargo del trabajador o pensionado", textPropLabelConceptCenter).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				text.NewCol(14, "Concepto de otros ingresos", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
				text.NewCol(7, "Valor recibido", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
				text.NewCol(7, "Valor Retenido", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
			),
			row.New().Add(
				text.NewCol(14, "Arrendamientos", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "56", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "63", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Honorarios, comisiones y servicios - 57 - 64
			row.New().Add(
				text.NewCol(14, "Honorarios, comisiones y servicios", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "57", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "64", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Intereses y rendimientos financieros - 58 - 65
			row.New().Add(
				text.NewCol(14, "Intereses y rendimientos financieros", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "58", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "65", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Enajenación de activos fijos - 59 - 66
			row.New().Add(
				text.NewCol(14, "Enajenación de activos fijos", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "59", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "66", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Loterías, rifas, apuestas y similares - 60 - 67
			row.New().Add(
				text.NewCol(14, "Loterías, rifas, apuestas y similares", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "60", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "67", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Otros 61 - 68
			row.New().Add(
				text.NewCol(14, "Otros", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "61", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "68", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Totales: (Valor recibido: Sume 56 a 61), (Valor retenido: Sume 63 a 68) - 62 - 69
			row.New().Add(
				text.NewCol(14, "Totales: (Valor recibido: Sume 56 a 61), (Valor retenido: Sume 63 a 68)", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "62", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(1, "69", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Total retenciones año gravable 2022 (Sume 55 + 69) - 70
			row.New().Add(
				text.NewCol(21, "Total retenciones año gravable 2022 (Sume 55 + 69)", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(1, "70", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(6, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			// Cuadro de identificación de los bienes poseidos
			row.New().Add(
				text.NewCol(1, "Item", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
				text.NewCol(20, "Identificación de los bienes poseídos", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
				text.NewCol(7, "72. Valor patrimonial", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
			),
			// 6 registros
			row.New().Add(
				text.NewCol(1, "1", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(1, "2", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(1, "3", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(1, "4", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(1, "5", textPropLabelConceptCenter).WithStyle(&cellStyleLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(1, "6", textPropLabelConceptCenter).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(20, "", textPropLabelConcept).WithStyle(&cellStyleBgLightBlueLeftBorder),
				text.NewCol(7, "", textPropLabelConceptRight).WithStyle(&cellStyleBgLightBlueLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "Deudas vigentes a 31 de diciembre de 2022", textPropLabelTitleBgBlue).WithStyle(&cellStyleBgBlueLeftBorder),
				text.NewCol(1, "73", textPropLabelConceptCenter).WithStyle(&cellStyleFullBorder),
				text.NewCol(6, "", textPropLabelConceptRight).WithStyle(&cellStyleFullBorder),
			),
			text.NewRow(0, "Identificación del dependiente económico de acuerdo al parágrafo 2 del artículo 387 del Estatuto Tributario", textPropLabelConceptCenter).WithStyle(&cellStyleBgLigthBlueFullBorder),
			row.New().Add(
				col.New(4).Add(
					text.New("74. Tipo documento", textPropLabelConceptCenter),
					text.New("", textPropDependent),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("75. No. Documento", textPropLabelConceptCenter),
					text.New("", textPropDependent),
				).WithStyle(&cellStyleFullBorder),
				col.New(16).Add(
					text.New("76. Apellidos y Nombres", textPropLabelConceptCenter),
					text.New("", textPropDependent),
				).WithStyle(&cellStyleFullBorder),
				col.New(4).Add(
					text.New("77. Parentesco", textPropLabelConceptCenter),
					text.New("", textPropDependent),
				).WithStyle(&cellStyleFullBorder),
			),
			row.New().Add(
				text.NewCol(21, "Certifico que durante el año gravable 2022:", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "Firma del Trabajador o Pensionado", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "1. Mi patrimonio bruto no excedió de 4.500 UVT ($171.018.000).", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "2. Mis ingresos brutos fueron inferiores a 1.400 UVT ($53.206.000).", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "3. No fui responsable del impuesto sobre las ventas a 31 de diciembre de 2022.", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "4. Mis consumos mediante tarjeta de crédito no excedieron la suma de 1.400 UVT ($53.206.000).", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "5. Que el total de mis compras y consumos no superaron la suma de 1.400 UVT ($53.206.000).", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "6. Que el valor total de mis consignaciones bancarias, depósitos o inversiones financieras no excedieron los 1.400 UVT ($53.206.000).", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				text.NewCol(21, "Por lo tanto, manifiesto que no estoy obligado a presentar declaración de renta y complementario por el año gravable 2022.", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				text.NewCol(7, "", textPropLabelDisclaimer).WithStyle(&cellStyleLeftBorder),
				col.New(0).WithStyle(&cellStyleLeftBorder),
			),
			row.New().Add(
				col.New(28).Add(
					text.New("Nota: este certificado sustituye para todos los efectos legales la declaración de Renta y Complementario para el trabajador o pensionado que lo firme.", textPropLabelDisclaimer),
					text.New("Para aquellos trabajadores independientes contribuyentes del impuesto unificado deberán presentar la declaración anual consolidada del Régimen Simple de Tributación (SIMPLE).", props.Text{Size: 5, Top: 3.5, Left: 1, Bottom: 1}),
				).WithStyle(&cellStyleFullBorder),
			),
		)

		mrt.AddPages(pageData)
	}

	return mrt
}
