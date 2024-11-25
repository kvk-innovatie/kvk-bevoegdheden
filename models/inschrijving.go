package models

import (
	"encoding/json"
)

type HrResponse struct {
	Envelope struct {
		Body struct {
			OphalenInschrijvingResponse OphalenInschrijvingResponse `json:"ophalenInschrijvingResponse"`
		} `json:"wstxns8:Body"`
	} `json:"soap:Envelope"`
}

type OphalenInschrijvingResponse struct {
	Peilmoment string     `json:"@peilmoment"`
	Meldingen  *Meldingen `json:"meldingen"`
	Product    struct {
		MaatschappelijkeActiviteit *MaatschappelijkeActiviteit `json:"maatschappelijkeActiviteit"`
	} `json:"product"`
}

type Meldingen struct {
	Informatie *Enumeratie `json:"informatie"`
	Fout       *Enumeratie `json:"fout"`
}

type MaatschappelijkeActiviteit struct {
	KvkNummer     string      `json:"kvkNummer"`
	Naam          string      `json:"naam"`
	Registratie   Registratie `json:"registratie"`
	BezoekLocatie struct {
		VolledigAdres string `json:"volledigAdres"`
	} `json:"bezoekLocatie"`
	Communicatiegegevens Communicatiegegevens `json:"communicatiegegevens"`
	SbiActiviteit        []SbiActiviteit      `json:"sbiActiviteit"`
	ManifesteertZichAls  struct {
		Onderneming Onderneming `json:"onderneming"`
	} `json:"manifesteertZichAls"`
	HeeftAlsEigenaar struct {
		NaamPersoon                     *NietNatuurlijkPersoon `json:"naamPersoon,omitempty"`
		Eenmanszaak                     *Eenmanszaak           `json:"natuurlijkPersoon,omitempty"`
		BuitenlandseVennootschap        *NietNatuurlijkPersoon `json:"buitenlandseVennootschap,omitempty"`
		EenmanszaakMetMeerdereEigenaren *NietNatuurlijkPersoon `json:"eenmanszaakMetMeerdereEigenaren,omitempty"`
		Rechtspersoon                   *NietNatuurlijkPersoon `json:"rechtspersoon,omitempty"`
		RechtspersoonInOprichting       *NietNatuurlijkPersoon `json:"rechtspersoonInOprichting,omitempty"`
		Samenwerkingsverband            *NietNatuurlijkPersoon `json:"samenwerkingsverband,omitempty"`
		AfgeslotenMoeder                *NietNatuurlijkPersoon `json:"afgeslotenMoeder,omitempty"`
	} `json:"heeftAlsEigenaar"`
}

func (u *MaatschappelijkeActiviteit) UnmarshalJSON(data []byte) error {
	type Alias MaatschappelijkeActiviteit
	aux := &struct {
		SbiActiviteit []SbiActiviteit `json:"sbiActiviteit"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		SbiActiviteit SbiActiviteit `json:"sbiActiviteit"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.SbiActiviteit = aux.SbiActiviteit
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.SbiActiviteit = []SbiActiviteit{
		auxOne.SbiActiviteit,
	}
	return nil
}

type Onderneming struct {
	SbiActiviteit []SbiActiviteit `json:"sbiActiviteit"`
	HandeltOnder  []HandeltOnder  `json:"handeltOnder"`
}

func (u *Onderneming) UnmarshalJSON(data []byte) error {
	type Alias Onderneming
	aux := &struct {
		SbiActiviteit []SbiActiviteit `json:"sbiActiviteit"`
		HandeltOnder  []HandeltOnder  `json:"handeltOnder"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOneC := &struct {
		SbiActiviteit SbiActiviteit  `json:"sbiActiviteit"`
		HandeltOnder  []HandeltOnder `json:"handeltOnder"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOneE := &struct {
		SbiActiviteit []SbiActiviteit `json:"sbiActiviteit"`
		HandeltOnder  HandeltOnder    `json:"handeltOnder"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		SbiActiviteit SbiActiviteit `json:"sbiActiviteit"`
		HandeltOnder  HandeltOnder  `json:"handeltOnder"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.SbiActiviteit = aux.SbiActiviteit
		u.HandeltOnder = aux.HandeltOnder
		return nil
	}
	err = json.Unmarshal(data, &auxOneC)
	if err == nil {
		u.SbiActiviteit = []SbiActiviteit{
			auxOneC.SbiActiviteit,
		}
		u.HandeltOnder = auxOneC.HandeltOnder
		return nil
	}
	err = json.Unmarshal(data, &auxOneE)
	if err == nil {
		u.SbiActiviteit = auxOneE.SbiActiviteit
		u.HandeltOnder = []HandeltOnder{
			auxOneE.HandeltOnder,
		}
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.SbiActiviteit = []SbiActiviteit{
		auxOne.SbiActiviteit,
	}
	u.HandeltOnder = []HandeltOnder{
		auxOne.HandeltOnder,
	}
	return nil
}

type SbiActiviteit struct {
	SbiCode struct {
		Code         string `json:"code"`
		Omschrijving string `json:"omschrijving"`
	} `json:"sbiCode"`
	IsHoofdactiviteit Enumeratie `json:"isHoofdactiviteit"`
}

type HandeltOnder struct {
	Handelsnaam struct {
		Naam string `json:"naam"`
	} `json:"handelsnaam"`
}

type Communicatiegegevens struct {
	EmailAdres         []string             `json:"emailAdres"`
	Communicatienummer []Communicatienummer `json:"communicatienummer"`
}

func (u *Communicatiegegevens) UnmarshalJSON(data []byte) error {
	type Alias Communicatiegegevens
	aux := &struct {
		Communicatienummer []Communicatienummer `json:"communicatienummer"`
		EmailAdres         []string             `json:"emailAdres"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOneC := &struct {
		Communicatienummer Communicatienummer `json:"communicatienummer"`
		EmailAdres         []string           `json:"emailAdres"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOneE := &struct {
		Communicatienummer []Communicatienummer `json:"communicatienummer"`
		EmailAdres         string               `json:"emailAdres"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		Communicatienummer Communicatienummer `json:"communicatienummer"`
		EmailAdres         string             `json:"emailAdres"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.Communicatienummer = aux.Communicatienummer
		u.EmailAdres = aux.EmailAdres
		return nil
	}
	err = json.Unmarshal(data, &auxOneC)
	if err == nil {
		u.Communicatienummer = []Communicatienummer{
			auxOneC.Communicatienummer,
		}
		u.EmailAdres = auxOneC.EmailAdres
		return nil
	}
	err = json.Unmarshal(data, &auxOneE)
	if err == nil {
		u.Communicatienummer = auxOneE.Communicatienummer
		u.EmailAdres = []string{
			auxOneE.EmailAdres,
		}
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.Communicatienummer = []Communicatienummer{
		auxOne.Communicatienummer,
	}
	u.EmailAdres = []string{
		auxOne.EmailAdres,
	}
	return nil
}

type Communicatienummer struct {
	Toegangscode string     `json:"toegangscode"`
	Nummer       string     `json:"nummer"`
	Soort        Enumeratie `json:"soort"`
}

type Eenmanszaak struct { // in KVK productstore this is 'natuurlijkPersoon' but that conflicts with NatuurlijkPersoon
	Registratie                Registratie                `json:"registratie"`
	PersoonRechtsvorm          string                     `json:"persoonRechtsvorm"`
	Geslachtsnaam              string                     `json:"geslachtsnaam"`
	Voornamen                  string                     `json:"voornamen"`
	VoorvoegselGeslachtsnaam   string                     `json:"voorvoegselGeslachtsnaam"`
	Geboortedatum              string                     `json:"geboortedatum"`
	Overlijdensdatum           string                     `json:"overlijdensdatum"`
	VolledigeNaam              string                     `json:"volledigeNaam"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
	Handlichting               Handlichting               `json:"handlichting"`
	Heeft                      []Functievervulling        `json:"heeft,omitempty"`
}

func (u *Eenmanszaak) UnmarshalJSON(data []byte) error {
	type Alias Eenmanszaak
	aux := &struct {
		Heeft []Functievervulling `json:"heeft"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		Heeft Functievervulling `json:"heeft"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.Heeft = aux.Heeft
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.Heeft = []Functievervulling{
		auxOne.Heeft,
	}
	return nil
}

type NietNatuurlijkPersoon struct {
	Rsin                       string                     `json:"rsin"`
	Registratie                Registratie                `json:"registratie"`
	DatumUitschrijving         string                     `json:"datumUitschrijving"`
	PersoonRechtsvorm          string                     `json:"persoonRechtsvorm"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
	Ontbinding                 Ontbinding                 `json:"ontbinding"`
	Heeft                      []Functievervulling        `json:"heeft,omitempty"`
	// LandVanVestiging           Enumeratie                  `json:"landVanVestiging"`
}

func (u *NietNatuurlijkPersoon) UnmarshalJSON(data []byte) error {
	type Alias NietNatuurlijkPersoon
	aux := &struct {
		Heeft []Functievervulling `json:"heeft"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		Heeft Functievervulling `json:"heeft"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.Heeft = aux.Heeft
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.Heeft = []Functievervulling{
		auxOne.Heeft,
	}
	return nil
}

// type NaamPersoon struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

// type BuitenlandseVennootschap struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// 	// LandVanVestiging           Enumeratie                  `json:"landVanVestiging"`
// }

// type EenmanszaakMetMeerdereEigenaren struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

// type Rechtspersoon struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

// type RechtspersoonInOprichting struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

// type Samenwerkingsverband struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

// type AfgeslotenMoeder struct {
// 	Registratie                Registratie                 `json:"registratie"`
// 	DatumUitschrijving         string                      `json:"datumUitschrijving"`
// 	PersoonRechtsvorm          string                      `json:"persoonRechtsvorm"`
// 	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
// 	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
// 	BuitenlandseRechtstoestand BuitenlandseRechtstoestand `json:"buitenlandseRechtstoestand"`
// 	Ontbinding                 Ontbinding                `json:"ontbinding"`
// 	Heeft                      []FunctieVervulling         `json:"heeft"`
// }

type Functievervulling struct {
	Aansprakelijke                       *FunctionarisOfGemachtigde `json:"aansprakelijke,omitempty"`
	Bestuursfunctie                      *FunctionarisOfGemachtigde `json:"bestuursfunctie,omitempty"`
	FunctionarisBijzondereRechtstoestand *FunctionarisOfGemachtigde `json:"functionarisBijzondereRechtstoestand,omitempty"`
	Gemachtigde                          *FunctionarisOfGemachtigde `json:"gemachtigde,omitempty"`
	OverigeFunctionaris                  *FunctionarisOfGemachtigde `json:"overigeFunctionaris,omitempty"`
	PubliekrechtelijkeFunctionaris       *FunctionarisOfGemachtigde `json:"publiekrechtelijkeFunctionaris,omitempty"`
}

type Door struct {
	NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon,omitempty"`
	Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon,omitempty"`
}

type FunctionarisOfGemachtigde struct {
	Functie      Enumeratie `json:"functie"`
	Functietitel struct {
		Titel string `json:"titel"`
	} `json:"functietitel"`
	Bevoegdheid  Bevoegdheid  `json:"bevoegdheid"`
	Volmacht     Volmacht     `json:"volmacht"`
	Handlichting Handlichting `json:"handlichting"`
	Schorsing    Schorsing    `json:"schorsing"`
	Door         Door         `json:"door"`
}

// type Gemachtigde struct {
// 	Functie      Enumeratie `json:"functie"`
// 	Functietitel struct {
// 		Titel string `json:"titel"`
// 	} `json:"functietitel"`
// 	Volmacht  Volmacht  `json:"volmacht"`
// 	Schorsing Schorsing `json:"schorsing"`
// 	Door      Door      `json:"door"`
// }

// type Aansprakelijke struct {
// 	Functie      Enumeratie    `json:"functie"`
// 	Bevoegdheid  Bevoegdheid   `json:"bevoegdheid"`
// 	Handlichting Handlichting `json:"handlichting"`
// 	Schorsing    Schorsing    `json:"schorsing"`
// 	Door         struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon"`
// 	} `json:"door"`
// }

// type Bestuursfunctie struct {
// 	Functie      Enumeratie `json:"functie"`
// 	Functietitel struct {
// 		Titel string `json:"titel"`
// 	} `json:"functietitel"`
// 	Bevoegdheid Bevoegdheid `json:"bevoegdheid"`
// 	Schorsing   Schorsing  `json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon"`
// 	} `json:"door"`
// }

// type FunctionarisBijzondereRechtstoestand struct {
// 	Functie   Enumeratie `json:"functie"`
// 	Schorsing Schorsing `json:"schorsing"`
// 	Door      struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon"`
// 	} `json:"door"`
// }

// type OverigeFunctionaris struct {
// 	Functie     Enumeratie  `json:"functie"`
// 	Bevoegdheid Bevoegdheid `json:"bevoegdheid"`
// 	Schorsing   Schorsing  `json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon"`
// 	} `json:"door"`
// }

// type PubliekrechtelijkeFunctionaris struct {
// 	Functie     Enumeratie  `json:"functie"`
// 	Bevoegdheid Bevoegdheid `json:"bevoegdheid"`
// 	Schorsing   Schorsing  `json:"schorsing"`
// 	Door        struct {
// 		NatuurlijkPersoon *NatuurlijkPersoon            `json:"natuurlijkPersoon"`
// 		Rechtspersoon     *RechtspersoonAlsFunctionaris `json:"rechtspersoon"`
// 	} `json:"door"`
// }

type Bevoegdheid struct {
	Soort            Enumeratie `json:"soort"`
	BeperkingInEuros struct {
		Waarde string     `json:"waarde"`
		Valuta Enumeratie `json:"valuta"`
	} `json:"beperkingInEuros"`
	OverigeBeperking           Enumeratie `json:"overigeBeperking"`
	IsBevoegdMetAnderePersonen Enumeratie `json:"isBevoegdMetAnderePersonen"`
}

type Volmacht struct {
	TypeVolmacht     Enumeratie       `json:"typeVolmacht"`
	BeperkteVolmacht BeperkteVolmacht `json:"beperkteVolmacht"`
}

type BeperkteVolmacht struct {
	BeperkingInHandeling []BeperkingInHandeling `json:"beperkingInHandeling"`
	BeperkingInGeld      struct {
		Waarde string     `json:"waarde"`
		Valuta Enumeratie `json:"valuta"`
	} `json:"beperkingInGeld"`
	MagOpgaveHandelsregisterDoen Enumeratie `json:"magOpgaveHandelsregisterDoen"`
	HeeftOverigeVolmacht         Enumeratie `json:"heeftOverigeVolmacht"`
	OmschrijvingOverigeVolmacht  string     `json:"omschrijvingOverigeVolmacht"`
}

type BeperkingInHandeling struct {
	SoortHandeling Enumeratie `json:"soortHandeling"`
}

func (u *BeperkteVolmacht) UnmarshalJSON(data []byte) error {
	type Alias BeperkteVolmacht
	aux := &struct {
		BeperkingInHandeling []BeperkingInHandeling `json:"beperkingInHandeling"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	auxOne := &struct {
		BeperkingInHandeling BeperkingInHandeling `json:"beperkingInHandeling"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	err := json.Unmarshal(data, &aux)
	if err == nil {
		u.BeperkingInHandeling = aux.BeperkingInHandeling
		return nil
	}
	err = json.Unmarshal(data, &auxOne)
	if err != nil {
		return err
	}
	u.BeperkingInHandeling = []BeperkingInHandeling{
		auxOne.BeperkingInHandeling,
	}
	return nil
}

type NatuurlijkPersoon struct {
	Geslachtsnaam              string                     `json:"geslachtsnaam"`
	VoorvoegselGeslachtsnaam   string                     `json:"voorvoegselGeslachtsnaam"`
	Voornamen                  string                     `json:"voornamen"`
	Geboortedatum              string                     `json:"geboortedatum"`
	Overlijdensdatum           string                     `json:"overlijdensdatum"`
	VolledigeNaam              string                     `json:"volledigeNaam"`
	BijzondereRechtstoestand   BijzondereRechtstoestand   `json:"bijzondereRechtstoestand"`
	BeperkingInRechtshandeling BeperkingInRechtshandeling `json:"beperkingInRechtshandeling"`
}

type RechtspersoonAlsFunctionaris struct {
	PersoonRechtsvorm string `json:"persoonRechtsvorm"`
	VolledigeNaam     string `json:"volledigeNaam"`
	IsEigenaarVan     struct {
		MaatschappelijkeActiviteit struct {
			KvkNummer string `json:"kvkNummer"`
		} `json:"maatschappelijkeActiviteit"`
	} `json:"isEigenaarVan"`
}

type BijzondereRechtstoestand struct {
	Registratie Registratie `json:"registratie"`
	Soort       Enumeratie  `json:"soort"`
}

type BeperkingInRechtshandeling struct {
	Registratie Registratie `json:"registratie"`
	Soort       Enumeratie  `json:"soort"`
}

type BuitenlandseRechtstoestand struct {
	Registratie  Registratie `json:"registratie"`
	Beschrijving string      `json:"beschrijving"`
}

type Handlichting struct {
	Registratie Registratie `json:"registratie"`
	IsVerleend  Enumeratie  `json:"isVerleend"`
}

type Ontbinding struct {
	Registratie Registratie `json:"registratie"`
	Aanleiding  Enumeratie  `json:"aanleiding"`
	Liquidatie  struct {
		Registratie Registratie `json:"registratie"`
	} `json:"liquidatie"`
}

type Schorsing struct {
	Registratie Registratie `json:"registratie"`
}

type Registratie struct {
	RegistratieTijdstip string `json:"registratieTijdstip"`
	DatumAanvang        string `json:"datumAanvang"`
	DatumEinde          string `json:"datumEinde"`
}

type Enumeratie struct {
	Code           string `json:"code"`
	Omschrijving   string `json:"omschrijving"`
	ReferentieType string `json:"referentieType"`
}
