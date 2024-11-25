package models

type IdentityNP struct {
	Geslachtsnaam            string `json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam string `json:"voorvoegselGeslachtsnaam"`
	Voornamen                string `json:"voornamen"`
	Geboortedatum            string `json:"geboortedatum"`
	VolledigeNaam            string `json:"volledigeNaam"`
}

type Interpretatie struct {
	HeeftBeperking string `json:"heeftBeperking"`
	IsBevoegd      string `json:"isBevoegd"`
	Reden          string `json:"reden"`
}

type BevoegdheidResponse struct {
	BevoegdheidUittreksel *BevoegdheidUittreksel      `json:"bevoegdheidUittreksel"`
	Inschrijving          *MaatschappelijkeActiviteit `json:"-"`
	Paths                 *Paths                      `json:"-"`
}

type BevoegdheidUittreksel struct {
	KvkNummer          string `json:"kvkNummer"`
	Rsin               string `json:"rsin"`
	Naam               string `json:"naam"`
	Handelsnamen       string `json:"handelsnamen"`
	TypeEigenaar       string `json:"typeEigenaar"`
	PersoonRechtsvorm  string `json:"persoonRechtsvorm"`
	Adres              string `json:"adres"`
	EmailAdres         string `json:"emailAdres"`
	Telefoon           string `json:"telefoon"`
	SbiActiviteit      string `json:"sbiActiviteit"`
	RegistratieAanvang string `json:"registratieAanvang"`

	DatumUitschrijving         string `json:"datumUitschrijving"`
	RegistratieEinde           string `json:"registratieEinde"`
	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand string `json:"buitenlandseRechtstoestand"`

	Peilmoment string `json:"peilmoment"`

	MatchedFunctionaris              *NatuurlijkPersoonFunctionaris  `json:"matchedFunctionaris,omitempty"`
	AlleFunctionarissen              []NatuurlijkPersoonFunctionaris `json:"functionarissen"`
	AlleRechtspersoonFunctionarissen []RechtspersoonFunctionaris     `json:"rechtspersoonFunctionarissen"`
}

type Functionaris struct {
	TypeFunctionaris string `json:"typeFunctionaris"`
	Functie          string `json:"functie"`
	Functietitel     string `json:"functietitel"`

	SoortBevoegdheid            string `json:"soortBevoegdheid"`
	BeperkingInEurosBevoegdheid string `json:"beperkingInEurosBevoegdheid"`
	OverigeBeperkingBevoegdheid string `json:"overigeBeperkingBevoegdheid"`
	IsBevoegdMetAnderePersonen  string `json:"isBevoegdMetAnderePersonen"`

	TypeVolmacht                 string `json:"typeVolmacht"`
	BeperkingInGeldVolmacht      string `json:"beperkingInGeldVolmacht"`
	BeperkingInHandelingVolmacht string `json:"beperkingInHandelingVolmacht"`
	HeeftOverigeVolmacht         string `json:"heeftOverigeVolmacht"`
	OmschrijvingOverigeVolmacht  string `json:"omschrijvingOverigeVolmacht"`
	MagOpgaveHandelsregisterDoen string `json:"magOpgaveHandelsregisterDoen"`

	SchorsingAanvang string `json:"schorsingAanvang"`
	SchorsingEinde   string `json:"schorsingEinde"`
	Handlichting     string `json:"handlichting"`

	Interpretatie Interpretatie `json:"interpretatie"`

	Importance int64 `json:"importance"`
}

type NatuurlijkPersoonFunctionaris struct {
	Geslachtsnaam            string `json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam string `json:"voorvoegselGeslachtsnaam"`
	Voornamen                string `json:"voornamen"`
	Geboortedatum            string `json:"geboortedatum"`
	Overlijdensdatum         string `json:"overlijdensdatum"`
	VolledigeNaam            string `json:"volledigeNaam"`

	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestandFunctionaris"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandelingFunctionaris"`

	Functionaris
}

type RechtspersoonFunctionaris struct {
	KvkNummer         string `json:"kvkNummer"`
	PersoonRechtsvorm string `json:"persoonRechtsvorm"`
	Naam              string `json:"naam"`

	Functionaris
}

type Paths struct {
	KvkNummer          string `json:"kvkNummer"`
	Rsin               string `json:"rsin"`
	Naam               string `json:"naam"`
	Handelsnamen       string `json:"handelsnamen"`
	TypeEigenaar       string `json:"typeEigenaar"`
	PersoonRechtsvorm  string `json:"persoonRechtsvorm"`
	Adres              string `json:"adres"`
	EmailAdres         string `json:"emailAdres"`
	Telefoon           string `json:"telefoon"`
	SbiActiviteit      string `json:"sbiActiviteit"`
	RegistratieAanvang string `json:"registratieAanvang"`

	DatumUitschrijving         string `json:"datumUitschrijving"`
	RegistratieEinde           string `json:"registratieEinde"`
	BijzondereRechtstoestand   string `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling string `json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand string `json:"buitenlandseRechtstoestand"`

	MatchedFunctionaris NatuurlijkPersoonFunctionaris `json:"matchedFunctionaris,omitempty"`
}
