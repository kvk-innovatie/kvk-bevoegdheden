package bevoegdheden

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
)

var ErrInvalidInput = errors.New("KVK nummer of persoonsgegevens kloppen niet")

var importanceTypeFunctionaris = map[string]int64{
	"Aansprakelijke":                       5,
	"Bestuursfunctie":                      5,
	"FunctionarisBijzondereRechtstoestand": 4,
	"PubliekrechtelijkeFunctionaris":       3,
	"Gemachtigde":                          2,
	"OverigeFunctionaris":                  1,
}

func convertDate(str string) string {
	if str == "" {
		return str
	}
	return str[6:] + "-" + str[4:6] + "-" + str[0:4]
}

func isSamePerson(p1 models.IdentityNP, p2 models.IdentityNP) bool {
	if p1.Voornamen == "" || p2.Voornamen == "" || p1.Geslachtsnaam == "" || p2.Geslachtsnaam == "" || p1.Geboortedatum == "" || p2.Geboortedatum == "" {
		return false
	}

	re := regexp.MustCompile(`\s*-\s*`)
	g1 := re.ReplaceAllString(p1.Geslachtsnaam, "-")
	g2 := re.ReplaceAllString(p2.Geslachtsnaam, "-")

	// surname from IRMA personal data card has prefix included. To compare add prefix to geslachtsnaam from KVK
	if p2.VoorvoegselGeslachtsnaam != "" {
		g2 = p2.VoorvoegselGeslachtsnaam + " " + g2
	}

	return (g1 == g2 && p1.Voornamen == p2.Voornamen && p1.Geboortedatum == p2.Geboortedatum)
}

func getFunctionaris(functionarisOfGemachtigde *models.FunctionarisOfGemachtigde, functionarisType string) models.Functionaris {
	bgd := functionarisOfGemachtigde.Bevoegdheid
	bvm := functionarisOfGemachtigde.Volmacht.BeperkteVolmacht
	beperkingInEuros := ""
	if bgd.BeperkingInEuros.Waarde != "" {
		beperkingInEuros = bgd.BeperkingInEuros.Waarde + " " + bgd.BeperkingInEuros.Valuta.Omschrijving
	}
	beperkingInGeld := ""
	if bvm.BeperkingInGeld.Waarde != "" {
		beperkingInGeld = bvm.BeperkingInGeld.Waarde + " " + bvm.BeperkingInGeld.Valuta.Omschrijving
	}
	soortBeperkingInHandeling := ""

	for i, bih := range bvm.BeperkingInHandeling {
		prefix := ", "
		if i == 0 {
			prefix = ""
		}
		if bih.SoortHandeling.Code != "" {
			soortBeperkingInHandeling = soortBeperkingInHandeling + prefix + bih.SoortHandeling.Omschrijving + ":" + bih.SoortHandeling.Code
		}
	}

	return models.Functionaris{
		TypeFunctionaris: functionarisType,
		Functie:          functionarisOfGemachtigde.Functie.Omschrijving,
		Functietitel:     functionarisOfGemachtigde.Functietitel.Titel,
		SchorsingAanvang: convertDate(functionarisOfGemachtigde.Schorsing.Registratie.DatumAanvang),
		SchorsingEinde:   convertDate(functionarisOfGemachtigde.Schorsing.Registratie.DatumEinde),
		Handlichting:     functionarisOfGemachtigde.Handlichting.IsVerleend.Code,

		SoortBevoegdheid:            bgd.Soort.Omschrijving,
		BeperkingInEurosBevoegdheid: beperkingInEuros,
		OverigeBeperkingBevoegdheid: bgd.OverigeBeperking.Omschrijving,
		IsBevoegdMetAnderePersonen:  bgd.IsBevoegdMetAnderePersonen.Omschrijving,

		TypeVolmacht:                 functionarisOfGemachtigde.Volmacht.TypeVolmacht.Omschrijving,
		BeperkingInGeldVolmacht:      beperkingInGeld,
		BeperkingInHandelingVolmacht: soortBeperkingInHandeling,
		HeeftOverigeVolmacht:         bvm.HeeftOverigeVolmacht.Omschrijving,
		OmschrijvingOverigeVolmacht:  bvm.OmschrijvingOverigeVolmacht,
		MagOpgaveHandelsregisterDoen: bvm.MagOpgaveHandelsregisterDoen.Omschrijving,

		Importance: importanceTypeFunctionaris[functionarisType],
	}
}

func getNatuurlijkePersoonFunctionaris(functionarisOfGemachtigde *models.FunctionarisOfGemachtigde, functionarisType string) models.NatuurlijkPersoonFunctionaris {
	np := functionarisOfGemachtigde.Door.NatuurlijkPersoon
	bijzondereRechtstoestand := ""
	if np.BijzondereRechtstoestand.Soort.Code != "" {
		bijzondereRechtstoestand = np.BijzondereRechtstoestand.Soort.Omschrijving + ":" + np.BijzondereRechtstoestand.Soort.Code
	}
	beperkingInRechtshandeling := ""
	if np.BeperkingInRechtshandeling.Soort.Code != "" {
		beperkingInRechtshandeling = np.BeperkingInRechtshandeling.Soort.Omschrijving + ":" + np.BeperkingInRechtshandeling.Soort.Code
	}
	functionaris := getFunctionaris(functionarisOfGemachtigde, functionarisType)
	npFunctionaris := models.NatuurlijkPersoonFunctionaris{
		Geslachtsnaam:              np.Geslachtsnaam,
		VoorvoegselGeslachtsnaam:   np.VoorvoegselGeslachtsnaam,
		Voornamen:                  np.Voornamen,
		Geboortedatum:              convertDate(np.Geboortedatum),
		Overlijdensdatum:           convertDate(np.Overlijdensdatum),
		VolledigeNaam:              np.VolledigeNaam,
		BijzondereRechtstoestand:   bijzondereRechtstoestand,
		BeperkingInRechtshandeling: beperkingInRechtshandeling,
		Functionaris:               functionaris,
	}

	return npFunctionaris
}

func getRechtspersoonFunctionaris(functionarisOfGemachtigde *models.FunctionarisOfGemachtigde, functionarisType string) models.RechtspersoonFunctionaris {
	rp := functionarisOfGemachtigde.Door.Rechtspersoon
	functionaris := getFunctionaris(functionarisOfGemachtigde, functionarisType)
	rpFunctionaris := models.RechtspersoonFunctionaris{
		KvkNummer:         rp.IsEigenaarVan.MaatschappelijkeActiviteit.KvkNummer,
		PersoonRechtsvorm: rp.PersoonRechtsvorm,
		Naam:              rp.VolledigeNaam,
		Functionaris:      functionaris,
	}

	return rpFunctionaris
}

func getFunctionarisPaths(basePath string) models.NatuurlijkPersoonFunctionaris {
	return models.NatuurlijkPersoonFunctionaris{
		Geslachtsnaam:              basePath + ".door.natuurlijkPersoon.geslachtsnaam",
		VoorvoegselGeslachtsnaam:   basePath + ".door.natuurlijkPersoon.voorvoegselGeslachtsnaam",
		Voornamen:                  basePath + ".door.natuurlijkPersoon.voornamen",
		Geboortedatum:              basePath + ".door.natuurlijkPersoon.geboortedatum",
		Overlijdensdatum:           basePath + ".door.natuurlijkPersoon.overlijdensdatum",
		VolledigeNaam:              basePath + ".door.natuurlijkPersoon.volledigeNaam",
		BijzondereRechtstoestand:   basePath + ".door.natuurlijkPersoon.bijzondereRechtstoestand.soort",
		BeperkingInRechtshandeling: basePath + ".door.natuurlijkPersoon.beperkingInRechtshandeling.soort",

		Functionaris: models.Functionaris{
			TypeFunctionaris: basePath,
			Functie:          basePath + ".functie.omschrijving",
			Functietitel:     basePath + ".functietitel.titel",
			SchorsingAanvang: basePath + ".schorsing.registratie.datumAanvang",
			SchorsingEinde:   basePath + ".schorsing.registratie.datumEinde",
			Handlichting:     basePath + ".handlichting.isVerleend.code",

			SoortBevoegdheid:            basePath + ".bevoegdheid.soort.omschrijving",
			BeperkingInEurosBevoegdheid: basePath + ".bevoegdheid.beperkingInEuros",
			OverigeBeperkingBevoegdheid: basePath + ".bevoegdheid.overigeBeperking.omschrijving",
			IsBevoegdMetAnderePersonen:  basePath + ".bevoegdheid.isBevoegdMetAnderePersonen.omschrijving",

			TypeVolmacht:                 basePath + ".volmacht.typeVolmacht.omschrijving",
			BeperkingInGeldVolmacht:      basePath + ".volmacht.beperkteVolmacht.beperkingInGeld",
			BeperkingInHandelingVolmacht: basePath + ".volmacht.beperkteVolmacht.beperkingInHandeling",
			HeeftOverigeVolmacht:         basePath + ".volmacht.beperkteVolmacht.heeftOverigeVolmacht.omschrijving",
			OmschrijvingOverigeVolmacht:  basePath + ".volmacht.beperkteVolmacht.omschrijvingOverigeVolmacht",
			MagOpgaveHandelsregisterDoen: basePath + ".volmacht.beperkteVolmacht.magOpgaveHandelsregisterDoen.omschrijving",
		},
	}
}

func addInterpretatieNP(bevoegdheidUittreksel *models.BevoegdheidUittreksel, functionaris *models.NatuurlijkPersoonFunctionaris) {
	namePerson := fmt.Sprintf("%s %s %s", functionaris.Voornamen, functionaris.VoorvoegselGeslachtsnaam, functionaris.Geslachtsnaam)
	interpretatie := &functionaris.Interpretatie
	interpretatie.HeeftBeperking = "Nee"

	addInterpretatie(bevoegdheidUittreksel, &functionaris.Functionaris, namePerson)

	if functionaris.Overlijdensdatum != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s staat geregistreerd als overleden op %s", namePerson, functionaris.Overlijdensdatum)
		return
	}

	if functionaris.BijzondereRechtstoestand != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s heeft een bijzondere rechtstoestand: %s", namePerson, functionaris.BijzondereRechtstoestand)
		return
	}

	if functionaris.BeperkingInRechtshandeling != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De persoon %s heeft een beperking in rechtshandeling: %s", namePerson, functionaris.BeperkingInRechtshandeling)
		return
	}

	if isMinderjarig(functionaris.Geboortedatum) {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"

		if bevoegdheidUittreksel.PersoonRechtsvorm == "Eenmanszaak" && functionaris.TypeFunctionaris == "Eigenaar" {
			if functionaris.Handlichting == "" {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is een minderjarige eigenaar eenmanszaak zonder handlichting bij inschrijving %s alleen bevoegd met schriftelijke toestemming van een wettelijke vertegenwoordiger.", namePerson, bevoegdheidUittreksel.KvkNummer)
			} else {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is een minderjarige eigenaar eenmanszaak met handlichting bij inschrijving %s. Raadpleeg het Handelsregister voor meer informatie.", namePerson, bevoegdheidUittreksel.KvkNummer)
			}
		} else {
			if functionaris.Handlichting == "" {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is minderjarig zonder handlichting bij inschrijving %s alleen bevoegd met schriftelijke toestemming van een wettelijke vertegenwoordiger.", namePerson, bevoegdheidUittreksel.KvkNummer)
			} else {
				interpretatie.Reden = fmt.Sprintf("De persoon %s is minderjarig met handlichting bij inschrijving %s. Raadpleeg het Handelsregister voor meer informatie.", namePerson, bevoegdheidUittreksel.KvkNummer)
			}
		}
		return
	}

	if bevoegdheidUittreksel.PersoonRechtsvorm == "Eenmanszaak" && functionaris.TypeFunctionaris == "Eigenaar" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("De persoon %s is eigenaar van een eenmanszaak", namePerson)
		return
	}
}

func addInterpretatie(bevoegdheidUittreksel *models.BevoegdheidUittreksel, functionaris *models.Functionaris, namePerson string) {
	interpretatie := &functionaris.Interpretatie
	interpretatie.HeeftBeperking = "Nee"

	if bevoegdheidUittreksel.DatumUitschrijving != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s staat niet meer ingeschreven: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.DatumUitschrijving)
		return
	}

	if bevoegdheidUittreksel.RegistratieEinde != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s is niet (meer) actief: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.RegistratieEinde)
		return
	}

	if bevoegdheidUittreksel.BijzondereRechtstoestand != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een bijzondere rechtstoestand: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BijzondereRechtstoestand)
		return
	}

	if bevoegdheidUittreksel.BeperkingInRechtshandeling != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een beperking in rechtshandeling: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BeperkingInRechtshandeling)
		return
	}

	if bevoegdheidUittreksel.BuitenlandseRechtstoestand != "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("De inschrijving %s heeft een buitenlandse rechtstoestand: %s", bevoegdheidUittreksel.KvkNummer, bevoegdheidUittreksel.BuitenlandseRechtstoestand)
		return
	}

	if functionaris.SchorsingAanvang != "" && functionaris.SchorsingEinde == "" {
		interpretatie.HeeftBeperking = "Ja"
		interpretatie.IsBevoegd = "Nee"
		interpretatie.Reden = fmt.Sprintf("%s is geschorst sinds: %s", namePerson, functionaris.SchorsingAanvang)
		return
	}

	if functionaris.SoortBevoegdheid == "Alleen/zelfstandig bevoegd" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, functionaris.Functie, functionaris.SoortBevoegdheid)
		return
	}
	if functionaris.SoortBevoegdheid == "Gezamenlijk bevoegd" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, functionaris.Functie, functionaris.SoortBevoegdheid)
		return
	}

	if functionaris.SoortBevoegdheid == "Onbeperkt bevoegd" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, functionaris.Functie, functionaris.SoortBevoegdheid)
		return
	}
	if functionaris.SoortBevoegdheid == "Beperkt bevoegd" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) is %s", namePerson, functionaris.Functie, functionaris.SoortBevoegdheid)
		return
	}

	if functionaris.TypeVolmacht == "Volledige volmacht" {
		interpretatie.IsBevoegd = "Ja"
		interpretatie.Reden = fmt.Sprintf("%s (%s) heeft %s", namePerson, functionaris.Functie, functionaris.TypeVolmacht)
		return
	}
	if functionaris.TypeVolmacht == "Beperkte volmacht" {
		interpretatie.IsBevoegd = "Niet vastgesteld"
		interpretatie.Reden = fmt.Sprintf("%s (%s) heeft %s", namePerson, functionaris.Functie, functionaris.TypeVolmacht)
		return
	}

	interpretatie.IsBevoegd = "Niet vastgesteld"
	interpretatie.Reden = "Geen reden gevonden"
}

func loopFunctionarissen(bevoegdheidUittreksel *models.BevoegdheidUittreksel, paths *models.Paths, identityNP models.IdentityNP, functievervullingen []models.Functievervulling, basePath string) {
	for i, heeft := range functievervullingen {
		var functionarisOfGemachtigde *models.FunctionarisOfGemachtigde
		path := basePath + ".heeft." + strconv.Itoa(i) + "."
		functionarisType := ""

		if heeft.Bestuursfunctie != nil {
			functionarisOfGemachtigde = heeft.Bestuursfunctie
			path = path + "bestuursfunctie"
			functionarisType = "Bestuursfunctie"
		} else if heeft.Aansprakelijke != nil {
			functionarisOfGemachtigde = heeft.Aansprakelijke
			path = path + "aansprakelijke"
			functionarisType = "Aansprakelijke"
		} else if heeft.FunctionarisBijzondereRechtstoestand != nil {
			functionarisOfGemachtigde = heeft.FunctionarisBijzondereRechtstoestand
			path = path + "functionarisBijzondereRechtstoestand"
			functionarisType = "FunctionarisBijzondereRechtstoestand"
		} else if heeft.OverigeFunctionaris != nil {
			functionarisOfGemachtigde = heeft.OverigeFunctionaris
			path = path + "overigeFunctionaris"
			functionarisType = "OverigeFunctionaris"
		} else if heeft.PubliekrechtelijkeFunctionaris != nil {
			functionarisOfGemachtigde = heeft.PubliekrechtelijkeFunctionaris
			path = path + "publiekrechtelijkeFunctionaris"
			functionarisType = "PubliekrechtelijkeFunctionaris"
		} else if heeft.Gemachtigde != nil {
			functionarisOfGemachtigde = heeft.Gemachtigde
			path = path + "gemachtigde"
			functionarisType = "Gemachtigde"
		}

		np := functionarisOfGemachtigde.Door.NatuurlijkPersoon
		if np == nil {
			rp := functionarisOfGemachtigde.Door.Rechtspersoon
			if rp == nil {
				continue
			}
			rechtsPersoonFunctionaris := getRechtspersoonFunctionaris(functionarisOfGemachtigde, functionarisType)
			addInterpretatie(bevoegdheidUittreksel, &rechtsPersoonFunctionaris.Functionaris, rechtsPersoonFunctionaris.Naam)
			bevoegdheidUittreksel.AlleRechtspersoonFunctionarissen = append(bevoegdheidUittreksel.AlleRechtspersoonFunctionarissen, rechtsPersoonFunctionaris)
			continue
		}

		npFunctionaris := getNatuurlijkePersoonFunctionaris(functionarisOfGemachtigde, functionarisType)
		addInterpretatieNP(bevoegdheidUittreksel, &npFunctionaris)
		bevoegdheidUittreksel.AlleFunctionarissen = append(bevoegdheidUittreksel.AlleFunctionarissen, npFunctionaris)

		identityFunctionaris := models.IdentityNP{
			Geslachtsnaam:            np.Geslachtsnaam,
			Voornamen:                np.Voornamen,
			VoorvoegselGeslachtsnaam: np.VoorvoegselGeslachtsnaam,
			Geboortedatum:            convertDate(np.Geboortedatum),
		}
		if isSamePerson(identityNP, identityFunctionaris) {
			if bevoegdheidUittreksel.MatchedFunctionaris == nil || bevoegdheidUittreksel.MatchedFunctionaris.Importance < npFunctionaris.Importance {
				bevoegdheidUittreksel.MatchedFunctionaris = &npFunctionaris
				paths.MatchedFunctionaris = getFunctionarisPaths(path)
			}
		}
	}
}

func eigenaarIsNatuurlijkPersoon(bevoegdheidUittreksel *models.BevoegdheidUittreksel, paths *models.Paths, identityNP models.IdentityNP, eenmanszaak *models.Eenmanszaak, typeEigenaar string, basePath string) {
	bevoegdheidUittreksel.TypeEigenaar = typeEigenaar
	bevoegdheidUittreksel.PersoonRechtsvorm = eenmanszaak.PersoonRechtsvorm
	bevoegdheidUittreksel.RegistratieEinde = convertDate(eenmanszaak.Registratie.DatumEinde)
	if eenmanszaak.BijzondereRechtstoestand.Soort.Code != "" {
		bevoegdheidUittreksel.BijzondereRechtstoestand = eenmanszaak.BijzondereRechtstoestand.Soort.Omschrijving + ":" + eenmanszaak.BijzondereRechtstoestand.Soort.Code
	}
	if eenmanszaak.BeperkingInRechtshandeling.Soort.Code != "" {
		bevoegdheidUittreksel.BeperkingInRechtshandeling = eenmanszaak.BeperkingInRechtshandeling.Soort.Omschrijving + ":" + eenmanszaak.BeperkingInRechtshandeling.Soort.Code
	}
	// bevoegdheidUittreksel.Handlichting = eenmanszaak.Handlichting.IsVerleend.Code

	paths.TypeEigenaar = basePath
	paths.PersoonRechtsvorm = basePath + ".persoonRechtsvorm"
	paths.RegistratieEinde = basePath + ".registratie.datumEinde"
	paths.BijzondereRechtstoestand = basePath + ".bijzondereRechtstoestand.soort"
	paths.BeperkingInRechtshandeling = basePath + ".beperkingInRechtshandeling.soort"
	// bevoegdheidUittreksel.Paths.Handlichting = basePath + ".handlichting.isVerleend.code"

	functionaris := models.NatuurlijkPersoonFunctionaris{
		Geslachtsnaam:            eenmanszaak.Geslachtsnaam,
		VoorvoegselGeslachtsnaam: eenmanszaak.VoorvoegselGeslachtsnaam,
		Voornamen:                eenmanszaak.Voornamen,
		Geboortedatum:            convertDate(eenmanszaak.Geboortedatum),
		Overlijdensdatum:         convertDate(eenmanszaak.Overlijdensdatum),
		VolledigeNaam:            eenmanszaak.VolledigeNaam,
		Functionaris: models.Functionaris{
			TypeFunctionaris: "Eigenaar",
			Importance:       5,
		},
	}
	addInterpretatieNP(bevoegdheidUittreksel, &functionaris)

	bevoegdheidUittreksel.AlleFunctionarissen = append(bevoegdheidUittreksel.AlleFunctionarissen, functionaris)

	identityFunctionaris := models.IdentityNP{
		Geslachtsnaam:            eenmanszaak.Geslachtsnaam,
		Voornamen:                eenmanszaak.Voornamen,
		VoorvoegselGeslachtsnaam: eenmanszaak.VoorvoegselGeslachtsnaam,
		Geboortedatum:            convertDate(eenmanszaak.Geboortedatum),
	}

	if isSamePerson(identityNP, identityFunctionaris) {
		bevoegdheidUittreksel.MatchedFunctionaris = &functionaris

		paths.MatchedFunctionaris = models.NatuurlijkPersoonFunctionaris{
			Geslachtsnaam:            basePath + ".geslachtsnaam",
			Voornamen:                basePath + ".voornamen",
			VoorvoegselGeslachtsnaam: basePath + ".voorvoegselGeslachtsnaam",
			Geboortedatum:            basePath + ".geboortedatum",
			Overlijdensdatum:         basePath + ".overlijdensdatum",
			VolledigeNaam:            basePath + ".volledigeNaam",
			Functionaris: models.Functionaris{
				Handlichting: basePath + ".handlichting.isVerleend.code",
			},
		}
	}
	loopFunctionarissen(bevoegdheidUittreksel, paths, identityNP, eenmanszaak.Heeft, basePath)
}

func eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel *models.BevoegdheidUittreksel, paths *models.Paths, identityNP models.IdentityNP, nnp *models.NietNatuurlijkPersoon, typeEigenaar string, basePath string) {
	bevoegdheidUittreksel.TypeEigenaar = typeEigenaar
	bevoegdheidUittreksel.Rsin = nnp.Rsin
	bevoegdheidUittreksel.PersoonRechtsvorm = nnp.PersoonRechtsvorm
	bevoegdheidUittreksel.DatumUitschrijving = convertDate(nnp.DatumUitschrijving)
	bevoegdheidUittreksel.RegistratieEinde = convertDate(nnp.Registratie.DatumEinde)
	if nnp.BijzondereRechtstoestand.Soort.Code != "" {
		bevoegdheidUittreksel.BijzondereRechtstoestand = nnp.BijzondereRechtstoestand.Soort.Omschrijving + ":" + nnp.BijzondereRechtstoestand.Soort.Code
	}
	if nnp.BeperkingInRechtshandeling.Soort.Code != "" {
		bevoegdheidUittreksel.BeperkingInRechtshandeling = nnp.BeperkingInRechtshandeling.Soort.Omschrijving + ":" + nnp.BeperkingInRechtshandeling.Soort.Code
	}
	bevoegdheidUittreksel.BuitenlandseRechtstoestand = nnp.BuitenlandseRechtstoestand.Beschrijving

	paths.TypeEigenaar = basePath
	paths.Rsin = basePath + ".rsin"
	paths.PersoonRechtsvorm = basePath + ".persoonRechtsvorm"
	paths.DatumUitschrijving = basePath + ".datumUitschrijving"
	paths.RegistratieEinde = basePath + ".registratie.datumEinde"
	paths.BijzondereRechtstoestand = basePath + ".bijzondereRechtstoestand.soort"
	paths.BeperkingInRechtshandeling = basePath + ".beperkingInRechtshandeling.soort"
	paths.BuitenlandseRechtstoestand = basePath + ".buitenlandseRechtstoestand.beschrijving"

	loopFunctionarissen(bevoegdheidUittreksel, paths, identityNP, nnp.Heeft, basePath)
}

func formatPeilMoment(pm string) string {
	if pm == "" {
		return pm
	}
	if len(pm) != 17 {
		return pm
	}
	return pm[:4] + "-" + pm[4:6] + "-" + pm[6:8] + " " + pm[8:10] + ":" + pm[10:12] + ":" + pm[12:14] + "." + pm[14:]
}

// GetBevoegdheidUittreksel TODO
func getBevoegdheidUittreksel(bevoegdheidUittreksel *models.BevoegdheidUittreksel, paths *models.Paths, ophalenInschrijvingResponse *models.OphalenInschrijvingResponse, identityNP models.IdentityNP) {
	bevoegdheidUittreksel.Peilmoment = formatPeilMoment(ophalenInschrijvingResponse.Peilmoment)

	ma := ophalenInschrijvingResponse.Product.MaatschappelijkeActiviteit
	bevoegdheidUittreksel.KvkNummer = ma.KvkNummer
	bevoegdheidUittreksel.Naam = ma.Naam
	bevoegdheidUittreksel.Adres = ma.BezoekLocatie.VolledigAdres
	bevoegdheidUittreksel.RegistratieAanvang = convertDate(ma.Registratie.DatumAanvang)

	paths.KvkNummer = "maatschappelijkeActiviteit.kvkNummer"
	paths.Naam = "maatschappelijkeActiviteit.naam"
	paths.Adres = "maatschappelijkeActiviteit.bezoekLocatie.volledigAdres"
	paths.RegistratieAanvang = "maatschappelijkeActiviteit.registratie.datumAanvang"

	if len(ma.Communicatiegegevens.EmailAdres) != 0 {
		bevoegdheidUittreksel.EmailAdres = ma.Communicatiegegevens.EmailAdres[0]
	}
	paths.EmailAdres = "maatschappelijkeActiviteit.communicatiegegevens.emailAdres.0"

	for i, nr := range ma.Communicatiegegevens.Communicatienummer {
		if nr.Soort.Code == "T" {
			bevoegdheidUittreksel.Telefoon = nr.Toegangscode + " " + nr.Nummer[1:]
			paths.Telefoon = "maatschappelijkeActiviteit.communicatiegegevens.communicatienummer." + strconv.Itoa(i)
			break
		}
	}

	for i, sbi := range ma.SbiActiviteit {
		if sbi.IsHoofdactiviteit.Code == "J" {
			bevoegdheidUittreksel.SbiActiviteit = sbi.SbiCode.Code + ", " + sbi.SbiCode.Omschrijving
			paths.SbiActiviteit = "maatschappelijkeActiviteit.sbiActiviteit." + strconv.Itoa(i) + ".sbiCode"
			break
		}
	}

	if bevoegdheidUittreksel.SbiActiviteit == "" {
		for i, sbi := range ma.ManifesteertZichAls.Onderneming.SbiActiviteit {
			if sbi.IsHoofdactiviteit.Code == "J" {
				bevoegdheidUittreksel.SbiActiviteit = sbi.SbiCode.Code + ", " + sbi.SbiCode.Omschrijving
				paths.SbiActiviteit = "maatschappelijkeActiviteit.manifesteertZichAls.onderneming.sbiActiviteit." + strconv.Itoa(i) + ".sbiCode"
				break
			}
		}
	}

	for i, handeltOnder := range ma.ManifesteertZichAls.Onderneming.HandeltOnder {
		prefix := ", "
		if i == 0 {
			prefix = ""
		}
		bevoegdheidUittreksel.Handelsnamen = bevoegdheidUittreksel.Handelsnamen + prefix + handeltOnder.Handelsnaam.Naam
	}
	paths.Handelsnamen = "maatschappelijkeActiviteit.manifesteertZichAls.onderneming.handeltOnder"

	if ma.HeeftAlsEigenaar.Eenmanszaak != nil {
		eigenaarIsNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.Eenmanszaak, "NatuurlijkPersoon", "maatschappelijkeActiviteit.heeftAlsEigenaar.natuurlijkPersoon")
	} else if ma.HeeftAlsEigenaar.NaamPersoon != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.NaamPersoon, "NaamPersoon", "maatschappelijkeActiviteit.heeftAlsEigenaar.naamPersoon")
	} else if ma.HeeftAlsEigenaar.BuitenlandseVennootschap != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.BuitenlandseVennootschap, "BuitenlandseVennootschap", "maatschappelijkeActiviteit.heeftAlsEigenaar.buitenlandseVennootschap")
	} else if ma.HeeftAlsEigenaar.EenmanszaakMetMeerdereEigenaren != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.EenmanszaakMetMeerdereEigenaren, "EenmanszaakMetMeerdereEigenaren", "maatschappelijkeActiviteit.heeftAlsEigenaar.eenmanszaakMetMeerdereEigenaren")
	} else if ma.HeeftAlsEigenaar.Rechtspersoon != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.Rechtspersoon, "Rechtspersoon", "maatschappelijkeActiviteit.heeftAlsEigenaar.rechtspersoon")
	} else if ma.HeeftAlsEigenaar.RechtspersoonInOprichting != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.RechtspersoonInOprichting, "RechtspersoonInOprichting", "maatschappelijkeActiviteit.heeftAlsEigenaar.rechtspersoonInOprichting")
	} else if ma.HeeftAlsEigenaar.Samenwerkingsverband != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.Samenwerkingsverband, "Samenwerkingsverband", "maatschappelijkeActiviteit.heeftAlsEigenaar.samenwerkingsverband")
	} else if ma.HeeftAlsEigenaar.AfgeslotenMoeder != nil {
		eigenaarIsNietNatuurlijkPersoon(bevoegdheidUittreksel, paths, identityNP, ma.HeeftAlsEigenaar.AfgeslotenMoeder, "AfgeslotenMoeder", "maatschappelijkeActiviteit.heeftAlsEigenaar.afgeslotenMoeder")
	}
}
