package main

import (
	"net"
	"strings"
)

type Language string

const (
	languageEnglish Language = "en"
	languageCzech   Language = "cs"
	languageGerman  Language = "de"
)

const (
	errorInvalidInvitation = "invalid_invitation"
	errorAllergies         = "allergies"
	errorAccommodation     = "accommodation"
	errorMessage           = "message"
	errorSave              = "save"
)

type Translation struct {
	LogoAlt                  string
	WeddingNav               string
	MapTitle                 string
	LandingTitle             string
	RSVPTitle                string
	SuccessTitle             string
	Welcome                  string
	GettingMarried           string
	LandingIntro             string
	InvitationCode           string
	OpenInvitation           string
	CodeHint                 string
	ExploreDetails           string
	InvalidInvitation        string
	DateRange                string
	WeekendTitle             string
	WeekendIntro             string
	AddCalendar              string
	CalendarTitle            string
	CalendarDescription      string
	NavWeekend               string
	NavLocation              string
	NavDetails               string
	NavRSVP                  string
	Plan                     string
	FridayDate               string
	SaturdayDate             string
	SundayDate               string
	Arrival                  string
	Ceremony                 string
	LocalTime                string
	Departure                string
	FindUs                   string
	Location                 string
	Country                  string
	GoogleMaps               string
	FewDetails               string
	Parking                  string
	ParkingDescription       string
	Accommodation            string
	AccommodationDescription string
	Witnesses                string
	WitnessNames             string
	Invited                  string
	GreetingPrefix           string
	NameSeparator            string
	GreetingSuffix           string
	ReservedSingle           string
	ReservedDual             string
	ReservedMultiple         string
	GuestJoinPrompt          string
	Coming                   string
	Allergies                string
	AllergyPlaceholder       string
	AccommodationOwn         string
	AccommodationHelp        string
	MessageLabel             string
	Optional                 string
	MessagePlaceholder       string
	UpdateRSVP               string
	SendRSVP                 string
	InvitationLabel          string
	ErrorAllergies           string
	ErrorAccommodation       string
	ErrorMessage             string
	ErrorSave                string
	RSVPReceived             string
	ThankYouCouple           string
	ThankYou                 string
	AttendingSuccess         string
	DeclinedSuccess          string
	EditRSVP                 string
}

func languageForHost(host string) Language {
	if hostname, _, err := net.SplitHostPort(host); err == nil {
		host = hostname
	}

	host = strings.TrimSuffix(strings.ToLower(host), ".")
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return languageEnglish
	}

	switch labels[0] {
	case "cs", "cz":
		return languageCzech
	case "de":
		return languageGerman
	case "en":
		return languageEnglish
	default:
		return languageEnglish
	}
}

func translationsFor(language Language) Translation {
	if translation, ok := translations[language]; ok {
		return translation
	}
	return translations[languageEnglish]
}

var translations = map[Language]Translation{
	languageEnglish: {
		LogoAlt:                  "Animated Luky and Lelaina wedding logo",
		WeddingNav:               "Wedding information",
		MapTitle:                 "Map to Statek Újezd u Plánice",
		LandingTitle:             "Luky & Lelaina · Wedding",
		RSVPTitle:                "RSVP · Luky & Lelaina",
		SuccessTitle:             "RSVP received · Wedding",
		Welcome:                  "Welcome to our wedding",
		GettingMarried:           "We're getting married",
		LandingIntro:             "We'd love to celebrate this special day with you. Enter the unique code from your invitation to RSVP.",
		InvitationCode:           "Invitation code",
		OpenInvitation:           "Open my invitation",
		CodeHint:                 "Your code can be found on your invitation.",
		ExploreDetails:           "Explore the wedding details",
		InvalidInvitation:        "That invitation code was not found. Please check the code and try again.",
		DateRange:                "27–29 August 2027",
		WeekendTitle:             "A weekend together",
		WeekendIntro:             "We can't wait to celebrate with you in Újezd. The venue is ours from Friday through Sunday.",
		AddCalendar:              "Add to calendar",
		CalendarTitle:            "Luky and Lelaina's Wedding",
		CalendarDescription:      "Ceremony at 11:00 local time.",
		NavWeekend:               "The weekend",
		NavLocation:              "Location",
		NavDetails:               "Good to know",
		NavRSVP:                  "RSVP",
		Plan:                     "The plan",
		FridayDate:               "Friday · 27 August",
		SaturdayDate:             "Saturday · 28 August",
		SundayDate:               "Sunday · 29 August",
		Arrival:                  "Arrive and settle in",
		Ceremony:                 "Ceremony at 11:00",
		LocalTime:                "Local time",
		Departure:                "Departure",
		FindUs:                   "Find us",
		Location:                 "The location",
		Country:                  "Czech Republic",
		GoogleMaps:               "Open in Google Maps",
		FewDetails:               "A few details",
		Parking:                  "Parking",
		ParkingDescription:       "On-site parking is in the field behind the venue. Shoes comfortable on grass are a good idea.",
		Accommodation:            "Accommodation",
		AccommodationDescription: "Please let us know in your RSVP if you would like help arranging a place to stay.",
		Witnesses:                "Our witnesses",
		WitnessNames:             "Honzí and Anna",
		Invited:                  "You're invited",
		GreetingPrefix:           "Dear",
		NameSeparator:            " & ",
		GreetingSuffix:           ",",
		ReservedSingle:           "We have reserved %d seat for you.",
		ReservedDual:             "We have reserved %d seats for you.",
		ReservedMultiple:         "We have reserved %d seats for you.",
		GuestJoinPrompt:          "Let us know who can join us.",
		Coming:                   "Coming",
		Allergies:                "Allergies or dietary needs",
		AllergyPlaceholder:       "e.g. nuts, gluten...",
		AccommodationOwn:         "I'll arrange my own accommodation",
		AccommodationHelp:        "I'd appreciate help arranging accommodation",
		MessageLabel:             "A message for us",
		Optional:                 "(optional)",
		MessagePlaceholder:       "Leave us a little note...",
		UpdateRSVP:               "Update my RSVP",
		SendRSVP:                 "Send my RSVP",
		InvitationLabel:          "Invitation:",
		ErrorAllergies:           "Please keep each guest's allergy note reasonably short.",
		ErrorAccommodation:       "Please tell us how you would like to handle accommodation.",
		ErrorMessage:             "Please keep your message reasonably short.",
		ErrorSave:                "There was a problem saving your RSVP. Please try again.",
		RSVPReceived:             "RSVP received",
		ThankYouCouple:           "Thank you, %s!",
		ThankYou:                 "Thank you!",
		AttendingSuccess:         "We're so happy you'll be joining us. Your RSVP has been saved.",
		DeclinedSuccess:          "We're sorry you won't be able to make it, but thank you for letting us know.",
		EditRSVP:                 "Edit my RSVP",
	},
	languageCzech: {
		LogoAlt:                  "Animované svatební logo Lukyho a Lelainy",
		WeddingNav:               "Informace o svatbě",
		MapTitle:                 "Mapa ke Statku Újezd u Plánice",
		LandingTitle:             "Luky & Lelaina · Svatba",
		RSVPTitle:                "Potvrzení účasti · Luky & Lelaina",
		SuccessTitle:             "Odpověď přijata · Svatba",
		Welcome:                  "Vítejte na naší svatbě",
		GettingMarried:           "Bereme se",
		LandingIntro:             "Budeme rádi, když s námi oslavíte tento výjimečný den. Pro potvrzení účasti zadejte jedinečný kód z pozvánky.",
		InvitationCode:           "Kód z pozvánky",
		OpenInvitation:           "Otevřít pozvánku",
		CodeHint:                 "Kód najdete ve své pozvánce.",
		ExploreDetails:           "Podrobnosti o svatbě",
		InvalidInvitation:        "Tento kód pozvánky se nepodařilo najít. Zkontrolujte ho prosím.",
		DateRange:                "27.–29. srpna 2027",
		WeekendTitle:             "Společný víkend",
		WeekendIntro:             "Těšíme se, že s vámi oslavíme naši svatbu v Újezdě. Statek máme od pátku do neděle.",
		AddCalendar:              "Přidat do kalendáře",
		CalendarTitle:            "Svatba Lukyho a Lelainy",
		CalendarDescription:      "Obřad začíná v 11:00 místního času.",
		NavWeekend:               "Víkend",
		NavLocation:              "Místo",
		NavDetails:               "Praktické informace",
		NavRSVP:                  "Účast",
		Plan:                     "Program",
		FridayDate:               "Pátek · 27. srpna",
		SaturdayDate:             "Sobota · 28. srpna",
		SundayDate:               "Neděle · 29. srpna",
		Arrival:                  "Příjezd a ubytování",
		Ceremony:                 "Obřad v 11:00",
		LocalTime:                "místního času",
		Departure:                "Odjezd",
		FindUs:                   "Kde nás najdete",
		Location:                 "Místo konání",
		Country:                  "Česká republika",
		GoogleMaps:               "Otevřít v Mapách Google",
		FewDetails:               "Praktické informace",
		Parking:                  "Parkování",
		ParkingDescription:       "Parkovat lze na louce za statkem. Doporučujeme pohodlnou obuv vhodnou do trávy.",
		Accommodation:            "Ubytování",
		AccommodationDescription: "Pokud budete chtít pomoci se zajištěním ubytování, napište nám to prosím do odpovědi.",
		Witnesses:                "Svědkové",
		WitnessNames:             "Honzí a Anna",
		Invited:                  "Srdečně vás zveme",
		GreetingPrefix:           "Milí",
		NameSeparator:            " a ",
		GreetingSuffix:           "!",
		ReservedSingle:           "Máme pro vás rezervováno %d místo.",
		ReservedDual:             "Máme pro vás rezervována %d místa.",
		ReservedMultiple:         "Máme pro vás rezervováno %d míst.",
		GuestJoinPrompt:          "Dejte nám vědět, kdo se může přidat.",
		Coming:                   "Dorazím",
		Allergies:                "Alergie a stravovací omezení",
		AllergyPlaceholder:       "např. ořechy, lepek...",
		AccommodationOwn:         "Ubytování si zařídím sám/sama",
		AccommodationHelp:        "Uvítám pomoc se zajištěním ubytování",
		MessageLabel:             "Vzkaz pro nás",
		Optional:                 "(nepovinné)",
		MessagePlaceholder:       "Napište nám pár slov...",
		UpdateRSVP:               "Aktualizovat odpověď",
		SendRSVP:                 "Odeslat odpověď",
		InvitationLabel:          "Pozvánka:",
		ErrorAllergies:           "Prosíme o stručný popis stravovacích požadavků každého hosta.",
		ErrorAccommodation:       "Prosíme, vyberte, jak chcete řešit ubytování.",
		ErrorMessage:             "Prosíme, zkraťte svou zprávu.",
		ErrorSave:                "Odpověď se nepodařilo uložit. Zkuste to prosím znovu.",
		RSVPReceived:             "Odpověď přijata",
		ThankYouCouple:           "Děkujeme, %s!",
		ThankYou:                 "Děkujeme!",
		AttendingSuccess:         "Máme radost, že se k nám přidáte. Vaši odpověď jsme uložili.",
		DeclinedSuccess:          "Mrzí nás, že se nemůžete zúčastnit. Děkujeme, že jste nám dali vědět.",
		EditRSVP:                 "Upravit odpověď",
	},
	languageGerman: {
		LogoAlt:                  "Animiertes Hochzeitslogo von Luky und Lelaina",
		WeddingNav:               "Hochzeitsinformationen",
		MapTitle:                 "Karte zum Statek Újezd u Plánice",
		LandingTitle:             "Luky & Lelaina · Hochzeit",
		RSVPTitle:                "Rückmeldung · Luky & Lelaina",
		SuccessTitle:             "Rückmeldung erhalten · Hochzeit",
		Welcome:                  "Willkommen zu unserer Hochzeit",
		GettingMarried:           "Wir heiraten",
		LandingIntro:             "Wir würden diesen besonderen Tag gern mit euch feiern. Gebt den persönlichen Code aus eurer Einladung ein, um zu- oder abzusagen.",
		InvitationCode:           "Einladungscode",
		OpenInvitation:           "Einladung öffnen",
		CodeHint:                 "Den Code findet ihr auf eurer Einladung.",
		ExploreDetails:           "Hochzeitsdetails ansehen",
		InvalidInvitation:        "Dieser Einladungscode wurde nicht gefunden. Bitte überprüft den Code.",
		DateRange:                "27.–29. August 2027",
		WeekendTitle:             "Ein Wochenende zusammen",
		WeekendIntro:             "Wir freuen uns darauf, mit euch in Újezd zu feiern. Das Gelände steht uns von Freitag bis Sonntag zur Verfügung.",
		AddCalendar:              "Zum Kalender hinzufügen",
		CalendarTitle:            "Hochzeit von Luky und Lelaina",
		CalendarDescription:      "Die Trauung beginnt um 11:00 Uhr Ortszeit.",
		NavWeekend:               "Wochenende",
		NavLocation:              "Ort",
		NavDetails:               "Gut zu wissen",
		NavRSVP:                  "Rückmeldung",
		Plan:                     "Der Ablauf",
		FridayDate:               "Freitag · 27. August",
		SaturdayDate:             "Samstag · 28. August",
		SundayDate:               "Sonntag · 29. August",
		Arrival:                  "Ankunft und Einrichten",
		Ceremony:                 "Trauung um 11:00",
		LocalTime:                "Ortszeit",
		Departure:                "Abreise",
		FindUs:                   "Hier findet ihr uns",
		Location:                 "Der Ort",
		Country:                  "Tschechien",
		GoogleMaps:               "In Google Maps öffnen",
		FewDetails:               "Gut zu wissen",
		Parking:                  "Parken",
		ParkingDescription:       "Parkplätze gibt es auf der Wiese hinter der Location. Bequeme Schuhe für Gras sind empfehlenswert.",
		Accommodation:            "Unterkunft",
		AccommodationDescription: "Gebt uns bitte in eurer Rückmeldung Bescheid, wenn wir euch bei der Suche nach einer Unterkunft helfen sollen.",
		Witnesses:                "Unsere Trauzeugen",
		WitnessNames:             "Honzí und Anna",
		Invited:                  "Ihr seid eingeladen",
		GreetingPrefix:           "Liebe",
		NameSeparator:            " & ",
		GreetingSuffix:           ",",
		ReservedSingle:           "Wir haben %d Platz für euch reserviert.",
		ReservedDual:             "Wir haben %d Plätze für euch reserviert.",
		ReservedMultiple:         "Wir haben %d Plätze für euch reserviert.",
		GuestJoinPrompt:          "Gebt uns Bescheid, wer mit uns feiern kann.",
		Coming:                   "Ich komme",
		Allergies:                "Allergien oder Ernährungswünsche",
		AllergyPlaceholder:       "z. B. Nüsse, Gluten ...",
		AccommodationOwn:         "Ich kümmere mich selbst um eine Unterkunft",
		AccommodationHelp:        "Ich freue mich über Hilfe bei der Unterkunft",
		MessageLabel:             "Eine Nachricht an uns",
		Optional:                 "(optional)",
		MessagePlaceholder:       "Hinterlasst uns eine kleine Nachricht ...",
		UpdateRSVP:               "Rückmeldung aktualisieren",
		SendRSVP:                 "Rückmeldung senden",
		InvitationLabel:          "Einladung:",
		ErrorAllergies:           "Bitte haltet die Angaben zu Allergien und Ernährungswünschen kurz.",
		ErrorAccommodation:       "Bitte gebt an, wie ihr die Unterkunft organisieren möchtet.",
		ErrorMessage:             "Bitte haltet eure Nachricht kurz.",
		ErrorSave:                "Eure Rückmeldung konnte nicht gespeichert werden. Bitte versucht es erneut.",
		RSVPReceived:             "Rückmeldung erhalten",
		ThankYouCouple:           "Vielen Dank, %s!",
		ThankYou:                 "Vielen Dank!",
		AttendingSuccess:         "Wie schön, dass ihr dabei seid. Eure Rückmeldung wurde gespeichert.",
		DeclinedSuccess:          "Schade, dass ihr nicht dabei sein könnt. Danke, dass ihr Bescheid gegeben habt.",
		EditRSVP:                 "Rückmeldung bearbeiten",
	},
}
